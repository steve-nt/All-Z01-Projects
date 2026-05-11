import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";

export const useProfileStore = defineStore("profile", {
  state: () => ({
    profile: null,
    loading: false,
    error: null,
    updatingPrivacy: false
  }),
  actions: {
    async fetchProfile(profileId) {
      this.loading = true;
      this.error = null;
      try {
        const data = await socialApi.getProfile(profileId);
        this.profile = data;
        return data;
      } catch (error) {
        this.profile = null;
        if (error?.status === 403) {
          this.error = { status: 403, message: "This profile is private" };
        } else {
          this.error = {
            status: error?.status || 500,
            message: error?.message || "Failed to load profile"
          };
        }
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async togglePrivacy(isPublic) {
      this.updatingPrivacy = true;
      this.error = null;
      try {
        const result = await socialApi.updatePrivacy(isPublic);
        if (this.profile) {
          this.profile = { ...this.profile, is_public: result.is_public };
        }
        return result;
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to update profile privacy"
        };
        throw error;
      } finally {
        this.updatingPrivacy = false;
      }
    },
    async updateProfile(payload) {
      this.error = null;
      try {
        await socialApi.updateProfile(payload);
        if (this.profile) {
          this.profile = { ...this.profile, ...payload };
        }
        return true;
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to update profile"
        };
        throw error;
      }
    },
    async uploadAvatar(file) {
      try {
        const result = await socialApi.uploadAvatar(file);
        if (this.profile && result?.avatarPath != null) {
          this.profile = { ...this.profile, avatar: result.avatarPath };
        }
        return result;
      } catch (error) {
        throw error;
      }
    },
    clearProfileState() {
      this.profile = null;
      this.loading = false;
      this.error = null;
      this.updatingPrivacy = false;
    }
  }
});
