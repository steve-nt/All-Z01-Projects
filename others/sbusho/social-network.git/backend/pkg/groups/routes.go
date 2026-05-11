package groups

import (
	"net/http"
	"social-network/backend/middleware"
)

// SetupGroupRoutes registers Groups, Group Posts, and Group Events API endpoints.
// Call this from main.go.
func SetupGroupRoutes() {
	// Groups
	http.HandleFunc("/api/groups", middleware.WrapHandler(GroupsHandler)) // GET browse, POST create (auth required for create)
	http.HandleFunc("/api/groups/view", middleware.WrapHandler(GroupViewHandler))

	// Invitations & join requests (auth)
	http.HandleFunc("/api/groups/invite", middleware.WrapHandler(middleware.RequireAuthJSON(InviteToGroupHandler)))
	http.HandleFunc("/api/groups/invitations/respond", middleware.WrapHandler(middleware.RequireAuthJSON(RespondGroupInviteHandler)))

	http.HandleFunc("/api/groups/join/request", middleware.WrapHandler(middleware.RequireAuthJSON(RequestJoinGroupHandler)))
	http.HandleFunc("/api/groups/join/requests", middleware.WrapHandler(middleware.RequireAuthJSON(ListJoinRequestsHandler)))
	http.HandleFunc("/api/groups/join/respond", middleware.WrapHandler(middleware.RequireAuthJSON(RespondJoinRequestHandler)))

	// Group posts & comments (auth)
	http.HandleFunc("/api/groups/posts", middleware.WrapHandler(middleware.RequireAuthJSON(GroupPostsHandler)))     // GET(list) + POST(create)
	http.HandleFunc("/api/groups/comments", middleware.WrapHandler(middleware.RequireAuthJSON(GroupCommentsHandler))) // GET(list) + POST(create)

	// Group events (auth)
	http.HandleFunc("/api/groups/events", middleware.WrapHandler(middleware.RequireAuthJSON(GroupEventsHandler))) // GET(list) + POST(create)
	http.HandleFunc("/api/groups/events/respond", middleware.WrapHandler(middleware.RequireAuthJSON(RespondEventHandler)))
}

