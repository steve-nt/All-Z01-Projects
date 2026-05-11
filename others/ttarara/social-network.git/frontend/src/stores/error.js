import { defineStore } from "pinia";

// Part 1: centralized error banner state
export const useErrorStore = defineStore("error", {
  state: () => ({
    message: "",
    type: ""
  }),
  actions: {
    // Part 1: set a user-facing error message
    setError(message, type = "error") {
      this.message = message;
      this.type = type;
    },
    // Part 1: clear error banner
    clear() {
      this.message = "";
      this.type = "";
    }
  }
});
