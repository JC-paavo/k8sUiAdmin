<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="2" y="4" width="10" height="8" rx="1.5" :stroke="color" stroke-width="1.6" />
    <rect x="11" y="4" width="10" height="8" rx="1.5" :stroke="color" stroke-width="1.6" />
    <rect x="20" y="4" width="10" height="8" rx="1.5" :stroke="color" stroke-width="1.6" />
    <line x1="7" y1="12" x2="7" y2="16" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="16" y1="12" x2="16" y2="16" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="25" y1="12" x2="25" y2="16" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <rect x="4" y="16" width="6" height="10" rx="1.5" :stroke="color" stroke-width="1.6" />
    <rect x="13" y="16" width="6" height="10" rx="1.5" :stroke="color" stroke-width="1.6" />
    <rect x="22" y="16" width="6" height="10" rx="1.5" :stroke="color" stroke-width="1.6" />
    <line x1="7" y1="20" x2="7" y2="23" :stroke="color" stroke-width="1.2" />
    <line x1="16" y1="20" x2="16" y2="23" :stroke="color" stroke-width="1.2" />
    <line x1="25" y1="20" x2="25" y2="23" :stroke="color" stroke-width="1.2" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#409EFF' }
})
</script>
