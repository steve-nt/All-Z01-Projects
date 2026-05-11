<template>
  <button class="button" type="button" :disabled="loading" @click="onClick">
    {{ label }}
  </button>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  state: {
    type: String,
    default: "follow"
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(["follow", "unfollow"]);

const label = computed(() => {
  if (props.loading) {
    return "Working...";
  }
  if (props.state === "following") {
    return "Unfollow";
  }
  if (props.state === "requested") {
    return "Requested";
  }
  return "Follow";
});

const onClick = () => {
  if (props.loading || props.state === "requested") {
    return;
  }
  if (props.state === "following") {
    emit("unfollow");
    return;
  }
  emit("follow");
};
</script>
