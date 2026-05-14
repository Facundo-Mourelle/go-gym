# AGENTS.md - Go Gym Development Guide

## Project Overview
Full-stack gym tracking app: Go backend (hexagonal architecture) + React/TypeScript frontend. Shared domain concepts (exercises, sessions, routines, metrics) span both layers.

**Key paths:**
- Backend: `cmd/api/main.go` (entry), `internal/` (domain, services, handlers, repos)
- Frontend: `gym-tracker-frontend/src/` (React components, stores, API client)
- Database: PostgreSQL 15 via Docker Compose (optional fallback to in-memory repos)

## Quick Start Commands

```bash
# Full stack (DB + backend + frontend)
./start.sh

# Individual services
./start_backend.sh    # go run cmd/api/main.go
./start_frontend.sh   # npm run dev (from gym-tracker-frontend/)
./stop.sh             # Stops all services

# Tests
go test ./...                    # Backend (standard testing package)
cd gym-tracker-frontend && npm test -- --run  # Frontend
```

## Architecture Notes

### Backend (Go 1.25)
- **Hexagonal**: Domain → Services → Handlers/Repos
- **Repos**: Dual-layer (PostgreSQL + in-memory fallback)
- **Auth**: JWT (golang-jwt/v5), 7-day expiration
- **Domain modules**: `calculator/` (volume, 1RM), `resistance/` (mechanical advantage profiles)
- **No migrations**: Schema changes require manual SQL + discussion

### Frontend (React 19 + TypeScript)
- **State**: Zustand (global stores)
- **Charts**: Recharts (progress visualization)
- **Build**: Vite + TypeScript, ESLint
- **API client**: Axios (baseURL: `http://localhost:8080`)

### Domain Concepts (Backend ↔ Frontend)
- **Exercise**: Equipment + movement patterns (primary/secondary) + ROM
- **Session**: Live workout tracking with sets (reps, weight, RPE/RIR)
- **Routine**: Reusable workout templates
- **Metrics**: Dashboard (weekly stats), progress (1RM trends)
- **Equipment**: Machines/free weights with resistance profiles

## Agent Workflow

### Domain-Based Agents
Organize by feature (not layer). Each domain agent owns backend + frontend changes:

1. **Auth Agent**: Register/login, JWT middleware, user context
2. **Exercise Agent**: Exercise CRUD, templates, patterns, equipment
3. **Session Agent**: Session lifecycle, set recording, completion
4. **Routine Agent**: Routine templates, exercise grouping
5. **Metrics Agent**: Dashboard, progress charts, calculations

### Before Starting Work
- **Backend changes**: Verify `go run cmd/api/main.go` starts without errors
- **Frontend changes**: Verify `npm run dev` (from `gym-tracker-frontend/`) starts without errors
- **Database schema**: If needed, ask user first. Explain why (data consistency, migration risk). Do NOT auto-migrate.

### After Code Changes
**Always run both:**
```bash
# Backend
go run cmd/api/main.go  # Verify startup
go test ./...           # Run tests (standard testing package)

# Frontend
cd gym-tracker-frontend
npm run lint            # ESLint
npm run build           # TypeScript check + Vite build
```

If either fails, fix before committing.

## Key Conventions

### Go
- Interfaces in `repository/` (UserRepository, ExerciseRepository, etc.)
- Services inject repos + domain logic
- Handlers use `Protected()` middleware for auth routes
- Error handling: Return errors, don't panic
- No comments unless explaining non-obvious logic

### React/TypeScript
- Components in `src/components/`
- Stores in `src/stores/` (Zustand)
- API calls in `src/api/` or within stores
- Types in `src/types/`
- No comments unless explaining non-obvious logic

### Environment
- Backend: `.env` (DATABASE_URL, JWT_SECRET, PORT)
- Frontend: `.env` (VITE_API_BASE_URL if needed)
- Both: Load via config/env files, not hardcoded

## Common Pitfalls

1. **Database fallback**: If `DATABASE_URL` is empty, backend uses in-memory repos. Data is lost on restart.
2. **CORS**: Already configured in `cmd/api/main.go` (allows all origins for dev).
3. **JWT expiration**: 7 days. Frontend must handle token refresh or re-login.
4. **Port conflicts**: Backend (8080), Frontend (5173), DB (5432). Check `stop.sh` if ports are stuck.
5. **Frontend API URL**: Ensure axios baseURL matches backend port (default `http://localhost:8080`).

## Testing Strategy

- **Backend**: Standard `testing` package. Write tests alongside features.
- **Frontend**: Jest/Vitest (if configured). Currently minimal test coverage.
- **Integration**: Manual testing via UI + API endpoints.

## File Structure Reference

```
go-gym/
├── cmd/api/main.go              # Backend entry, service wiring
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTP handlers (auth, exercise, session, etc.)
│   │   ├── middleware/          # JWT auth middleware
│   │   └── router.go            # Route registration
│   ├── domain/                  # Pure domain logic
│   │   ├── calculator/          # Volume, 1RM formulas
│   │   ├── resistance/          # Mechanical advantage profiles
│   │   └── *.go                 # Exercise, Session, Routine, etc.
│   ├── service/                 # Business logic (cases of use)
│   ├── repository/              # Data access interfaces + implementations
│   │   ├── postgres/            # PostgreSQL implementations
│   │   └── memory/              # In-memory fallback
│   └── config/                  # Config loading
├── gym-tracker-frontend/
│   ├── src/
│   │   ├── components/          # React components
│   │   ├── stores/              # Zustand stores
│   │   ├── api/                 # API client
│   │   └── types/               # TypeScript types
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml           # PostgreSQL service
├── start.sh / stop.sh           # Service orchestration
└── .env                         # Local config (not committed)
```

## Debugging Tips

- **Backend logs**: `backend.log` (from `start_backend.sh`)
- **Frontend logs**: `frontend.log` (from `start_frontend.sh`)
- **DB connection**: Check `DATABASE_URL` in `.env`. If empty, in-memory mode is active.
- **API health**: `curl http://localhost:8080/health`
- **Frontend dev**: `http://localhost:5173` (Vite HMR enabled)
