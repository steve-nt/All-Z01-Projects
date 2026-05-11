<template>
  <label class="privacy-toggle">
    <input
      :checked="isPublic"
      :disabled="loading"
      type="checkbox"
      @change="onChange"
    />
    <span>{{ isPublic ? "Public profile" : "Private profile" }}</span>
  </label>
</template>

<script setup>
const props = defineProps({
  isPublic: {
    type: Boolean,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(["change"]);

const onChange = (event) => {
  emit("change", event.target.checked);
};
</script>
