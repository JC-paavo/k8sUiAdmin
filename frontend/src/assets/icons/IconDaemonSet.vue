<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M16 4L5 10V16C5 22.63 9.93 28.75 16 31C22.07 28.75 27 22.63 27 16V10L16 4Z" :stroke="color" stroke-width="1.8" stroke-linejoin="round" />
    <circle cx="12" cy="15" r="2" :stroke="color" stroke-width="1.5" />
    <circle cx="16" cy="21" r="2" :stroke="color" stroke-width="1.5" />
    <circle cx="20" cy="15" r="2" :stroke="color" stroke-width="1.5" />
    <line x1="12" y1="17" x2="14.5" y2="19.5" :stroke="color" stroke-width="1.3" />
    <line x1="20" y1="17" x2="17.5" y2="19.5" :stroke="color" stroke-width="1.3" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#E6A23C' }
})
</script>
