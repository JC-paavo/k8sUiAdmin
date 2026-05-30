<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <circle cx="14" cy="16" r="5" :stroke="color" stroke-width="1.8" />
    <path d="M18 19L26 27" :stroke="color" stroke-width="1.8" stroke-linecap="round" />
    <line x1="22" y1="23" x2="25" y2="26" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <path d="M18 14C18 12 18 10 22 10C24 10 24 11.5 24 13" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <rect x="22" y="13" width="5" height="4" rx="0.5" :stroke="color" stroke-width="1.3" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#E6A23C' }
})
</script>
