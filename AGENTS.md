# Agent Guide

This is a fullstack monorepo. Two agents work independently — one on the frontend, one on the backend — and coordinate through `.coordination/`.

---

## Agents

### Frontend Agent
**Lives in:** `frontend/`
**Stack:** Nuxt 4, shadcn-vue, Tailwind CSS 4, Pinia, VeeValidate + Zod, @vueuse/motion, PocketBase JS SDK
**Responsibilities:**
- Build pages (`app/pages/`), layouts (`app/layouts/`), components (`app/components/`)
- Manage client state with Pinia stores (`app/stores/`)
- Call the backend via PocketBase SDK (auth, collections) or `$fetch` for Gin endpoints
- Own all UI, routing, animations, and form validation
- Detailed rules: `docs/frontend-guidelines.md`

### Backend Agent
**Lives in:** `backend/`
**Stack:** Go, PocketBase (embedded, :8090), Gin (:8313)
**Responsibilities:**
- Define database schema as Go migrations (`pb_migrations/`) — never through the admin UI
- Add Gin routes in `internal/router/router.go` grouped by resource
- Write handlers (`internal/handlers/`) that stay thin — business logic goes in `internal/services/`)
- Add middleware (`internal/middleware/`) for auth, roles, and other common checks
- All responses via `models.Respond(c, status, message, data)`
- Detailed rules: `docs/backend-guidelines.md`

---

## Coordination

**At the start of every session:**
1. `ls .coordination/` — get the current filename (it's a timestamp, e.g. `20260412T120000Z.md`)
2. If different from your last remembered timestamp → read the file, it's short
3. Read your guideline doc and check `skills/` before writing anything from scratch

**After finishing work:**
1. Write `.coordination/<new-timestamp>.md` (keep it under 15 lines)
2. Delete the old timestamped file

Only include what the other agent must know: API contract changes, new collections, decisions that affect the other side.

---

## Common commands

All commands run from the repo root via `make`:

| Command | What it does |
|---|---|
| `make dev-frontend` | Start Nuxt dev server on :3000 |
| `make dev-backend` | Start Go + PocketBase on :8313/:8090 |
| `make install` | `pnpm install` inside `frontend/` |
| `make tidy` | `go mod tidy` inside `backend/` |
| `make migrate name=<name>` | Create a new backend migration file |
| `make up` | Build and start all services via Docker Compose |
| `make down` | Stop all Docker services |

---

## Conventions at a glance

| | Frontend | Backend |
|---|---|---|
| Language | TypeScript / Vue SFC | Go |
| New page | `app/pages/<name>.vue` | — |
| New route | — | group in `internal/router/router.go` |
| New collection | — | new file in `pb_migrations/` |
| Response shape | `{ message, data }` | `models.Respond(c, status, msg, data)` |
| Auth | PocketBase SDK → Pinia store | `middleware.RequireAuth(app)` |
| Validation | VeeValidate + Zod | `c.ShouldBindJSON` + explicit checks |

## Ports

| Service | Port |
|---|---|
| Frontend | 3000 |
| Gin API | 8313 |
| PocketBase admin + DB | 8090 |

## Skills

Copy-paste templates in `skills/` — use these before writing from scratch:

| Skill | Path |
|---|---|
| Auth pages (login + register) | `skills/frontend/auth-pages.md` |
| Store manager pattern + API wrapper | `skills/frontend/auth-manager.md` |
| Dashboard sidebar layout | `skills/frontend/dashboard-shell.md` |
| Data table | `skills/frontend/data-table.md` |
| CRUD form component | `skills/frontend/crud-form.md` |
| CRUD handler | `skills/backend/crud-handler.md` |
| PocketBase service | `skills/backend/pb-service.md` |
| Auth + role middleware | `skills/backend/auth-middleware.md` |
| Schema migration | `skills/backend/migration.md` |
