import { defineStore } from "pinia";
import {
  getNotifications,
  markNotificationRead as apiMarkRead,
  markAllNotificationsRead as apiMarkAllRead
} from "../services/notificationsApi";

export const useNotificationsStore = defineStore("notifications", {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    isLoading: false,
    loaded: false
  }),
  actions: {
    async fetchNotifications(options = {}) {
      this.isLoading = true;
      try {
        const res = await getNotifications(options);
        this.notifications = res.notifications ?? [];
        this.unreadCount = res.unread_count ?? 0;
        this.loaded = true;
      } finally {
        this.isLoading = false;
      }
    },
    async markAsRead(notificationId) {
      await apiMarkRead(notificationId);
      const n = this.notifications.find((x) => x.id === notificationId);
      if (n) n.is_read = true;
      if (this.unreadCount > 0) this.unreadCount -= 1;
    },
    async markAllAsRead() {
      await apiMarkAllRead();
      this.notifications.forEach((n) => (n.is_read = true));
      this.unreadCount = 0;
    },
    /** Mark only follow_request notifications as read */
    async markFollowRequestsAsRead() {
      const unreadFollowRequestIds = this.notifications
        .filter((n) => n.type === "follow_request" && !n.is_read)
        .map((n) => n.id);
      if (unreadFollowRequestIds.length === 0) return;

      await Promise.allSettled(
        unreadFollowRequestIds.map((id) => apiMarkRead(id))
      );

      // Refresh counts so the nav badge becomes correct immediately.
      await this.fetchNotifications().catch(() => {});
    },
    /** Add a notification from real-time WebSocket (prepend, update unread count) */
    addNotification(notification) {
      const normalized = this._normalize(notification);
      const exists = this.notifications.some((n) => n.id === normalized.id);
      if (!exists) {
        this.notifications.unshift(normalized);
        if (!normalized.is_read) this.unreadCount += 1;
      }
    },
    _normalize(n) {
      const id = n.notification_id ?? n.id;
      return {
        id: id ?? `rt-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
        type: n.type ?? "notification",
        message: n.message ?? "",
        related_user_id: n.related_user_id ?? null,
        related_group_id: n.related_group_id ?? null,
        related_post_id: n.related_post_id ?? null,
        related_comment_id: n.related_comment_id ?? null,
        related_event_id: n.related_event_id ?? null,
        is_read: n.is_read ?? false,
        created_at: n.created_at ?? new Date().toISOString()
      };
    },
    reset() {
      this.notifications = [];
      this.unreadCount = 0;
      this.loaded = false;
    }
  }
});
