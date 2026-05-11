import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    currentUser: null,
    isChecking: false,
    sessionChecked: false
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.currentUser?.userID),
    loggedIn() {
      return this.isLoggedIn;
    },
    userId() {
      return this.currentUser?.userID ?? null;
    },
    nickname() {
      return this.currentUser?.nickname ?? null;
    }
  },
  actions: {
    setSession(session) {
      if (session?.loggedIn && session?.userID) {
        this.currentUser = {
          userID: session.userID,
          nickname: session.nickname || null
        };
      } else {
        this.currentUser = null;
      }
      this.sessionChecked = true;
    },
    clearSession() {
      this.currentUser = null;
      this.sessionChecked = true;
    },
    async loadAuthStatus() {
      if (this.isChecking) {
        return;
      }
      this.isChecking = true;
      try {
        const data = await socialApi.getAuthStatus();
        this.setSession(data || {});
      } catch (error) {
        this.clearSession();
      } finally {
        this.isChecking = false;
      }
    },
    async checkSession() {
      await this.loadAuthStatus();
    },
    async logout() {
      try {
        await fetch("/logout", {
          method: "POST",
          credentials: "include",
          headers: { Accept: "application/json" }
        });
      } finally {
        this.clearSession();
      }
    }
  }
});
