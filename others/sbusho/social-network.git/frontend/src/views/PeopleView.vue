<template>
  <div class="page-card">
    <h1>Find People</h1>

    <div class="field">
      <label class="label" for="people-search">Search by nickname</label>
      <input
        id="people-search"
        v-model.trim="query"
        class="input"
        type="search"
        placeholder="Type a name"
      />
      <p class="muted">Results update automatically.</p>
    </div>

    <p v-if="peopleStore.loading" class="muted">Searching users...</p>
    <p v-else-if="peopleStore.error" class="error">{{ peopleStore.error.message }}</p>
    <p v-if="followError" class="error">{{ followError }}</p>
    <p v-else-if="peopleStore.users.length === 0" class="muted">No users found.</p>

    <ul v-else class="user-list">
      <li v-for="user in peopleStore.users" :key="user.user_id" class="people-item">
        <RouterLink
          :to="`/profile/${user.user_id}`"
          class="people-summary"
        >
          <img
            v-if="user.avatar"
            :src="user.avatar"
            alt="User avatar"
            class="people-avatar"
          />
          <div v-else class="people-avatar people-avatar-placeholder">?</div>
          <span>{{ displayName(user.user_name, user.user_id) }}</span>
        </RouterLink>

        <button
          class="button"
          type="button"
          :disabled="
            followStore.isActionPending(`request:${user.user_id}`) ||
            Boolean(followStates[user.user_id])
          "
          @click="onFollow(user.user_id)"
        >
          {{
            followStates[user.user_id] === "requested"
              ? "Requested"
              : followStates[user.user_id] === "following"
                ? "Following"
                : followStore.isActionPending(`request:${user.user_id}`)
                  ? "Following..."
                  : "Follow"
          }}
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import { RouterLink } from "vue-router";
import { useFollowStore } from "../stores/follow";
import { usePeopleStore } from "../stores/people";

const peopleStore = usePeopleStore();
const followStore = useFollowStore();
const query = ref("");
const followError = ref("");
const followStates = ref({});

/** Show username when possible; if name looks like email use part before @, else "User #id". */
function displayName(name, userId) {
  if (!name || !String(name).trim()) return "User #" + (userId ?? "");
  const n = String(name).trim();
  if (n.includes("@")) return n.split("@")[0];
  return n;
}

let debounceTimer = null;

const runSearch = async () => {
  await peopleStore.search(query.value, 20).catch(() => {});
};

watch(query, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer);
  }
  debounceTimer = setTimeout(() => {
    runSearch();
  }, 300);
});

onMounted(runSearch);

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer);
  }
});

const onFollow = async (userId) => {
  followError.value = "";
  const parsedUserId = Number(userId);
  try {
    const response = await followStore.requestFollow(parsedUserId);
    followStates.value = {
      ...followStates.value,
      [parsedUserId]: response?.status === "accepted" ? "following" : "requested"
    };
  } catch (error) {
    followError.value =
      error?.message || followStore.error?.message || "Failed to send follow request.";
  }
};
</script>

<style scoped>
.people-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}

.people-summary {
  display: flex;
  align-items: center;
  gap: 10px;
}

.people-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--border);
}

.people-avatar-placeholder {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: var(--text);
  background: rgba(255, 255, 255, 0.08);
}
</style>
