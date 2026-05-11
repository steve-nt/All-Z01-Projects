<template>
  <!-- Part 1: shared layout shell with navigation -->
  <div class="app-shell">
    <ToastContainer />
    <nav v-if="!isStandaloneAuthPage" class="nav">
      <div>
        <div class="nav-brand" role="banner" aria-label="AthensZone Social">
          <span class="nav-brand-main">AthensZone</span>
          <span class="nav-brand-tag">Social</span>
        </div>
      </div>
      <div class="nav-actions">
        <template v-if="authStore.sessionChecked && authStore.loggedIn">
          <!-- Auth-only navigation -->
          <RouterLink to="/feed">Feed</RouterLink>
          <RouterLink to="/profile/me">My Profile</RouterLink>
          <RouterLink to="/groups">Groups</RouterLink>
          <RouterLink to="/groups/invitations/list">Group invitations</RouterLink>
          <RouterLink to="/people">Find People</RouterLink>
          <div class="nav-notifications-wrap" ref="notificationsWrapRef">
            <button
              type="button"
              class="nav-notifications-trigger"
              :class="{ 'nav-notifications-trigger--open': showNotificationsDropdown, 'nav-notifications-trigger--has-unread': notificationsStore.unreadCount > 0 }"
              aria-haspopup="true"
              :aria-expanded="showNotificationsDropdown"
              @click="toggleNotificationsDropdown"
            >
              Notifications
              <span v-if="notificationsStore.unreadCount > 0" class="nav-badge nav-badge--unread">
                {{ notificationsStore.unreadCount > 99 ? "99+" : notificationsStore.unreadCount }}
              </span>
            </button>
            <div
              v-if="showNotificationsDropdown"
              class="nav-notifications-dropdown"
              role="menu"
            >
              <div class="nav-notifications-dropdown-header">
                <span>Recent</span>
                <RouterLink to="/notifications" @click="closeNotificationsDropdown">See all</RouterLink>
              </div>
              <div v-if="recentNotifications.length === 0" class="nav-notifications-empty">
                No new notifications
              </div>
              <ul v-else class="nav-notifications-dropdown-list">
                <li
                  v-for="n in recentNotifications"
                  :key="n.id"
                  class="nav-notifications-dropdown-item"
                  :class="{ 'nav-notifications-dropdown-item--unread': !n.is_read }"
                >
                  <span class="nav-notifications-dropdown-msg">{{ notificationMessageDisplay(n.message) }}</span>
                  <span class="nav-notifications-dropdown-meta">{{ formatNotificationTime(n.created_at) }}</span>
                </li>
              </ul>
            </div>
          </div>
          <RouterLink
            to="/follow/requests"
            class="nav-follow-requests-link"
          >
            Follow Requests
            <span
              v-if="followRequestsUnreadCount > 0"
              class="nav-follow-requests-dot"
              aria-hidden="true"
            />
          </RouterLink>
          <RouterLink to="/messages">Messages</RouterLink>
          <button class="button secondary" type="button" @click="logout">
            Logout
          </button>
        </template>
        <template v-else>
          <!-- No nav links when logged out; Login/Register are on the welcome page -->
        </template>
      </div>
    </nav>

    <main :class="['main-content', { 'main-content--full': isStandaloneAuthPage }]">
      <!-- Part 1: global error banner (403 / network) -->
      <div v-if="!isStandaloneAuthPage && errorStore.message" class="banner">
        <div>
          <strong v-if="errorStore.type === 'forbidden'">Forbidden:</strong>
          <strong v-else-if="errorStore.type === 'network'">Network:</strong>
          <strong v-else>Notice:</strong>
          {{ errorStore.message }}
        </div>
        <button class="button secondary" type="button" @click="errorStore.clear">
          Dismiss
        </button>
      </div>
      <!-- Part 1: initial session loading state -->
      <div v-if="authStore.isChecking && !isStandaloneAuthPage" class="page-card">
        <p class="muted">Checking session...</p>
      </div>
      <RouterView v-else />
    </main>

    <footer class="app-footer">
      <p class="app-footer-quote">"Connection is why we're here. It's what gives purpose and meaning to our lives."</p>
      <p class="app-footer-credit">© 2026 AthensZone Social. Created with care.</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { RouterLink, RouterView, useRouter } from "vue-router";
import ToastContainer from "./components/ToastContainer.vue";
import { useAuthStore } from "./stores/auth";
import { useErrorStore } from "./stores/error";
import { useNotificationsStore } from "./stores/notifications";
import { useMessagesStore } from "./stores/messages";
import { wsService } from "./services/websocket";

const authStore = useAuthStore();
const errorStore = useErrorStore();
const notificationsStore = useNotificationsStore();
const messagesStore = useMessagesStore();
const router = useRouter();

const showNotificationsDropdown = ref(false);
const notificationsWrapRef = ref(null);

const isStandaloneAuthPage = computed(() => {
  const name = router.currentRoute.value.name;
  return name === "login" || name === "register";
});

const recentNotifications = computed(() =>
  notificationsStore.notifications.slice(0, 7)
);

const followRequestsUnreadCount = computed(() =>
  notificationsStore.notifications.filter(
    (n) => n.type === "follow_request" && !n.is_read
  ).length
);

function formatNotificationTime(createdAt) {
  if (!createdAt) return "";
  const d = new Date(createdAt);
  const now = new Date();
  const diffMs = now - d;
  if (diffMs < 60 * 1000) return "just now";
  if (diffMs < 60 * 60 * 1000) return `${Math.floor(diffMs / 60000)}m ago`;
  if (diffMs < 24 * 60 * 60 * 1000) return `${Math.floor(diffMs / 3600000)}h ago`;
  return d.toLocaleDateString();
}

/** Show username instead of email in notification text (e.g. "from user@mail.com" -> "from user") */
function notificationMessageDisplay(message) {
  if (!message || typeof message !== "string") return message || "";
  return message.replace(/\S+@\S+\.\S+/g, (email) => email.split("@")[0]);
}

function toggleNotificationsDropdown() {
  showNotificationsDropdown.value = !showNotificationsDropdown.value;
  if (showNotificationsDropdown.value && !notificationsStore.loaded) {
    notificationsStore.fetchNotifications().catch(() => {});
  }
}

function closeNotificationsDropdown() {
  showNotificationsDropdown.value = false;
}

function handleClickOutside(event) {
  if (
    showNotificationsDropdown.value &&
    notificationsWrapRef.value &&
    !notificationsWrapRef.value.contains(event.target)
  ) {
    showNotificationsDropdown.value = false;
  }
}

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});
onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside);
});

const logout = async () => {
  wsService.disconnect();
  notificationsStore.reset();
  messagesStore.clearPresence();
  await authStore.logout();
  if (router.currentRoute.value.name !== "login") {
    await router.push({ name: "login" });
  }
};
</script>

<style scoped>
.main-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  width: min(1120px, 92%);
  margin: 24px auto;
}
.main-content--full {
  width: 100%;
  max-width: none;
  margin: 0;
  padding: 48px 24px;
  min-height: 0;
  background: var(--bg);
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-footer {
  flex-shrink: 0;
  padding: 14px 24px;
  background: #000;
  color: rgba(255, 255, 255, 0.85);
  border-top: 1px solid var(--border-inverse);
  text-align: center;
}

.app-footer-quote {
  margin: 0 0 6px;
  font-size: 0.875rem;
  font-style: italic;
  max-width: 520px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.4;
}

.app-footer-credit {
  margin: 0;
  font-size: 0.75rem;
  opacity: 0.75;
}

.nav-notifications-wrap {
  position: relative;
}
.nav-notifications-trigger {
  display: inline-flex;
  align-items: center;
  padding: 0;
  margin: 0;
  border: none;
  background: transparent;
  color: inherit;
  font: inherit;
  cursor: pointer;
  text-decoration: none;
}
.nav-notifications-trigger:hover {
  text-decoration: underline;
}
.nav-notifications-trigger--open {
  text-decoration: underline;
}
.nav-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  margin-left: 4px;
  font-size: 11px;
  font-weight: 700;
  color: #000;
  background: #fff;
  border-radius: 9px;
  vertical-align: middle;
}
.nav-badge--unread {
  background: #dc2626;
  color: #fff;
  animation: nav-badge-pulse 1.5s ease-in-out infinite;
}
@keyframes nav-badge-pulse {
  50% { opacity: 0.9; }
}
.nav-notifications-trigger--has-unread {
  font-weight: 700;
}

.nav-follow-requests-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.nav-follow-requests-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #dc2626;
  display: inline-block;
}
.nav-notifications-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 6px;
  min-width: 280px;
  max-width: 360px;
  max-height: 70vh;
  overflow: auto;
  background: var(--surface);
  color: var(--surface-text);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  box-shadow: var(--shadow);
  z-index: 1000;
}
.nav-notifications-dropdown-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid var(--border);
  font-weight: 600;
  font-size: 14px;
}
.nav-notifications-dropdown-header a {
  color: var(--link);
  font-weight: 500;
  font-size: 13px;
}
.nav-notifications-empty {
  padding: 20px 14px;
  opacity: 0.7;
  font-size: 14px;
}
.nav-notifications-dropdown-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.nav-notifications-dropdown-item {
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
  font-size: 14px;
}
.nav-notifications-dropdown-item--unread {
  background: rgba(0, 0, 0, 0.04);
}
.nav-notifications-dropdown-msg {
  display: block;
}
.nav-notifications-dropdown-meta {
  display: block;
  margin-top: 2px;
  font-size: 12px;
  opacity: 0.7;
}
</style>
