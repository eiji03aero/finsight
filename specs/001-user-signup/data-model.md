# Data Model: User Signup

**Feature**: User Signup
**Date**: 2025-12-27
**ORM**: Ent (https://entgo.io/)

## Overview

This document defines the data model for the user signup feature using Ent schema definitions. It covers three main entities: User, Workspace, and their many-to-many relationship.

## Entity Definitions

### 1. User Entity

**Purpose**: Represents a user account in the system.

**Ent Schema Location**: `backend/internal/infrastructure/ent/schema/user.go`

**Fields**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | `int` | Primary Key, Auto-increment | Unique identifier for the user |
| `email` | `string` | Unique, Not Empty, Lowercase | User's email address (used for login) |
| `password_hash` | `string` | Not Empty | Bcrypt-hashed password |
| `created_at` | `time.Time` | Auto-set on create | Timestamp of account creation |
| `updated_at` | `time.Time` | Auto-update | Timestamp of last update |

**Edges** (Relationships):

| Edge Name | Type | Target Entity | Description |
|-----------|------|---------------|-------------|
| `workspaces` | Many-to-Many | `Workspace` | Workspaces associated with this user |

**Indexes**:
- Unique index on `email` (enforced by Ent's `Unique()` constraint)

**Hooks**:
- `BeforeCreate`: Normalize email to lowercase
- `BeforeUpdate`: Normalize email to lowercase (if email is being updated)

**Validation Rules**:
- Email must match regex pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
- Password hash must not be empty (passwords are hashed before storage)

**Ent Schema Example**:
```go
// Fields returns the fields of the User schema.
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

// Edges returns the edges of the User schema.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("workspaces", Workspace.Type).
            Ref("users"),
    }
}

// Hooks returns the hooks of the User schema.
func (User) Hooks() []ent.Hook {
    return []ent.Hook{
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                    if email, exists := m.Field("email"); exists {
                        emailStr := email.(string)
                        m.SetField("email", strings.ToLower(emailStr))
                    }
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
        ),
    }
}
```

---

### 2. Workspace Entity

**Purpose**: Represents a workspace (isolated environment for users).

**Ent Schema Location**: `backend/internal/infrastructure/ent/schema/workspace.go`

**Fields**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | `int` | Primary Key, Auto-increment | Unique identifier for the workspace |
| `name` | `string` | Not Empty | Workspace name (user-provided) |
| `created_at` | `time.Time` | Auto-set on create | Timestamp of workspace creation |
| `updated_at` | `time.Time` | Auto-update | Timestamp of last update |

**Edges** (Relationships):

| Edge Name | Type | Target Entity | Description |
|-----------|------|---------------|-------------|
| `users` | Many-to-Many | `User` | Users who have access to this workspace |

**Indexes**:
- None required (workspace names are not unique across the system)

**Validation Rules**:
- Name must not be empty
- Name can contain any characters (including special characters, emojis, etc.)

**Ent Schema Example**:
```go
// Fields returns the fields of the Workspace schema.
func (Workspace) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            NotEmpty(),
        field.Time("created_at").
            Default(time.Now).
            Immutable(),
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now),
    }
}

// Edges returns the edges of the Workspace schema.
func (Workspace) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("users", User.Type),
    }
}
```

---

### 3. User-Workspace Relationship

**Type**: Many-to-Many

**Implementation**: Ent automatically creates a join table for many-to-many relationships.

**Join Table Details** (auto-generated by Ent):

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `user_id` | `int` | Foreign Key (users.id), Not Null | Reference to user |
| `workspace_id` | `int` | Foreign Key (workspaces.id), Not Null | Reference to workspace |

**Primary Key**: Composite key on (`user_id`, `workspace_id`)

**Indexes**:
- Index on `user_id` (for efficient user → workspaces queries)
- Index on `workspace_id` (for efficient workspace → users queries)

**Notes**:
- Ent handles the join table automatically via the edge definitions
- For signup flow: one user is linked to one workspace initially
- Future: can support multiple users per workspace and multiple workspaces per user

---

## Database Schema (Generated SQL)

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for email uniqueness (auto-created by UNIQUE constraint)
CREATE UNIQUE INDEX users_email_key ON users (email);
```

### Workspaces Table

```sql
CREATE TABLE workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### User-Workspace Join Table

```sql
CREATE TABLE user_workspaces (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workspace_id INTEGER NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, workspace_id)
);

-- Indexes for efficient lookups
CREATE INDEX user_workspaces_user_id_idx ON user_workspaces (user_id);
CREATE INDEX user_workspaces_workspace_id_idx ON user_workspaces (workspace_id);
```

---

## State Transitions

### User Entity States

```
[Non-existent]
    → (Signup) → [Active User]
```

**Notes**:
- For MVP: no account activation, suspension, or deletion states
- Future: can add `status` field with values: `active`, `suspended`, `deleted`

### Workspace Entity States

```
[Non-existent]
    → (User Signup) → [Active Workspace]
```

**Notes**:
- Workspace is created simultaneously with user during signup
- For MVP: no workspace archival or deletion
- Future: can add `status` field or `deleted_at` for soft deletes

---

## Relationships & Cardinality

```
User ◄──────────► Workspace
     (many-to-many)

- One User can have many Workspaces
- One Workspace can have many Users
- For signup: One User is initially linked to One Workspace
```

**Example Scenarios**:

1. **Signup Flow**:
   - User signs up → Creates 1 User + 1 Workspace + 1 link

2. **Future Multi-Workspace**:
   - User creates additional workspace → Creates 1 Workspace + 1 link
   - User invited to workspace → Creates 1 link (no new User or Workspace)

---

## Data Validation Summary

### Backend (Ent)

| Field | Validation | Enforced By |
|-------|-----------|-------------|
| `users.email` | Format (regex) | Ent validator |
| `users.email` | Uniqueness | Database constraint + Ent |
| `users.email` | Lowercase | Ent hook |
| `users.password_hash` | Not empty | Ent validator |
| `workspaces.name` | Not empty | Ent validator |

### Frontend (Zod)

| Field | Validation | Error Message |
|-------|-----------|---------------|
| `email` | Format (email) | "Invalid email format" |
| `password` | Min length (8) | "Password must be at least 8 characters" |
| `workspaceName` | Not empty | "Workspace name is required" |

---

## Migration Files

Based on the schema, three migration files will be created:

1. **`001_create_users_table.sql`**:
   - Creates `users` table with all fields
   - Creates unique index on `email`

2. **`002_create_workspaces_table.sql`**:
   - Creates `workspaces` table with all fields

3. **`003_create_user_workspaces_table.sql`**:
   - Creates `user_workspaces` join table
   - Creates foreign key constraints
   - Creates indexes on `user_id` and `workspace_id`

---

## TypeScript Types (Frontend)

**Location**: `client/src/entities/user/types.ts`

```typescript
export interface User {
  id: number
  email: string
  createdAt: string // ISO 8601 timestamp
  updatedAt: string
}

export interface Workspace {
  id: number
  name: string
  createdAt: string
  updatedAt: string
}

export interface SignupRequest {
  email: string
  password: string
  workspaceName: string
}

export interface SignupResponse {
  user: User
  workspace: Workspace
  message: string
}
```

**Notes**:
- `password_hash` is never exposed to frontend
- Timestamps are serialized as ISO 8601 strings in JSON

---

## Query Patterns

### Common Queries (Ent)

**1. Create User with Workspace** (Signup):
```go
tx, _ := client.Tx(ctx)

user, _ := tx.User.Create().
    SetEmail("user@example.com").
    SetPasswordHash(hashedPassword).
    Save(ctx)

workspace, _ := tx.Workspace.Create().
    SetName("My Workspace").
    AddUsers(user).
    Save(ctx)

tx.Commit()
```

**2. Check Email Existence**:
```go
exists, _ := client.User.Query().
    Where(user.Email(strings.ToLower(email))).
    Exist(ctx)
```

**3. Get User with Workspaces**:
```go
u, _ := client.User.Query().
    Where(user.ID(userID)).
    WithWorkspaces().
    Only(ctx)
```

---

## Security Considerations

1. **Password Storage**:
   - Never store plain-text passwords
   - Use bcrypt with cost factor ≥ 10
   - `password_hash` field marked as `Sensitive()` in Ent (not logged)

2. **Email Uniqueness**:
   - Database constraint prevents duplicates
   - Case-insensitive via hook (all emails stored lowercase)
   - Prevents race conditions

3. **Foreign Key Constraints**:
   - `ON DELETE CASCADE` ensures cleanup if user/workspace is deleted
   - Maintains referential integrity

4. **Input Validation**:
   - Both frontend (Zod) and backend (Ent) validate
   - Defense in depth approach

---

## Next Steps

With the data model defined, proceed to:
1. Create Ent schemas (`user.go`, `workspace.go`)
2. Run `ent generate` to generate Go code
3. Generate migration SQL files
4. Define API contracts in `contracts/signup-api.yaml`
