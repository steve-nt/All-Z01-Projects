<template>
  <div class="page-card">
    <div class="notifications-header">
      <h1>Notifications</h1>
      <button
        v-if="notificationsStore.unreadCount > 0"
        type="button"
        class="button secondary"
        :disabled="notificationsStore.isLoading"
        @click="markAllRead"
      >
        Mark all as read
      </button>
    </div>

    <p v-if="notificationsStore.isLoading && !notificationsStore.loaded" class="muted">
      Loading notifications...
    </p>
    <p v-else-if="notificationsStore.notifications.length === 0" class="muted">
      No notifications yet.
    </p>
    <ul v-else class="notifications-list">
      <li
        v-for="n in notificationsStore.notifications"
        :key="n.id"
        class="notification-item"
        :class="{
          'notification-item--unread': !n.is_read,
          'notification-item--message': n.type === 'message' || n.type === 'group_message'
        }"
      >
        <div class="notification-body">
          <span v-if="n.type" class="notification-type-label">{{ typeLabel(n.type) }}</span>
          <span class="notification-message">{{ notificationMessageDisplay(n.message) }}</span>
          <span class="notification-meta">
            {{ formatTime(n.created_at) }}
          </span>
          <div class="notification-actions">
            <RouterLink
              v-if="n.type === 'message'"
              :to="{ name: 'messages' }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              Open Messages
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'group_message'"
              :to="{ name: 'messages' }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              Open Messages
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'follow_request'"
              :to="{ name: 'follow-requests' }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View follow requests
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'new_follower'"
              :to="{ name: 'profile-followers', params: { id: String(myUserId) } }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View followers
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'group_invitation'"
              :to="{ name: 'group-invitations' }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View invitations
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'group_join_request' && n.related_group_id"
              :to="{ name: 'group', params: { id: String(n.related_group_id) } }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View group
            </RouterLink>
            <RouterLink
              v-else-if="(n.type === 'group_event_created' || n.type === 'group_event_response') && n.related_group_id"
              :to="{ name: 'group', params: { id: String(n.related_group_id) } }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View group
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'group_invitation_response'"
              :to="{ name: 'groups' }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View groups
            </RouterLink>
            <RouterLink
              v-else-if="n.type === 'group_join_response'"
              :to="{ name: 'groups' }"
              class="notification-action-link"
              @click="markAsRead(n.id)"
            >
              View groups
            </RouterLink>
          </div>
        </div>
        <button
          v-if="!n.is_read"
          type="button"
          class="button secondary notification-mark-read"
          @click="markAsRead(n.id)"
        >
          Mark read
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { onMounted, computed } from "vue";
import { useAuthStore } from "../stores/auth";
import { useNotificationsStore } from "../stores/notifications";

const authStore = useAuthStore();
const notificationsStore = useNotificationsStore();
const myUserId = computed(() => authStore.userId ?? 0);

onMounted(() => {
  notificationsStore.fetchNotifications();
});

function markAsRead(id) {
  notificationsStore.markAsRead(id);
}

async function markAllRead() {
  await notificationsStore.markAllAsRead();
}

function typeLabel(type) {
  const labels = {
    follow_request: "Follow request",
    new_follower: "New follower",
    group_invitation: "Group invitation",
    group_join_request: "Join request",
    group_join_response: "Join response",
    group_invitation_response: "Invitation response",
    group_event_created: "New event",
    group_event_response: "Event response",
    message: "Private message",
    group_message: "Group message",
    welcome: "Welcome"
  };
  return labels[type] ?? type;
}

function formatTime(createdAt) {
  if (!createdAt) return "";
  const d = new Date(createdAt);
  const now = new Date();
  const diffMs = now - d;
  if (diffMs < 60 * 1000) return "just now";
  if (diffMs < 60 * 60 * 1000) return `${Math.floor(diffMs / 60000)}m ago`;
  if (diffMs < 24 * 60 * 60 * 1000) return `${Math.floor(diffMs / 3600000)}h ago`;
  return d.toLocaleDateString();
}

/** Show username instead of email in notification text */
function notificationMessageDisplay(message) {
  if (!message || typeof message !== "string") return message || "";
  return message.replace(/\S+@\S+\.\S+/g, (email) => email.split("@")[0]);
}
</script>

<style scoped>
.notifications-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}
.notifications-header h1 {
  margin: 0;
}

.notifications-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.notification-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 0;
  border-bottom: 1px solid var(--border);
}
.notification-item:last-child {
  border-bottom: none;
}
.notification-item--unread {
  background: linear-gradient(
    90deg,
    rgba(0, 0, 0, 0.04) 0%,
    transparent 100%
  );
  margin: 0 -24px;
  padding-left: 24px;
  padding-right: 24px;
}
.notification-body {
  flex: 1;
  min-width: 0;
}
.notification-type-label {
  display: inline-block;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.02em;
  opacity: 0.7;
  margin-bottom: 4px;
}
.notification-item--message .notification-type-label {
  opacity: 0.7;
}
.notification-message {
  display: block;
  font-size: 15px;
}
.notification-meta {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  opacity: 0.7;
}
.notification-actions {
  margin-top: 8px;
}
.notification-action-link {
  font-size: 13px;
  color: var(--link);
  font-weight: 500;
}
.notification-action-link:hover {
  text-decoration: underline;
}
.notification-item--message .notification-action-link {
  color: var(--link);
}
.notification-mark-read {
  flex-shrink: 0;
  font-size: 13px;
}
</style>
