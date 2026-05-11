<template>
  <div class="page-card">
    <h1>Following</h1>
    <p v-if="followStore.loading" class="muted">Loading following...</p>
    <p v-if="localError" class="error">{{ localError }}</p>

    <ul v-if="!followStore.loading && following.length" class="user-list">
      <li v-for="user in following" :key="user.user_id" class="user-list-item">
        <RouterLink :to="`/profile/${user.user_id}`" class="user-link">
          <img
            v-if="user.avatar"
            :src="avatarSrc(user.avatar)"
            alt="User avatar"
            class="user-avatar"
          />
          <div v-else class="user-avatar user-avatar--placeholder">?</div>
          <span>{{ user.nickname || `User #${user.user_id}` }}</span>
        </RouterLink>
      </li>
    </ul>

    <p v-else-if="!followStore.loading" class="muted">No following users found.</p>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { useFollowStore } from "../stores/follow";

const route = useRoute();
const followStore = useFollowStore();
const localError = ref("");

const userId = computed(() => Number(route.params.id || 0));
const following = computed(() => followStore.following || []);

onMounted(async () => {
  localError.value = "";
  try {
    await followStore.fetchFollowing(userId.value);
  } catch (error) {
    localError.value = error?.message || "Failed to load following.";
  }
});

function avatarSrc(path) {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  return path.startsWith("/") ? path : `/${path}`;
}
</script>

<style scoped>
.user-list-item {
  align-items: center;
}

.user-link {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--border);
}

.user-avatar--placeholder {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  background: rgba(0, 0, 0, 0.08);
  color: var(--surface-text);
}
</style>
