<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="5" y="5" width="22" height="22" rx="3" :stroke="color" stroke-width="1.8" />
    <circle cx="16" cy="12" r="3" :stroke="color" stroke-width="1.8" />
    <circle cx="11" cy="21" r="2" :stroke="color" stroke-width="1.8" />
    <circle cx="21" cy="21" r="2" :stroke="color" stroke-width="1.8" />
    <line x1="14" y1="14.5" x2="12" y2="19" :stroke="color" stroke-width="1.5" />
    <line x1="18" y1="14.5" x2="20" y2="19" :stroke="color" stroke-width="1.5" />
    <path d="M10 5V3" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <path d="M16 5V3" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <path d="M22 5V3" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#409EFF' }
})
</script>
