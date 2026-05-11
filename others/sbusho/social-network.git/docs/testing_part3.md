# Part 3 — Testing Guide (Posts & Groups)

This file contains **copy/paste** commands to test Part 3 features end-to-end:
- Posts (privacy + images + comments)
- Groups (create, invite, join request, membership)
- Group posts/comments
- Group events + responses
- Notifications DB inserts

Assumptions:
- Backend runs on `http://localhost:8080`
- You run commands from any directory (paths are absolute where needed)

---

## 0) Start from scratch (clean DB)

```bash
rm -f /home/tarara-x/Desktop/history/social-network/backend/data/social_network.db
```

Start the server (terminal 1):

```bash
cd /home/tarara-x/Desktop/history/social-network/backend
go run main.go
```

Quick health check (terminal 2):

```bash
curl -i http://localhost:8080/health
```

---

## 1) Register 2 users (Alice=A, Bob=B)

```bash
curl -i -X POST http://localhost:8080/register \
  -F "email=a@test.com" \
  -F "password=Password1!" \
  -F "confirmPassword=Password1!" \
  -F "first_name=Alice" \
  -F "last_name=Test" \
  -F "date_of_birth=2000-01-01" \
  -F "nickname=alice" \
  -F "about_me=hi" \
  -F "is_public=true" \
  -F "is_active=true"

curl -i -X POST http://localhost:8080/register \
  -F "email=b@test.com" \
  -F "password=Password1!" \
  -F "confirmPassword=Password1!" \
  -F "first_name=Bob" \
  -F "last_name=Test" \
  -F "date_of_birth=2000-01-01" \
  -F "nickname=bob" \
  -F "about_me=yo" \
  -F "is_public=true" \
  -F "is_active=true"
```

Expected:
- HTTP `303 See Other` redirect to `/login?success=registration`

---

## 2) Login and store cookies

```bash
curl -i -c /tmp/cookiesA.txt -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data "email=a@test.com&password=Password1!"

curl -i -c /tmp/cookiesB.txt -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data "email=b@test.com&password=Password1!"
```

Expected:
- HTTP `303 See Other` with `Set-Cookie: session=...`

Get user IDs:

```bash
curl -s -b /tmp/cookiesA.txt http://localhost:8080/api/auth/status | python3 -m json.tool
curl -s -b /tmp/cookiesB.txt http://localhost:8080/api/auth/status | python3 -m json.tool
```

Expected:
- Alice: `"userID": 1`
- Bob: `"userID": 2`

---

## 3) Posts — privacy test

### 3.1 Create public + almost_private posts as Alice

```bash
curl -i -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/posts/create \
  -H "Content-Type: application/json" \
  -d '{"content":"PUBLIC post","privacy":"public"}'

curl -i -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/posts/create \
  -H "Content-Type: application/json" \
  -d '{"content":"ALMOST_PRIVATE post","privacy":"almost_private"}'
```

### 3.2 Bob views Alice posts BEFORE following

```bash
curl -s -b /tmp/cookiesB.txt "http://localhost:8080/api/posts?user_id=1" | python3 -m json.tool
```

Expected:
- Bob should see **only** `privacy="public"` posts at this point.

### 3.3 Bob follows Alice (public profile ⇒ auto-accepted)

```bash
curl -i -b /tmp/cookiesB.txt -X POST http://localhost:8080/api/follow/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'
```

### 3.4 Bob views Alice posts AFTER following

```bash
curl -s -b /tmp/cookiesB.txt "http://localhost:8080/api/posts?user_id=1" | python3 -m json.tool
```

Expected:
- Bob now sees `public` **and** `almost_private`.

### 3.5 Create private post visible only to Bob

```bash
curl -i -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/posts/create \
  -H "Content-Type: application/json" \
  -d '{"content":"PRIVATE post","privacy":"private","visible_to":[2]}'
```

Verify:
- Guest should **not** see private/almost_private:

```bash
curl -s "http://localhost:8080/api/posts?user_id=1" | python3 -m json.tool
```

- Bob should see all 3:

```bash
curl -s -b /tmp/cookiesB.txt "http://localhost:8080/api/posts?user_id=1" | python3 -m json.tool
```

---

## 4) Posts — image/GIF test (2-step)

### 4.1 Upload image as Alice

Replace `/ABSOLUTE/PATH/TO/file.gif` with a real file path.

```bash
curl -s -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/upload-image \
  -F "image=@/ABSOLUTE/PATH/TO/file.gif" | python3 -m json.tool
```

Copy the returned `filename`.

### 4.2 Create post with `image_filename`

Replace `<FILENAME>` with the value from the previous response.

```bash
curl -i -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/posts/create \
  -H "Content-Type: application/json" \
  -d '{"content":"POST with image","privacy":"public","image_filename":"<FILENAME>"}'
```

Verify feed includes `"image_url"`:

```bash
curl -s -b /tmp/cookiesA.txt "http://localhost:8080/api/posts?user_id=1" | python3 -m json.tool
```

---

## 5) Post comments test

### 5.1 List comments for a post (example: post_id=1)

```bash
curl -s -b /tmp/cookiesA.txt "http://localhost:8080/api/posts/comments?post_id=1" | python3 -m json.tool
```

### 5.2 Create comment (Alice)

```bash
curl -i -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/posts/comments \
  -H "Content-Type: application/json" \
  -d '{"post_id":1,"content":"Nice post!"}'
```

Verify:

```bash
curl -s -b /tmp/cookiesA.txt "http://localhost:8080/api/posts/comments?post_id=1" | python3 -m json.tool
```

---

## 6) Groups — invite flow

### 6.1 Alice creates a group

```bash
curl -s -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/groups \
  -H "Content-Type: application/json" \
  -d '{"group_name":"G1","description":"test group"}' | python3 -m json.tool
```

Expected:
- `{ "group_id": 1 }` on a fresh DB.

### 6.2 Alice invites Bob to group 1

```bash
curl -s -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/groups/invite \
  -H "Content-Type: application/json" \
  -d '{"group_id": 1, "user_id": 2}' | python3 -m json.tool
```

Expected:
- `{ "invitation_id": 1, "status": "pending" }` on a fresh DB.

### 6.3 Bob accepts invitation

```bash
curl -i -b /tmp/cookiesB.txt -X POST http://localhost:8080/api/groups/invitations/respond \
  -H "Content-Type: application/json" \
  -d '{"invitation_id": 1, "response":"accepted"}'
```

Verify membership-sensitive endpoint works:

```bash
curl -s -b /tmp/cookiesB.txt "http://localhost:8080/api/groups/posts?group_id=1" | python3 -m json.tool
```

---

## 7) Groups — group posts and comments

### 7.1 Bob creates a group post

```bash
curl -s -b /tmp/cookiesB.txt -X POST http://localhost:8080/api/groups/posts \
  -H "Content-Type: application/json" \
  -d '{"group_id": 1, "content":"hello group"}' | python3 -m json.tool
```

Expected:
- `{ "group_post_id": 1 }` on a fresh DB.

### 7.2 Bob comments on group post 1

```bash
curl -s -b /tmp/cookiesB.txt -X POST http://localhost:8080/api/groups/comments \
  -H "Content-Type: application/json" \
  -d '{"group_post_id": 1, "content":"first comment"}' | python3 -m json.tool
```

List comments:

```bash
curl -s -b /tmp/cookiesB.txt "http://localhost:8080/api/groups/comments?group_post_id=1" | python3 -m json.tool
```

---

## 8) Groups — events + responses

### 8.1 Alice creates an event

```bash
curl -s -b /tmp/cookiesA.txt -X POST http://localhost:8080/api/groups/events \
  -H "Content-Type: application/json" \
  -d '{"group_id": 1, "title":"Meet","description":"desc","event_datetime":"2026-01-22T20:30:00Z"}' | python3 -m json.tool
```

Expected:
- `{ "event_id": 1 }` on a fresh DB.

### 8.2 Bob responds

```bash
curl -i -b /tmp/cookiesB.txt -X POST http://localhost:8080/api/groups/events/respond \
  -H "Content-Type: application/json" \
  -d '{"event_id": 1, "response":"going"}'
```

---

## 9) Verify Notifications were inserted (sqlite)

```bash
sqlite3 /home/tarara-x/Desktop/history/social-network/backend/data/social_network.db \
"select id,type,user_id,related_user_id,related_group_id,related_event_id,message from Notifications order by id desc limit 50;"
```

Expected to see (among others):
- `group_invitation` (to Bob)
- `group_invitation_response` (to Alice)
- `group_event_created` (to Bob)
- `group_event_response` (to Alice)

---

## Troubleshooting

### `{"error":"invalid_json","message":"invalid JSON body"}`
- Your JSON is malformed. Common cause: using placeholders like `<GROUP_ID>` inside JSON.

### `{"error":"unauthorized",...}` / HTTP 401
- Missing/invalid cookie jar (`-b /tmp/cookiesA.txt`).

### HTTP 403
- You are logged in but not allowed (privacy/membership rules).

