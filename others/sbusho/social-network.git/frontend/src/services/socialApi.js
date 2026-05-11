import { apiRequest } from "./apiClient";

function withUserId(path, userId) {
  const params = new URLSearchParams({ user_id: String(userId) });
  return `${path}?${params.toString()}`;
}

export const socialApi = {
  getAuthStatus() {
    return apiRequest("/api/auth/status");
  },
  getProfile(userId) {
    return apiRequest(withUserId("/api/profile", userId));
  },
  updatePrivacy(isPublic) {
    return apiRequest("/api/profile/privacy", {
      method: "POST",
      body: JSON.stringify({ is_public: Boolean(isPublic) })
    });
  },

  updateProfile(payload) {
    return apiRequest("/api/profile/update", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  uploadAvatar(file) {
    const formData = new FormData();
    formData.append("avatar", file);
    return apiRequest("/api/upload-avatar", {
      method: "POST",
      body: formData
    });
  },
  requestFollow(userId) {
    return apiRequest("/api/follow/request", {
      method: "POST",
      body: JSON.stringify({ user_id: userId })
    });
  },
  acceptFollowRequest(userId) {
    return apiRequest("/api/follow/accept", {
      method: "POST",
      body: JSON.stringify({ user_id: userId })
    });
  },
  declineFollowRequest(userId) {
    return apiRequest("/api/follow/decline", {
      method: "POST",
      body: JSON.stringify({ user_id: userId })
    });
  },
  unfollow(userId) {
    return apiRequest("/api/follow/unfollow", {
      method: "POST",
      body: JSON.stringify({ user_id: userId })
    });
  },
  getFollowers(userId) {
    return apiRequest(withUserId("/api/followers", userId));
  },
  getFollowing(userId) {
    return apiRequest(withUserId("/api/following", userId));
  },
  getPendingFollowRequests() {
    return apiRequest("/api/follow/requests");
  },
  searchUsers(q, limit = 20) {
    const params = new URLSearchParams();
    if (typeof q === "string") {
      params.set("q", q);
    }
    params.set("limit", String(limit));
    return apiRequest(`/api/users/search?${params.toString()}`);
  },
  getGroups() {
    return apiRequest("/api/groups");
  },
  createGroup(payload) {
    return apiRequest("/api/groups", {
      method: "POST",
      body: JSON.stringify({
        group_name: payload.group_name,
        description: payload.description
      })
    });
  },
  getGroupDetails(groupId) {
    const params = new URLSearchParams({ group_id: String(groupId) });
    return apiRequest(`/api/groups/view?${params.toString()}`);
  },

  getGroupInvitations() {
    return apiRequest("/api/groups/invitations");
  },

  inviteToGroup(groupId, userId) {
    return apiRequest("/api/groups/invite", {
      method: "POST",
      body: JSON.stringify({ group_id: Number(groupId), user_id: Number(userId) })
    });
  },

  respondToGroupInvitation(invitationId, response) {
    const r = String(response).toLowerCase();
    if (r !== "accepted" && r !== "declined") {
      return Promise.reject(new Error("response must be accepted or declined"));
    }
    return apiRequest("/api/groups/invitations/respond", {
      method: "POST",
      body: JSON.stringify({ invitation_id: Number(invitationId), response: r })
    });
  },

  getGroupJoinRequests(groupId) {
    const params = new URLSearchParams({ group_id: String(groupId) });
    return apiRequest(`/api/groups/join/requests?${params.toString()}`);
  },

  requestToJoinGroup(groupId) {
    return apiRequest("/api/groups/join/request", {
      method: "POST",
      body: JSON.stringify({ group_id: Number(groupId) })
    });
  },

  respondToJoinRequest(requestId, response) {
    const r = String(response).toLowerCase();
    if (r !== "accepted" && r !== "declined") {
      return Promise.reject(new Error("response must be accepted or declined"));
    }
    return apiRequest("/api/groups/join/respond", {
      method: "POST",
      body: JSON.stringify({ request_id: Number(requestId), response: r })
    });
  },

  getGroupPosts(groupId) {
    const params = new URLSearchParams({ group_id: String(groupId) });
    return apiRequest(`/api/groups/posts?${params.toString()}`);
  },

  createGroupPost(payload) {
    const body = {
      group_id: Number(payload.group_id),
      content: (payload.content != null ? payload.content : "").toString().trim()
    };
    if (payload.image_url != null && payload.image_url !== "") {
      body.image_url = payload.image_url;
    }
    return apiRequest("/api/groups/posts", {
      method: "POST",
      body: JSON.stringify(body)
    });
  },

  getGroupComments(groupPostId) {
    const params = new URLSearchParams({ group_post_id: String(groupPostId) });
    return apiRequest(`/api/groups/comments?${params.toString()}`);
  },

  createGroupComment(payload) {
    const body = {
      group_post_id: Number(payload.group_post_id),
      content: payload.content.trim()
    };
    if (payload.image_url != null && payload.image_url !== "") {
      body.image_url = payload.image_url;
    }
    return apiRequest("/api/groups/comments", {
      method: "POST",
      body: JSON.stringify(body)
    });
  },

  getGroupEvents(groupId) {
    const params = new URLSearchParams({ group_id: String(groupId) });
    return apiRequest(`/api/groups/events?${params.toString()}`);
  },

  createGroupEvent(payload) {
    return apiRequest("/api/groups/events", {
      method: "POST",
      body: JSON.stringify({
        group_id: Number(payload.group_id),
        title: payload.title.trim(),
        description: payload.description.trim(),
        event_datetime: payload.event_datetime
      })
    });
  },

  respondToGroupEvent(eventId, response) {
    const r = String(response).toLowerCase().replace(/\s+/g, " ");
    if (r !== "going" && r !== "not going") {
      return Promise.reject(new Error("response must be going or not going"));
    }
    return apiRequest("/api/groups/events/respond", {
      method: "POST",
      body: JSON.stringify({ event_id: Number(eventId), response: r })
    });
  },

  // Posts (feed)
  getPosts(options = {}) {
    const params = new URLSearchParams();
    if (options.user_id != null) {
      params.set("user_id", String(options.user_id));
    }
    if (options.limit != null) {
      params.set("limit", String(options.limit));
    }
    if (options.offset != null) {
      params.set("offset", String(options.offset));
    }
    const query = params.toString();
    return apiRequest(query ? `/api/posts?${query}` : "/api/posts");
  },

  createPost(payload) {
    const body = {
      content: payload.content,
      privacy: payload.privacy
    };
    if (payload.privacy === "private" && Array.isArray(payload.visible_to)) {
      body.visible_to = payload.visible_to;
    }
    if (payload.image_filename != null && payload.image_filename !== "") {
      body.image_filename = payload.image_filename;
    }
    return apiRequest("/api/posts/create", {
      method: "POST",
      body: JSON.stringify(body)
    });
  },

  uploadPostImage(file) {
    const formData = new FormData();
    formData.append("image", file);
    return apiRequest("/api/upload-image", {
      method: "POST",
      body: formData
    });
  },

  // Post comments
  getPostComments(postId) {
    const params = new URLSearchParams({ post_id: String(postId) });
    return apiRequest(`/api/posts/comments?${params.toString()}`);
  },

  createPostComment(payload) {
    const body = {
      post_id: payload.post_id,
      content: payload.content.trim()
    };
    if (payload.image_url != null && payload.image_url !== "") {
      body.image_url = payload.image_url;
    }
    return apiRequest("/api/posts/comments", {
      method: "POST",
      body: JSON.stringify(body)
    });
  },

  // Private & group messages
  getConversations() {
    return apiRequest("/api/messages/conversations");
  },
  getMessageContacts() {
    return apiRequest("/api/messages/contacts");
  },
  getMessages(userId, limit = 50) {
    const params = new URLSearchParams({ user_id: String(userId) });
    if (limit != null) params.set("limit", String(limit));
    return apiRequest(`/api/messages?${params.toString()}`);
  },
  sendMessage(recipientId, content) {
    return apiRequest("/api/messages/send", {
      method: "POST",
      body: JSON.stringify({ recipient_id: Number(recipientId), content: String(content).trim() })
    });
  },
  markMessagesRead(payload) {
    return apiRequest("/api/messages/read", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },
  getGroupMessages(groupId, limit = 50) {
    const params = new URLSearchParams({ group_id: String(groupId) });
    if (limit != null) params.set("limit", String(limit));
    return apiRequest(`/api/messages/group?${params.toString()}`);
  },
  sendGroupMessage(groupId, content) {
    return apiRequest("/api/messages/group/send", {
      method: "POST",
      body: JSON.stringify({ group_id: Number(groupId), content: String(content).trim() })
    });
  }
};
