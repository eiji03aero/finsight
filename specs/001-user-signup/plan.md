# Implementation Plan: User Signup

**Branch**: `001-user-signup` | **Date**: 2025-12-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-user-signup/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement user signup functionality that allows new users to create an account by providing email, password, and workspace name. The system will create both user and workspace records, enforce email uniqueness at both application and database levels, hash passwords securely, automatically log in the user, and redirect them to the top page. This is the foundational feature for user acquisition and account management.

**Technical Approach**:
- Backend: RESTful API endpoint using Go/Gin framework with clean architecture pattern
- Frontend: React form with TanStack Form for validation and TanStack Router for navigation
- Database: PostgreSQL with unique constraint on email field
- Security: Password hashing using bcrypt, case-insensitive email handling
- Session: Session-based authentication with secure session storage

## Technical Context

**Language/Version**:
- Backend: Go 1.24.0
- Frontend: TypeScript with React 19.2.0, Node.js (latest LTS)

**Primary Dependencies**:
- Backend: Gin (v1.11.0), bcrypt (for password hashing), PostgreSQL driver, session management library
- Frontend: TanStack Router (v1.132.0), TanStack Form (v1.0.0), TanStack Query (v5.66.5), Zod (v4.2.1), Radix UI components

**Storage**:
- PostgreSQL 18.1 (via Docker)
- Database: `finsight_dev`
- Tables needed: `users`, `workspaces`, `user_workspaces` (junction table)

**Testing**:
- Backend: Go testing (`go test`), integration tests for API endpoints
- Frontend: Vitest (v3.0.5), React Testing Library (v16.2.0) - only for functions and custom hooks (no component tests)

**Target Platform**:
- Backend: Linux server (Docker container)
- Frontend: Modern web browsers (Chrome, Firefox, Safari, Edge)
- Deployment: Docker Compose orchestration

**Project Type**: Web application (separate frontend + backend)

**Performance Goals**:
- Signup completion within 30 seconds (per SC-001)
- Handle 100+ concurrent signups without degradation (per SC-004)
- API response time < 500ms for signup endpoint
- Database query optimization for email uniqueness checks

**Constraints**:
- Email uniqueness must be enforced at database level (unique constraint)
- Passwords must be hashed (bcrypt, minimum cost factor 10)
- Case-insensitive email comparison required
- Session security: HttpOnly, Secure cookies in production
- CORS configuration for frontend-backend communication

**Scale/Scope**:
- Initial implementation for MVP
- Expected user base: 1,000-10,000 users initially
- 3 database tables involved
- 2-3 API endpoints (signup, session check, redirect)
- 1 frontend page (signup form)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Status**: ⚠️  Constitution file is currently a template. No specific principles are defined yet, so this check is marked as PASSED by default. When the project constitution is established, this section should be re-evaluated.

**Assumptions for this feature**:
- No complexity violations identified based on standard web application practices
- Clean architecture pattern in backend is already established
- Feature-Sliced Design pattern in frontend is already established
- Standard REST API design (no unusual architectural decisions)

**Re-evaluation Required**: After Phase 1 design completion, verify that:
- API contracts follow established patterns
- Data model follows existing conventions
- No unnecessary abstractions introduced
- Test coverage meets project standards

## Project Structure

### Documentation (this feature)

```text
specs/001-user-signup/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── signup-api.yaml  # OpenAPI specification for signup endpoints
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── model/
│   │   │   ├── user.go           # User entity
│   │   │   └── workspace.go      # Workspace entity
│   │   └── service/
│   │       └── auth_service.go   # Password hashing, email validation
│   ├── application/
│   │   └── usecase/
│   │       └── signup_usecase.go # Signup business logic
│   └── infrastructure/
│       ├── repositories/
│       │   ├── user_repository.go
│       │   └── workspace_repository.go
│       ├── http/
│       │   ├── handler/
│       │   │   └── signup_handler.go
│       │   └── router/
│       │       └── router.go (update)
│       └── session/
│           └── session_manager.go

client/
└── src/
    ├── app/
    │   └── routes/
    │       └── signup.tsx        # Route configuration for signup page
    ├── pages/
    │   └── Signup/
    │       ├── index.tsx         # Exports SignupPage component
    │       └── ui/
    │           ├── SignupPage.tsx    # Page component
    │           └── SignupForm.tsx    # Form component (no test)
    ├── features/
    │   └── auth/
    │       ├── ui/
    │       │   └── (shared auth UI components if needed)
    │       └── model/
    │           ├── useSignup.ts      # TanStack Query hook for signup API
    │           └── useSignup.test.ts # Test for useSignup hook
    ├── entities/
    │   └── user/
    │       └── types.ts          # User and Workspace TypeScript types
    └── shared/
        ├── api/
        │   └── client.ts         # Base API client configuration
        └── lib/
            ├── validation.ts     # Zod schemas for form validation
            └── validation.test.ts # Test for validation functions

db/
└── migrations/
    ├── 001_create_users_table.sql
    ├── 002_create_workspaces_table.sql
    └── 003_create_user_workspaces_table.sql
```

**Structure Decision**:

This project follows a **Web Application** architecture with separate frontend and backend:

1. **Backend (Go)**: Clean Architecture pattern
   - `domain/`: Core business entities and domain services (User, Workspace)
   - `application/`: Use cases orchestrating domain logic (SignupUseCase)
   - `infrastructure/`:
     - `repositories/`: Data access implementations
     - `http/`: HTTP handlers and routing
     - `session/`: Session management
   - Clear separation of concerns, dependency inversion principle

2. **Frontend (React/TypeScript)**: Feature-Sliced Design (FSD) with project-specific conventions
   - `app/routes/`: Route definitions linking URLs to pages
   - `pages/`: Page components organized by feature
     - `Signup/index.tsx`: Exports the page component
     - `Signup/ui/SignupPage.tsx`: Main page component
     - `Signup/ui/SignupForm.tsx`: Form component
   - `features/`: Feature modules with `ui/` (components) and `model/` (hooks/logic)
   - `entities/`: Business entities (user types)
   - `shared/`: Shared utilities and configurations
   - **Testing strategy**: Colocation pattern (tests next to source files), only for functions and custom hooks (no component tests)

3. **Database**: Managed via Docker Compose with migration support
   - PostgreSQL 18.1
   - Migration files in `db/migrations/` (mounted to database container)
   - One SQL file per purpose (3 files for 3 tables)
   - Migrations executed automatically on container startup

This structure supports:
- Independent development and testing of frontend/backend
- Clear boundaries between layers
- Easy addition of new features without coupling
- Scalable team collaboration
- Automated database schema management via Docker

## Complexity Tracking

> **No violations identified** - this section is not applicable for this feature.

This feature follows established patterns and does not introduce any architectural complexity that violates standard practices:
- Uses existing clean architecture in backend
- Uses existing FSD pattern in frontend with project conventions
- Standard REST API design
- Standard session-based authentication
- No new frameworks or unusual patterns introduced
