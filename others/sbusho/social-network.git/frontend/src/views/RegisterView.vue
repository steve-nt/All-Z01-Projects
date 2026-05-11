<template>
  <div class="page-card auth-card auth-card--wide">
    <h1>Register</h1>

    <p v-if="formError" class="error">
      {{ formError }}
    </p>

    <form class="form" @submit.prevent="onSubmit">
      <div class="form-grid">
        <div class="field">
          <label class="label" for="first_name">First name</label>
          <input
            id="first_name"
            v-model.trim="firstName"
            class="input"
            autocomplete="given-name"
            required
          />
        </div>

        <div class="field">
          <label class="label" for="last_name">Last name</label>
          <input
            id="last_name"
            v-model.trim="lastName"
            class="input"
            autocomplete="family-name"
            required
          />
        </div>
      </div>

      <div class="field">
        <label class="label" for="email">Email</label>
        <input
          id="email"
          v-model.trim="email"
          class="input"
          type="email"
          autocomplete="email"
          required
        />
      </div>

      <div class="field">
        <label class="label" for="date_of_birth">Date of birth</label>
        <input
          id="date_of_birth"
          v-model="dateOfBirth"
          class="input"
          type="date"
          required
        />
      </div>

      <div class="form-grid">
        <div class="field">
          <label class="label" for="password">Password</label>
          <input
            id="password"
            v-model="password"
            class="input"
            type="password"
            autocomplete="new-password"
            required
          />
          <div class="hint muted">
            8+ chars, uppercase, lowercase, number, symbol
          </div>
        </div>

        <div class="field">
          <label class="label" for="confirmPassword">Confirm password</label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            class="input"
            type="password"
            autocomplete="new-password"
            required
          />
        </div>
      </div>

      <div class="form-grid">
        <div class="field">
          <label class="label" for="nickname">Nickname (optional)</label>
          <input
            id="nickname"
            v-model.trim="nickname"
            class="input"
            autocomplete="nickname"
          />
        </div>

        <div class="field checkbox-field">
          <label class="checkbox">
            <input v-model="isPublic" type="checkbox" />
            Public profile
          </label>
          <div class="hint muted">
            You can change this later in your profile settings.
          </div>
        </div>
      </div>

      <div class="field">
        <label class="label" for="about_me">About me (optional)</label>
        <textarea
          id="about_me"
          v-model.trim="aboutMe"
          class="input textarea"
          rows="4"
        />
      </div>

      <div class="actions">
        <button class="button" type="submit" :disabled="isSubmitting">
          {{ isSubmitting ? "Creating account..." : "Create account" }}
        </button>
        <RouterLink class="muted" to="/login">Already have an account?</RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import { apiRequest } from "../services/apiClient";

const router = useRouter();

const firstName = ref("");
const lastName = ref("");
const email = ref("");
const dateOfBirth = ref("");
const password = ref("");
const confirmPassword = ref("");
const nickname = ref("");
const aboutMe = ref("");
const isPublic = ref(true);

const isSubmitting = ref(false);
const formError = ref("");

const onSubmit = async () => {
  formError.value = "";

  if (password.value !== confirmPassword.value) {
    formError.value = "Passwords do not match";
    return;
  }

  isSubmitting.value = true;
  try {
    const body = new FormData();
    body.set("email", email.value);
    body.set("password", password.value);
    body.set("confirmPassword", confirmPassword.value);
    body.set("first_name", firstName.value);
    body.set("last_name", lastName.value);
    body.set("date_of_birth", dateOfBirth.value);
    body.set("nickname", nickname.value);
    body.set("about_me", aboutMe.value);
    body.set("is_public", isPublic.value ? "true" : "false");
    body.set("is_active", "true");

    await apiRequest("/register", {
      method: "POST",
      body,
      headers: { Accept: "application/json" }
    });

    await router.push({ name: "login", query: { registered: "1" } });
  } catch (error) {
    formError.value = error?.message || "Registration failed";
  } finally {
    isSubmitting.value = false;
  }
};
</script>
