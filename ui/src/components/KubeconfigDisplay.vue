<script setup lang="ts">
const props = defineProps<{
  yaml: string | null
}>()
</script>

<template>
  <div v-if="yaml" class="w-full min-h-5/6 border-2 border-gray-600 rounded-md p-4">
    <pre>{{ yaml }}</pre>
  </div>
</template>
