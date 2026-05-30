<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="10" y="4" width="12" height="8" rx="2" :stroke="color" stroke-width="1.8" />
    <path d="M12 10L7 16L10 16L10 28L14 28L14 21C14 19.9 14.9 19 16 19C17.1 19 18 19.9 18 21L18 28L22 28L22 16L25 16L20 10L16 5L12 10Z" :stroke="color" stroke-width="1.8" stroke-linejoin="round" />
    <path d="M16 12V16" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <circle cx="16" cy="17" r="1" :stroke="color" stroke-width="1.2" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#909399' }
})
</script>
