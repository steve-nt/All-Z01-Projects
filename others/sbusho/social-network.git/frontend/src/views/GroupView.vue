<template>
  <div class="page-card">
    <p v-if="groupsStore.loading" class="muted">Loading group...</p>
    <p v-else-if="groupsStore.error" class="error">{{ groupsStore.error.message }}</p>

    <template v-else-if="groupData">
      <h1>{{ groupData.group_name }}</h1>
      <p class="muted">Group ID: {{ groupData.group_id }}</p>
      <p class="muted">Creator ID: {{ groupData.creator_id }}</p>
      <p>{{ groupData.description }}</p>
      <p class="muted">Member status: {{ isMember ? "Member" : "Not a member" }}</p>

      <!-- Request to join (non-members only) -->
      <section v-if="!isMember" class="join-request-section">
        <button
          type="button"
          class="button"
          :disabled="groupsStore.requestingJoinGroupId === groupId"
          @click="onRequestToJoin"
        >
          {{ groupsStore.requestingJoinGroupId === groupId ? "Sending..." : "Request to join" }}
        </button>
        <p v-if="joinRequestMessage" :class="joinRequestError ? 'error' : 'muted'">
          {{ joinRequestMessage }}
        </p>
      </section>

      <!-- Join requests (creator only) -->
      <section v-if="isCreator" class="join-requests-section">
        <h2 class="section-heading">Join requests</h2>
        <p v-if="groupsStore.loadingJoinRequests" class="muted">Loading...</p>
        <p v-else-if="groupsStore.joinRequests.length === 0" class="muted">
          No pending join requests.
        </p>
        <ul v-else class="join-requests-list">
          <li
            v-for="req in groupsStore.joinRequests"
            :key="req.request_id"
            class="join-request-item"
          >
            <span class="requester-name">{{ req.requester_name || "User #" + req.requester_id }}</span>
            <div class="join-request-actions">
              <button
                type="button"
                class="button button-small"
                :disabled="groupsStore.respondingJoinRequestId === req.request_id"
                @click="onRespondJoinRequest(req.request_id, 'accepted')"
              >
                {{ groupsStore.respondingJoinRequestId === req.request_id ? "..." : "Accept" }}
              </button>
              <button
                type="button"
                class="button secondary button-small"
                :disabled="groupsStore.respondingJoinRequestId === req.request_id"
                @click="onRespondJoinRequest(req.request_id, 'declined')"
              >
                Decline
              </button>
            </div>
          </li>
        </ul>
      </section>

      <!-- Invite (members only) -->
      <section v-if="isMember" class="invite-section">
        <h2 class="invite-heading">Invite to group</h2>
        <div class="field">
          <label class="label" for="invite-search">Search by name or email</label>
          <input
            id="invite-search"
            v-model.trim="inviteQuery"
            class="input"
            type="search"
            placeholder="Type to search users..."
          />
        </div>
        <p v-if="inviteSearching" class="muted">Searching...</p>
        <p v-else-if="inviteQuery && inviteResults.length === 0" class="muted">
          No users found.
        </p>
        <ul v-else-if="inviteResults.length > 0" class="invite-list">
          <li v-for="user in inviteResults" :key="user.user_id" class="invite-item">
            <span class="invite-user-name">{{ user.nickname || user.email || "User #" + user.user_id }}</span>
            <button
              type="button"
              class="button button-small"
              :disabled="groupsStore.invitingUserId === user.user_id"
              @click="onInvite(user.user_id)"
            >
              {{ groupsStore.invitingUserId === user.user_id ? "Inviting..." : "Invite" }}
            </button>
          </li>
        </ul>
        <p v-if="inviteMessage" :class="inviteError ? 'error' : 'muted'">
          {{ inviteMessage }}
        </p>
      </section>

      <!-- Group posts (members only) -->
      <section v-if="isMember" class="group-posts-section">
        <h2 class="section-heading">Group posts</h2>
        <form class="form create-post-form" @submit.prevent="onSubmitGroupPost">
          <div class="field">
            <label class="label" for="group-post-content">Post to the group</label>
            <div class="emoji-bar">
              <button
                v-for="emoji in emojiList"
                :key="emoji"
                type="button"
                class="emoji-btn"
                :title="'Add ' + emoji"
                @click="insertGroupPostEmoji(emoji)"
              >
                {{ emoji }}
              </button>
            </div>
            <textarea
              id="group-post-content"
              ref="groupPostContentRef"
              v-model="groupPostContent"
              class="input textarea"
              rows="3"
              placeholder="Write something or add an image below (or both)"
            />
          </div>
          <div class="field">
            <label class="label" for="group-post-image">Image (optional)</label>
            <input
              id="group-post-image"
              type="file"
              accept="image/*"
              class="input"
              @change="onGroupPostFileSelect"
            />
          </div>
          <div class="actions">
            <button
              type="submit"
              class="button"
              :disabled="groupsStore.creatingGroupPost || (!(groupPostContent && groupPostContent.trim()) && !groupPostFile)"
            >
              {{ groupsStore.creatingGroupPost ? "Posting..." : "Post" }}
            </button>
          </div>
        </form>
        <p v-if="groupsStore.loadingGroupPosts" class="muted">Loading posts...</p>
        <p v-else-if="groupsStore.groupPosts.length === 0" class="muted">
          No posts yet. Be the first to post!
        </p>
        <ul v-else class="group-posts-list">
          <li
            v-for="post in groupsStore.groupPosts"
            :key="post.group_post_id"
            class="group-post-item"
          >
            <div class="group-post-meta">
              <span class="group-post-author">{{ post.author || "User #" + post.user_id }}</span>
              <span class="muted"> · {{ formatDate(post.created_at) }}</span>
            </div>
            <p class="group-post-content">{{ post.content }}</p>
            <img
              v-if="post.image_url"
              :src="post.image_url"
              alt="Post attachment"
              class="group-post-image"
            />
            <!-- Comments -->
            <div class="group-post-comments">
              <p v-if="commentsLoading(post.group_post_id)" class="muted comments-loading">
                Loading comments...
              </p>
              <template v-else>
                <ul class="comments-list">
                  <li
                    v-for="c in getGroupComments(post.group_post_id)"
                    :key="c.group_comment_id"
                    class="comment-item"
                  >
                    <span class="comment-author">{{ c.author || "User #" + c.user_id }}</span>
                    <span class="muted"> · {{ formatDate(c.created_at) }}</span>
                    <p class="comment-content">{{ c.content }}</p>
                    <img
                      v-if="c.image_url"
                      :src="c.image_url"
                      alt="Comment attachment"
                      class="comment-image"
                    />
                  </li>
                </ul>
                <form
                  class="comment-form"
                  @submit.prevent="onSubmitGroupComment(post.group_post_id)"
                >
                  <div class="emoji-bar">
                    <button
                      v-for="emoji in emojiList"
                      :key="emoji"
                      type="button"
                      class="emoji-btn"
                      :title="'Add ' + emoji"
                      @click="appendGroupCommentEmoji(post.group_post_id, emoji)"
                    >
                      {{ emoji }}
                    </button>
                  </div>
                  <textarea
                    v-model.trim="commentContentByPostId[post.group_post_id]"
                    class="input textarea comment-input"
                    rows="2"
                    placeholder="Write a comment..."
                  />
                  <div class="comment-form-row">
                    <label class="comment-file-label">
                      <input
                        type="file"
                        accept="image/*"
                        class="input-hidden"
                        :data-comment-file="post.group_post_id"
                        @change="(e) => onGroupCommentFileSelect(post.group_post_id, e)"
                      />
                      <span class="muted">Image (optional)</span>
                    </label>
                    <button
                      type="submit"
                      class="button button-small"
                      :disabled="isCreatingComment(post.group_post_id) || !(commentContentByPostId[post.group_post_id] || '').trim()"
                    >
                      {{ isCreatingComment(post.group_post_id) ? "Sending..." : "Comment" }}
                    </button>
                  </div>
                </form>
              </template>
            </div>
          </li>
        </ul>
      </section>

      <!-- Group events (members only) -->
      <section v-if="isMember" class="group-events-section">
        <h2 class="section-heading">Events</h2>
        <form class="form create-event-form" @submit.prevent="onSubmitEvent">
          <div class="field">
            <label class="label" for="event-title">Title</label>
            <input
              id="event-title"
              v-model.trim="eventTitle"
              class="input"
              type="text"
              placeholder="Event title"
              required
            />
          </div>
          <div class="field">
            <label class="label" for="event-description">Description</label>
            <textarea
              id="event-description"
              v-model.trim="eventDescription"
              class="input textarea"
              rows="2"
              placeholder="Description"
              required
            />
          </div>
          <div class="field">
            <label class="label" for="event-datetime">Date & time</label>
            <input
              id="event-datetime"
              v-model="eventDateTime"
              class="input"
              type="datetime-local"
              required
            />
          </div>
          <div class="actions">
            <button
              type="submit"
              class="button"
              :disabled="groupsStore.creatingGroupEvent || !eventTitle || !eventDescription || !eventDateTime"
            >
              {{ groupsStore.creatingGroupEvent ? "Creating..." : "Create event" }}
            </button>
          </div>
        </form>
        <p v-if="groupsStore.loadingGroupEvents" class="muted">Loading events...</p>
        <p v-else-if="groupsStore.groupEvents.length === 0" class="muted">
          No events yet.
        </p>
        <ul v-else class="group-events-list">
          <li
            v-for="ev in groupsStore.groupEvents"
            :key="ev.event_id"
            class="group-event-item"
          >
            <div class="event-title-row">
              <strong>{{ ev.title }}</strong>
              <span class="muted">{{ formatEventDate(ev.event_datetime) }}</span>
            </div>
            <p class="event-description">{{ ev.description }}</p>
            <div class="event-actions">
              <button
                type="button"
                class="button button-small"
                :disabled="groupsStore.respondingEventId === ev.event_id"
                @click="onEventResponse(ev.event_id, 'going')"
              >
                {{ groupsStore.respondingEventId === ev.event_id ? "..." : "Going" }}
              </button>
              <button
                type="button"
                class="button secondary button-small"
                :disabled="groupsStore.respondingEventId === ev.event_id"
                @click="onEventResponse(ev.event_id, 'not going')"
              >
                Not going
              </button>
            </div>
          </li>
        </ul>
      </section>
    </template>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { useGroupsStore } from "../stores/groups";
import { socialApi } from "../services/socialApi";

const route = useRoute();
const groupsStore = useGroupsStore();
const authStore = useAuthStore();

const groupId = computed(() => Number(route.params.id || 0));
const groupData = computed(() => groupsStore.selectedGroup?.group || null);
const isMember = computed(() => Boolean(groupsStore.selectedGroup?.is_member));
const isCreator = computed(
  () => groupData.value && authStore.userId === groupData.value.creator_id
);

const joinRequestMessage = ref("");
const joinRequestError = ref(false);

const inviteQuery = ref("");
const inviteResults = ref([]);
const inviteSearching = ref(false);
const inviteMessage = ref("");
const inviteError = ref(false);
let inviteDebounceTimer = null;

const groupPostContent = ref("");
const groupPostFile = ref(null);

const eventTitle = ref("");
const eventDescription = ref("");
const eventDateTime = ref("");

const commentContentByPostId = ref({});
const commentFileByPostId = ref({});
const groupPostContentRef = ref(null);

const emojiList = ["😀", "😊", "👍", "❤️", "😂", "🎉", "🙏", "👋", "✨", "🔥"];

function insertGroupPostEmoji(emoji) {
  const ta = groupPostContentRef.value;
  if (ta) {
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const text = groupPostContent.value || "";
    groupPostContent.value = text.slice(0, start) + emoji + text.slice(end);
    nextTick(() => {
      ta.focus();
      ta.selectionStart = ta.selectionEnd = start + emoji.length;
    });
  } else {
    groupPostContent.value = (groupPostContent.value || "") + emoji;
  }
}

function appendGroupCommentEmoji(groupPostId, emoji) {
  const current = commentContentByPostId.value[groupPostId] || "";
  commentContentByPostId.value = {
    ...commentContentByPostId.value,
    [groupPostId]: current + emoji
  };
}

function getGroupComments(groupPostId) {
  return groupsStore.commentsByGroupPostId[groupPostId] || [];
}

function commentsLoading(groupPostId) {
  return groupsStore.loadingCommentsByGroupPostId[groupPostId] === true;
}

function isCreatingComment(groupPostId) {
  return groupsStore.creatingCommentPostId === groupPostId;
}

function onGroupCommentFileSelect(groupPostId, event) {
  const file = event.target?.files?.[0];
  commentFileByPostId.value = { ...commentFileByPostId.value, [groupPostId]: file || null };
}

async function onSubmitGroupComment(groupPostId) {
  const text = (commentContentByPostId.value[groupPostId] || "").trim();
  if (!text) return;

  groupsStore.clearError();
  const payload = { group_post_id: groupPostId, content: text };

  const file = commentFileByPostId.value[groupPostId];
  if (file) {
    try {
      const uploadResult = await socialApi.uploadPostImage(file);
      const imageUrl = uploadResult?.imageUrl ?? uploadResult?.image_url;
      if (imageUrl) payload.image_url = imageUrl;
    } catch (err) {
      groupsStore.error = {
        status: err?.status || 500,
        message: err?.message || "Image upload failed."
      };
      return;
    }
  }

  try {
    await groupsStore.createGroupComment(payload);
    commentContentByPostId.value = { ...commentContentByPostId.value, [groupPostId]: "" };
    commentFileByPostId.value = { ...commentFileByPostId.value, [groupPostId]: null };
    const input = document.querySelector(`input[data-comment-file="${groupPostId}"]`);
    if (input) input.value = "";
  } catch (_) {
    // error in store or banner
  }
}

function formatDate(createdAt) {
  if (!createdAt) return "";
  try {
    const d = new Date(createdAt);
    return d.toLocaleDateString(undefined, {
      dateStyle: "short",
      timeStyle: "short"
    });
  } catch {
    return String(createdAt);
  }
}

function formatEventDate(eventDatetime) {
  if (!eventDatetime) return "";
  try {
    const d = new Date(eventDatetime);
    return d.toLocaleDateString(undefined, {
      dateStyle: "medium",
      timeStyle: "short"
    });
  } catch {
    return String(eventDatetime);
  }
}

async function onSubmitEvent() {
  if (!eventTitle.value.trim() || !eventDescription.value.trim() || !eventDateTime.value) return;
  groupsStore.clearError();
  const dt = new Date(eventDateTime.value).toISOString();
  try {
    await groupsStore.createGroupEvent({
      group_id: groupId.value,
      title: eventTitle.value.trim(),
      description: eventDescription.value.trim(),
      event_datetime: dt
    });
    eventTitle.value = "";
    eventDescription.value = "";
    eventDateTime.value = "";
  } catch (_) {
    // error in store or banner
  }
}

async function onEventResponse(eventId, response) {
  groupsStore.clearError();
  try {
    await groupsStore.respondToEvent(eventId, response, groupId.value);
  } catch (_) {
    // error in store or banner
  }
}

function onGroupPostFileSelect(event) {
  groupPostFile.value = event.target?.files?.[0] || null;
}

async function onSubmitGroupPost() {
  const hasText = groupPostContent.value && groupPostContent.value.trim().length > 0;
  const hasImage = Boolean(groupPostFile.value);
  if (!hasText && !hasImage) {
    groupsStore.error = { status: 400, message: "Write something or add an image to post." };
    return;
  }
  groupsStore.clearError();
  const payload = {
    group_id: groupId.value,
    content: (groupPostContent.value && groupPostContent.value.trim()) || ""
  };
  if (groupPostFile.value) {
    try {
      const uploadResult = await socialApi.uploadPostImage(groupPostFile.value);
      const imageUrl = uploadResult?.imageUrl ?? uploadResult?.image_url;
      if (imageUrl) payload.image_url = imageUrl;
    } catch (err) {
      groupsStore.error = {
        status: err?.status || 500,
        message: err?.message || "Image upload failed."
      };
      return;
    }
  }
  try {
    await groupsStore.createGroupPost(payload);
    groupPostContent.value = "";
    groupPostFile.value = null;
    const input = document.getElementById("group-post-image");
    if (input) input.value = "";
  } catch (_) {
    // error in store or banner
  }
}

async function runInviteSearch() {
  const q = inviteQuery.value.trim();
  if (!q) {
    inviteResults.value = [];
    return;
  }
  inviteSearching.value = true;
  inviteMessage.value = "";
  try {
    const data = await socialApi.searchUsers(q, 20);
    const users = data?.users || [];
    const myId = authStore.userId;
    inviteResults.value = myId != null ? users.filter((u) => u.user_id !== myId) : users;
  } catch {
    inviteResults.value = [];
  } finally {
    inviteSearching.value = false;
  }
}

watch(inviteQuery, () => {
  if (inviteDebounceTimer) clearTimeout(inviteDebounceTimer);
  inviteDebounceTimer = setTimeout(runInviteSearch, 300);
});

async function onInvite(userId) {
  groupsStore.clearError();
  inviteMessage.value = "";
  inviteError.value = false;
  try {
    await groupsStore.inviteUser(groupId.value, userId);
    inviteMessage.value = "Invitation sent.";
  } catch (err) {
    inviteError.value = true;
    inviteMessage.value = err?.message || "Failed to send invitation.";
  }
}

async function onRequestToJoin() {
  groupsStore.clearError();
  joinRequestMessage.value = "";
  joinRequestError.value = false;
  try {
    await groupsStore.requestToJoin(groupId.value);
    joinRequestMessage.value = "Request sent. The group creator will review it.";
  } catch (err) {
    joinRequestError.value = true;
    joinRequestMessage.value = err?.message || "Failed to send request.";
  }
}

async function onRespondJoinRequest(requestId, response) {
  groupsStore.clearError();
  try {
    await groupsStore.respondToJoinRequest(requestId, response);
  } catch {
    // error in store
  }
}

const load = async () => {
  if (!groupId.value) {
    groupsStore.error = { status: 400, message: "Invalid group id" };
    return;
  }
  await groupsStore.fetchGroupDetails(groupId.value).catch(() => {});
  if (groupsStore.selectedGroup?.is_member && groupData.value?.creator_id === authStore.userId) {
    await groupsStore.fetchJoinRequests(groupId.value).catch(() => {});
  } else {
    groupsStore.joinRequests = [];
  }
  if (groupsStore.selectedGroup?.is_member) {
    await groupsStore.fetchGroupPosts(groupId.value).catch(() => {});
    await Promise.all([
      ...groupsStore.groupPosts.map((p) => groupsStore.fetchGroupComments(p.group_post_id)),
      groupsStore.fetchGroupEvents(groupId.value)
    ]);
  } else {
    groupsStore.groupPosts = [];
    groupsStore.commentsByGroupPostId = {};
    groupsStore.groupEvents = [];
  }
};

onMounted(load);
watch(() => route.params.id, load);
</script>

<style scoped>
.join-request-section,
.join-requests-section,
.invite-section,
.group-posts-section,
.group-events-section {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}
.create-post-form,
.create-event-form {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
}
.group-events-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.group-event-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}
.event-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 6px;
}
.event-description {
  margin: 0 0 10px;
  white-space: pre-wrap;
  font-size: 0.95rem;
}
.event-actions {
  display: flex;
  gap: 8px;
}
.group-posts-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.group-post-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}
.group-post-meta {
  margin-bottom: 6px;
}
.group-post-author {
  font-weight: 600;
}
.group-post-content {
  margin: 0 0 8px;
  white-space: pre-wrap;
}
.group-post-image {
  max-width: 100%;
  max-height: 400px;
  border-radius: 8px;
  object-fit: contain;
  border: 1px solid var(--border);
}
.group-post-comments {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}
.comments-loading {
  margin: 0 0 8px;
  font-size: 0.9rem;
}
.comments-list {
  list-style: none;
  padding: 0;
  margin: 0 0 12px;
}
.comment-item {
  padding: 8px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.comment-item:last-child {
  border-bottom: none;
}
.comment-author {
  font-weight: 600;
  font-size: 0.9rem;
}
.comment-content {
  margin: 6px 0 0;
  white-space: pre-wrap;
  font-size: 0.95rem;
}
.comment-image {
  max-width: 100%;
  max-height: 200px;
  border-radius: 6px;
  object-fit: contain;
  margin-top: 6px;
  border: 1px solid var(--border);
}
.comment-form {
  margin-top: 8px;
}
.comment-input {
  width: 100%;
  margin-bottom: 8px;
  resize: vertical;
}
.comment-form-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}
.comment-file-label {
  cursor: pointer;
  font-size: 0.9rem;
}
.input-hidden {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
}
.section-heading,
.invite-heading {
  font-size: 1.1rem;
  margin: 0 0 12px;
}
.join-request-section .button {
  margin-bottom: 8px;
}
.join-requests-list {
  list-style: none;
  padding: 0;
  margin: 12px 0 0;
}
.join-request-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  padding: 8px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.requester-name {
  font-weight: 500;
}
.join-request-actions {
  display: flex;
  gap: 8px;
}
.invite-list {
  list-style: none;
  padding: 0;
  margin: 12px 0 0;
}
.invite-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.invite-user-name {
  font-weight: 500;
}
.button-small {
  padding: 6px 12px;
  font-size: 0.9rem;
}

.emoji-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 8px;
}
.emoji-btn {
  background: var(--surface-2, #f4f4f4);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 4px 8px;
  font-size: 1.25rem;
  cursor: pointer;
  line-height: 1.2;
}
.emoji-btn:hover {
  background: var(--border);
}
</style>
