<template>
  <div class="page-card">
    <h1>Group invitations</h1>
    <p class="muted">Accept or decline invitations to join groups.</p>

    <p v-if="groupsStore.loadingInvitations" class="muted">Loading invitations...</p>
    <p v-else-if="groupsStore.error" class="error">{{ groupsStore.error.message }}</p>
    <p v-else-if="groupsStore.groupInvitations.length === 0" class="muted">
      You have no pending group invitations.
    </p>

    <ul v-else class="invitations-list">
      <li
        v-for="inv in groupsStore.groupInvitations"
        :key="inv.invitation_id"
        class="invitation-item"
      >
        <div class="invitation-info">
          <strong>{{ inv.group_name }}</strong>
          <span class="muted"> — {{ inv.inviter_name }} invited you</span>
        </div>
        <div class="invitation-actions">
          <button
            type="button"
            class="button"
            :disabled="respondingId === inv.invitation_id"
            @click="onRespond(inv.invitation_id, 'accepted')"
          >
            {{ respondingId === inv.invitation_id ? "..." : "Accept" }}
          </button>
          <button
            type="button"
            class="button secondary"
            :disabled="respondingId === inv.invitation_id"
            @click="onRespond(inv.invitation_id, 'declined')"
          >
            Decline
          </button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useGroupsStore } from "../stores/groups";

const groupsStore = useGroupsStore();
const router = useRouter();
const respondingId = ref(null);

async function onRespond(invitationId, response) {
  groupsStore.clearError();
  const inv = groupsStore.groupInvitations.find(
    (i) => i.invitation_id === invitationId
  );
  const groupId = inv?.group_id;
  respondingId.value = invitationId;
  try {
    await groupsStore.respondToInvitation(invitationId, response);
    if (response === "accepted" && groupId) {
      router.push(`/groups/${groupId}`);
    }
  } catch (err) {
    // error in store
  } finally {
    respondingId.value = null;
  }
}

onMounted(() => {
  groupsStore.fetchGroupInvitations().catch(() => {});
});
</script>

<style scoped>
.invitations-list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0;
}
.invitation-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}
.invitation-info {
  flex: 1;
  min-width: 0;
}
.invitation-actions {
  display: flex;
  gap: 8px;
}
</style>
