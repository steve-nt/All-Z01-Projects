import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";
import { useAuthStore } from "./auth";
import { wsService } from "../services/websocket";

const ONLINE_THRESHOLD_MS = 60 * 1000; // 60s

export const useMessagesStore = defineStore("messages", {
  state: () => ({
    conversations: [],
    contacts: [],
    groupChats: [],
    messages: [],
    selectedChat: null,
    loadingConversations: false,
    loadingContacts: false,
    loadingMessages: false,
    sending: false,
    error: null,
    lastSeenByUser: {},
    groupUnread: {},
    /** User IDs currently connected via WebSocket (from server presence) */
    connectedUserIds: {}
  }),

  getters: {
    isUserOnline(state) {
      return (userId) => {
        const auth = useAuthStore();
        if (userId === auth.userId) {
          return wsService.isOpen();
        }
        if (state.connectedUserIds[userId]) return true;
        const t = state.lastSeenByUser[userId];
        return t != null && Date.now() - t < ONLINE_THRESHOLD_MS;
      };
    },
    unreadCountForGroup: (state) => (groupId) => state.groupUnread[groupId] || 0,
    selectedConversation() {
      if (!this.selectedChat) return null;
      if (this.selectedChat.type === "private") {
        return this.conversations.find((c) => c.user_id === this.selectedChat.id) || null;
      }
      return this.groupChats.find((g) => g.group_id === this.selectedChat.id) || null;
    }
  },

  actions: {
    setLastSeen(userId) {
      this.lastSeenByUser = { ...this.lastSeenByUser, [userId]: Date.now() };
    },

    /** Update presence from WebSocket (user_id came online or went offline) */
    setPresence(userId, online) {
      if (online) {
        this.connectedUserIds = { ...this.connectedUserIds, [userId]: true };
      } else {
        const next = { ...this.connectedUserIds };
        delete next[userId];
        this.connectedUserIds = next;
      }
    },

    setSelectedChat(type, id) {
      this.selectedChat = id != null ? { type, id } : null;
      this.messages = [];
    },

    async fetchConversations() {
      this.loadingConversations = true;
      this.error = null;
      try {
        const data = await socialApi.getConversations();
        this.conversations = data?.conversations || [];
        return this.conversations;
      } catch (e) {
        this.error = e?.message || "Failed to load conversations";
        this.conversations = [];
        throw e;
      } finally {
        this.loadingConversations = false;
      }
    },

    async fetchContacts() {
      this.loadingContacts = true;
      try {
        const data = await socialApi.getMessageContacts();
        this.contacts = data?.contacts || [];
        return this.contacts;
      } catch (e) {
        this.contacts = [];
        return [];
      } finally {
        this.loadingContacts = false;
      }
    },

    async fetchGroupChats() {
      try {
        const data = await socialApi.getGroups();
        const groups = data?.groups || [];
        this.groupChats = groups.filter((g) => g.is_member === true);
        return this.groupChats;
      } catch (e) {
        this.groupChats = [];
        return [];
      }
    },

    async fetchPrivateMessages(userId) {
      if (!userId) {
        this.messages = [];
        return [];
      }
      this.loadingMessages = true;
      this.error = null;
      try {
        const data = await socialApi.getMessages(userId);
        this.messages = data?.messages || [];
        return this.messages;
      } catch (e) {
        this.error = e?.message || "Failed to load messages";
        this.messages = [];
        throw e;
      } finally {
        this.loadingMessages = false;
      }
    },

    async fetchGroupMessages(groupId) {
      if (!groupId) {
        this.messages = [];
        return [];
      }
      this.loadingMessages = true;
      this.error = null;
      try {
        const data = await socialApi.getGroupMessages(groupId);
        this.messages = data?.messages || [];
        return this.messages;
      } catch (e) {
        this.error = e?.message || "Failed to load group messages";
        this.messages = [];
        throw e;
      } finally {
        this.loadingMessages = false;
      }
    },

    async sendPrivateMessage(recipientId, content) {
      const auth = useAuthStore();
      const myId = auth.userId;
      const myName = auth.nickname || "Me";
      this.sending = true;
      this.error = null;
      try {
        const result = await socialApi.sendMessage(recipientId, content);
        const msg = {
          id: result?.message_id,
          sender_id: myId,
          sender_name: myName,
          recipient_id: recipientId,
          content,
          created_at: new Date().toISOString()
        };
        this.messages = [...this.messages, msg];
        return msg;
      } catch (e) {
        this.error = e?.message || "Failed to send message";
        throw e;
      } finally {
        this.sending = false;
      }
    },

    async sendGroupMessage(groupId, content) {
      const auth = useAuthStore();
      const myId = auth.userId;
      const myName = auth.nickname || "Me";
      this.sending = true;
      this.error = null;
      try {
        const result = await socialApi.sendGroupMessage(groupId, content);
        const msg = {
          id: result?.message_id,
          sender_id: myId,
          sender_name: myName,
          group_id: groupId,
          content,
          created_at: new Date().toISOString()
        };
        this.messages = [...this.messages, msg];
        return msg;
      } catch (e) {
        this.error = e?.message || "Failed to send message";
        throw e;
      } finally {
        this.sending = false;
      }
    },

    async markAsRead(senderId) {
      try {
        await socialApi.markMessagesRead({ sender_id: senderId });
      } catch (_) {}
    },

    appendPrivateMessage(payload) {
      const { sender_id, recipient_id, message_id, content, created_at, sender_name } = payload;
      if (this.selectedChat?.type !== "private" || this.selectedChat?.id !== sender_id) {
        return;
      }
      const exists = this.messages.some((m) => m.id === message_id);
      if (exists) return;
      this.messages = [
        ...this.messages,
        {
          id: message_id,
          sender_id,
          sender_name: sender_name || "User",
          recipient_id,
          content,
          created_at: created_at || new Date().toISOString()
        }
      ];
      this.setLastSeen(sender_id);
    },

    appendGroupMessage(payload) {
      const { group_id, sender_id, message_id, content, created_at, sender_name } = payload;
      if (this.selectedChat?.type !== "group" || this.selectedChat?.id !== group_id) {
        return;
      }
      const exists = this.messages.some((m) => m.id === message_id);
      if (exists) return;
      this.messages = [
        ...this.messages,
        {
          id: message_id,
          sender_id,
          sender_name: sender_name || "User",
          group_id,
          content,
          created_at: created_at || new Date().toISOString()
        }
      ];
      this.setLastSeen(sender_id);
    },

    incrementPrivateUnread(senderId) {
      const idx = this.conversations.findIndex((c) => c.user_id === senderId);
      if (idx === -1) return;
      const list = [...this.conversations];
      const c = { ...list[idx], unread_count: (list[idx].unread_count || 0) + 1 };
      list[idx] = c;
      this.conversations = list;
    },

    incrementGroupUnread(groupId) {
      this.groupUnread = {
        ...this.groupUnread,
        [groupId]: (this.groupUnread[groupId] || 0) + 1
      };
    },

    clearGroupUnread(groupId) {
      this.groupUnread = { ...this.groupUnread, [groupId]: 0 };
    },

    /** Clear presence (e.g. on logout) so reconnecting shows fresh state */
    clearPresence() {
      this.connectedUserIds = {};
    }
  }
});
