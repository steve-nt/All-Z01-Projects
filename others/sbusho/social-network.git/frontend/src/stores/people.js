import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";

export const usePeopleStore = defineStore("people", {
  state: () => ({
    query: "",
    users: [],
    loading: false,
    error: null
  }),
  actions: {
    async search(query, limit = 20) {
      this.query = query;
      this.loading = true;
      this.error = null;
      try {
        const data = await socialApi.searchUsers(query, limit);
        this.users = data?.users || [];
        return this.users;
      } catch (error) {
        this.users = [];
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to search users"
        };
        throw error;
      } finally {
        this.loading = false;
      }
    }
  }
});
