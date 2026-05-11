<template>
  <div class="profile-settings" :class="{ embedded }">
    <template v-if="!embedded">
      <h2>Profile settings</h2>
      <p class="muted">Update your profile information. Only you can see this section.</p>
    </template>

    <div class="settings-section">
      <label class="settings-label">Profile photo</label>
      <div class="avatar-upload">
        <img
          class="avatar-preview"
          :src="profile.avatar || fallbackAvatar"
          alt="Avatar"
        />
        <div class="avatar-upload-actions">
          <input
            ref="fileInputRef"
            type="file"
            accept="image/jpeg,image/png,image/gif"
            class="file-input"
            @change="onAvatarChange"
          />
          <button type="button" class="button secondary upload-image-btn" @click="triggerFileInput">
            {{ avatarUploading ? "Uploading…" : "Upload image" }}
          </button>
        </div>
      </div>
    </div>

    <form class="settings-form" @submit.prevent="saveProfile">
      <div class="settings-section">
        <label class="settings-label" for="nickname">Nickname</label>
        <input
          id="nickname"
          v-model="form.nickname"
          type="text"
          class="settings-input"
          placeholder="Display name"
        />
      </div>

      <div class="settings-section">
        <label class="settings-label" for="about">About me</label>
        <textarea
          id="about"
          v-model="form.about_me"
          class="settings-input settings-textarea"
          rows="3"
          placeholder="Tell others about yourself"
        />
      </div>

      <div class="settings-section">
        <label class="settings-label" for="dob">Date of birth</label>
        <input
          id="dob"
          v-model="form.date_of_birth"
          type="date"
          class="settings-input"
        />
        <span v-if="ageDisplay" class="muted age-display">{{ ageDisplay }}</span>
      </div>

      <div class="settings-section">
        <label class="settings-label" for="relationship">Relationship status</label>
        <select id="relationship" v-model="form.relationship_status" class="settings-input">
          <option value="">Prefer not to say</option>
          <option value="single">Single</option>
          <option value="in_relationship">In a relationship</option>
          <option value="married">Married</option>
          <option value="engaged">Engaged</option>
          <option value="complicated">It's complicated</option>
          <option value="other">Other</option>
        </select>
      </div>

      <div class="settings-section">
        <label class="settings-label" for="hobbies">Hobbies</label>
        <textarea
          id="hobbies"
          v-model="form.hobbies"
          class="settings-input settings-textarea"
          rows="2"
          placeholder="e.g. Reading, hiking, music"
        />
      </div>

      <p v-if="settingsError" class="error">{{ settingsError }}</p>
      <button type="submit" class="button primary" :disabled="saving">
        {{ saving ? "Saving…" : "Save profile" }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { useProfileStore } from "../stores/profile";

const props = defineProps({
  profile: {
    type: Object,
    required: true
  },
  /** When true, hide the main heading (e.g. when used inside a modal that has its own title) */
  embedded: {
    type: Boolean,
    default: false
  }
});

const profileStore = useProfileStore();
const fileInputRef = ref(null);
const avatarUploading = ref(false);
const saving = ref(false);
const settingsError = ref("");

const fallbackAvatar =
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='96' height='96'%3E%3Crect width='96' height='96' fill='%23e5e7eb'/%3E%3Ctext x='48' y='52' text-anchor='middle' fill='%236b7280' font-size='12'%3EAvatar%3C/text%3E%3C/svg%3E";

const form = ref({
  nickname: "",
  about_me: "",
  date_of_birth: "",
  relationship_status: "",
  hobbies: ""
});

function syncFormFromProfile() {
  if (!props.profile) return;
  form.value = {
    nickname: props.profile.nickname || "",
    about_me: props.profile.about_me || "",
    date_of_birth: props.profile.date_of_birth || "",
    relationship_status: props.profile.relationship_status || "",
    hobbies: props.profile.hobbies || ""
  };
}

watch(
  () => props.profile,
  (p) => {
    if (p) syncFormFromProfile();
  },
  { immediate: true }
);

const ageDisplay = computed(() => {
  const dob = form.value.date_of_birth || props.profile?.date_of_birth;
  if (!dob) return "";
  const birth = new Date(dob);
  if (Number.isNaN(birth.getTime())) return "";
  const today = new Date();
  let age = today.getFullYear() - birth.getFullYear();
  const m = today.getMonth() - birth.getMonth();
  if (m < 0 || (m === 0 && today.getDate() < birth.getDate())) age--;
  return age >= 0 && age < 150 ? `${age} years old` : "";
});

function triggerFileInput() {
  fileInputRef.value?.click();
}

async function onAvatarChange(event) {
  const file = event.target.files?.[0];
  if (!file) return;
  avatarUploading.value = true;
  settingsError.value = "";
  try {
    await profileStore.uploadAvatar(file);
  } catch (err) {
    settingsError.value = err?.message || "Failed to upload image.";
  } finally {
    avatarUploading.value = false;
    event.target.value = "";
  }
}

async function saveProfile() {
  saving.value = true;
  settingsError.value = "";
  const payload = {};
  if (form.value.nickname !== (props.profile?.nickname ?? "")) payload.nickname = form.value.nickname;
  if (form.value.about_me !== (props.profile?.about_me ?? "")) payload.about_me = form.value.about_me;
  if (form.value.date_of_birth !== (props.profile?.date_of_birth ?? "")) payload.date_of_birth = form.value.date_of_birth || "";
  if (form.value.relationship_status !== (props.profile?.relationship_status ?? "")) payload.relationship_status = form.value.relationship_status || "";
  if (form.value.hobbies !== (props.profile?.hobbies ?? "")) payload.hobbies = form.value.hobbies || "";

  if (Object.keys(payload).length === 0) {
    saving.value = false;
    return;
  }

  try {
    await profileStore.updateProfile(payload);
  } catch (err) {
    settingsError.value = err?.message || "Failed to save profile.";
  } finally {
    saving.value = false;
  }
}
</script>

<style scoped>
.profile-settings {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border);
}
.profile-settings.embedded {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
}
.profile-settings h2 {
  font-size: 1.25rem;
  margin-bottom: 0.25rem;
}
.settings-section {
  margin-bottom: 1rem;
}
.settings-label {
  display: block;
  font-weight: 500;
  margin-bottom: 0.25rem;
}
.settings-input {
  width: 100%;
  max-width: 24rem;
  padding: 0.5rem;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: var(--surface);
  color: var(--surface-text);
  outline: none;
}
.settings-input:focus {
  border-color: var(--border-strong);
  box-shadow: 0 0 0 3px var(--ring-surface);
}
.settings-textarea {
  resize: vertical;
  min-height: 4rem;
}
.avatar-upload {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.avatar-preview {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--border);
}
.avatar-upload-actions {
  position: relative;
  display: inline-block;
}

.file-input {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}

.upload-image-btn {
  position: relative;
  z-index: 1;
  /* Always visible: explicit colors so it doesn’t blend into the modal */
  background: var(--surface-2, #f4f4f4) !important;
  color: var(--surface-text, #0b0b0b) !important;
  border: 1px solid var(--border-strong, rgba(0, 0, 0, 0.2)) !important;
}

.upload-image-btn:hover:not(:disabled) {
  background: var(--border, rgba(0, 0, 0, 0.12)) !important;
}
.age-display {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.875rem;
}
.settings-form .button {
  margin-top: 0.5rem;
}
.error {
  margin: 0.5rem 0;
}
</style>
