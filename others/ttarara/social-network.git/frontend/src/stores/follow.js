import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";

export const useFollowStore = defineStore("follow", {
  state: () => ({
    followers: [],
    following: [],
    pendingRequests: [],
    loading: false,
    error: null,
    pendingActions: {}
  }),
  getters: {
    isActionPending: (state) => (key) => Boolean(state.pendingActions[key])
  },
  actions: {
    setActionPending(key, value) {
      this.pendingActions = { ...this.pendingActions, [key]: value };
    },
    clearError() {
      this.error = null;
    },
    async requestFollow(targetUserId) {
      const parsedTargetUserId = Number(targetUserId);
      if (!Number.isInteger(parsedTargetUserId) || parsedTargetUserId <= 0) {
        this.error = {
          status: 400,
          message: "Invalid user id for follow request"
        };
        throw new Error("Invalid user id for follow request");
      }

      const actionKey = `request:${parsedTargetUserId}`;
      this.setActionPending(actionKey, true);
      this.clearError();
      try {
        return await socialApi.requestFollow(parsedTargetUserId);
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to send follow request"
        };
        throw error;
      } finally {
        this.setActionPending(actionKey, false);
      }
    },
    async acceptRequest(userId) {
      const actionKey = `accept:${userId}`;
      this.setActionPending(actionKey, true);
      this.clearError();
      try {
        const response = await socialApi.acceptFollowRequest(userId);
        this.pendingRequests = this.pendingRequests.filter(
          (request) => request.user_id !== userId
        );
        return response;
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to accept follow request"
        };
        throw error;
      } finally {
        this.setActionPending(actionKey, false);
      }
    },
    async declineRequest(userId) {
      const actionKey = `decline:${userId}`;
      this.setActionPending(actionKey, true);
      this.clearError();
      try {
        const response = await socialApi.declineFollowRequest(userId);
        this.pendingRequests = this.pendingRequests.filter(
          (request) => request.user_id !== userId
        );
        return response;
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to decline follow request"
        };
        throw error;
      } finally {
        this.setActionPending(actionKey, false);
      }
    },
    async unfollow(targetUserId) {
      const actionKey = `unfollow:${targetUserId}`;
      this.setActionPending(actionKey, true);
      this.clearError();
      try {
        return await socialApi.unfollow(targetUserId);
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to unfollow user"
        };
        throw error;
      } finally {
        this.setActionPending(actionKey, false);
      }
    },
    async fetchFollowers(userId) {
      this.loading = true;
      this.clearError();
      try {
        const data = await socialApi.getFollowers(userId);
        this.followers = data?.followers || [];
        return this.followers;
      } catch (error) {
        this.followers = [];
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to load followers"
        };
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async fetchFollowing(userId) {
      this.loading = true;
      this.clearError();
      try {
        const data = await socialApi.getFollowing(userId);
        this.following = data?.following || [];
        return this.following;
      } catch (error) {
        this.following = [];
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to load following"
        };
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async fetchPendingRequests() {
      this.loading = true;
      this.clearError();
      try {
        const data = await socialApi.getPendingFollowRequests();
        this.pendingRequests = data?.requests || [];
        return this.pendingRequests;
      } catch (error) {
        this.pendingRequests = [];
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to load pending requests"
        };
        throw error;
      } finally {
        this.loading = false;
      }
    }
  }
});
