<script setup lang="ts">
import { ref } from 'vue'
import YAML from 'yaml'
import { copyToClipboard } from '../utils/clipboard'
import type { Kubeconfig } from '../types/Kubeconfig'

const props = defineProps<{
  kubeconfigs: Kubeconfig[]
}>()

const selectedKubeconfig = ref<string | null>(null)

function selectKubeconfig(kubeconfig: object) {
  selectedKubeconfig.value = YAML.stringify(kubeconfig)
}

function copyYaml() {
  if (selectedKubeconfig.value) {
    copyToClipboard(selectedKubeconfig.value)
  }
}
</script>

<template>
  <div class="flex flex-col space-y-8 my-8 max-w-max">
    <div
      v-for="kubeconfig in kubeconfigs"
      :key="kubeconfig.name"
      class="text-lg py-8 px-16 rounded-md bg-gray-600 border-2 border-gray-600"
      @click="selectKubeconfig(kubeconfig.kubeconfig)"
    >
      {{ kubeconfig.name }}
    </div>

    <!-- <div v-if="selectedKubeconfig" class="kubeconfig-box">
      <pre>{{ selectedKubeconfig }}</pre>
      <button class="copy-button" @click="copyYaml">Copy</button>
    </div> -->
  </div>
</template>
