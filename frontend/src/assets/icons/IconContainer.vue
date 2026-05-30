<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="5" y="4" width="22" height="20" rx="3" :stroke="color" stroke-width="1.8" />
    <line x1="5" y1="12" x2="27" y2="12" :stroke="color" stroke-width="1.5" />
    <line x1="13" y1="4" x2="13" y2="24" :stroke="color" stroke-width="1.3" />
    <line x1="19" y1="4" x2="19" y2="24" :stroke="color" stroke-width="1.3" />
    <rect x="8" y="14" width="3" height="3" rx="0.5" :stroke="color" stroke-width="1.2" />
    <rect x="14" y="14" width="3" height="3" rx="0.5" :stroke="color" stroke-width="1.2" />
    <rect x="20" y="14" width="3" height="3" rx="0.5" :stroke="color" stroke-width="1.2" />
    <rect x="8" y="18.5" width="3" height="3" rx="0.5" :stroke="color" stroke-width="1.2" />
    <rect x="14" y="18.5" width="3" height="3" rx="0.5" :stroke="color" stroke-width="1.2" />
    <rect x="20" y="18.5" width="3" height="3" rx="0.5" :stroke="color" stroke-width="1.2" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#409EFF' }
})
</script>
