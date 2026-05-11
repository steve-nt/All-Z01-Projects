# Vue.js Frontend Plan (SPA) — Project Split into 5 Parts

This document explains the recommended approach to build a **Vue.js frontend** for this repo’s Go backend, and proposes a clear **5-part** division of work aligned with the existing backend Parts 1–3 and the planned real-time work (Part 4 notes).

---

## What to build (recommended)

- Build a **Vue SPA (Single Page App)** using:
  - **Vite** (dev server + build tooling)
  - **Vue Router** (routes/pages)
  - **Pinia** (state management)
- Treat the Go backend as a **JSON API** (the Vue app should not rely on server-rendered HTML templates).

### Do we need HTML pages?

- **No** (not as multiple HTML pages).
- A Vue SPA typically has a single entry HTML file (e.g. `index.html`) managed by Vite.
- “Pages” are implemented as **Vue routes + components**:
  - `/login`, `/register`, `/feed`, `/profile/:id`, `/groups/:id`, etc.

Note: The backend currently includes minimal placeholder templates under `backend/frontend/templates/` (e.g. `login.html`, `register.html`). You can keep them as placeholders/legacy, but the Vue SPA does not need them.

---

## Backend integration note (CORS + cookies)

If the Vue frontend runs on a different origin/port than the backend, browsers enforce:
- **CORS** for HTTP API calls
- **Cookie restrictions** (SameSite / Secure) for session cookies

### Recommended for development (least friction)

Use a **Vite dev proxy** so your frontend can call the backend without fighting CORS/cookie issues during development.

Why:
- Avoids common preflight/CORS headaches
- Avoids cross-origin cookie issues (especially around `SameSite=None; Secure`)

### Alternative (fine, but more moving parts)

Add proper **CORS middleware** on the Go backend (allow specific frontend origins + allow credentials) and adjust cookie settings for cross-origin cookies.

Use this when:
- You intentionally want frontend and backend on different origins in dev
- You’re preparing production deployment with separate domains

---

## Project split into 5 parts (frontend deliverables)

This split is designed to be incremental, testable, and aligned with the current backend state and docs.

Iana
### Part 1) Frontend foundation (skeleton)

Deliverables:
- Project scaffolding (Vue + Vite)
- Routing structure (Vue Router)
- Global state (Pinia): auth/session state, user/profile caching
- API client conventions:
  - Base URL strategy (proxy vs direct backend URL)
  - Standard error handling
  - Automatic redirect/handling for `401 Unauthorized`
- Basic layout/navigation shell

---

## How to start Part 1 (Frontend foundation) — if you’re new to Vue

### Goal of Part 1

By the end of Part 1 you should have a running Vue app with:
- Routes + layout (navbar + empty pages)
- An auth/session store (even before the login UI exists)
- A consistent way to call the backend API (and handle `401/403` cleanly)
- A development setup that avoids CORS pain (recommended: Vite proxy)

### What to learn first (minimum Vue knowledge)

- Components: template + script + style, **props** and **emits**
- Reactivity: `ref`, `reactive`, computed values
- Lifecycle: `onMounted` (for “check session on page load”)
- Vue Router: routes + route params (e.g. `/profile/:id`)
- Pinia: global state (auth/session, cached data)

### Recommended execution order (do these in sequence)

#### 1) Decide your dev integration approach (important)

- Recommended: use a **Vite dev proxy** so the browser doesn’t treat your API calls as cross-origin during development.
- Alternative: run frontend and backend on different origins and implement backend **CORS + cookie cross-origin** rules.

#### 2) Create the frontend skeleton (no features yet)

- Create the Vue app (Vite)
- Confirm you can run the dev server and see the app in the browser

#### 3) Define routes/pages early (placeholders are fine)

Create empty “page” views for:
- Public: `Login`, `Register`
- Auth-only (later parts): `Feed`, `MyProfile`, `GroupsList`, `GroupView`, `Notifications`, `Messages`

Tip: keep page components thin; put logic in stores/services.

#### 4) Add a global layout + navigation

- Top navigation that changes by auth state:
  - Logged out: Login/Register
  - Logged in: Feed/Profile/Groups + Logout
- Add an “initial session loading” state on app startup.

#### 5) Add Pinia stores (foundation stores only)

Start with:
- Auth store:
  - state: `loggedIn`, `userID`, `nickname`
  - action: “check session” via `GET /api/auth/status`
  - later: “login”, “logout”
- (Optional) Profile cache store: cache profiles by `user_id` for later parts

#### 6) Standardize API calls for the team

Agree on conventions now:
- One API client wrapper (base URL/proxy path in one place)
- Shared error handling rules:
  - `401` → clear auth state + redirect to login
  - `403` → show “private/forbidden” UI (don’t redirect)
  - network errors → show toast/banner
- Choose either `fetch` or Axios (either is fine; consistency is the priority)

#### 7) Add router guards

- Protect auth-only routes (redirect to `/login` if logged out)
- Handle first-load correctly:
  - run “check session” once before deciding redirects

### Definition of done (Part 1)

Part 1 is complete when:
- App loads with a shared layout and working navigation
- Routes exist and switching routes works
- Refreshing the page keeps you logged in (when backend session cookie is valid) by using `GET /api/auth/status`
- Visiting a protected route while logged out redirects to login
- Everyone can run the frontend using the same documented steps

Andriana
### Part 2) Authentication UI (maps to backend Part 1)

Deliverables:
- Register UI (support for multipart/avatar can be added later if needed)
- Login / Logout
- Session restore on refresh via `GET /api/auth/status`
- Route guards (auth-only routes)

Sofia
### Part 3) Profiles + Follow system UI (maps to backend Part 2)

Deliverables:
- Profile view (handle privacy rules; show friendly UI on `403 Forbidden`)
- Privacy toggle (for the logged-in user)
- Follow flows:
  - request follow
  - accept/decline follow request
  - unfollow
- Followers/following lists
- Pending follow requests view


Charoula
### Part 4) Posts + Groups + Events UI (maps to backend Part 3)

Deliverables:
- Feed and post creation
  - privacy: `public`, `almost_private`, `private`
  - private allow-list UI (must be accepted followers)
- Image/GIF flow (2-step): upload → create post referencing returned filename
- Post comments: list/create
- Groups:
  - list/create/view
  - invite flow
  - join request flow
  - membership enforcement UX (handle `403` cleanly)
- Group posts/comments
- Group events + responses


Georgia
### Part 5) Real-time layer + notifications/messaging UI (matches Part 4 phases doc)

Deliverables (recommended order):
- WebSocket connection foundation:
  - connection lifecycle + reconnection
  - authenticate the connection using existing session system
- Real-time notifications UI:
  - toasts + notifications panel
  - “new notification” updates without refresh
- Private messaging UI (1-to-1)
- Group chat UI
- Polish:
  - emoji/unicode handling
  - UX refinements (unread badges, message timestamps, etc.)

---

## Notes about current backend readiness (from repo docs)

- Backend **Part 1 (auth)** and **Part 2 (profiles/follows)** are implemented and testable.
- Backend **Part 3 (posts/groups/events)** is documented and includes endpoints and test flows.
- Real-time work is planned as phased steps; notifications are already being inserted into the DB, but a full notifications API/UI is a logical next step for Part 5.


