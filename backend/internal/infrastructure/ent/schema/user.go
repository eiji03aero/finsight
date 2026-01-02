package schema

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Unique().
			NotEmpty().
			Validate(func(s string) error {
				emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !emailRegex.MatchString(s) {
					return fmt.Errorf("invalid email format")
				}
				return nil
			}),
		field.String("password_hash").
			NotEmpty().
			Sensitive(), // Prevents password hash from being logged
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("workspaces", Workspace.Type).
			Ref("users"),
	}
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		// Hook to normalize email to lowercase
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if email, exists := m.Field("email"); exists {
					emailStr := email.(string)
					if err := m.SetField("email", strings.ToLower(emailStr)); err != nil {
						return nil, err
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
