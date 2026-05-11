<template>
  <li class="pending-item">
    <div class="pending-user">
      <div class="pending-avatar-wrap">
        <img
          v-if="request.avatar"
          :src="request.avatar"
          alt="User avatar"
          class="pending-avatar"
        />
        <div v-else class="pending-avatar pending-avatar--placeholder">?</div>
      </div>
      <div class="pending-meta">
        <p class="pending-name">
          {{ request.nickname || `User #${request.user_id}` }}
        </p>
        <p class="muted">Requested at {{ request.created_at }}</p>
      </div>
    </div>
    <div class="actions">
      <button
        class="button"
        type="button"
        :disabled="loading"
        @click="$emit('accept', request.user_id)"
      >
        Accept
      </button>
      <button
        class="button secondary"
        type="button"
        :disabled="loading"
        @click="$emit('decline', request.user_id)"
      >
        Decline
      </button>
    </div>
  </li>
</template>

<script setup>
defineProps({
  request: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
});

defineEmits(["accept", "decline"]);
</script>

<style scoped>
.pending-user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pending-avatar-wrap {
  flex-shrink: 0;
}

.pending-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--border);
}

.pending-avatar--placeholder {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  background: rgba(0, 0, 0, 0.08);
  color: var(--surface-text);
}

.pending-name {
  margin: 0 0 2px;
  font-weight: 600;
}
</style>
