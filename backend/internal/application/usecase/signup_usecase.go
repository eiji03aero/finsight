package usecase

import (
	"context"
	"fmt"

	"backend/internal/domain/model"
	"backend/internal/domain/service"
	"backend/internal/infrastructure/ent"
	"backend/internal/infrastructure/repositories"
)

type SignupUseCase struct {
	userRepo      *repositories.UserRepository
	workspaceRepo *repositories.WorkspaceRepository
	client        *ent.Client
}

func NewSignupUseCase(
	userRepo *repositories.UserRepository,
	workspaceRepo *repositories.WorkspaceRepository,
	client *ent.Client,
) *SignupUseCase {
	return &SignupUseCase{
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		client:        client,
	}
}

// Execute orchestrates the signup process
func (uc *SignupUseCase) Execute(
	ctx context.Context,
	email string,
	password string,
	workspaceName string,
) (*model.User, *model.Workspace, error) {
	// Validate email format
	if err := service.ValidateEmail(email); err != nil {
		return nil, nil, fmt.Errorf("email validation failed: %w", err)
	}

	// Validate password length
	if len(password) < 8 {
		return nil, nil, fmt.Errorf("password must be at least 8 characters")
	}

	// Validate workspace name
	if workspaceName == "" {
		return nil, nil, fmt.Errorf("workspace name is required")
	}

	// Check if email already exists
	exists, err := uc.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, nil, fmt.Errorf("email already registered")
	}

	// Hash password
	passwordHash, err := service.HashPassword(password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user and workspace in a transaction
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Rollback helper
	rollback := func(tx *ent.Tx, err error) error {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v, original error: %w", rbErr, err)
		}
		return err
	}

	// Create user
	entUser, err := tx.User.
		Create().
		SetEmail(email).
		SetPasswordHash(passwordHash).
		Save(ctx)
	if err != nil {
		return nil, nil, rollback(tx, fmt.Errorf("failed to create user: %w", err))
	}

	// Create workspace and link to user
	entWorkspace, err := tx.Workspace.
		Create().
		SetName(workspaceName).
		AddUsers(entUser).
		Save(ctx)
	if err != nil {
		return nil, nil, rollback(tx, fmt.Errorf("failed to create workspace: %w", err))
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to domain models
	user := &model.User{
		ID:           entUser.ID,
		Email:        entUser.Email,
		PasswordHash: entUser.PasswordHash,
		CreatedAt:    entUser.CreatedAt,
		UpdatedAt:    entUser.UpdatedAt,
	}

	workspace := &model.Workspace{
		ID:        entWorkspace.ID,
		Name:      entWorkspace.Name,
		CreatedAt: entWorkspace.CreatedAt,
		UpdatedAt: entWorkspace.UpdatedAt,
	}

	return user, workspace, nil
}
