<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M4 8C4 6.9 4.9 6 6 6H12L14.5 8.5H26C27.1 8.5 28 9.4 28 10.5V24C28 25.1 27.1 26 26 26H6C4.9 26 4 25.1 4 24V8Z" :stroke="color" stroke-width="1.8" stroke-linejoin="round" />
    <path d="M4 13H28" :stroke="color" stroke-width="1.5" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#909399' }
})
</script>
