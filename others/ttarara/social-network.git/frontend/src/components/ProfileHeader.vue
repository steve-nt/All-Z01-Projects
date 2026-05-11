<template>
  <div class="profile-header">
    <img
      class="avatar"
      :src="profile.avatar || fallbackAvatar"
      alt="Profile avatar"
    />
    <div>
      <h1 class="profile-name">{{ fullName }}</h1>
      <p class="muted">@{{ profile.nickname || "no-nickname" }}</p>
      <p v-if="profile.email" class="muted">{{ profile.email }}</p>
      <p v-if="ageDisplay" class="muted">{{ ageDisplay }}</p>
      <p v-if="relationshipLabel" class="muted">{{ relationshipLabel }}</p>
      <p v-if="profile.about_me" class="about">{{ profile.about_me }}</p>
      <p v-if="profile.hobbies" class="hobbies">{{ profile.hobbies }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  profile: {
    type: Object,
    required: true
  }
});

const fallbackAvatar =
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='96' height='96'%3E%3Crect width='96' height='96' fill='%23e5e7eb'/%3E%3Ctext x='48' y='52' text-anchor='middle' fill='%236b7280' font-size='12'%3EAvatar%3C/text%3E%3C/svg%3E";

const fullName = computed(() => {
  const first = props.profile.first_name || "";
  const last = props.profile.last_name || "";
  return `${first} ${last}`.trim() || "User";
});

const ageDisplay = computed(() => {
  const dob = props.profile.date_of_birth;
  if (!dob) return "";
  const birth = new Date(dob);
  if (Number.isNaN(birth.getTime())) return "";
  const today = new Date();
  let age = today.getFullYear() - birth.getFullYear();
  const m = today.getMonth() - birth.getMonth();
  if (m < 0 || (m === 0 && today.getDate() < birth.getDate())) age--;
  return age >= 0 && age < 150 ? `${age} years old` : "";
});

const relationshipLabels = {
  single: "Single",
  in_relationship: "In a relationship",
  married: "Married",
  engaged: "Engaged",
  complicated: "It's complicated",
  other: "Other"
};

const relationshipLabel = computed(() => {
  const s = props.profile.relationship_status;
  if (!s) return "";
  return relationshipLabels[s] || s;
});
</script>
