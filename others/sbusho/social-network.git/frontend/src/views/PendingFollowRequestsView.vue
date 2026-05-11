<template>
  <div class="page-card">
    <h1>Pending follow requests</h1>

    <p v-if="followStore.loading" class="muted">Loading pending requests...</p>
    <p v-if="localError" class="error">{{ localError }}</p>

    <p v-if="!followStore.loading && followStore.pendingRequests.length === 0" class="muted">
      No pending follow requests.
    </p>

    <ul v-else class="pending-list">
      <PendingRequestItem
        v-for="request in followStore.pendingRequests"
        :key="request.user_id"
        :request="request"
        :loading="isRequestLoading(request.user_id)"
        @accept="onAccept"
        @decline="onDecline"
      />
    </ul>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import PendingRequestItem from "../components/PendingRequestItem.vue";
import { useFollowStore } from "../stores/follow";

const followStore = useFollowStore();
const localError = ref("");

const isRequestLoading = (userId) => {
  return (
    followStore.isActionPending(`accept:${userId}`) ||
    followStore.isActionPending(`decline:${userId}`)
  );
};

const load = async () => {
  localError.value = "";
  try {
    await followStore.fetchPendingRequests();
  } catch (error) {
    localError.value = error?.message || "Failed to load pending requests.";
  }
};

const onAccept = async (userId) => {
  localError.value = "";
  try {
    await followStore.acceptRequest(userId);
  } catch (error) {
    localError.value = error?.message || "Failed to accept request.";
  }
};

const onDecline = async (userId) => {
  localError.value = "";
  try {
    await followStore.declineRequest(userId);
  } catch (error) {
    localError.value = error?.message || "Failed to decline request.";
  }
};

onMounted(load);
</script>
