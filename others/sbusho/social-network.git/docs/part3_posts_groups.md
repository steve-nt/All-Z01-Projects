# Part 3 — Posts & Groups (Backend Implementation Walkthrough)

This doc is intended to be read **top-to-bottom** by someone who wants to understand what was implemented for Part 3:
- Files added/changed
- Endpoints and the correct calling order
- Authentication and authorization rules (privacy + group membership)
- Database tables touched and what gets written where
- Manual verification steps (curl/sqlite)

Code locations:
- Posts: `backend/pkg/posts/*`
- Groups + Events: `backend/pkg/groups/*`
- Routes wired in: `backend/main.go`
- Image upload (pre-existing): `backend/pkg/authentication/avatar_image.go` (`/api/upload-image`)

---

## 0) What changed in the codebase

### New packages added
- **Posts**: `backend/pkg/posts`
  - `routes.go`: registers routes under `/api/posts/*`
  - `handlers.go`: privacy filtering, post CRUD (create/list/view), post comments
- **Groups + Events**: `backend/pkg/groups`
  - `routes.go`: registers routes under `/api/groups/*`
  - `handlers.go`: group CRUD, invites, join requests, group posts/comments, events + responses

### Route wiring (minimal change)
`backend/main.go` now calls:
- `posts.SetupPostRoutes()`
- `groups.SetupGroupRoutes()`

No existing auth/profile/follow handlers were modified; Part 3 only adds new endpoints.

---

## 1) Authentication model used (shared by Part 3)

The project uses sessions via:
- Cookie name: `session`
- DB table: `Sessions`
- Helpers:
  - `utils.IsValidSession(cookieValue)`
  - `utils.GetUserIDFromSession(cookieValue)`

Part 3 endpoints enforce auth either with:
- `middleware.RequireAuthJSON` (route wrapper), or
- optional cookie checks inside handlers (for endpoints where auth is optional, like `/api/posts`).

## Posts

### Features implemented
- Create posts (authenticated)
- Post privacy levels:
  - `public`: visible to everyone
  - `almost_private`: visible only to **accepted followers** of the post author (plus the author)
  - `private`: visible only to users included in `Post_Visibility` (plus the author)
- List posts (feed) with privacy-aware filtering
- View single post with privacy check
- Comments on posts (list + create) with privacy check
- Optional image/GIF support via existing upload endpoint (`/api/upload-image`)

### DB tables used
- `Posts`
- `Post_Visibility` (used only for `privacy="private"`)
- `Posts_Images` (links uploaded filename to a post)
- `Comments`
- `Comments_Images`
- `Followers` (for `almost_private`)
- `Users` (author display name)

### Privacy filtering logic (core rule set)

For a post `(Posts.post_id, Posts.user_id, Posts.privacy)`:
- **Author can always see** their own posts.
- **`public`**: anyone can see.
- **`almost_private`**:
  - viewer must be logged in
  - require accepted follower relationship:
    - `Followers.follower_id = viewer_id`
    - `Followers.following_id = post_author_id`
    - `Followers.status = 'accepted'`
- **`private`**:
  - viewer must be logged in
  - require explicit allow-list entry:
    - `Post_Visibility.post_id = post_id`
    - `Post_Visibility.user_id = viewer_id`

This is applied in two places:
- The feed query in `GET /api/posts`
- The per-post permission check used by `GET /api/posts/view` and comments

### Endpoints

#### GET `/api/posts`
List posts with privacy filtering.

Query params:
- `user_id` (optional): filter by author user ID
- `limit` (optional, default 20, max 100)
- `offset` (optional, default 0)

Auth:
- Optional (uses cookie `session` if present)

Response (example):
```json
{
  "viewer_id": 2,
  "limit": 20,
  "offset": 0,
  "posts": [
    {
      "post_id": 1,
      "user_id": 1,
      "author": "alice",
      "content": "PUBLIC post",
      "privacy": "public",
      "created_at": "2026-01-22T13:14:25Z",
      "image_url": "/frontend/uploads/images/..."
    }
  ]
}
```

Implementation notes (what happens internally):
- Determines `viewer_id` from cookie (if present), otherwise `0`.
- Runs a **privacy-aware SQL query**:
  - guest: `WHERE privacy='public'`
  - logged-in: `WHERE (own OR public OR (almost_private + accepted follower) OR (private + Post_Visibility))`
- Optional `user_id` query adds `AND p.user_id = ?`.
- Adds a single `image_url` if the post has an entry in `Posts_Images` (first image only).

#### GET `/api/posts/view?post_id=`
Fetch a single post by ID with privacy enforcement.

Auth:
- Optional (privacy applies)

Implementation notes:
- Validates `post_id`.
- Checks whether the current viewer can see the post using the same privacy rules.
- Returns `403` if not allowed.

#### POST `/api/posts/create`
Create a post.

Auth:
- Required (cookie `session`)

Body:
```json
{
  "content": "hello",
  "privacy": "public | almost_private | private",
  "visible_to": [2,3],
  "image_filename": "123_1700000000_photo.gif"
}
```

Rules:
- `content` required
- `privacy` required
- If `privacy="private"`, `visible_to` must contain **at least 1 accepted follower** (not yourself)
- `image_filename` is optional and should be taken from the response of `/api/upload-image`

Response:
```json
{ "post_id": 123 }
```

Implementation notes (step-by-step):
- Validates JSON and required fields (`content`, `privacy`).
- Inserts into `Posts`.
- If `privacy="private"`:
  - requires `visible_to`
  - each user in `visible_to` must be an **accepted follower** of the post author (checked in `Followers`)
  - inserts into `Post_Visibility`
  - rejects if after validation there are 0 valid recipients
- If `image_filename` is provided:
  - stores a row in `Posts_Images` with `image_url = /frontend/uploads/images/<filename>`
  - file type is inferred from filename extension (jpg/jpeg/png/gif) to satisfy the migration CHECK constraint

#### GET/POST `/api/posts/comments`

##### GET `/api/posts/comments?post_id=`
List comments for a post (privacy enforced).

##### POST `/api/posts/comments`
Create a comment for a post (privacy enforced).

Auth:
- Required for POST

Body:
```json
{
  "post_id": 1,
  "content": "nice!",
  "image_url": "/frontend/uploads/images/x.png"
}
```

---

## Image/GIF support

Image upload endpoint is **pre-existing**:

#### POST `/api/upload-image`
Auth:
- Required

Multipart:
- field name: `image`

Response includes `filename`, which is then used in:
- `POST /api/posts/create` as `image_filename`

Notes:
- Upload supports JPEG/PNG/GIF (validated).
- Files stored under `backend/frontend/uploads/images/`.
- For **posts**, link the uploaded file by sending `image_filename` to `POST /api/posts/create`.
- For **post comments / group posts / group comments**, the handlers accept an `image_url`/`image_path` string (example: `/frontend/uploads/images/x.gif`) and store it in the relevant `*_Images` table.

---

## Groups

### Features implemented
- Create group (authenticated)
- Browse/list groups (public)
- View group details (public, includes membership flag if logged in)
- Invite user to group (member-only)
- Accept/decline invitation (invitee-only)
- Request to join group (authenticated, non-member)
- Accept/decline join request (creator-only)
- Group posts + comments (members-only, list + create)

### DB tables used
- `Groups`
- `Group_Members`
- `Group_Invitations`
- `Group_Join_Requests`
- `Group_Posts`
- `Group_Posts_Images`
- `Group_Comments`
- `Group_Comments_Images`
- `Users` (author display name)

### Membership & permission rules
- **Viewing/creating group posts, group comments, group events, event responses**:
  - user must be authenticated
  - user must be a member: `Group_Members(group_id, user_id)`
- **Inviting users to a group**:
  - inviter must be a member
  - invitee must not already be a member
- **Responding to join requests**:
  - only the group creator can accept/decline (`Groups.creator_id`)

### Endpoints

#### GET/POST `/api/groups`

##### GET `/api/groups`
List all groups.

##### POST `/api/groups`
Create a group.

Auth:
- Required for POST

Body:
```json
{ "group_name": "G1", "description": "test group" }
```

Response:
```json
{ "group_id": 1 }
```

Creator is automatically inserted into `Group_Members` with role `creator`.

Implementation notes:
- Uses a DB transaction:
  - insert into `Groups`
  - insert creator into `Group_Members` with role `creator`

#### GET `/api/groups/view?group_id=`
Group details + membership flag (if logged in).

Implementation notes:
- Fetches group from `Groups`.
- If a session cookie exists: checks membership and returns `is_member`.

#### POST `/api/groups/invite`
Invite a user to a group.

Auth:
- Required

Rules:
- Inviter must be a group member.
- Invitee must not already be a member.

Body:
```json
{ "group_id": 1, "user_id": 2 }
```

Response:
```json
{ "invitation_id": 1, "status": "pending" }
```

Implementation notes:
- Inserts `Group_Invitations` (status `pending`).
- Adds a notification for the invitee:
  - `Notifications.type = 'group_invitation'`
  - `user_id = invitee_id`
  - `related_user_id = inviter_id`
  - `related_group_id = group_id`

#### POST `/api/groups/invitations/respond`
Accept/decline an invitation (invitee-only).

Auth:
- Required

Body:
```json
{ "invitation_id": 1, "response": "accepted | declined" }
```

Response:
```json
{ "status": "accepted" }
```

Implementation notes:
- Validates the invitation belongs to the current user and is `pending`.
- On `accepted`:
  - inserts into `Group_Members` (role `member`)
- Adds notification to inviter:
  - `type = 'group_invitation_response'`

#### POST `/api/groups/join/request`
Request to join a group (non-member).

Auth:
- Required

Body:
```json
{ "group_id": 1 }
```

#### POST `/api/groups/join/respond`
Accept/decline a join request (creator-only).

Auth:
- Required

Body:
```json
{ "request_id": 1, "response": "accepted | declined" }
```

#### GET/POST `/api/groups/posts`
Members-only group posts.

Auth:
- Required

GET query:
- `group_id` required

POST body:
```json
{ "group_id": 1, "content": "hello group", "image_url": "/frontend/uploads/images/x.gif" }
```

#### GET/POST `/api/groups/comments`
Members-only group comments.

Auth:
- Required

GET query:
- `group_post_id` required

POST body:
```json
{ "group_post_id": 1, "content": "nice", "image_url": "/frontend/uploads/images/x.png" }
```

---

## Group Events

### Features implemented
- Create events in group (members-only)
- List events in group (members-only)
- Respond to an event: `going` / `not going` (members-only)

### DB tables used
- `Group_Events`
- `Group_Event_Responses`
- `Group_Members` (membership enforcement)

### Endpoints

#### GET/POST `/api/groups/events`
Auth:
- Required

GET query:
- `group_id` required

POST body:
```json
{
  "group_id": 1,
  "title": "Meet",
  "description": "desc",
  "event_datetime": "2026-01-22T20:30:00Z"
}
```

Response:
```json
{ "event_id": 1 }
```

#### POST `/api/groups/events/respond`
Auth:
- Required

Body:
```json
{ "event_id": 1, "response": "going | not going" }
```

Response:
```json
{ "status": "going" }
```

Implementation notes:
- Ensures user is a member of the group (derived from `Group_Events.group_id`).
- Upserts into `Group_Event_Responses` so each user has one response per event.
- Notifies the event creator (if not the same as responder):
  - `type = 'group_event_response'`
  - `user_id = creator_id`
  - `related_user_id = responder_id`
  - `related_group_id = group_id`
  - `related_event_id = event_id`

---

## Notifications (DB inserts)

This Part 3 implementation inserts notifications rows (best-effort) into:
- `Notifications`

Current notification types used:
- `group_invitation`
- `group_invitation_response`
- `group_join_request`
- `group_join_response`
- `group_event_created`
- `group_event_response`

Note:
- Notification **API/UI** to list/read notifications is not part of Part 3 backend here; only inserts are created.

---

## 2) Happy-path flows (recommended call order)

### A) Create a post (no image)
1. Login → get `session` cookie
2. `POST /api/posts/create` with `content` + `privacy`
3. `GET /api/posts` or `GET /api/posts?user_id=<author>` to verify

### B) Create a post with image/GIF (2-step)
1. Login → get `session` cookie
2. Upload image:
   - `POST /api/upload-image` (multipart form field `image`)
   - Read the returned `filename`
3. Create post:
   - `POST /api/posts/create` with `image_filename="<filename>"`
4. Verify feed returns `image_url`

### C) Private post (allow-list)
Precondition: the users in `visible_to` must be **accepted followers** of the author.
1. Make follower relationship accepted (via `/api/follow/request` + possibly accept if profile is private)
2. `POST /api/posts/create` with `privacy="private"` and `visible_to=[...]`
3. Verify:
   - guest cannot see the post
   - allowed user can see the post

### D) Invite flow (group management)
1. A creates group: `POST /api/groups`
2. A invites B: `POST /api/groups/invite`
3. B accepts: `POST /api/groups/invitations/respond` with `accepted`
4. B can now:
   - create group posts: `POST /api/groups/posts`
   - comment: `POST /api/groups/comments`
   - see events: `GET /api/groups/events`

### E) Join request flow (alternative way to join)
1. User requests to join: `POST /api/groups/join/request`
2. Creator accepts/declines: `POST /api/groups/join/respond`
3. On accepted, requester becomes a member.

### F) Event flow
1. Member creates event: `POST /api/groups/events`
2. Members list events: `GET /api/groups/events?group_id=...`
3. Member responds: `POST /api/groups/events/respond` with `going` or `not going`
4. Creator receives `group_event_response` notification row.

---

## 3) Common error cases (what they mean)

### `401 Unauthorized`
- You called an endpoint that requires auth without a valid `session` cookie.

### `403 Forbidden`
- You are authenticated, but you are not allowed:
  - Posts: you don’t satisfy privacy rules (not follower / not allow-listed)
  - Groups: you are not a group member (members-only endpoints)
  - Join request respond: you are not the group creator

### `400 Bad Request`
- Invalid JSON body (malformed JSON)
- Missing required fields (e.g. `group_id`, `post_id`, `content`)
- Invalid `privacy` value
- For private posts: `visible_to` is empty or contains users that are not accepted followers

---

## Verification (manual)

### Start from scratch (clean DB)
```bash
rm -f /home/tarara-x/Desktop/history/social-network/backend/data/social_network.db
cd /home/tarara-x/Desktop/history/social-network/backend
go run main.go
```

### Register + login 2 users
Use `/register` then `/login` to get cookie jars:
```bash
curl -i -c /tmp/cookiesA.txt -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data "email=a@test.com&password=Password1!"
```

### Validate privacy behavior
Check what B can see:
```bash
curl -s -b /tmp/cookiesB.txt "http://localhost:8080/api/posts?user_id=1" | python3 -m json.tool
```

### Validate invitations and events (and notifications)
Use sqlite to inspect:
```bash
sqlite3 /home/tarara-x/Desktop/history/social-network/backend/data/social_network.db \
"select id,type,user_id,related_user_id,related_group_id,related_event_id,message from Notifications order by id desc limit 20;"
```

---

## Known gaps (outside the Part 3 deliverables list)

These are not required by the “Posts & Groups Responsibility” list above, but they are common next steps:
- Add endpoints to **list pending invitations** for the current user (instead of fetching via sqlite during manual testing).
- Add endpoints to **list pending join requests** for a group creator.
- Add Notifications API/UI (Part 4+), since Part 3 only inserts notification rows.

