# Tasks: User Signup

**Input**: Design documents from `/specs/001-user-signup/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are included for functions and custom hooks only (no component tests per project conventions)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/`
- **Frontend**: `client/src/`
- **Database**: `db/migrations/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and Ent ORM setup

- [ ] T001 Install Ent ORM dependencies in backend: `docker-compose exec api go get entgo.io/ent/cmd/ent`
- [ ] T002 [P] Install Gorilla Sessions in backend: `docker-compose exec api go get github.com/gorilla/sessions`
- [ ] T003 [P] Install bcrypt in backend: `docker-compose exec api go get golang.org/x/crypto/bcrypt`
- [ ] T004 Run `docker-compose exec api go mod tidy` to update dependencies
- [ ] T005 Initialize Ent schemas in `backend/internal/infrastructure/ent/schema/`: `docker-compose exec api go run entgo.io/ent/cmd/ent init --target internal/infrastructure/ent/schema User Workspace`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Database Schema Setup

- [ ] T006 Define Ent User schema in `backend/internal/infrastructure/ent/schema/user.go` with fields: email (unique, lowercase via hook), password_hash, created_at, updated_at
- [ ] T007 Define Ent Workspace schema in `backend/internal/infrastructure/ent/schema/workspace.go` with fields: name, created_at, updated_at
- [ ] T008 Add many-to-many edge between User and Workspace in both schema files
- [ ] T009 Add email validation hook and lowercase normalization hook to User schema
- [ ] T010 Generate Ent client code: `docker-compose exec api go run entgo.io/ent/cmd/ent generate ./internal/infrastructure/ent/schema`

### Database Migrations

- [ ] T011 [P] Create migration file `db/migrations/001_create_users_table.sql` with users table schema
- [ ] T012 [P] Create migration file `db/migrations/002_create_workspaces_table.sql` with workspaces table schema
- [ ] T013 [P] Create migration file `db/migrations/003_create_user_workspaces_table.sql` with join table schema
- [ ] T014 Restart database to apply migrations: `docker-compose restart db`
- [ ] T015 Verify migrations applied successfully: `docker-compose exec db psql -U finsight_user -d finsight_dev -c "\dt"`

### Session Management Infrastructure

- [ ] T016 Create session manager in `backend/internal/infrastructure/session/session_manager.go` with InitStore, SetSession, GetSession functions
- [ ] T017 Initialize session store in `backend/cmd/server/main.go` with secret key from environment

### Frontend Types & Validation

- [ ] T018 [P] Create TypeScript types in `client/src/entities/user/types.ts`: User, Workspace, SignupRequest, SignupResponse interfaces
- [ ] T019 [P] Create Zod validation schema in `client/src/shared/lib/validation.ts` for signup form (email, password min 8 chars, workspaceName)
- [ ] T020 [P] Create test file `client/src/shared/lib/validation.test.ts` for validation schema
- [ ] T021 Create base API client in `client/src/shared/api/client.ts` with fetch wrapper and credentials: 'include'

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - New User Account Creation with Workspace (Priority: P1) 🎯 MVP

**Goal**: Enable new users to sign up with email, password, and workspace name, automatically log in, and redirect to top page

**Independent Test**: Navigate to `/signup`, enter valid credentials and workspace name, submit form, verify user and workspace are created in database, user is logged in (session cookie set), and redirected to `/` (top page)

### Backend Implementation for US1

- [ ] T022 [P] [US1] Create auth service in `backend/internal/domain/service/auth_service.go` with HashPassword and ValidateEmail functions
- [ ] T023 [P] [US1] Create test file `backend/internal/domain/service/auth_service_test.go` for auth service functions
- [ ] T024 [P] [US1] Create user repository in `backend/internal/infrastructure/repositories/user_repository.go` with CreateUser, GetUserByEmail, EmailExists methods
- [ ] T025 [P] [US1] Create workspace repository in `backend/internal/infrastructure/repositories/workspace_repository.go` with CreateWorkspace method
- [ ] T026 [US1] Create signup usecase in `backend/internal/application/usecase/signup_usecase.go` with Execute method (orchestrates: validate email → check email exists → hash password → create user + workspace in transaction → return user and workspace)
- [ ] T027 [US1] Create test file `backend/internal/application/usecase/signup_usecase_test.go` for signup usecase
- [ ] T028 [US1] Create signup handler in `backend/internal/infrastructure/http/handler/signup_handler.go` with Signup and GetSession methods
- [ ] T029 [US1] Add routes in `backend/internal/infrastructure/http/router/router.go`: POST /api/auth/signup, GET /api/auth/session
- [ ] T030 [US1] Update CORS middleware in `backend/internal/infrastructure/http/middleware/cors.go` to allow credentials from localhost:3000

### Frontend Implementation for US1

- [ ] T031 [P] [US1] Create useSignup hook in `client/src/features/auth/model/useSignup.ts` using TanStack Query's useMutation for POST /api/auth/signup
- [ ] T032 [P] [US1] Create test file `client/src/features/auth/model/useSignup.test.ts` for useSignup hook
- [ ] T033 [US1] Create SignupForm component in `client/src/pages/Signup/ui/SignupForm.tsx` with TanStack Form, Zod validation, and useSignup hook
- [ ] T034 [US1] Create SignupPage component in `client/src/pages/Signup/ui/SignupPage.tsx` that renders SignupForm
- [ ] T035 [US1] Export SignupPage from `client/src/pages/Signup/index.tsx`
- [ ] T036 [US1] Create signup route in `client/src/app/routes/signup.tsx` linking `/signup` to SignupPage
- [ ] T037 [US1] Add navigation logic in useSignup success callback to redirect to `/` using TanStack Router

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Users can sign up, get logged in, and be redirected to the top page.

---

## Phase 4: User Story 2 & 3 - Secure Password Handling & Email Uniqueness (Priority: P2)

**Note**: US2 and US3 are combined as they enhance the security layer of US1 rather than adding new user flows

**Goal**: Ensure passwords are validated (min 8 chars), hashed with bcrypt, and stored securely. Ensure email uniqueness is enforced at both application and database levels with case-insensitive checking.

**Independent Test**:
- **US2**: Attempt to create account with password < 8 chars (should fail). Create valid account and verify password is hashed in database.
- **US3**: Attempt to create two accounts with same email (should fail on second attempt). Attempt with different capitalization (test@example.com vs Test@Example.com) and verify both are rejected.

### Backend Enhancements for US2 & US3

- [ ] T038 [P] [US2] [US3] Verify bcrypt cost factor is set to 10 or higher in `backend/internal/domain/service/auth_service.go` HashPassword function
- [ ] T039 [P] [US2] [US3] Add password minimum length validation (8 chars) to signup usecase in `backend/internal/application/usecase/signup_usecase.go`
- [ ] T040 [US2] [US3] Verify email uniqueness check is case-insensitive in user repository `backend/internal/infrastructure/repositories/user_repository.go` EmailExists method
- [ ] T041 [US2] [US3] Add integration test in `backend/tests/integration/signup_test.go` for: weak password rejection, duplicate email rejection, case-insensitive email rejection, password hashing verification

### Frontend Enhancements for US2 & US3

- [ ] T042 [P] [US2] [US3] Verify password min length validation (8 chars) is in Zod schema in `client/src/shared/lib/validation.ts`
- [ ] T043 [P] [US2] Ensure password input field in `client/src/pages/Signup/ui/SignupForm.tsx` has type="password" (masked characters)
- [ ] T044 [US3] Add error handling in `client/src/features/auth/model/useSignup.ts` to display "Email already registered" error when API returns 409 Conflict

**Checkpoint**: At this point, User Stories 1, 2, and 3 should all work together. Password security and email uniqueness are fully enforced.

---

## Phase 5: User Story 4 - Workspace Name Customization (Priority: P3)

**Goal**: Allow users to provide a custom workspace name during signup and have it properly stored and displayed

**Independent Test**: Create account with workspace name "My Project", verify workspace is created with exact name "My Project" in database, and verify workspace name is displayed on the top page after login

### Backend Implementation for US4

- [ ] T045 [US4] Verify workspace name validation (not empty) in Ent schema `backend/internal/infrastructure/ent/schema/workspace.go`
- [ ] T046 [US4] Verify workspace name is set from request in signup usecase `backend/internal/application/usecase/signup_usecase.go`
- [ ] T047 [US4] Update signup handler in `backend/internal/infrastructure/http/handler/signup_handler.go` to accept workspaceName from request body

### Frontend Implementation for US4

- [ ] T048 [US4] Verify workspaceName field is included in SignupForm component `client/src/pages/Signup/ui/SignupForm.tsx`
- [ ] T049 [US4] Verify workspaceName validation (not empty) is in Zod schema `client/src/shared/lib/validation.ts`
- [ ] T050 [US4] Ensure workspaceName is included in SignupRequest type in `client/src/entities/user/types.ts`

**Checkpoint**: All user stories should now be independently functional. Users can sign up with custom workspace names.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and final validation

- [ ] T051 [P] Add error logging to signup usecase in `backend/internal/application/usecase/signup_usecase.go`
- [ ] T052 [P] Add success logging to signup handler in `backend/internal/infrastructure/http/handler/signup_handler.go`
- [ ] T053 [P] Add loading state UI to SignupForm component in `client/src/pages/Signup/ui/SignupForm.tsx`
- [ ] T054 [P] Add error message display UI to SignupForm component in `client/src/pages/Signup/ui/SignupForm.tsx`
- [ ] T055 [P] Add environment variable SESSION_SECRET to `.env` file and update session_manager.go to use it
- [ ] T056 Run full backend test suite: `docker-compose exec api go test ./...`
- [ ] T057 Run full frontend test suite: `docker-compose exec client npm test`
- [ ] T058 Perform manual end-to-end test following quickstart.md validation scenarios
- [ ] T059 [P] Update CLAUDE.md with any new backend commands or patterns
- [ ] T060 Code review and refactoring pass

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion - Core MVP functionality
- **User Story 2 & 3 (Phase 4)**: Depends on Foundational completion - Can run in parallel with US1 or after US1
- **User Story 4 (Phase 5)**: Depends on Foundational completion - Can run in parallel with other stories or last
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories - **RECOMMENDED MVP**
- **User Story 2 & 3 (P2)**: Can start after Foundational (Phase 2) - Enhances US1 but independently testable
- **User Story 4 (P3)**: Can start after Foundational (Phase 2) - Extends US1 but independently testable

### Within Each User Story

- Backend: Auth service & repositories can be built in parallel → Usecase → Handler → Router
- Frontend: Types, validation, and API client can be built in parallel → Custom hook → Components → Routes
- Tests for functions/hooks should be created alongside their implementation

### Parallel Opportunities

**Setup Phase (Phase 1)**:
- T002, T003 (Gorilla Sessions, bcrypt) can run in parallel

**Foundational Phase (Phase 2)**:
- T011, T012, T013 (migration files) can run in parallel
- T018, T019, T020, T021 (frontend types and validation) can run in parallel after T010

**User Story 1 (Phase 3)**:
- T022, T023, T024, T025 (auth service and repositories) can run in parallel
- T031, T032 (useSignup hook and test) can run in parallel with backend tasks

**User Story 2 & 3 (Phase 4)**:
- T038, T039, T040 (backend enhancements) can run in parallel
- T042, T043, T044 (frontend enhancements) can run in parallel

**Polish Phase (Phase 6)**:
- T051, T052, T053, T054, T055, T059 can run in parallel

---

## Parallel Example: User Story 1 Backend

```bash
# Launch all backend foundation tasks together:
Task T022: "Create auth service in backend/internal/domain/service/auth_service.go"
Task T023: "Create auth service test in backend/internal/domain/service/auth_service_test.go"
Task T024: "Create user repository in backend/internal/infrastructure/repositories/user_repository.go"
Task T025: "Create workspace repository in backend/internal/infrastructure/repositories/workspace_repository.go"
```

## Parallel Example: User Story 1 Frontend

```bash
# Launch frontend tasks together (after backend usecase is done):
Task T031: "Create useSignup hook in client/src/features/auth/model/useSignup.ts"
Task T032: "Create useSignup test in client/src/features/auth/model/useSignup.test.ts"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (~10 min)
2. Complete Phase 2: Foundational (CRITICAL - ~1-2 hours)
3. Complete Phase 3: User Story 1 (~2-3 hours)
4. **STOP and VALIDATE**: Test User Story 1 independently using quickstart.md
5. Deploy/demo if ready - you now have a working signup flow!

**Total MVP Effort**: ~4-6 hours for a fully functional signup system

### Incremental Delivery

1. **Foundation (Phases 1-2)**: Ent schemas, migrations, session management, types → Foundation ready
2. **MVP (Phase 3 - US1)**: Full signup flow → Test independently → **Deploy/Demo MVP!**
3. **Security Layer (Phase 4 - US2 & US3)**: Password & email validation → Test independently → Deploy/Demo
4. **Workspace Customization (Phase 5 - US4)**: Workspace name input → Test independently → Deploy/Demo
5. **Polish (Phase 6)**: Logging, error handling, final validation → Deploy final version

### Parallel Team Strategy

With multiple developers:

1. **Together**: Complete Setup + Foundational (Phases 1-2)
2. **Once Foundational is done**:
   - Developer A: Backend for User Story 1 (T022-T030)
   - Developer B: Frontend for User Story 1 (T031-T037)
   - Developer C: User Story 2 & 3 enhancements (if starting early)
3. Stories complete and integrate independently

---

## Task Summary

**Total Tasks**: 60
**MVP Tasks (Setup + Foundational + US1)**: T001-T037 (37 tasks)
**Parallel Opportunities**: 21 tasks marked [P]

**Breakdown by Phase**:
- Phase 1 (Setup): 5 tasks
- Phase 2 (Foundational): 16 tasks
- Phase 3 (US1 - MVP): 16 tasks
- Phase 4 (US2 & US3): 7 tasks
- Phase 5 (US4): 6 tasks
- Phase 6 (Polish): 10 tasks

**Breakdown by User Story**:
- US1 (New User Account Creation): 16 tasks (MVP)
- US2 & US3 (Security & Email Uniqueness): 7 tasks
- US4 (Workspace Customization): 6 tasks

**Test Coverage**:
- Backend: 3 test files (auth_service_test.go, signup_usecase_test.go, signup_test.go integration)
- Frontend: 3 test files (validation.test.ts, useSignup.test.ts, validation functions)
- No component tests (per project conventions)

**Independent Test Criteria Met**: Each user story phase includes clear independent test criteria that can be validated without other stories being complete.

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Test files for functions and hooks are created alongside implementation (colocation pattern)
- No component tests per project testing strategy
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Docker commands must be run with `docker-compose exec` prefix per CLAUDE.md
- All paths are absolute or relative to project root
- Ent code generation (T010) is required before writing repository code