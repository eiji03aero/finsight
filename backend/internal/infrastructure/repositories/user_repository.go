package repositories

import (
	"context"
	"strings"

	"backend/internal/domain/model"
	"backend/internal/infrastructure/ent"
	"backend/internal/infrastructure/ent/user"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*model.User, error) {
	entUser, err := r.client.User.
		Create().
		SetEmail(email).
		SetPasswordHash(passwordHash).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toUserModel(entUser), nil
}

// GetUserByEmail retrieves a user by email (case-insensitive)
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(user.Email(strings.ToLower(email))).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toUserModel(entUser), nil
}

// EmailExists checks if an email already exists (case-insensitive)
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return r.client.User.
		Query().
		Where(user.Email(strings.ToLower(email))).
		Exist(ctx)
}

// toModel converts ent.User to domain model User
func toUserModel(entUser *ent.User) *model.User {
	return &model.User{
		ID:           entUser.ID,
		Email:        entUser.Email,
		PasswordHash: entUser.PasswordHash,
		CreatedAt:    entUser.CreatedAt,
		UpdatedAt:    entUser.UpdatedAt,
	}
}
