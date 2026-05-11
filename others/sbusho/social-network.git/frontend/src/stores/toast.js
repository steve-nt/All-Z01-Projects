import { defineStore } from "pinia";

let nextId = 0;
const DEFAULT_DURATION_MS = 5000;

export const useToastStore = defineStore("toast", {
  state: () => ({
    toasts: []
  }),
  actions: {
    add(message, type = "info", durationMs = DEFAULT_DURATION_MS) {
      const id = ++nextId;
      this.toasts.push({ id, message, type, durationMs });
      return id;
    },
    remove(id) {
      this.toasts = this.toasts.filter((t) => t.id !== id);
    },
    /** Show a general notification toast (follow request, group invite, etc.) */
    notification(message) {
      return this.add(message, "notification", DEFAULT_DURATION_MS);
    },
    /** Show a private message toast (distinct from other notifications) */
    privateMessage(message) {
      return this.add(message, "message", DEFAULT_DURATION_MS);
    }
  }
});
