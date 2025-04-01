<script setup lang="ts">
import { ref } from 'vue'
import YAML from 'yaml'
import type { Kubeconfig } from '../types/Kubeconfig'

defineProps<{
  kubeconfigs: Kubeconfig[]
}>()

const emit = defineEmits<{
  (e: 'kubeconfig-selected', kubeconfig: string): void
}>()

const selectedKubeconfig = ref<number | null>(null)

function selectKubeconfig(kubeconfig: object, index: number) {
  selectedKubeconfig.value = index
  emit('kubeconfig-selected', YAML.stringify(kubeconfig))
}
</script>

<template>
  <div
    class="flex-none flex flex-col gap-4 max-w-max min-w-min"
  >
    <button
      v-for="(kubeconfig, index) in kubeconfigs"
      :key="index"
      class="text-lg py-6 px-12 rounded-md border-2 border-gray-600 cursor-pointer whitespace-nowrap"
      :class="{
        'bg-accent text-primary-950': selectedKubeconfig === index,
        'bg-gray-700': selectedKubeconfig !== index,
      }"
      @click="selectKubeconfig(kubeconfig.kubeconfig, index)"
    >
      {{ kubeconfig.name }}
    </button>
  </div>
</template>
