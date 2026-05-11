## Georgia Part 2 Summary

This summary is based on the repository state under `backend/` and `docs/`.

### Implemented (routes and behavior)
- `GET /api/profile?user_id=<int>`: Returns public profile data; enforces privacy for private profiles via accepted follower status.
- `POST /api/profile/privacy`: Toggles `Users.is_public` for the authenticated user.
- `POST /api/follow/request`: Creates a follow record; auto-accepts when target is public, otherwise marks as `pending` and creates a follow request notification.
- `POST /api/follow/accept`: Accepts a pending follow request.
- `POST /api/follow/decline`: Declines a pending follow request (keeps record with `declined`).
- `POST /api/follow/unfollow`: Deletes the follow relationship.
- `GET /api/followers?user_id=<int>`: Lists accepted followers with privacy enforcement.
- `GET /api/following?user_id=<int>`: Lists accepted following with privacy enforcement.
- `GET /api/follow/requests`: Lists pending follow requests for the authenticated user.

### Profile privacy (public/private)
- Private profiles block non-followers: viewing, followers list, and following list require owner or accepted follower access.

### Follow request notifications
- On pending follow request, a notification is inserted with:
  - Columns: `user_id`, `type`, `related_user_id`, `message`
  - `type` is `follow_request`, `message` is `"New follow request"`.

### Welcome notification validation
- `registerHandler` inserts into `Notifications` with `(user_id, type, message, created_at)`.
- Schema includes `user_id`, `type`, `message`, `created_at`.
- Result: welcome notification is supported and not broken.

### User activity endpoints
- Not implemented in Part 2.
- Suggested future placeholders (no code added):
  - `GET /api/activity`
  - `GET /api/activity?user_id=<int>`
  - `POST /api/activity/read`

### Manual testing (curl)
```bash
# Register + login user A (cookie jar a.txt)
curl -i -c a.txt -X POST http://localhost:8080/register \
  -F "email=alice@example.com" -F "password=P@ssword1!" -F "confirmPassword=P@ssword1!" \
  -F "first_name=Alice" -F "last_name=Anderson" -F "date_of_birth=1990-01-01" \
  -F "nickname=alice" -F "about_me=I am Alice" -F "is_public=true" \
  --next -i -c a.txt -b a.txt -X POST http://localhost:8080/login \
  -d "email=alice@example.com" -d "password=P@ssword1!"

# Register + login user B (cookie jar b.txt)
curl -i -c b.txt -X POST http://localhost:8080/register \
  -F "email=bob@example.com" -F "password=P@ssword1!" -F "confirmPassword=P@ssword1!" \
  -F "first_name=Bob" -F "last_name=Brown" -F "date_of_birth=1991-02-02" \
  -F "nickname=bob" -F "about_me=I am Bob" -F "is_public=true" \
  --next -i -c b.txt -b b.txt -X POST http://localhost:8080/login \
  -d "email=bob@example.com" -d "password=P@ssword1!"

# Toggle privacy (make Bob private)
curl -i -b b.txt -X POST http://localhost:8080/api/profile/privacy \
  -H "Content-Type: application/json" -d '{"is_public":false}'

# Follow request (Alice -> Bob)
curl -i -b a.txt -X POST http://localhost:8080/api/follow/request \
  -H "Content-Type: application/json" -d '{"user_id":2}'

# List follow requests (Bob)
curl -i -b b.txt http://localhost:8080/api/follow/requests

# Accept follow request (Bob accepts Alice)
curl -i -b b.txt -X POST http://localhost:8080/api/follow/accept \
  -H "Content-Type: application/json" -d '{"user_id":1}'

# Followers and following lists
curl -i -b b.txt "http://localhost:8080/api/followers?user_id=2" \
  "http://localhost:8080/api/following?user_id=1"

# Unfollow (Alice -> Bob)
curl -i -b a.txt -X POST http://localhost:8080/api/follow/unfollow \
  -H "Content-Type: application/json" -d '{"user_id":2}'
```




NEED TO FIX:
Below is a **copy-paste curl cookbook** for **every endpoint** related to **register, login, profiles, privacy, follow system**.
Assume server is running on `http://localhost:8080`.

I’ll use:

* `a.txt` → cookies for user A
* `b.txt` → cookies for user B
  (Change emails if needed.)

---

# AUTH / ACCOUNT

## Register

```bash
curl -i -c a.txt -X POST http://localhost:8080/register \
  -F "email=a@example.com" \
  -F "password=P@ssword1!" \
  -F "confirmPassword=P@ssword1!" \
  -F "first_name=A" \
  -F "last_name=A" \
  -F "date_of_birth=1990-01-01" \
  -F "nickname=userA" \
  -F "about_me=hi" \
  -F "is_public=false"
```

✅ Expect: `303 See Other`

---

## Login

```bash
curl -i -c a.txt -X POST http://localhost:8080/login \
  -d "email=a@example.com" \
  -d "password=P@ssword1!"
```

✅ Expect: `303` + `Set-Cookie: session=...`

---

## Auth status

```bash
curl -i -b a.txt http://localhost:8080/api/auth/status
```

✅ Example:

```json
{"loggedIn":true,"userID":2,"nickname":"userA"}
```

---

# PROFILES (Georgia)

## View profile (public or allowed private)

```bash
curl -i http://localhost:8080/api/profile?user_id=2
```

or (authenticated):

```bash
curl -i -b a.txt http://localhost:8080/api/profile?user_id=2
```

---

## Toggle privacy

```bash
curl -i -b a.txt -X POST http://localhost:8080/api/profile/privacy \
  -H "Content-Type: application/json" \
  -d '{"is_public":true}'
```

or

```bash
-d '{"is_public":false}'
```

---

# FOLLOW SYSTEM (Georgia)

## Send follow request

```bash
curl -i -b b.txt -X POST http://localhost:8080/api/follow/request \
  -H "Content-Type: application/json" \
  -d '{"user_id":2}'
```

* Public target → auto-accepted
* Private target → pending

---

## View pending follow requests (target user)

```bash
curl -i -b a.txt http://localhost:8080/api/follow/requests
```

Example:

```json
{
  "user_id": 2,
  "requests": [
    {
      "user_id": 1,
      "nickname": "userB",
      "status": "pending"
    }
  ]
}
```

---

## Accept follow request

```bash
curl -i -b a.txt -X POST http://localhost:8080/api/follow/accept \
  -H "Content-Type: application/json" \
  -d '{"user_id":1}'
```

---

## Decline follow request

```bash
curl -i -b a.txt -X POST http://localhost:8080/api/follow/decline \
  -H "Content-Type: application/json" \
  -d '{"user_id":1}'
```

---

## Unfollow

```bash
curl -i -b b.txt -X POST http://localhost:8080/api/follow/unfollow \
  -H "Content-Type: application/json" \
  -d '{"user_id":2}'
```

---

# LISTS (Georgia)

## Followers list

```bash
curl -i http://localhost:8080/api/followers?user_id=2
```

or (private profile, authenticated):

```bash
curl -i -b a.txt http://localhost:8080/api/followers?user_id=2
```

---

## Following list

```bash
curl -i http://localhost:8080/api/following?user_id=1
```

---

# PRIVACY ENFORCEMENT CHECKS

## Stranger blocked from private profile

```bash
curl -i http://localhost:8080/api/profile?user_id=2
```

✅ Expect: `403 Forbidden`

---

## Accepted follower can view private profile

```bash
curl -i -b b.txt http://localhost:8080/api/profile?user_id=2
```

✅ Expect: `200 OK`

---

# FINAL VERDICT ✅

If **all commands above work**, then:

✔ Profiles
✔ Privacy toggles
✔ Follow requests
✔ Accept / decline
✔ Followers / following
✔ Auth-based access control

👉 **Georgia Part 2 is DONE.**

If you want next:

* a **README testing section**
* or a **“Georgia Part 2 verification checklist”** for audit
  just say so.
