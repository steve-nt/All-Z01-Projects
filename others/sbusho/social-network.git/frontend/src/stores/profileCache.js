import { defineStore } from "pinia";

// Part 1: simple in-memory cache for profiles
export const useProfileCacheStore = defineStore("profileCache", {
  state: () => ({
    profilesById: {}
  }),
  actions: {
    // Part 1: read cached profile by user id
    getProfile(userId) {
      return this.profilesById[userId] || null;
    },
    // Part 1: set profile cache
    setProfile(userId, profile) {
      this.profilesById[userId] = profile;
    },
    // Part 1: clear single cached profile
    clearProfile(userId) {
      delete this.profilesById[userId];
    },
    // Part 1: clear all cached profiles
    clearAll() {
      this.profilesById = {};
    }
  }
});
