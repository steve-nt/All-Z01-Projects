import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import { useAuthStore } from "./stores/auth";
import { useErrorStore } from "./stores/error";
import { useNotificationsStore } from "./stores/notifications";
import { useToastStore } from "./stores/toast";
import { useMessagesStore } from "./stores/messages";
import {
  setForbiddenHandler,
  setNetworkErrorHandler,
  setUnauthorizedHandler
} from "./services/apiClient";
import { wsService } from "./services/websocket";
import "./assets/base.css";

// Part 1: initialize Vue app and global stores
const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

const authStore = useAuthStore();
const errorStore = useErrorStore();
// Part 1: global HTTP error handling (shared conventions)
setUnauthorizedHandler(() => {
  wsService.disconnect();
  const notificationsStore = useNotificationsStore();
  const messagesStore = useMessagesStore();
  notificationsStore.reset();
  messagesStore.clearPresence();
  authStore.clearSession();
  if (router.currentRoute.value.path !== "/login") {
    router.push("/login");
  }
});
setForbiddenHandler((path, payload) => {
  const isGroupPath =
    typeof path === "string" && path.includes("/api/groups/");
  const message = isGroupPath
    ? "Members only. Join or request to join this group to see content."
    : (payload && (payload.message || payload.error)) ||
      "You do not have access to this resource.";
  errorStore.setError(message, "forbidden");
});
setNetworkErrorHandler(() => {
  errorStore.setError(
    "Network error. Make sure the backend is running.",
    "network"
  );
});

// Part 1: wait for router + session check before mounting
router.isReady().then(async () => {
  await authStore.checkSession();
  const notificationsStore = useNotificationsStore();
  if (authStore.loggedIn) {
    wsService.connect();
    notificationsStore.fetchNotifications().catch(() => {});
  }
  app.mount("#app");

  // Real-time: show toast, update notifications in header, update unread in Messages sidebar
  const toastStore = useToastStore();
  const messagesStore = useMessagesStore();
  wsService.on("message", (data) => {
    if (data.type === "private_message") {
      const payload = data.data ?? data;
      messagesStore.appendPrivateMessage({
        sender_id: payload.sender_id,
        recipient_id: payload.recipient_id,
        message_id: payload.message_id,
        content: payload.content,
        created_at: payload.created_at,
        sender_name: payload.sender_name
      });
      const selected = messagesStore.selectedChat;
      if (!selected || selected.type !== "private" || selected.id !== payload.sender_id) {
        messagesStore.incrementPrivateUnread(payload.sender_id);
      }
      toastStore.privateMessage(payload?.content ?? "New message");
      // Add to header notifications (same as group) so it's always visible
      const notifMessage = `New message from ${payload.sender_name || "Someone"}`;
      notificationsStore.addNotification({
        type: "message",
        message: notifMessage,
        related_user_id: payload.sender_id,
        created_at: payload.created_at || new Date().toISOString()
      });
      return;
    }
    if (data.type === "group_message") {
      const payload = data.data ?? data;
      messagesStore.appendGroupMessage({
        group_id: payload.group_id,
        sender_id: payload.sender_id,
        message_id: payload.message_id,
        content: payload.content,
        created_at: payload.created_at,
        sender_name: payload.sender_name
      });
      const selected = messagesStore.selectedChat;
      if (!selected || selected.type !== "group" || selected.id !== payload.group_id) {
        messagesStore.incrementGroupUnread(payload.group_id);
      }
      toastStore.notification(
        `New message from ${payload.sender_name || "Someone"} (group)`
      );
      // Add to header notifications (same as private)
      notificationsStore.addNotification({
        type: "group_message",
        message: `New message from ${payload.sender_name || "Someone"} (group)`,
        related_group_id: payload.group_id,
        related_user_id: payload.sender_id,
        created_at: payload.created_at || new Date().toISOString()
      });
      return;
    }
    if (data.type === "presence") {
      const payload = data.data ?? data;
      const userId = payload.user_id;
      const online = payload.online === true;
      if (userId != null) messagesStore.setPresence(userId, online);
      return;
    }
    if (data.type === "notification") {
      const payload = data.data ?? data;
      const message = payload?.message ?? "New notification";
      const isPrivateMessage = payload?.type === "message";
      const isGroupMessage = payload?.type === "group_message";
      // Toast: only for non-message types (we already show toast in private_message / group_message handlers)
      if (!isPrivateMessage && !isGroupMessage) {
        toastStore.notification(message);
      }
      // Skip adding to header for message/group_message - we already added from private_message/group_message handlers
      if (isPrivateMessage || isGroupMessage) return;
      notificationsStore.addNotification(payload);
    }
  });
});
