# Quickstart Guide: User Signup Feature

**Feature**: User Signup
**Branch**: `001-user-signup`
**Last Updated**: 2025-12-27

## Overview

This guide helps developers get started with implementing and testing the user signup feature. It covers setup, development workflow, and testing procedures.

## Prerequisites

Before starting, ensure you have:

- [x] Docker and Docker Compose installed
- [x] Git repository cloned
- [x] `.env` file configured (copy from `.env.example`)
- [x] Branch `001-user-signup` checked out

## Quick Start (5 minutes)

```bash
# 1. Ensure you're on the feature branch
git checkout 001-user-signup

# 2. Start all services
docker-compose up -d

# 3. Install Ent (first time only)
docker-compose exec api go get entgo.io/ent/cmd/ent

# 4. Install dependencies (first time only)
docker-compose exec api go mod tidy

# 5. Verify services are running
docker-compose ps
```

Expected output:
```
NAME                COMMAND                  SERVICE   STATUS    PORTS
finsight-api        "air"                    api       Up        0.0.0.0:28080->8080/tcp
finsight-client     "npm run dev"            client    Up        0.0.0.0:3000->3000/tcp
finsight-db         "docker-entrypoint.s…"   db        Up        0.0.0.0:5432->5432/tcp
```

---

## Development Workflow

### Phase 1: Backend Setup

#### 1.1 Initialize Ent Schemas

```bash
# Create Ent schemas (if not exists)
docker-compose exec api go run entgo.io/ent/cmd/ent init --target internal/infrastructure/ent/schema User Workspace
```

This creates:
- `backend/internal/infrastructure/ent/schema/user.go`
- `backend/internal/infrastructure/ent/schema/workspace.go`

#### 1.2 Define Ent Schemas

Edit the generated files according to [data-model.md](./data-model.md):

**`backend/internal/infrastructure/ent/schema/user.go`**:
- Add fields: `email`, `password_hash`, `created_at`, `updated_at`
- Add edge: `workspaces` (many-to-many to Workspace)
- Add hooks for email normalization
- Add validators for email format

**`backend/internal/infrastructure/ent/schema/workspace.go`**:
- Add fields: `name`, `created_at`, `updated_at`
- Add edge: `users` (many-to-many to User)

#### 1.3 Generate Ent Code

```bash
# Generate Ent client code
docker-compose exec api go run entgo.io/ent/cmd/ent generate ./internal/infrastructure/ent/schema
```

This generates:
- `backend/internal/infrastructure/ent/*.go` (client, mutations, queries, etc.)

#### 1.4 Create Database Migrations

**Option A: Manual SQL** (Recommended for first iteration):

Create three files in `db/migrations/`:

1. `001_create_users_table.sql`
2. `002_create_workspaces_table.sql`
3. `003_create_user_workspaces_table.sql`

See [data-model.md](./data-model.md) for SQL schema.

**Option B: Atlas Migration** (For schema changes):

```bash
# Install Atlas CLI (on host machine)
brew install ariga/tap/atlas  # macOS
# or
curl -sSf https://atlasgo.sh | sh  # Linux

# Generate migration from Ent schema
atlas migrate diff create_user_signup \
  --dir "file://db/migrations" \
  --to "ent://backend/internal/infrastructure/ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

#### 1.5 Apply Migrations

```bash
# Restart database to apply migrations
docker-compose restart db

# Verify tables were created
docker-compose exec db psql -U finsight_user -d finsight_dev -c "\dt"
```

Expected tables:
- `users`
- `workspaces`
- `user_workspaces`

#### 1.6 Install Backend Dependencies

```bash
# Install Gorilla Sessions
docker-compose exec api go get github.com/gorilla/sessions

# Install bcrypt (should already be in deps)
docker-compose exec api go get golang.org/x/crypto/bcrypt

# Tidy dependencies
docker-compose exec api go mod tidy
```

---

### Phase 2: Backend Implementation

Implement files in this order:

#### 2.1 Domain Layer

1. **`internal/domain/model/user.go`** (if needed for business logic)
2. **`internal/domain/model/workspace.go`** (if needed)
3. **`internal/domain/service/auth_service.go`**:
   - `HashPassword(password string) (string, error)`
   - `ValidateEmail(email string) error`

#### 2.2 Infrastructure Layer

1. **`internal/infrastructure/repositories/user_repository.go`**:
   - `CreateUser(ctx, email, passwordHash) (*ent.User, error)`
   - `GetUserByEmail(ctx, email) (*ent.User, error)`
   - `EmailExists(ctx, email) (bool, error)`

2. **`internal/infrastructure/repositories/workspace_repository.go`**:
   - `CreateWorkspace(ctx, name) (*ent.Workspace, error)`

3. **`internal/infrastructure/session/session_manager.go`**:
   - `InitStore(secretKey string) *sessions.CookieStore`
   - `SetSession(c *gin.Context, userID, workspaceID int, email string) error`
   - `GetSession(c *gin.Context) (userID, workspaceID int, email string, err error)`

#### 2.3 Application Layer

1. **`internal/application/usecase/signup_usecase.go`**:
   - `Execute(ctx, email, password, workspaceName) (user, workspace, error)`
   - Orchestrates: validation → hash password → create user+workspace in transaction → return

#### 2.4 HTTP Layer

1. **`internal/infrastructure/http/handler/signup_handler.go`**:
   - `Signup(c *gin.Context)`
   - `GetSession(c *gin.Context)`

2. **Update `internal/infrastructure/http/router/router.go`**:
   - Add routes: `POST /api/auth/signup`, `GET /api/auth/session`

---

### Phase 3: Frontend Implementation

#### 3.1 Setup Types

1. **`client/src/entities/user/types.ts`**:
   - Define `User`, `Workspace`, `SignupRequest`, `SignupResponse` interfaces

#### 3.2 Setup Validation

1. **`client/src/shared/lib/validation.ts`**:
   - Define Zod schema for signup form
   - Export `signupSchema`

#### 3.3 Setup API Client

1. **`client/src/shared/api/client.ts`** (if not exists):
   - Create base fetch wrapper with `credentials: 'include'`

2. **`client/src/features/auth/model/useSignup.ts`**:
   - Create `useSignup` hook using TanStack Query's `useMutation`
   - Call `POST /api/auth/signup`

#### 3.4 Create UI Components

1. **`client/src/pages/Signup/ui/SignupForm.tsx`**:
   - Form with TanStack Form
   - Fields: email, password, workspaceName
   - Use Zod validation
   - Call `useSignup` on submit

2. **`client/src/pages/Signup/ui/SignupPage.tsx`**:
   - Page wrapper component
   - Renders SignupForm

3. **`client/src/pages/Signup/index.tsx`**:
   - Export SignupPage

#### 3.5 Setup Routing

1. **`client/src/app/routes/signup.tsx`**:
   - Define route for `/signup`
   - Import SignupPage from `pages/Signup`

---

## Testing

### Backend Tests

#### Unit Tests

```bash
# Run all backend tests
docker-compose exec api go test ./...

# Run specific package tests
docker-compose exec api go test ./internal/application/usecase
```

**Key test files to create**:
- `internal/application/usecase/signup_usecase_test.go`
- `internal/domain/service/auth_service_test.go`

**Test scenarios**:
- Valid signup creates user and workspace
- Duplicate email returns error
- Invalid email format returns error
- Weak password returns error
- Transaction rollback on workspace creation failure

#### Integration Tests

```bash
# Run integration tests
docker-compose exec api go test -tags=integration ./tests/integration
```

**Test scenarios**:
- Full signup flow (API → DB)
- Session creation and retrieval
- Email uniqueness enforcement (database level)

### Frontend Tests

```bash
# Run frontend tests
docker-compose exec client npm test

# Run specific test file
docker-compose exec client npm test -- useSignup.test.ts
```

**Key test files to create**:
- `client/src/features/auth/model/useSignup.test.ts`
- `client/src/shared/lib/validation.test.ts`

**Test scenarios**:
- `useSignup` hook handles success
- `useSignup` hook handles errors
- Validation schema catches invalid inputs

### Manual Testing

#### 1. Test Signup Flow

```bash
# Using curl
curl -X POST http://localhost:28080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "workspaceName": "Test Workspace"
  }' \
  -c cookies.txt \
  -v
```

Expected: `201 Created` with session cookie

#### 2. Test Session Retrieval

```bash
curl http://localhost:28080/api/auth/session \
  -b cookies.txt \
  -v
```

Expected: `200 OK` with user and workspace data

#### 3. Test Email Uniqueness

```bash
# Try to signup with same email again
curl -X POST http://localhost:28080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "workspaceName": "Another Workspace"
  }' \
  -v
```

Expected: `409 Conflict` with "Email already registered"

#### 4. Test Frontend

1. Open browser: http://localhost:3000/signup
2. Fill form with valid data
3. Submit
4. Verify redirect to top page (`/`)
5. Verify session persists (refresh page → still logged in)

---

## Common Issues & Troubleshooting

### Database Connection Issues

```bash
# Check if database is ready
docker-compose exec db pg_isready

# Check database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Ent Code Generation Errors

```bash
# Make sure schemas are valid Go code
docker-compose exec api go build ./internal/infrastructure/ent/schema

# Regenerate with verbose output
docker-compose exec api go run entgo.io/ent/cmd/ent generate --verbose ./internal/infrastructure/ent/schema
```

### Migration Not Applied

```bash
# Check migration files exist
ls -la db/migrations/

# Check Docker volume mount
docker-compose exec db ls -la /docker-entrypoint-initdb.d/

# If migrations don't run, recreate database volume
docker-compose down -v
docker-compose up -d
```

### Session Cookie Not Set

- Verify `credentials: 'include'` in frontend fetch calls
- Check CORS settings in backend allow credentials
- Ensure cookie domain/path match
- Check `SameSite` attribute (should be `Lax` for dev)

### CORS Errors

Update `backend/internal/infrastructure/http/middleware/cors.go`:
```go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

---

## Environment Variables

### Backend (`.env`)

```bash
# Database
POSTGRES_DB=finsight_dev
POSTGRES_USER=finsight_user
POSTGRES_PASSWORD=finsight_password
DB_PORT=5432

# Server
GIN_MODE=debug

# Session (add these)
SESSION_SECRET=your-secret-key-change-in-production
SESSION_MAX_AGE=604800  # 7 days in seconds
```

### Frontend (`client/.env.local`)

```bash
VITE_API_URL=http://localhost:28080
```

---

## Development Commands Reference

### Docker

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f api      # Backend
docker-compose logs -f client   # Frontend
docker-compose logs -f db       # Database

# Restart a service
docker-compose restart api
```

### Backend (Go)

```bash
# Run commands in API container
docker-compose exec api <command>

# Examples:
docker-compose exec api go test ./...
docker-compose exec api go mod tidy
docker-compose exec api go build ./cmd/server
```

### Frontend (Node)

```bash
# Run commands in client container
docker-compose exec client <command>

# Examples:
docker-compose exec client npm test
docker-compose exec client npm run build
docker-compose exec client npm run lint
```

---

## Next Steps

After completing this feature:

1. Run full test suite (backend + frontend)
2. Test manually with different scenarios
3. Update `spec.md` if requirements changed
4. Create pull request from `001-user-signup` to `main`
5. Request code review
6. Address review comments
7. Merge to main

---

## Additional Resources

- [Ent Documentation](https://entgo.io/docs/getting-started)
- [TanStack Query Documentation](https://tanstack.com/query/latest)
- [TanStack Form Documentation](https://tanstack.com/form/latest)
- [Gorilla Sessions Documentation](http://www.gorillatoolkit.org/pkg/sessions)
- [OpenAPI Specification](./contracts/signup-api.yaml)
- [Data Model](./data-model.md)
- [Research Document](./research.md)