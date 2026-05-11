<template>
  <div class="follow-preview">
    <RouterLink :to="`/profile/${profileId}/followers`">
      <strong>{{ followersCount }}</strong> followers
    </RouterLink>
    <RouterLink :to="`/profile/${profileId}/following`">
      <strong>{{ followingCount }}</strong> following
    </RouterLink>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";

defineProps({
  profileId: {
    type: Number,
    required: true
  },
  followersCount: {
    type: Number,
    default: 0
  },
  followingCount: {
    type: Number,
    default: 0
  }
});
</script>
