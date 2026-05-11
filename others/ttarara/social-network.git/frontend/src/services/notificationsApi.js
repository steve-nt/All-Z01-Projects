import { apiRequest } from "./apiClient";

/**
 * GET /api/notifications
 * @param {{ unread_only?: boolean, limit?: number }} options
 * @returns {Promise<{ notifications: Array, unread_count: number, total: number }>}
 */
export function getNotifications(options = {}) {
  const params = new URLSearchParams();
  if (options.unread_only === true) params.set("unread_only", "true");
  if (options.limit != null) params.set("limit", String(options.limit));
  const query = params.toString();
  return apiRequest(query ? `/api/notifications?${query}` : "/api/notifications");
}

/**
 * POST /api/notifications/read - mark one as read
 * @param {number} notificationId
 */
export function markNotificationRead(notificationId) {
  return apiRequest("/api/notifications/read", {
    method: "POST",
    body: JSON.stringify({ notification_id: notificationId })
  });
}

/**
 * POST /api/notifications/read - mark all as read
 */
export function markAllNotificationsRead() {
  return apiRequest("/api/notifications/read", {
    method: "POST",
    body: JSON.stringify({ mark_all_read: true })
  });
}
