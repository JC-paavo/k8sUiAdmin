<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="7" y="4" width="18" height="24" rx="3" :stroke="color" stroke-width="1.8" />
    <circle cx="16" cy="11" r="2.5" :stroke="color" stroke-width="1.5" />
    <circle cx="16" cy="19" r="2.5" :stroke="color" stroke-width="1.5" />
    <line x1="16" y1="13.5" x2="16" y2="16.5" :stroke="color" stroke-width="1.5" />
    <circle cx="16" cy="26" r="1.5" :stroke="color" stroke-width="1.3" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#409EFF' }
})
</script>
