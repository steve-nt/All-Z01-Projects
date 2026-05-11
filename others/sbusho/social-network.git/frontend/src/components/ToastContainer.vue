<template>
  <div class="toast-container" aria-live="polite">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toastStore.toasts"
        :key="toast.id"
        class="toast"
        :class="`toast--${toast.type}`"
        role="alert"
      >
        <span class="toast-message">{{ toast.message }}</span>
        <button
          type="button"
          class="toast-dismiss"
          aria-label="Dismiss"
          @click="dismiss(toast)"
        >
          ×
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { watch, onUnmounted } from "vue";
import { useToastStore } from "../stores/toast";

const toastStore = useToastStore();
const timeouts = new Map();

function dismiss(toast) {
  const t = timeouts.get(toast.id);
  if (t) clearTimeout(t);
  timeouts.delete(toast.id);
  toastStore.remove(toast.id);
}

function scheduleDismiss(toast) {
  if (timeouts.has(toast.id)) return;
  const id = setTimeout(() => {
    timeouts.delete(toast.id);
    toastStore.remove(toast.id);
  }, toast.durationMs);
  timeouts.set(toast.id, id);
}

watch(
  () => toastStore.toasts,
  (toasts) => {
    toasts.forEach(scheduleDismiss);
  },
  { deep: true }
);

onUnmounted(() => {
  timeouts.forEach((id) => clearTimeout(id));
  timeouts.clear();
});
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-width: min(380px, calc(100vw - 32px));
  pointer-events: none;
}
.toast-container > * {
  pointer-events: auto;
}

.toast {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 8px;
  box-shadow: var(--shadow);
  background: var(--surface);
  color: var(--text);
  border: 1px solid var(--border);
  border-left: 4px solid #000;
  font-size: 14px;
}
.toast--notification {
  border-left-style: solid;
}
.toast--message {
  border-left-style: dashed;
}
.toast--info {
  border-left-style: solid;
}
.toast--success {
  border-left-style: double;
}
.toast--error {
  border-left-style: dotted;
}

.toast-message {
  flex: 1;
  min-width: 0;
}
.toast-dismiss {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  padding: 0;
  border: none;
  background: transparent;
  color: inherit;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  opacity: 0.85;
  border-radius: 4px;
}
.toast-dismiss:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.06);
}

.toast-enter-active,
.toast-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
.toast-move {
  transition: transform 0.2s ease;
}
</style>
