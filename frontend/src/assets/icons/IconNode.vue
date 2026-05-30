<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="6" y="3" width="20" height="14" rx="2" :stroke="color" stroke-width="1.8" />
    <rect x="9" y="6" width="6" height="4" rx="0.5" :stroke="color" stroke-width="1.3" />
    <rect x="17" y="6" width="6" height="4" rx="0.5" :stroke="color" stroke-width="1.3" />
    <rect x="9" y="11.5" width="14" height="3" rx="0.5" :stroke="color" stroke-width="1.3" />
    <line x1="12" y1="20" x2="12" y2="24" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="20" y1="20" x2="20" y2="24" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="8" y1="24" x2="24" y2="24" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="8" y1="27" x2="24" y2="27" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="16" y1="24" x2="16" y2="27" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#67C23A' }
})
</script>
