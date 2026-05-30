<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <circle cx="16" cy="16" r="4" :stroke="color" stroke-width="1.8" />
    <circle cx="5" cy="8" r="2.5" :stroke="color" stroke-width="1.5" />
    <circle cx="27" cy="8" r="2.5" :stroke="color" stroke-width="1.5" />
    <circle cx="5" cy="24" r="2.5" :stroke="color" stroke-width="1.5" />
    <circle cx="27" cy="24" r="2.5" :stroke="color" stroke-width="1.5" />
    <line x1="7.5" y1="9.5" x2="12.5" y2="13.5" :stroke="color" stroke-width="1.5" />
    <line x1="24.5" y1="9.5" x2="19.5" y2="13.5" :stroke="color" stroke-width="1.5" />
    <line x1="7.5" y1="22.5" x2="12.5" y2="18.5" :stroke="color" stroke-width="1.5" />
    <line x1="24.5" y1="22.5" x2="19.5" y2="18.5" :stroke="color" stroke-width="1.5" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#409EFF' }
})
</script>
