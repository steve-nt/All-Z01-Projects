<template>
  <form class="form create-post-form" @submit.prevent="onSubmit">
    <div class="field">
      <label class="label" for="post-content">What's on your mind?</label>
      <div class="emoji-bar">
        <button
          v-for="emoji in emojiList"
          :key="emoji"
          type="button"
          class="emoji-btn"
          :title="'Add ' + emoji"
          @click="insertEmoji(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
      <textarea
        id="post-content"
        ref="postContentRef"
        v-model="content"
        class="input textarea"
        rows="3"
        placeholder="Write a post or add an image below (or both)"
      />
    </div>
    <div class="field">
      <label class="label" for="post-privacy">Privacy</label>
      <select id="post-privacy" v-model="privacy" class="input">
        <option value="public">Public</option>
        <option value="almost_private">Followers only</option>
        <option value="private">Private (select people)</option>
      </select>
    </div>
    <div v-if="privacy === 'private'" class="field">
      <label class="label">Visible to (accepted followers)</label>
      <p v-if="loadingFollowers" class="muted">Loading followers...</p>
      <div v-else-if="followersList.length === 0" class="muted">
        You have no followers yet. Choose "Public" or "Followers only" instead.
      </div>
      <div v-else class="visibility-list">
        <label
          v-for="f in followersList"
          :key="f.user_id"
          class="checkbox"
          :title="f.nickname || ('User #' + f.user_id)"
        >
          <input
            v-model="visibleTo"
            type="checkbox"
            :value="f.user_id"
          />
          {{ f.nickname || "User #" + f.user_id }}
        </label>
      </div>
    </div>
    <div class="field">
      <label class="label" for="post-image">Image (optional)</label>
      <input
        id="post-image"
        type="file"
        accept="image/*"
        class="input"
        @change="onFileSelect"
      />
    </div>
    <div class="actions">
      <button class="button" type="submit" :disabled="creating">
        {{ creating ? "Posting..." : "Post" }}
      </button>
    </div>
  </form>
</template>

<script setup>
import { nextTick, ref, watch } from "vue";
import { useAuthStore } from "../stores/auth";
import { usePostsStore } from "../stores/posts";
import { socialApi } from "../services/socialApi";

const props = defineProps({
  /** When true, form is disabled (e.g. not own profile) */
  disabled: { type: Boolean, default: false }
});

const emit = defineEmits(["submitted"]);

const postsStore = usePostsStore();
const authStore = useAuthStore();

const content = ref("");
const privacy = ref("public");
const visibleTo = ref([]);
const selectedFile = ref(null);
const followersList = ref([]);
const loadingFollowers = ref(false);

const creating = ref(false);
const postContentRef = ref(null);

const emojiList = ["😀", "😊", "👍", "❤️", "😂", "🎉", "🙏", "👋", "✨", "🔥"];

function insertEmoji(emoji) {
  const ta = postContentRef.value;
  if (ta) {
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const text = content.value || "";
    content.value = text.slice(0, start) + emoji + text.slice(end);
    nextTick(() => {
      ta.focus();
      ta.selectionStart = ta.selectionEnd = start + emoji.length;
    });
  } else {
    content.value = (content.value || "") + emoji;
  }
}

async function loadFollowers() {
  const uid = authStore.userId;
  if (!uid) return;
  loadingFollowers.value = true;
  try {
    const data = await socialApi.getFollowers(uid);
    followersList.value = data?.followers || [];
  } catch {
    followersList.value = [];
  } finally {
    loadingFollowers.value = false;
  }
}

watch(privacy, (val) => {
  if (val === "private") {
    loadFollowers();
  } else {
    visibleTo.value = [];
  }
});

function onFileSelect(event) {
  const file = event.target?.files?.[0];
  selectedFile.value = file || null;
}

async function onSubmit() {
  const hasText = (content.value && content.value.trim().length > 0);
  const hasImage = Boolean(selectedFile.value);
  if (props.disabled || creating.value || (!hasText && !hasImage)) {
    if (!hasText && !hasImage) {
      postsStore.error = { status: 400, message: "Write something or add an image to post." };
    }
    return;
  }

  postsStore.clearError();
  const payload = {
    content: (content.value && content.value.trim()) || "",
    privacy: privacy.value
  };
  if (privacy.value === "private" && visibleTo.value.length > 0) {
    payload.visible_to = [...visibleTo.value];
  }
  if (privacy.value === "private" && visibleTo.value.length === 0) {
    postsStore.error = { status: 400, message: "Select at least one follower for private posts." };
    return;
  }

  if (selectedFile.value) {
    try {
      const uploadResult = await socialApi.uploadPostImage(selectedFile.value);
      const filename = uploadResult?.filename;
      if (filename) {
        payload.image_filename = filename;
      }
    } catch (err) {
      postsStore.error = {
        status: err?.status || 500,
        message: err?.message || "Image upload failed."
      };
      return;
    }
  }

  creating.value = true;
  try {
    await postsStore.createPost(payload);
    content.value = "";
    visibleTo.value = [];
    selectedFile.value = null;
    const input = document.getElementById("post-image");
    if (input) input.value = "";
    emit("submitted");
  } catch (_) {
    // error in store
  } finally {
    creating.value = false;
  }
}
</script>

<style scoped>
.create-post-form {
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border);
}
.visibility-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 200px;
  overflow-y: auto;
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
