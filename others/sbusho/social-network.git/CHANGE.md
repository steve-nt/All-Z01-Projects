# Changes since 49a3eb50af954f4d2668843d0c727caf88e153ce

- Baseline: `49a3eb50af954f4d2668843d0c727caf88e153ce`
- Compared to: `c16be046f2afd947383f70cf08d0e79f3e755525`
- Date generated: 2026-02-12

## Summary (high level)

- Implemented SPA Profile + Follow system UI (profile view, privacy toggle, follow/unfollow, pending requests, followers/following views).
- Added backend user discovery endpoint `GET /api/users/search` and wired it into auth routes.
- Added People discovery page (`/people`) with debounced search and follow actions.
- Replaced placeholder groups pages with functional groups browse/detail flows and wired group API calls.
- Added create-group form on `/groups` with validation and post-create navigation.
- Expanded profile page to show reachable followers/following links and own follow-request count.
- Standardized frontend API access through `socialApi` and Pinia stores for profile/follow/groups/people.
- Added reusable Vue components for profile/follow UI pieces and pending-request item actions.

## Uncommitted changes

- Working tree is clean (`git status`: nothing to commit).
- No staged or unstaged local changes are present beyond `HEAD`.

## Backend changes

### New/changed API endpoints

1. `GET /api/users/search`
- Handler: `UsersSearchHandler` in `backend/pkg/profile/profileHandler.go`
- Route registration: `backend/pkg/authentication/routes.go` via `SetupAuthRoutes()`
- Auth: required (`middleware.RequireAuthJSON`)
- Query params:
  - `q` optional (trimmed)
  - `limit` optional (default `20`, max `50`, invalid/non-positive => `400`)
- Response JSON:
```json
{
  "users": [
    { "user_id": 1, "nickname": "alice", "avatar": "/uploads/a.png" }
  ]
}
```
- Status codes:
  - `200` success
  - `400` invalid `limit`
  - `401` unauthorized
  - `500` query/scan failures
- Notes:
  - Excludes current authenticated user from results (`user_id != currentUserID`).
  - Matches on `nickname`, `email`, `first_name`, `last_name` with NULL-safe `COALESCE(... LIKE ...)` and substring wildcard.

2. Existing follow/profile endpoints consumed by new frontend work (no route-path changes)
- `GET /api/profile?user_id=...`
- `POST /api/profile/privacy`
- `POST /api/follow/request`
- `POST /api/follow/accept`
- `POST /api/follow/decline`
- `POST /api/follow/unfollow`
- `GET /api/followers?user_id=...`
- `GET /api/following?user_id=...`
- `GET /api/follow/requests`

### Database / migrations

- No new migration files added.
- No schema/index DDL changes in this commit range.

### Other backend changes

- `backend/pkg/authentication/routes.go`
  - Added `/api/users/search` registration to `SetupAuthRoutes()`.
  - Updated route list comment block to include the new endpoint.
- `backend/pkg/profile/profileHandler.go`
  - Added `UsersSearchHandler` implementation.
  - Added `strings` import for query normalization.

## Frontend changes

### New pages/routes

- `/people` → `frontend/src/views/PeopleView.vue`
  - Auth-guarded people discovery with debounced search and follow actions.
- `/groups` now uses `frontend/src/views/GroupsBrowseView.vue`
  - Lists groups, supports local filter, and includes create-group form.
- `/groups/:id` → `frontend/src/views/GroupView.vue` (upgraded from placeholder)
  - Loads group details and membership state from backend.
- `/follow/requests` → `frontend/src/views/PendingFollowRequestsView.vue`
- `/profile/:id/followers` → `frontend/src/views/FollowersListView.vue`
- `/profile/:id/following` → `frontend/src/views/FollowingListView.vue`

### New stores/services/components

- `frontend/src/services/socialApi.js`
  - Added methods for profile/follow/groups/people APIs:
    - `getProfile`, `updatePrivacy`, `requestFollow`, `acceptFollowRequest`, `declineFollowRequest`, `unfollow`
    - `getFollowers`, `getFollowing`, `getPendingFollowRequests`
    - `getGroups`, `createGroup`, `getGroupDetails`
    - `searchUsers`
- `frontend/src/stores/profile.js`
  - Profile load/update state with explicit 403 private-profile handling.
- `frontend/src/stores/follow.js`
  - Follow/follower/following/pending-request state and actions.
  - Pending-action tracking for button disabling.
  - `requestFollow` now validates and normalizes numeric target user id.
- `frontend/src/stores/groups.js`
  - Group list/detail/create state/actions.
- `frontend/src/stores/people.js`
  - Search query/result state + async search action.
- `frontend/src/stores/auth.js`
  - Refactored to `currentUser` model with getters (`isLoggedIn`, `userId`, `nickname`), `loadAuthStatus`.
- New reusable components:
  - `frontend/src/components/FollowButton.vue`
  - `frontend/src/components/FollowersFollowingPreview.vue`
  - `frontend/src/components/PendingRequestItem.vue`
  - `frontend/src/components/PrivacyToggle.vue`
  - `frontend/src/components/ProfileHeader.vue`

### UI/UX behavior changes

- `frontend/src/views/ProfileView.vue`
  - Replaced placeholder with real profile flow:
    - fetch profile by route id (including `/profile/me` handling)
    - show friendly private-profile state on 403
    - show follow/unfollow actions for non-self views
    - show privacy toggle for own profile
    - fetch followers/following counts and render direct links
    - show own pending follow-request count/link
    - graceful 403 handling for followers/following lists as `Private`
- `frontend/src/views/PeopleView.vue`
  - Debounced search input (300ms), follow button pending states, follow result labels, visible error messaging.
- `frontend/src/views/PendingFollowRequestsView.vue`
  - Accept/decline actions remove handled entries from list.
- `frontend/src/views/GroupsBrowseView.vue`
  - Create-group UX with required fields and inline error feedback.
- `frontend/src/App.vue`
  - Added nav links: `Find People`, `Follow Requests` (and updated groups/profile navigation behavior from prior placeholder phase).
- `frontend/src/router/index.js`
  - Expanded auth-guarded route map for new views.
- `frontend/src/assets/base.css`
  - Added styles for profile header, follow previews, pending/user/group lists, and groups browse/search visuals.

## DevOps / tooling changes

- No intentional build-system or runtime infra changes detected (no Docker/env/script additions in this range).
- One generated file changed:
  - `frontend/node_modules/.vite/deps/_metadata.json` (Vite dependency metadata)
- Notes / Risks:
  - `frontend/node_modules/.vite/deps/_metadata.json` is environment-generated and typically should not be committed in source control unless intentionally tracked.

## File-level change list

Source of truth: `git diff --name-status 49a3eb50af954f4d2668843d0c727caf88e153ce..HEAD`

### Added

- `frontend/src/components/FollowButton.vue`
- `frontend/src/components/FollowersFollowingPreview.vue`
- `frontend/src/components/PendingRequestItem.vue`
- `frontend/src/components/PrivacyToggle.vue`
- `frontend/src/components/ProfileHeader.vue`
- `frontend/src/services/socialApi.js`
- `frontend/src/stores/follow.js`
- `frontend/src/stores/groups.js`
- `frontend/src/stores/people.js`
- `frontend/src/stores/profile.js`
- `frontend/src/views/FollowersListView.vue`
- `frontend/src/views/FollowingListView.vue`
- `frontend/src/views/GroupsBrowseView.vue`
- `frontend/src/views/PendingFollowRequestsView.vue`
- `frontend/src/views/PeopleView.vue`

### Modified

- `backend/pkg/authentication/routes.go`
- `backend/pkg/profile/profileHandler.go`
- `frontend/node_modules/.vite/deps/_metadata.json`
- `frontend/src/App.vue`
- `frontend/src/assets/base.css`
- `frontend/src/router/index.js`
- `frontend/src/stores/auth.js`
- `frontend/src/views/GroupView.vue`
- `frontend/src/views/ProfileView.vue`

### Deleted

- None

## Breaking changes / required actions

- Backend must be restarted to pick up `/api/users/search` route and handler changes.
- Frontend should be rebuilt/restarted (`npm run dev` or `npm run build`) to load new routes/views/stores.
- No DB migration step required for this range.

## How to test (smoke tests)

### Run backend

```bash
cd backend
go run .
```

### Run frontend

```bash
cd frontend
npm run dev
```

### Curl smoke tests

1. Health check
```bash
curl -i http://localhost:8080/health
```

2. Auth status (anonymous)
```bash
curl -i http://localhost:8080/api/auth/status
```

3. Login (capture cookie)
```bash
curl -i -c cookie.txt -X POST http://localhost:8080/login \
  -F "email=<email>" \
  -F "password=<password>"
```

4. User search (auth required)
```bash
curl -i -b cookie.txt "http://localhost:8080/api/users/search?q=ali&limit=10"
```

5. Follow request
```bash
curl -i -b cookie.txt -X POST http://localhost:8080/api/follow/request \
  -H 'Content-Type: application/json' \
  -d '{"user_id":2}'
```

6. Pending follow requests
```bash
curl -i -b cookie.txt http://localhost:8080/api/follow/requests
```

7. Group browse
```bash
curl -i -b cookie.txt http://localhost:8080/api/groups
```

8. Create group + view group
```bash
curl -i -b cookie.txt -X POST http://localhost:8080/api/groups \
  -H 'Content-Type: application/json' \
  -d '{"group_name":"Smoke Group","description":"Created by smoke test"}'

curl -i -b cookie.txt "http://localhost:8080/api/groups/view?group_id=1"
```
