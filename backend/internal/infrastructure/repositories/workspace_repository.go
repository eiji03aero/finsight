package repositories

import (
	"context"

	"backend/internal/domain/model"
	"backend/internal/infrastructure/ent"
)

type WorkspaceRepository struct {
	client *ent.Client
}

func NewWorkspaceRepository(client *ent.Client) *WorkspaceRepository {
	return &WorkspaceRepository{client: client}
}

// CreateWorkspace creates a new workspace
func (r *WorkspaceRepository) CreateWorkspace(ctx context.Context, name string) (*model.Workspace, error) {
	entWorkspace, err := r.client.Workspace.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toWorkspaceModel(entWorkspace), nil
}

// toModel converts ent.Workspace to domain model Workspace
func toWorkspaceModel(entWorkspace *ent.Workspace) *model.Workspace {
	return &model.Workspace{
		ID:        entWorkspace.ID,
		Name:      entWorkspace.Name,
		CreatedAt: entWorkspace.CreatedAt,
		UpdatedAt: entWorkspace.UpdatedAt,
	}
}
