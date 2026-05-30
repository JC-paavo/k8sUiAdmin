<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M5 6C5 4.9 5.9 4 7 4H20L27 11V26C27 27.1 26.1 28 25 28H7C5.9 28 5 27.1 5 26V6Z" :stroke="color" stroke-width="1.8" stroke-linejoin="round" />
    <path d="M20 4V11H27" :stroke="color" stroke-width="1.8" stroke-linejoin="round" />
    <line x1="9" y1="14" x2="23" y2="14" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="9" y1="18" x2="23" y2="18" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="9" y1="22" x2="18" y2="22" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
    <line x1="13" y1="10" x2="17" y2="10" :stroke="color" stroke-width="1.5" stroke-linecap="round" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#F56C6C' }
})
</script>
