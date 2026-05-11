package messages

import (
	"net/http"
	"social-network/backend/middleware"
)

// SetupMessageRoutes registers private messaging and group chat API endpoints
func SetupMessageRoutes() {
	// Private Messages
	// POST /api/messages/send - Send a private message
	http.HandleFunc("/api/messages/send", middleware.WrapHandler(middleware.RequireAuthJSON(SendMessageHandler)))

	// GET /api/messages?user_id=X - Get conversation history with a user
	http.HandleFunc("/api/messages", middleware.WrapHandler(middleware.RequireAuthJSON(GetMessagesHandler)))

	// GET /api/messages/conversations - Get list of conversations
	http.HandleFunc("/api/messages/conversations", middleware.WrapHandler(middleware.RequireAuthJSON(GetConversationsHandler)))

	// GET /api/messages/contacts - Get users you can message (followers/following/public)
	http.HandleFunc("/api/messages/contacts", middleware.WrapHandler(middleware.RequireAuthJSON(GetContactsHandler)))

	// POST /api/messages/read - Mark messages as read
	http.HandleFunc("/api/messages/read", middleware.WrapHandler(middleware.RequireAuthJSON(MarkReadHandler)))

	// Group Chat
	// POST /api/messages/group/send - Send a message to group chat
	http.HandleFunc("/api/messages/group/send", middleware.WrapHandler(middleware.RequireAuthJSON(SendGroupMessageHandler)))

	// GET /api/messages/group?group_id=X - Get group chat history
	http.HandleFunc("/api/messages/group", middleware.WrapHandler(middleware.RequireAuthJSON(GetGroupMessagesHandler)))
}

