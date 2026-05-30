<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <ellipse cx="16" cy="8" rx="9" ry="3.5" :stroke="color" stroke-width="1.8" />
    <path d="M7 8V24C7 25.93 11.03 27.5 16 27.5C20.97 27.5 25 25.93 25 24V8" :stroke="color" stroke-width="1.8" />
    <line x1="7" y1="14" x2="25" y2="14" :stroke="color" stroke-width="1.5" />
    <line x1="7" y1="20" x2="25" y2="20" :stroke="color" stroke-width="1.5" />
    <circle cx="16" cy="10" r="1.5" :stroke="color" stroke-width="1.2" />
    <circle cx="16" cy="16" r="1.5" :stroke="color" stroke-width="1.2" />
    <circle cx="16" cy="22" r="1.5" :stroke="color" stroke-width="1.2" />
  </svg>
</template>

<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#67C23A' }
})
</script>
