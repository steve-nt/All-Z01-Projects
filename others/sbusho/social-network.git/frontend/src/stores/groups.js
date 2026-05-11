import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";

const GROUP_FORBIDDEN_MESSAGE =
  "Members only. Join or request to join this group to see content.";

function messageForGroupError(error, fallback) {
  return error?.status === 403 ? GROUP_FORBIDDEN_MESSAGE : (error?.message || fallback);
}

export const useGroupsStore = defineStore("groups", {
  state: () => ({
    groups: [],
    selectedGroup: null,
    groupInvitations: [],
    joinRequests: [],
    groupPosts: [],
    groupEvents: [],
    commentsByGroupPostId: {},
    loadingCommentsByGroupPostId: {},
    creatingCommentPostId: null,
    loadingGroupEvents: false,
    creatingGroupEvent: false,
    respondingEventId: null,
    loading: false,
    loadingInvitations: false,
    loadingJoinRequests: false,
    loadingGroupPosts: false,
    creatingGroupPost: false,
    error: null,
    invitingUserId: null,
    requestingJoinGroupId: null,
    respondingJoinRequestId: null
  }),
  actions: {
    clearError() {
      this.error = null;
    },
    async fetchGroupPosts(groupId) {
      const id = Number(groupId);
      if (!id) {
        this.groupPosts = [];
        return [];
      }
      this.loadingGroupPosts = true;
      try {
        const data = await socialApi.getGroupPosts(id);
        this.groupPosts = data?.posts || [];
        return this.groupPosts;
      } catch (error) {
        this.groupPosts = [];
        throw error;
      } finally {
        this.loadingGroupPosts = false;
      }
    },
    async createGroupPost(payload) {
      this.creatingGroupPost = true;
      this.clearError();
      try {
        await socialApi.createGroupPost(payload);
        await this.fetchGroupPosts(payload.group_id);
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: messageForGroupError(error, "Failed to create group post")
        };
        throw error;
      } finally {
        this.creatingGroupPost = false;
      }
    },
    async fetchGroupComments(groupPostId) {
      const id = Number(groupPostId);
      if (!id) return [];
      this.loadingCommentsByGroupPostId = { ...this.loadingCommentsByGroupPostId, [id]: true };
      try {
        const data = await socialApi.getGroupComments(id);
        const comments = data?.comments || [];
        this.commentsByGroupPostId = { ...this.commentsByGroupPostId, [id]: comments };
        return comments;
      } catch {
        this.commentsByGroupPostId = { ...this.commentsByGroupPostId, [id]: [] };
        return [];
      } finally {
        this.loadingCommentsByGroupPostId = { ...this.loadingCommentsByGroupPostId, [id]: false };
      }
    },
    async createGroupComment(payload) {
      this.creatingCommentPostId = payload.group_post_id;
      try {
        await socialApi.createGroupComment(payload);
        await this.fetchGroupComments(payload.group_post_id);
      } catch (error) {
        throw error;
      } finally {
        this.creatingCommentPostId = null;
      }
    },
    async fetchGroupEvents(groupId) {
      const id = Number(groupId);
      if (!id) {
        this.groupEvents = [];
        return [];
      }
      this.loadingGroupEvents = true;
      try {
        const data = await socialApi.getGroupEvents(id);
        this.groupEvents = data?.events || [];
        return this.groupEvents;
      } catch (error) {
        this.groupEvents = [];
        throw error;
      } finally {
        this.loadingGroupEvents = false;
      }
    },
    async createGroupEvent(payload) {
      this.creatingGroupEvent = true;
      this.clearError();
      try {
        await socialApi.createGroupEvent(payload);
        await this.fetchGroupEvents(payload.group_id);
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: messageForGroupError(error, "Failed to create event")
        };
        throw error;
      } finally {
        this.creatingGroupEvent = false;
      }
    },
    async respondToEvent(eventId, response, groupId) {
      this.respondingEventId = eventId;
      try {
        await socialApi.respondToGroupEvent(eventId, response);
        const gid = groupId ?? this.groupEvents.find((e) => e.event_id === eventId)?.group_id;
        if (gid) await this.fetchGroupEvents(gid);
      } catch (error) {
        throw error;
      } finally {
        this.respondingEventId = null;
      }
    },
    async fetchJoinRequests(groupId) {
      const id = Number(groupId);
      if (!id) {
        this.joinRequests = [];
        return [];
      }
      this.loadingJoinRequests = true;
      try {
        const data = await socialApi.getGroupJoinRequests(id);
        this.joinRequests = data?.join_requests || [];
        return this.joinRequests;
      } catch (error) {
        this.joinRequests = [];
        throw error;
      } finally {
        this.loadingJoinRequests = false;
      }
    },
    async requestToJoin(groupId) {
      this.requestingJoinGroupId = groupId;
      try {
        await socialApi.requestToJoinGroup(groupId);
      } finally {
        this.requestingJoinGroupId = null;
      }
    },
    async respondToJoinRequest(requestId, response) {
      this.respondingJoinRequestId = requestId;
      try {
        await socialApi.respondToJoinRequest(requestId, response);
        this.joinRequests = this.joinRequests.filter(
          (r) => r.request_id !== requestId
        );
      } finally {
        this.respondingJoinRequestId = null;
      }
    },
    async fetchGroupInvitations() {
      this.loadingInvitations = true;
      try {
        const data = await socialApi.getGroupInvitations();
        this.groupInvitations = data?.invitations || [];
        return this.groupInvitations;
      } catch (error) {
        this.groupInvitations = [];
        throw error;
      } finally {
        this.loadingInvitations = false;
      }
    },
    async inviteUser(groupId, userId) {
      this.invitingUserId = userId;
      try {
        await socialApi.inviteToGroup(groupId, userId);
      } finally {
        this.invitingUserId = null;
      }
    },
    async respondToInvitation(invitationId, response) {
      await socialApi.respondToGroupInvitation(invitationId, response);
      this.groupInvitations = this.groupInvitations.filter(
        (inv) => inv.invitation_id !== invitationId
      );
    },
    async fetchGroups() {
      this.loading = true;
      this.clearError();
      try {
        const data = await socialApi.getGroups();
        this.groups = data?.groups || [];
        return this.groups;
      } catch (error) {
        this.groups = [];
        this.error = {
          status: error?.status || 500,
          message: messageForGroupError(error, "Failed to load groups")
        };
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async fetchGroupDetails(groupId) {
      this.loading = true;
      this.clearError();
      try {
        const data = await socialApi.getGroupDetails(groupId);
        this.selectedGroup = data || null;
        return this.selectedGroup;
      } catch (error) {
        this.selectedGroup = null;
        this.error = {
          status: error?.status || 500,
          message: messageForGroupError(error, "Failed to load group details")
        };
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async createGroup(payload) {
      this.loading = true;
      this.clearError();
      try {
        const created = await socialApi.createGroup(payload);
        await this.fetchGroups();
        return created;
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: messageForGroupError(error, "Failed to create group")
        };
        throw error;
      } finally {
        this.loading = false;
      }
    }
  }
});
