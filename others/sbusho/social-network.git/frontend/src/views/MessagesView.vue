<template>
  <div class="messages-page">
    <div class="messages-layout">
      <!-- Left sidebar: conversations + groups -->
      <aside class="messages-sidebar">
        <div class="messages-sidebar-header">
          <h1 class="messages-sidebar-title">Messages</h1>
        </div>

        <div class="messages-sidebar-tabs">
          <button
            type="button"
            class="messages-tab"
            :class="{ 'messages-tab--active': activeTab === 'private' }"
            @click="activeTab = 'private'"
          >
            Private
          </button>
          <button
            type="button"
            class="messages-tab"
            :class="{ 'messages-tab--active': activeTab === 'groups' }"
            @click="activeTab = 'groups'"
          >
            Groups
          </button>
        </div>

        <div v-if="activeTab === 'private'" class="messages-conversation-list">
          <div v-if="messagesStore.loadingConversations || messagesStore.loadingContacts" class="messages-loading">Loading…</div>
          <template v-else-if="privateList.length === 0">
            <p class="messages-empty">No one to message yet. Follow someone or accept a follow request to start chatting.</p>
          </template>
          <button
            v-for="c in privateList"
            :key="c.user_id"
            type="button"
            class="messages-conversation-item"
            :class="{ 'messages-conversation-item--active': isSelectedPrivate(c.user_id) }"
            @click="selectPrivate(c)"
          >
            <div class="messages-conv-avatar-wrap">
              <img
                v-if="c.user_avatar"
                :src="avatarUrl(c.user_avatar)"
                alt=""
                class="messages-conv-avatar"
              />
              <div v-else class="messages-conv-avatar messages-conv-avatar--placeholder">
                {{ displayName(c.user_name, c.user_id).charAt(0).toUpperCase() }}
              </div>
              <span
                class="messages-status-dot"
                :class="messagesStore.isUserOnline(c.user_id) ? 'messages-status-dot--online' : 'messages-status-dot--offline'"
                :title="messagesStore.isUserOnline(c.user_id) ? 'Online' : 'Offline'"
              />
            </div>
            <div class="messages-conv-body">
              <span class="messages-conv-name">{{ displayName(c.user_name, c.user_id) }}</span>
              <span class="messages-conv-meta">
                {{ c.last_message_at ? formatTime(c.last_message_at) : (messagesStore.isUserOnline(c.user_id) ? 'Online' : 'Offline') }}
              </span>
            </div>
            <span v-if="(c.unread_count || 0) > 0" class="messages-unread-badge">
              {{ (c.unread_count || 0) > 99 ? '99+' : (c.unread_count || 0) }}
            </span>
          </button>
        </div>

        <div v-else class="messages-conversation-list">
          <div v-if="messagesStore.groupChats.length === 0" class="messages-empty">
            You are not in any groups. Join a group to see its chat.
          </div>
          <button
            v-for="g in messagesStore.groupChats"
            :key="g.group_id"
            type="button"
            class="messages-conversation-item messages-conversation-item--group"
            :class="{ 'messages-conversation-item--active': isSelectedGroup(g.group_id) }"
            @click="selectGroup(g)"
          >
            <div class="messages-conv-avatar messages-conv-avatar--group">
              {{ (g.group_name || 'G').charAt(0).toUpperCase() }}
            </div>
            <div class="messages-conv-body">
              <span class="messages-conv-name">{{ g.group_name || 'Group' }}</span>
            </div>
            <span v-if="messagesStore.unreadCountForGroup(g.group_id) > 0" class="messages-unread-badge">
              {{ messagesStore.unreadCountForGroup(g.group_id) > 99 ? '99+' : messagesStore.unreadCountForGroup(g.group_id) }}
            </span>
          </button>
        </div>
      </aside>

      <!-- Right: chat panel -->
      <section class="messages-chat-panel">
        <template v-if="messagesStore.selectedChat">
          <header class="messages-chat-header">
            <div class="messages-chat-header-info">
              <template v-if="messagesStore.selectedChat.type === 'private'">
                <span class="messages-chat-title">{{ selectedPrivateName }}</span>
                <span
                  class="messages-chat-status"
                  :class="messagesStore.isUserOnline(messagesStore.selectedChat.id) ? 'messages-chat-status--online' : 'messages-chat-status--offline'"
                >
                  {{ messagesStore.isUserOnline(messagesStore.selectedChat.id) ? 'Online' : 'Offline' }}
                </span>
                <!-- onlineTick forces re-evaluation of isUserOnline when time passes -->
                <span v-if="messagesStore.selectedChat.type === 'private'" class="sr-only">{{ onlineTick }}</span>
              </template>
              <template v-else>
                <span class="messages-chat-title">{{ selectedGroupName }}</span>
                <span class="messages-chat-status messages-chat-status--group">Group chat</span>
              </template>
            </div>
          </header>

          <div ref="messagesEndRef" class="messages-list-wrap">
            <div v-if="messagesStore.loadingMessages" class="messages-loading-inline">Loading messages…</div>
            <ul v-else class="messages-list">
              <li
                v-for="m in messagesStore.messages"
                :key="m.id"
                class="messages-list-item"
                :class="{ 'messages-list-item--sent': isSent(m) }"
              >
                <div class="messages-bubble">
                  <span class="messages-bubble-sender">{{ messageSenderName(m) }}</span>
                  <span class="messages-bubble-content">{{ m.content }}</span>
                  <span class="messages-bubble-time">{{ formatMessageTime(m.created_at) }}</span>
                </div>
              </li>
            </ul>
          </div>

          <div class="messages-input-wrap">
            <div v-if="messagesStore.error" class="messages-input-error">{{ messagesStore.error }}</div>
            <div class="messages-input-row">
              <button
                type="button"
                class="messages-emoji-btn"
                aria-label="Insert emoji"
                @click="showEmojiPicker = !showEmojiPicker"
              >
                🙂
              </button>
              <div v-if="showEmojiPicker" class="messages-emoji-picker">
                <button
                  v-for="emoji in emojiList"
                  :key="emoji"
                  type="button"
                  class="messages-emoji-item"
                  @click="insertEmoji(emoji)"
                >
                  {{ emoji }}
                </button>
              </div>
              <input
                v-model="inputText"
                type="text"
                class="messages-input"
                placeholder="Type a message…"
                :disabled="messagesStore.sending"
                @keydown.enter.prevent="sendMessage"
              />
              <button
                type="button"
                class="button messages-send-btn"
                :disabled="!inputText.trim() || messagesStore.sending"
                @click="sendMessage"
              >
                {{ messagesStore.sending ? 'Sending…' : 'Send' }}
              </button>
            </div>
          </div>
        </template>
        <div v-else class="messages-chat-placeholder">
          <p>Select a conversation or group to start chatting.</p>
          <p class="muted">Private chat is only with users you follow or who follow you. Group chat is available for groups you are a member of.</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from "vue";
import { useAuthStore } from "../stores/auth";
import { useMessagesStore } from "../stores/messages";

const authStore = useAuthStore();
const messagesStore = useMessagesStore();

const activeTab = ref("private");
const inputText = ref("");
const showEmojiPicker = ref(false);
const messagesEndRef = ref(null);
// Force re-render every 12s so online/offline status updates (60s threshold in store)
const onlineTick = ref(0);
let onlineTickInterval;
onMounted(() => {
  onlineTickInterval = setInterval(() => { onlineTick.value++; }, 12000);
});
onUnmounted(() => {
  if (onlineTickInterval) clearInterval(onlineTickInterval);
});

const emojiList = [
  "😀", "😃", "😄", "😁", "😅", "😂", "🤣", "😊", "😇", "🙂",
  "😉", "😍", "🥰", "😘", "👍", "👋", "❤️", "🔥", "✨", "🎉",
  "💯", "🙏", "👏", "😢", "😭", "😤", "🤔", "😎", "🥳", "💪"
];

function avatarUrl(path) {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  return path.startsWith("/") ? path : `/${path}`;
}

/** Show username (nickname) when possible; never show full email – use part before @ or "User #id". */
function displayName(name, userId) {
  if (!name || !name.trim()) return "User #" + (userId ?? "");
  const n = name.trim();
  if (n.includes("@")) return n.split("@")[0];
  return n;
}

function formatTime(createdAt) {
  if (!createdAt) return "";
  const d = new Date(createdAt);
  const now = new Date();
  const diffMs = now - d;
  if (diffMs < 60 * 1000) return "Just now";
  if (diffMs < 60 * 60 * 1000) return `${Math.floor(diffMs / 60000)}m`;
  if (diffMs < 24 * 60 * 60 * 1000) return `${Math.floor(diffMs / 3600000)}h`;
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

function formatMessageTime(createdAt) {
  if (!createdAt) return "";
  const d = new Date(createdAt);
  return d.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" });
}

const selectedPrivateName = computed(() => {
  if (messagesStore.selectedChat?.type !== "private") return "";
  const conv = messagesStore.conversations.find((x) => x.user_id === messagesStore.selectedChat.id);
  if (conv) return displayName(conv.user_name, conv.user_id);
  const contact = messagesStore.contacts.find((x) => x.user_id === messagesStore.selectedChat.id);
  return contact ? displayName(contact.user_name, contact.user_id) : "User #" + messagesStore.selectedChat.id;
});

/** Merged list: conversations first (with last message), then contacts not yet in conversations. Shows online/offline. */
const privateList = computed(() => {
  const convIds = new Set(messagesStore.conversations.map((c) => c.user_id));
  const contactsOnly = messagesStore.contacts
    .filter((c) => !convIds.has(c.user_id))
    .map((c) => ({
      user_id: c.user_id,
      user_name: c.user_name,
      user_avatar: c.user_avatar,
      last_message_at: null,
      unread_count: 0
    }));
  return [...messagesStore.conversations, ...contactsOnly];
});

const selectedGroupName = computed(() => {
  if (messagesStore.selectedChat?.type !== "group") return "";
  const g = messagesStore.groupChats.find((x) => x.group_id === messagesStore.selectedChat.id);
  return g?.group_name || "Group";
});

function isSelectedPrivate(userId) {
  return messagesStore.selectedChat?.type === "private" && messagesStore.selectedChat?.id === userId;
}

function isSelectedGroup(groupId) {
  return messagesStore.selectedChat?.type === "group" && messagesStore.selectedChat?.id === groupId;
}

function isSent(m) {
  return m.sender_id === authStore.userId;
}

/** Display name for the sender of a message (shown above every bubble). Never show email – use username or part before @. */
function messageSenderName(m) {
  if (isSent(m)) {
    return displayName(authStore.nickname, authStore.userId) || "Me";
  }
  return displayName(m.sender_name, m.sender_id) || "User";
}

async function selectPrivate(c) {
  messagesStore.setSelectedChat("private", c.user_id);
  await messagesStore.fetchPrivateMessages(c.user_id);
  if (c.unread_count > 0) messagesStore.markAsRead(c.user_id);
  scrollToBottom();
}

async function selectGroup(g) {
  messagesStore.setSelectedChat("group", g.group_id);
  await messagesStore.fetchGroupMessages(g.group_id);
  messagesStore.clearGroupUnread(g.group_id);
  scrollToBottom();
}

function insertEmoji(emoji) {
  inputText.value += emoji;
}

async function sendMessage() {
  const text = inputText.value.trim();
  if (!text || !messagesStore.selectedChat) return;
  try {
    if (messagesStore.selectedChat.type === "private") {
      await messagesStore.sendPrivateMessage(messagesStore.selectedChat.id, text);
      await messagesStore.fetchConversations();
    } else {
      await messagesStore.sendGroupMessage(messagesStore.selectedChat.id, text);
    }
    inputText.value = "";
    await nextTick();
    scrollToBottom();
  } catch (_) {}
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesEndRef.value) {
      messagesEndRef.value.scrollTop = messagesEndRef.value.scrollHeight;
    }
  });
}

watch(
  () => messagesStore.messages.length,
  () => scrollToBottom()
);

onMounted(async () => {
  await Promise.all([
    messagesStore.fetchConversations(),
    messagesStore.fetchContacts(),
    messagesStore.fetchGroupChats()
  ]);
});
</script>

<style scoped>
.messages-page {
  background: var(--surface);
  color: var(--surface-text);
  border-radius: var(--radius);
  border: 1px solid var(--border);
  box-shadow: var(--shadow);
  overflow: hidden;
  height: 70vh;
  max-height: 70vh;
  max-width: 1000px;
  margin-left: auto;
  margin-right: auto;
}

.messages-layout {
  display: grid;
  grid-template-columns: 320px 1fr;
  height: 100%;
  max-height: 100%;
}

@media (max-width: 720px) {
  .messages-layout {
    grid-template-columns: 1fr;
  }
}

.messages-sidebar {
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border);
  background: var(--surface-2);
}

.messages-sidebar-header {
  padding: 16px;
  border-bottom: 1px solid var(--border);
}

.messages-sidebar-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
}

.messages-sidebar-tabs {
  display: flex;
  padding: 8px;
  gap: 4px;
}

.messages-tab {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  color: var(--surface-text);
}

.messages-tab:hover {
  background: var(--surface-2);
}

.messages-tab--active {
  background: #000;
  color: #fff;
  border-color: #000;
}

.messages-conversation-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.messages-loading,
.messages-empty {
  padding: 16px;
  margin: 0;
  font-size: 14px;
  color: var(--muted);
}

.messages-conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: inherit;
  text-align: left;
  cursor: pointer;
  font: inherit;
  margin-bottom: 4px;
}

.messages-conversation-item:hover {
  background: rgba(0, 0, 0, 0.06);
}

.messages-conversation-item--active {
  background: rgba(0, 0, 0, 0.1);
}

.messages-conv-avatar-wrap {
  position: relative;
  flex-shrink: 0;
}

.messages-conv-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  background: var(--surface-2);
}

.messages-conv-avatar--placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--surface-text);
}

.messages-conv-avatar--group {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  font-weight: 700;
  background: #333;
  color: #fff;
}

.messages-status-dot {
  position: absolute;
  bottom: 2px;
  right: 2px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid var(--surface-2);
}
.messages-status-dot--online {
  background: #22c55e;
}
.messages-status-dot--offline {
  background: #6b7280;
}

.messages-conv-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.messages-conv-name {
  font-weight: 600;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.messages-conv-meta {
  font-size: 12px;
  color: var(--muted);
}

.messages-unread-badge {
  flex-shrink: 0;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  background: #000;
  color: #fff;
  border-radius: 10px;
}

.messages-chat-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.messages-chat-header {
  padding: 14px 20px;
  border-bottom: 1px solid var(--border);
  background: var(--surface-2);
}

.messages-chat-header-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.messages-chat-title {
  font-weight: 700;
  font-size: 1rem;
}

.messages-chat-status {
  font-size: 13px;
  color: var(--muted);
}

.messages-chat-status--online {
  color: #22c55e;
}

.messages-chat-status--offline {
  color: var(--muted);
}

.messages-chat-status--group {
  color: var(--muted);
}

.messages-list-wrap {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  min-height: 200px;
}

.messages-loading-inline {
  padding: 16px;
  color: var(--muted);
  font-size: 14px;
}

.messages-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.messages-list-item {
  display: flex;
  justify-content: flex-start;
}

.messages-list-item--sent {
  justify-content: flex-end;
}

.messages-bubble {
  max-width: 75%;
  padding: 10px 14px;
  border-radius: 16px;
  background: var(--surface-2);
  border: 1px solid var(--border);
}

.messages-list-item--sent .messages-bubble {
  background: #000;
  color: #fff;
  border-color: #000;
}

.messages-bubble-sender {
  display: block;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 4px;
  color: var(--surface-text);
}

.messages-list-item--sent .messages-bubble-sender {
  color: rgba(255, 255, 255, 0.9);
}

.messages-bubble-content {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 14px;
}

.messages-bubble-time {
  display: block;
  font-size: 11px;
  margin-top: 4px;
  opacity: 0.8;
}

.messages-chat-placeholder {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  color: var(--muted);
}

.messages-chat-placeholder p {
  margin: 0 0 8px;
}

.messages-input-wrap {
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  background: var(--surface-2);
}

.messages-input-error {
  font-size: 13px;
  color: #b91c1c;
  margin-bottom: 8px;
}

.messages-input-row {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.messages-emoji-btn {
  width: 40px;
  height: 40px;
  padding: 0;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
  font-size: 1.25rem;
  cursor: pointer;
  flex-shrink: 0;
}

.messages-emoji-btn:hover {
  background: var(--surface-2);
}

.messages-emoji-picker {
  position: absolute;
  bottom: 100%;
  left: 0;
  margin-bottom: 8px;
  padding: 8px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow);
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 4px;
  max-height: 200px;
  overflow-y: auto;
  z-index: 10;
}

.messages-emoji-item {
  width: 36px;
  height: 36px;
  padding: 0;
  border: none;
  border-radius: 6px;
  background: transparent;
  font-size: 1.25rem;
  cursor: pointer;
}

.messages-emoji-item:hover {
  background: var(--surface-2);
}

.messages-input {
  flex: 1;
  min-width: 0;
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 14px;
  background: var(--surface);
  color: var(--surface-text);
}

.messages-input:focus {
  outline: none;
  border-color: var(--border-strong);
  box-shadow: 0 0 0 2px var(--ring-surface);
}

.messages-send-btn {
  flex-shrink: 0;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
