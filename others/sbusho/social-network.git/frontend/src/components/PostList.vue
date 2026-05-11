<template>
  <ul v-if="posts.length > 0" class="feed-list">
    <li v-for="post in posts" :key="post.post_id" class="feed-item">
      <div class="feed-meta">
        <span class="feed-author">{{ displayName(post.author, post.user_id) }}</span>
        <span class="muted"> · {{ formatDate(post.created_at) }} · {{ post.privacy }}</span>
      </div>
      <p class="feed-content">{{ post.content }}</p>
      <img
        v-if="post.image_url"
        :src="post.image_url"
        alt="Post attachment"
        class="feed-image"
      />
      <section class="comments-section">
        <p
          v-if="commentsLoading(post.post_id)"
          class="muted comments-loading"
        >
          Loading comments...
        </p>
        <template v-else>
          <ul class="comments-list">
            <li
              v-for="c in getComments(post.post_id)"
              :key="c.comment_id"
              class="comment-item"
            >
              <div class="comment-bubble">
                <span class="comment-meta">
                  <span class="comment-author">{{ displayName(c.author, c.user_id) }}</span>
                  <span class="muted"> · {{ formatDate(c.created_at) }}</span>
                </span>
                <p class="comment-content">{{ c.content }}</p>
                <img
                  v-if="c.image_url"
                  :src="c.image_url"
                  alt="Comment attachment"
                  class="comment-image"
                />
              </div>
            </li>
          </ul>
          <form
            class="comment-form"
            @submit.prevent="onSubmitComment(post.post_id)"
          >
            <div class="emoji-bar">
              <button
                v-for="emoji in emojiList"
                :key="emoji"
                type="button"
                class="emoji-btn"
                :title="'Add ' + emoji"
                @click="appendCommentEmoji(post.post_id, emoji)"
              >
                {{ emoji }}
              </button>
            </div>
            <textarea
              v-model.trim="commentContentByPostId[post.post_id]"
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
                  :data-comment-file="post.post_id"
                  @change="(e) => onCommentFileSelect(post.post_id, e)"
                />
                <span class="muted">Image (optional)</span>
              </label>
              <button
                type="submit"
                class="button button-small"
                :disabled="isCreatingComment(post.post_id) || !commentContentByPostId[post.post_id]"
              >
                {{ isCreatingComment(post.post_id) ? "Sending..." : "Comment" }}
              </button>
            </div>
          </form>
        </template>
      </section>
    </li>
  </ul>
</template>

<script setup>
import { computed, ref } from "vue";
import { usePostsStore } from "../stores/posts";
import { socialApi } from "../services/socialApi";

const props = defineProps({
  posts: {
    type: Array,
    default: () => []
  }
});

const postsStore = usePostsStore();
const commentContentByPostId = ref({});
const commentFileByPostId = ref({});

const posts = computed(() => props.posts);

function getComments(postId) {
  return postsStore.commentsByPostId[postId] || [];
}

function commentsLoading(postId) {
  return postsStore.loadingCommentsByPostId[postId] === true;
}

function isCreatingComment(postId) {
  return postsStore.creatingCommentPostId === postId;
}

function onCommentFileSelect(postId, event) {
  const file = event.target?.files?.[0];
  commentFileByPostId.value = { ...commentFileByPostId.value, [postId]: file || null };
}

async function onSubmitComment(postId) {
  const text = commentContentByPostId.value[postId];
  if (!text || !text.trim()) return;

  postsStore.clearError();
  const payload = { post_id: postId, content: text.trim() };

  const file = commentFileByPostId.value[postId];
  if (file) {
    try {
      const uploadResult = await socialApi.uploadPostImage(file);
      const imageUrl = uploadResult?.imageUrl ?? uploadResult?.image_url;
      if (imageUrl) payload.image_url = imageUrl;
    } catch (err) {
      postsStore.error = {
        status: err?.status || 500,
        message: err?.message || "Image upload failed."
      };
      return;
    }
  }

  try {
    await postsStore.createComment(payload);
    commentContentByPostId.value = { ...commentContentByPostId.value, [postId]: "" };
    commentFileByPostId.value = { ...commentFileByPostId.value, [postId]: null };
    const input = document.querySelector(`input[data-comment-file="${postId}"]`);
    if (input) input.value = "";
  } catch (_) {}
}

const emojiList = ["😀", "😊", "👍", "❤️", "😂", "🎉", "🙏", "👋", "✨", "🔥"];

function appendCommentEmoji(postId, emoji) {
  const current = commentContentByPostId.value[postId] || "";
  commentContentByPostId.value = {
    ...commentContentByPostId.value,
    [postId]: current + emoji
  };
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

/** Show username (nickname) when possible; never show full email – use part before @ or "User #id". */
function displayName(name, userId) {
  if (!name || !String(name).trim()) return "User #" + (userId ?? "");
  const n = String(name).trim();
  if (n.includes("@")) return n.split("@")[0];
  return n;
}
</script>

<style scoped>
.feed-list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0;
}
.feed-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}
.feed-meta {
  margin-bottom: 6px;
}
.feed-author {
  font-weight: 600;
}
.feed-content {
  margin: 0 0 8px;
  white-space: pre-wrap;
}
.feed-image {
  max-width: 100%;
  max-height: 400px;
  border-radius: 8px;
  object-fit: contain;
  border: 1px solid var(--border);
}

.comments-section {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}
.comments-section::before {
  content: "Comments";
  display: block;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--muted);
  margin-bottom: 10px;
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
  padding: 6px 0;
}
.comment-item:last-child {
  padding-bottom: 0;
}
.comment-bubble {
  background: var(--surface-2, rgba(0, 0, 0, 0.04));
  border-left: 3px solid var(--border-strong, rgba(0, 0, 0, 0.2));
  border-radius: 0 8px 8px 0;
  padding: 10px 12px;
  margin-left: 0;
}
.comment-meta {
  display: block;
  margin-bottom: 4px;
}
.comment-author {
  font-weight: 600;
  font-size: 0.875rem;
}
.comment-meta .muted {
  font-size: 0.8rem;
}
.comment-content {
  margin: 0;
  white-space: pre-wrap;
  font-size: 0.9rem;
  line-height: 1.4;
}
.comment-image {
  max-width: 100%;
  max-height: 200px;
  border-radius: 6px;
  object-fit: contain;
  margin-top: 8px;
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
.button-small {
  padding: 6px 12px;
  font-size: 0.9rem;
}

.emoji-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 6px;
}
.emoji-btn {
  background: var(--surface-2, rgba(0, 0, 0, 0.04));
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 2px 6px;
  font-size: 1.1rem;
  cursor: pointer;
  line-height: 1.2;
}
.emoji-btn:hover {
  background: var(--border);
}
</style>
