# Setup + Part 2 Authentication UI (Vue) — Steps We Did

This document summarizes **exactly what was done** to get the project running locally and to implement **Part 2: Authentication UI (maps to backend Part 1)**.

## Part 2 checklist (deliverables)

- ✅ **Register UI**
  - Frontend form: `frontend/src/views/RegisterView.vue`
  - Backend endpoint: `POST /register` (SPA uses `Accept: application/json`)
- ✅ **Login UI**
  - Frontend form: `frontend/src/views/LoginView.vue`
  - Backend endpoint: `POST /login` (SPA uses `Accept: application/json`)
- ✅ **Logout**
  - UI button: `frontend/src/App.vue`
  - Store action: `frontend/src/stores/auth.js`
  - Backend endpoint: `POST /logout` (SPA uses `Accept: application/json`)
- ✅ **Session restore on refresh**
  - Backend endpoint: `GET /api/auth/status`
  - Called at startup in: `frontend/src/main.js`
  - Stored in: `frontend/src/stores/auth.js`
- ✅ **Route guards (auth-only routes)**
  - Guard logic in: `frontend/src/router/index.js`
  - Behavior:
    - Protected routes redirect to `/login?redirect=<originalPath>`
    - Logged-in users visiting `/login` or `/register` redirect to `/feed`

## Prerequisites

- **Go**: project uses a Go toolchain declared in `backend/go.mod` (Go 1.24.x toolchain).
- **Node.js + npm**: frontend is a Vite/Vue app.
- **Linux build deps for SQLite** (needed for `github.com/mattn/go-sqlite3`):

```bash
sudo apt update
sudo apt install -y build-essential libsqlite3-dev
```

## Run the project (local dev)

### Backend (Go + SQLite migrations)

```bash
cd /home/tarara-x/Desktop/social/social-network/backend
go mod download
go run main.go
```

- **Backend URL**: `http://localhost:8080`
- **Health check**: `http://localhost:8080/health`
- **SQLite DB path** (auto-created): `backend/data/social_network.db`
- **Migrations** auto-run on startup from: `backend/pkg/db/migrations/sqlite/`

### Frontend (Vue + Vite)

```bash
cd /home/tarara-x/Desktop/social/social-network/frontend
npm install
npm run dev
```

- **Frontend URL** (default): `http://localhost:5173`
- Dev proxy is already configured in `frontend/vite.config.js` so these routes go to the backend:
  - `/api/*` → `http://localhost:8080`
  - `/login`, `/register`, `/logout`, `/health`
  - `/ws` (WebSocket)

## Fix we did: `vite: Permission denied`

When running `npm run dev`, we hit:

```
sh: 1: vite: Permission denied
```

Cause: `node_modules/.bin/vite` was present but **not executable** (mode `664`).

Fix (one-time):

```bash
cd /home/tarara-x/Desktop/social/social-network/frontend
chmod +x node_modules/.bin/vite node_modules/vite/bin/vite.js
```

## Fix we did: Node 18 vs “latest” Vite

After fixing permissions, Vite complained because the project had `"vite": "latest"` and the installed Vite required **Node 20+**.

Fix we implemented: **pin versions** in `frontend/package.json` so it works on **Node 18**:

- Pinned **Vite** to `^5.x` and `@vitejs/plugin-vue` to `^5.x`
- Pinned `vue`, `pinia`, `vue-router` to stable versions
- Added `"engines": { "node": ">=18.0.0" }`

Then we did a clean reinstall:

```bash
cd /home/tarara-x/Desktop/social/social-network/frontend
rm -rf node_modules package-lock.json
npm install
```

Result: `npm run dev` starts successfully on Node 18.

## Part 2 deliverables implemented (Authentication UI)

### Deliverable: Register UI

- Implemented a real register form in:
  - `frontend/src/views/RegisterView.vue`
- It submits **form data** to the backend:
  - `POST /register`
- Field names match backend handler:
  - `email`, `password`, `confirmPassword`, `first_name`, `last_name`, `date_of_birth`
  - optional: `nickname`, `about_me`
  - `is_public` and `is_active`

On success, the UI navigates to:
- `/login?registered=1` (shows a “Registration successful” message)

### Deliverable: Login UI

- Implemented a real login form in:
  - `frontend/src/views/LoginView.vue`
- It submits to:
  - `POST /login`
- Uses cookies (`credentials: "include"`) so the backend session cookie is stored in the browser.

On success:
- The frontend calls `authStore.checkSession()` (GET `/api/auth/status`)
- Then redirects to:
  - `?redirect=...` if present (from route guards), otherwise `/feed`

### Deliverable: Logout

- Logout button is already present in the global nav:
  - `frontend/src/App.vue`
- Logout action calls:
  - `POST /logout`
- Then clears auth store and routes user to `/login`

### Deliverable: Session restore on refresh

Already present and verified:
- On app start (after router ready), frontend does:
  - `GET /api/auth/status`
  - in `frontend/src/main.js`
- Auth store tracks:
  - `sessionChecked`, `loggedIn`, etc. in `frontend/src/stores/auth.js`

### Deliverable: Route guards (auth-only routes)

Already present and improved:
- In `frontend/src/router/index.js`
- Behavior:
  - If route has `meta.requiresAuth` and user is not logged in → redirect to login with:
    - `?redirect=<originalPath>`
  - If user is logged in and tries to visit `/login` or `/register` → redirect to `/feed`

## Backend change we made to support SPA (important)

The backend auth endpoints originally behaved like classic HTML routes (render templates / redirect).
For the Vue SPA, we added **JSON responses when the request includes** `Accept: application/json`.

Updated endpoints:
- `POST /register`
  - HTML flow: redirects to `/login?success=registration`
  - SPA flow (Accept JSON): returns `201` JSON `{ ok: true, userID: <id> }`
- `POST /login`
  - HTML flow: redirects to `/`
  - SPA flow (Accept JSON): returns `200` JSON `{ ok: true }` (and sets the session cookie)
- `POST /logout`
  - HTML flow: redirects to `/`
  - SPA flow (Accept JSON): returns `200` JSON `{ ok: true }`

Implementation files:
- `backend/pkg/authentication/registerHandler.go`
- `backend/pkg/authentication/loginHandler.go`
- `backend/pkg/authentication/authHandler.go`

## Smoke test commands (optional)

You can test auth without the UI:

```bash
EMAIL="user$(date +%s)@example.com"
PASS="Abcdef1!"

# Register (JSON)
curl -i -H "Accept: application/json" \
  -F "email=$EMAIL" \
  -F "password=$PASS" \
  -F "confirmPassword=$PASS" \
  -F "first_name=Test" \
  -F "last_name=User" \
  -F "date_of_birth=2000-01-01" \
  -F "is_public=true" \
  -F "is_active=true" \
  http://localhost:8080/register

# Login (JSON) with cookie jar
curl -i -c /tmp/social_cookiejar.txt -b /tmp/social_cookiejar.txt \
  -H "Accept: application/json" \
  -d "email=$EMAIL&password=$PASS" \
  -X POST http://localhost:8080/login

# Auth status
curl -i -b /tmp/social_cookiejar.txt -H "Accept: application/json" \
  http://localhost:8080/api/auth/status

# Logout (JSON)
curl -i -c /tmp/social_cookiejar.txt -b /tmp/social_cookiejar.txt \
  -H "Accept: application/json" \
  -X POST http://localhost:8080/logout
```

## Notes / gotchas

- If the backend was already running before the JSON changes, you must **restart it** to pick up the new behavior:

```bash
ss -ltnp | grep ':8080'
kill <pid>
cd /home/tarara-x/Desktop/social/social-network/backend
go run main.go
```

- For development, the Vue app should call backend endpoints using relative paths (like `/api/...`, `/login`) so Vite’s proxy handles routing and you avoid CORS issues.

