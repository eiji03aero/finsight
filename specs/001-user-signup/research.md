# Research & Technical Decisions: User Signup

**Feature**: User Signup
**Date**: 2025-12-27
**Status**: Complete

## Overview

This document captures technical research and decisions made during the planning phase for the user signup feature. It resolves all "NEEDS CLARIFICATION" items from the Technical Context and documents the rationale for key technology choices.

## Research Areas

### 1. Password Hashing Library (Backend)

**Decision**: Use `golang.org/x/crypto/bcrypt`

**Rationale**:
- Already included in project dependencies (found in go.mod)
- Industry-standard library for password hashing in Go
- Implements bcrypt algorithm with configurable cost factor
- Part of Go's official extended packages (maintained by Go team)
- Automatic salting and secure defaults
- Well-tested and audited for security

**Implementation Details**:
- Use `bcrypt.GenerateFromPassword()` with cost factor 10 (minimum secure level)
- Use `bcrypt.CompareHashAndPassword()` for login verification
- Store hashed password as string in database (bcrypt output is already base64 encoded)

**Alternatives Considered**:
- `golang.org/x/crypto/argon2`: More modern but requires manual salt management
- `golang.org/x/crypto/scrypt`: Good alternative but bcrypt is more widely adopted in Go community

**References**:
- [Go crypto/bcrypt documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetsheatcheat.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)

---

### 2. ORM & Database Access (Backend)

**Decision**: Use `ent` (https://entgo.io/)

**Rationale**:
- Project requirement (specified by user)
- Modern, type-safe ORM for Go with excellent code generation
- Schema-as-code approach: entities defined in Go, migrations auto-generated
- Built-in features we need:
  - Unique constraints (email uniqueness)
  - Hooks for data transformation (lowercase email)
  - Transaction support (atomic user+workspace creation)
  - Edge definitions (user-workspace relationships)
- Strong type safety: compile-time checks for queries
- Excellent integration with PostgreSQL
- Active development and good documentation

**Implementation Details**:
- Define schemas in `backend/internal/infrastructure/ent/schema/`:
  - `user.go`: User entity with email, password_hash fields
  - `workspace.go`: Workspace entity with name field
  - Define edges: User -> Workspaces (many-to-many)
- Use `ent generate` to create type-safe client code
- Migrations: `ent migrate diff` to generate SQL migration files
- Unique constraints and indexes defined in schema:
  ```go
  field.String("email").Unique().NotEmpty()
  ```
- Hooks for email normalization:
  ```go
  func (u *User) BeforeCreate() {
      u.Email = strings.ToLower(u.Email)
  }
  ```

**Migration Workflow**:
1. Define/update ent schema
2. Run `ent generate ./internal/infrastructure/ent/schema`
3. Generate migration: `atlas migrate diff migration_name --dir "file://db/migrations" --to "ent://internal/infrastructure/ent/schema" --dev-url "docker://postgres/15/test?search_path=public"`
4. Migration files stored in `db/migrations/`
5. Apply migrations via Docker init or CLI

**Alternatives Considered**:
- `gorm`: Popular but less type-safe, reflection-based
- `sqlx`: Lighter but no code generation or schema management
- `sqlc`: Good for raw SQL but ent provides better schema-first workflow

**References**:
- [Ent Documentation](https://entgo.io/docs/getting-started)
- [Ent Schema Definition](https://entgo.io/docs/schema-def)
- [Ent Migrations](https://entgo.io/docs/versioned-migrations)

---

### 3. Session Management (Backend)

**Decision**: Use `github.com/gorilla/sessions` with cookie store

**Rationale**:
- De facto standard for session management in Go web applications
- Integrates seamlessly with Gin framework
- Supports multiple storage backends (starting with cookie store for MVP)
- Built-in security features:
  - Secure cookie signing and encryption
  - HttpOnly flag support
  - Secure flag support (HTTPS-only in production)
  - SameSite attribute support
- Simple API for getting/setting session values
- Easy to migrate to Redis/database store later if needed

**Implementation Details**:
- Use cookie-based sessions for MVP (stateless, no infrastructure needed)
- Session data stored: user ID, email, workspace ID
- Cookie settings:
  - HttpOnly: true (prevent XSS attacks)
  - Secure: true (HTTPS only in production)
  - SameSite: Lax (CSRF protection)
  - MaxAge: 7 days (configurable)
- Cookie name: "finsight_session"

**Migration Path**:
- Can easily switch to Redis store (`github.com/gorilla/sessions/redistore`) when scaling

**Alternatives Considered**:
- `github.com/gin-contrib/sessions`: Gin-specific wrapper, but gorilla/sessions is more flexible
- JWT tokens: Stateless but harder to invalidate, overkill for MVP

**References**:
- [Gorilla Sessions GitHub](https://github.com/gorilla/sessions)
- [Gorilla Sessions with Gin example](https://github.com/gin-gonic/contrib/tree/master/sessions)

---

### 4. Database Migration Management

**Decision**: Use Ent's built-in migration system with Atlas

**Rationale**:
- Ent provides first-class migration support via Atlas (https://atlasgo.io/)
- Automatic migration generation from schema changes
- Versioned migration files (SQL)
- Type-safe: migrations match ent schema definitions
- Supports migration testing and dry-runs
- Can still use manual SQL migrations if needed

**Implementation Details**:
- Use Atlas CLI for migration generation
- Migration files stored in `db/migrations/`
- Each migration is a SQL file with up/down migrations
- Mount migrations to PostgreSQL container for automatic execution
- Naming convention: `<timestamp>_<description>.sql`
- For this feature (manual creation for first iteration):
  - `001_create_users_table.sql`
  - `002_create_workspaces_table.sql`
  - `003_create_user_workspaces_table.sql`

**Docker Compose Configuration**:
```yaml
db:
  volumes:
    - ./db/migrations:/docker-entrypoint-initdb.d
```

**Workflow**:
1. Define ent schema
2. Generate migrations: `atlas migrate diff --env local`
3. Review generated SQL
4. Commit migrations to `db/migrations/`
5. Docker automatically applies on startup

**Alternatives Considered**:
- Manual SQL only: Harder to keep in sync with code
- `golang-migrate/migrate`: Good but ent+Atlas integration is seamless

**References**:
- [Ent Versioned Migrations](https://entgo.io/docs/versioned-migrations)
- [Atlas Documentation](https://atlasgo.io/getting-started)

---

### 5. Email Validation (Backend)

**Decision**: Use Ent validators + custom validation in domain service

**Rationale**:
- Ent supports field validators out of the box
- Can combine regex pattern validation with custom logic
- Domain service layer adds business logic validation (email format, uniqueness)
- Client-side validation (Zod) provides first line of defense

**Implementation Details**:
- Ent schema validation:
  ```go
  field.String("email").
      Unique().
      NotEmpty().
      Validate(func(s string) error {
          // Regex validation
          emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
          if !emailRegex.MatchString(s) {
              return fmt.Errorf("invalid email format")
          }
          return nil
      })
  ```
- Case normalization via ent hooks (before create)
- Uniqueness enforced at database level via unique constraint

**Alternatives Considered**:
- `github.com/badoux/checkmail`: Third-party library not needed with ent validators
- Email verification workflow: Out of scope for MVP (can add later)

---

### 6. Frontend Form Management

**Decision**: Use `@tanstack/react-form` with Zod validation

**Rationale**:
- Already in project dependencies (`package.json`)
- Type-safe form state management
- Excellent integration with Zod for validation schemas
- Supports async validation (for checking email uniqueness)
- Good performance (minimal re-renders)
- Modern React patterns (hooks-based)

**Implementation Details**:
- Define Zod schema for signup form in `shared/lib/validation.ts`
- Use TanStack Form for form state and submission
- Form fields: email, password, workspace name
- Client-side validation before API call

**Validation Schema** (Zod):
```typescript
const signupSchema = z.object({
  email: z.string().email("Invalid email format"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  workspaceName: z.string().min(1, "Workspace name is required"),
})
```

**Alternatives Considered**:
- `react-hook-form`: Equally good, but TanStack Form is already in dependencies
- Native form validation: Insufficient for complex validation requirements

**References**:
- [TanStack Form Documentation](https://tanstack.com/form/latest)
- [Zod Documentation](https://zod.dev/)

---

### 7. Frontend API Client

**Decision**: Use `@tanstack/react-query` with custom hook

**Rationale**:
- Already in project dependencies
- Industry standard for data fetching in React
- Built-in features:
  - Loading/error states
  - Retry logic
  - Request deduplication
  - Cache management
- Excellent developer experience with DevTools
- Type-safe with TypeScript

**Implementation Details**:
- Create `useSignup` hook in `features/auth/model/useSignup.ts`
- Use `useMutation` for signup POST request
- Handle success: session is set via cookie, redirect to top page
- Handle error: display error message in form

**API Client Setup**:
- Base client in `shared/api/client.ts`
- Use `fetch` for HTTP requests
- Configure base URL from environment variable (`VITE_API_URL`)
- Set credentials: 'include' for cookie-based sessions

**Example Hook**:
```typescript
export const useSignup = () => {
  return useMutation({
    mutationFn: async (data: SignupFormData) => {
      const response = await fetch(`${API_URL}/api/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data),
      })
      if (!response.ok) throw new Error('Signup failed')
      return response.json()
    },
  })
}
```

**Alternatives Considered**:
- `swr`: Good alternative but React Query is more feature-complete
- Plain `fetch`: Too low-level, would need to rebuild caching/loading states

**References**:
- [TanStack Query Documentation](https://tanstack.com/query/latest)
- [TanStack Query Best Practices](https://tkdodo.eu/blog/practical-react-query)

---

### 8. Frontend Routing & Navigation

**Decision**: Use `@tanstack/react-router` for navigation after signup

**Rationale**:
- Already in project dependencies
- Type-safe routing with full TypeScript support
- Programmatic navigation via `useNavigate` hook
- Supports route-based code splitting
- Modern routing solution with excellent DX

**Implementation Details**:
- Define `/signup` route in `app/routes/signup.tsx`
- After successful signup, navigate to `/` (top page)
- Use `router.navigate({ to: '/' })`

**Alternatives Considered**:
- React Router: TanStack Router provides better type safety

**References**:
- [TanStack Router Documentation](https://tanstack.com/router/latest)

---

### 9. Case-Insensitive Email Handling

**Decision**: Use Ent hooks to normalize email to lowercase + database unique constraint

**Rationale**:
- Ensures "test@example.com" === "Test@Example.com"
- Database unique constraint works on lowercase emails
- Ent hooks provide a clean way to enforce normalization
- Simple and performant

**Implementation Details**:

**Ent Hook** (in schema):
```go
func (User) Hooks() []ent.Hook {
    return []ent.Hook{
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                    if email, ok := m.Field("email"); ok {
                        m.SetField("email", strings.ToLower(email.(string)))
                    }
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpCreate,
        ),
    }
}
```

**Database Constraint** (generated by ent):
```sql
CREATE UNIQUE INDEX users_email_key ON users (email);
```

**Alternatives Considered**:
- CITEXT extension: PostgreSQL-specific, not needed with hook approach
- Manual normalization in handler: Ent hook is cleaner and centralized

---

### 10. Transaction Handling for User + Workspace Creation

**Decision**: Use Ent's transaction API

**Rationale**:
- Ent provides built-in transaction support
- Type-safe transaction operations
- Ensures user and workspace creation are atomic
- Prevents orphaned records
- Clean API with context-based transactions

**Implementation Details**:
```go
tx, err := client.Tx(ctx)
if err != nil {
    return err
}
defer func() {
    if v := recover(); v != nil {
        tx.Rollback()
        panic(v)
    }
}()

// Create user
user, err := tx.User.Create().
    SetEmail(email).
    SetPasswordHash(hashedPassword).
    Save(ctx)
if err != nil {
    return rollback(tx, err)
}

// Create workspace
workspace, err := tx.Workspace.Create().
    SetName(workspaceName).
    AddUsers(user).
    Save(ctx)
if err != nil {
    return rollback(tx, err)
}

return tx.Commit()
```

**Edge Case Handling**:
- If workspace creation fails → rollback user creation
- If edge creation fails → rollback both
- Return appropriate error to client

**Alternatives Considered**:
- Manual transaction with database/sql: Ent provides cleaner API
- Saga pattern: Over-engineered for single-database operation

---

## Summary of Key Technologies

| Component | Technology | Version | Reason |
|-----------|-----------|---------|--------|
| ORM | `entgo.io/ent` | Latest | Type-safe ORM, schema-as-code, excellent PostgreSQL support |
| Password Hashing | `golang.org/x/crypto/bcrypt` | Latest | Industry standard, already in deps |
| Session Management | `github.com/gorilla/sessions` | Latest | De facto standard for Go web apps |
| Database Migrations | Ent + Atlas | Latest | Auto-generated from schema, versioned SQL files |
| Frontend Forms | `@tanstack/react-form` | 1.0.0 | Already in deps, type-safe |
| Frontend API | `@tanstack/react-query` | 5.66.5 | Already in deps, industry standard |
| Frontend Routing | `@tanstack/react-router` | 1.132.0 | Already in deps, type-safe |
| Validation | Zod | 4.2.1 | Already in deps, shared FE/BE types |

## Ent Setup Requirements

Before implementing this feature, ensure Ent is properly set up:

1. **Install Ent**:
   ```bash
   docker-compose exec api go get entgo.io/ent/cmd/ent
   ```

2. **Initialize Ent** (if not already done):
   ```bash
   docker-compose exec api go run entgo.io/ent/cmd/ent init --target internal/infrastructure/ent/schema User Workspace
   ```

3. **Install Atlas CLI** (for migrations):
   ```bash
   # On macOS
   brew install ariga/tap/atlas

   # On Linux
   curl -sSf https://atlasgo.sh | sh
   ```

4. **Configure Atlas** (create `atlas.hcl` in backend/):
   ```hcl
   env "local" {
     url = "postgres://finsight_user:finsight_password@localhost:5432/finsight_dev?sslmode=disable"
     dev = "docker://postgres/15/dev"
   }
   ```

## Next Steps

All technical unknowns have been resolved. Proceed to Phase 1:
1. Generate `data-model.md` with Ent entity definitions
2. Create API contracts in `contracts/signup-api.yaml`
3. Generate `quickstart.md` for development setup
4. Update agent context with new technologies (Ent, Atlas, Gorilla Sessions)
