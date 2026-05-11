<template>
  <div class="page-card">
    <h1>Feed</h1>
    <p class="muted">Posts from everyone. Create posts from your profile.</p>

    <p v-if="postsStore.loading" class="muted">Loading feed...</p>
    <p v-else-if="postsStore.error" class="error">{{ postsStore.error.message }}</p>
    <p v-else-if="postsStore.posts.length === 0" class="muted">No posts yet. Be the first to post from your profile!</p>
    <PostList v-else :posts="postsStore.posts" />
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import PostList from "../components/PostList.vue";
import { usePostsStore } from "../stores/posts";

const postsStore = usePostsStore();

onMounted(async () => {
  await postsStore.fetchPosts().catch(() => {});
  await Promise.all(
    postsStore.posts.map((p) => postsStore.fetchComments(p.post_id))
  );
});
</script>

<style scoped>
/* List styles are in PostList.vue */
</style>
