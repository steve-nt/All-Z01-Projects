<template>
  <div class="page-card auth-card">
    <h1>Login</h1>

    <p v-if="successMessage" class="success">
      {{ successMessage }}
    </p>

    <p v-if="formError" class="error">
      {{ formError }}
    </p>

    <form class="form" @submit.prevent="onSubmit">
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
        <label class="label" for="password">Password</label>
        <input
          id="password"
          v-model="password"
          class="input"
          type="password"
          autocomplete="current-password"
          required
        />
      </div>

      <div class="actions">
        <button class="button" type="submit" :disabled="isSubmitting">
          {{ isSubmitting ? "Logging in..." : "Login" }}
        </button>
        <RouterLink class="muted" to="/register">Need an account?</RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { apiRequest } from "../services/apiClient";
import { useAuthStore } from "../stores/auth";
import { useNotificationsStore } from "../stores/notifications";
import { wsService } from "../services/websocket";

const authStore = useAuthStore();
const notificationsStore = useNotificationsStore();
const route = useRoute();
const router = useRouter();

const email = ref("");
const password = ref("");
const isSubmitting = ref(false);
const formError = ref("");

const successMessage = computed(() => {
  if (route.query.registered === "1") {
    return "Registration successful. Please log in.";
  }
  if (typeof route.query.message === "string" && route.query.message.trim()) {
    return route.query.message.trim();
  }
  return "";
});

const onSubmit = async () => {
  formError.value = "";
  isSubmitting.value = true;
  try {
    const body = new FormData();
    body.set("email", email.value);
    body.set("password", password.value);

    await apiRequest("/login", {
      method: "POST",
      body,
      headers: { Accept: "application/json" }
    });

    await authStore.checkSession();
    if (!authStore.loggedIn) {
      throw new Error("Login failed");
    }
    wsService.connect();
    await notificationsStore.fetchNotifications().catch(() => {});

    const redirect =
      typeof route.query.redirect === "string" ? route.query.redirect : "";
    await router.push(redirect || { name: "feed" });
  } catch (error) {
    formError.value = error?.message || "Login failed";
  } finally {
    isSubmitting.value = false;
  }
};
</script>
