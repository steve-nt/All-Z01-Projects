Files Changed/Added

profileHandler.go (added): introduces profile view/privacy logic plus follow request, followers/following, and request list endpoints.

routes.go: registers new Part 2 API routes for profile and follow flows.

registerHandler.go: captures is_public/is_active, stores profile fields, and writes a welcome notification on signup.

utils.go: FileServiceWithAuth now injects UserID and Nickname from session; session helpers now rely on shared sqlite.GetDB() without closing per call.

middleware.go: updates imports and simplifies HandlerFunc alias (used for auth enforcement wrappers).

loginHandler.go, authHandler.go, avatar_image.go: update imports to backend/utils package path; no new Part 2 behavior in diff besides import alignment.

login.html, register.html (added): placeholder templates.

testing_part2.md (added): test plan and curl recipes for Part 2 endpoints.

README.md (added): high-level auth/profile notes and Part 2 API description.

.gitignore: ignores *.txt (for cookie jars in testing).

server.go: placeholder to keep package compiling.

New Routes and Endpoints

Registered in routes.go:

GET /api/profile?user_id=<int> → profile.ProfileViewHandler (no auth required; privacy enforced) routes.go (line 31), profileHandler.go (line 15)
POST /api/profile/privacy → profile.ProfilePrivacyHandler (auth required) routes.go (line 33), profileHandler.go (line 124)
POST /api/follow/request → profile.FollowRequestHandler (auth required) routes.go (line 34), profileHandler.go (line 177)
POST /api/follow/accept → profile.FollowAcceptHandler (auth required) routes.go (line 35), profileHandler.go (line 278)
POST /api/follow/decline → profile.FollowDeclineHandler (auth required) routes.go (line 36), profileHandler.go (line 347)
GET /api/followers?user_id=<int> → profile.FollowersHandler (privacy enforced) routes.go (line 37), profileHandler.go (line 416)
GET /api/following?user_id=<int> → profile.FollowingHandler (privacy enforced) routes.go (line 38), profileHandler.go (line 504)
POST /api/follow/unfollow → profile.UnfollowHandler (auth required) routes.go (line 39), profileHandler.go (line 592)
GET /api/follow/requests → profile.FollowRequestsHandler (auth required) routes.go (line 40), profileHandler.go (line 642)


Profile System (Public/Private)

ProfileViewHandler fetches user profile fields from Users and enforces privacy using is_public and the viewer’s follow status profileHandler.go (line 15).
Full profile JSON is returned if is_public is true or viewer is the owner (viewerID == userID) profileHandler.go (line 63).
For private profiles, unauthenticated viewers get 403 with {error:"forbidden"} profileHandler.go (line 81).
Authenticated viewers must be an accepted follower (Followers with status='accepted') to see full fields; otherwise 403 profileHandler.go (line 90).
Returned fields on full view: user_id, email, first_name, last_name, nickname, avatar, about_me, is_public profileHandler.go (line 65).


Follow Request Flow (Pending/Accepted/Declined)

FollowRequestHandler creates follower records; if target is public, it auto-accepts (status='accepted'), otherwise status='pending' profileHandler.go (line 244).
It blocks self-follow and duplicate relationships with 400 profileHandler.go (line 201), profileHandler.go (line 224).
FollowAcceptHandler transitions a pending request to accepted and returns previous_status profileHandler.go (line 277).
FollowDeclineHandler transitions pending to declined (record retained for history) profileHandler.go (line 397).
UnfollowHandler deletes the Followers row for the authenticated follower/target pair profileHandler.go (line 592).
Followers/Following Lists

FollowersHandler lists accepted followers (Followers.status='accepted') with user_id, nickname, avatar profileHandler.go (line 469).
FollowingHandler lists accepted following targets similarly profileHandler.go (line 557).
Both handlers enforce privacy: if the target user is private and the viewer is neither the owner nor an accepted follower, they return 403 profileHandler.go (line 451), profileHandler.go (line 539).
Privacy Toggle

ProfilePrivacyHandler updates Users.is_public for the authenticated user via JSON { "is_public": true|false } profileHandler.go (line 124).
Validation requires POST and a non-nil is_public field; responds with {user_id, is_public} on success profileHandler.go (line 141).
Follow Request Notifications

When a follow request is pending, FollowRequestHandler inserts a Notifications record with type follow_request, related_user_id set to the requester, and message "New follow request" profileHandler.go (line 258).
Separate: RegisterHandler writes a welcome notification on signup registerHandler.go (line 173).
Auth Enforcement in Part 2

middleware.RequireAuthJSON is applied to /api/profile/privacy, /api/follow/request, /api/follow/accept, /api/follow/decline, /api/follow/unfollow, /api/follow/requests in routes.go (lines 33-40).
GET /api/profile, /api/followers, /api/following are public but enforce privacy using utils.CheckAuth and Followers status inside the handlers profileHandler.go (line 33), profileHandler.go (line 433), profileHandler.go (line 521).
Session identity continues to be derived from the session cookie via utils.CheckAuth and utils.GetUserIDFromSession utils.go (line 215).

