<template>
  <div class="page-card">
    <h1>Groups</h1>

    <form class="form" @submit.prevent="onCreateGroup">
      <div class="field">
        <label class="label" for="group-name">Group name</label>
        <input
          id="group-name"
          v-model.trim="groupName"
          class="input"
          type="text"
          required
        />
      </div>
      <div class="field">
        <label class="label" for="group-description">Description</label>
        <textarea
          id="group-description"
          v-model.trim="description"
          class="input textarea"
          rows="3"
          required
        />
      </div>
      <div class="actions">
        <button class="button" type="submit" :disabled="isSubmitting">
          {{ isSubmitting ? "Creating..." : "Create group" }}
        </button>
      </div>
      <p v-if="formError" class="error">{{ formError }}</p>
    </form>

    <div class="field group-search">
      <label class="label" for="group-search">Search groups</label>
      <input
        id="group-search"
        v-model.trim="search"
        class="input"
        type="search"
        placeholder="Search by group name"
      />
      <p class="muted">Search is local on fetched results (backend has no search query endpoint).</p>
    </div>

    <p v-if="groupsStore.loading" class="muted">Loading groups...</p>
    <p v-else-if="groupsStore.error" class="error">{{ groupsStore.error.message }}</p>
    <p v-else-if="filteredGroups.length === 0" class="muted">No groups found.</p>

    <ul v-else class="group-list">
      <li v-for="group in filteredGroups" :key="group.group_id" class="group-list-item">
        <div>
          <RouterLink :to="`/groups/${group.group_id}`" class="group-title">
            {{ group.group_name }}
          </RouterLink>
          <p class="muted">{{ group.description }}</p>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import { useGroupsStore } from "../stores/groups";

const groupsStore = useGroupsStore();
const router = useRouter();
const search = ref("");
const groupName = ref("");
const description = ref("");
const formError = ref("");
const isSubmitting = ref(false);

const filteredGroups = computed(() => {
  const q = search.value.trim().toLowerCase();
  if (!q) {
    return groupsStore.groups;
  }
  return groupsStore.groups.filter((group) =>
    String(group.group_name || "").toLowerCase().includes(q)
  );
});

onMounted(async () => {
  await groupsStore.fetchGroups().catch(() => {});
});

const onCreateGroup = async () => {
  formError.value = "";

  if (!groupName.value) {
    formError.value = "Group name is required.";
    return;
  }
  if (!description.value) {
    formError.value = "Description is required.";
    return;
  }

  isSubmitting.value = true;
  try {
    const created = await groupsStore.createGroup({
      group_name: groupName.value,
      description: description.value
    });
    groupName.value = "";
    description.value = "";
    if (created?.group_id) {
      await router.push(`/groups/${created.group_id}`);
    }
  } catch (error) {
    formError.value = error?.message || "Failed to create group.";
  } finally {
    isSubmitting.value = false;
  }
};
</script>
