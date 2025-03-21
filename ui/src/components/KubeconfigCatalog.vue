<script setup lang="ts">
import { ref } from 'vue'
import YAML from 'yaml'
import { copyToClipboard } from '../utils/clipboard'
import type { Kubeconfig } from '../types/Kubeconfig';

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
  <div>
    <div v-if="kubeconfigs.length" class="catalog">
      <div
        v-for="kubeconfig in kubeconfigs"
        :key="kubeconfig.name"
        class="catalog-item"
        @click="selectKubeconfig(kubeconfig.kubeconfig)"
      >
        {{ kubeconfig.name }}
      </div>
    </div>
    <p v-else class="empty-message">No kubeconfigs available.</p>

    <div v-if="selectedKubeconfig" class="kubeconfig-box">
      <pre>{{ selectedKubeconfig }}</pre>
      <button class="copy-button" @click="copyYaml">Copy</button>
    </div>
  </div>
</template>

<style scoped>
.catalog {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.catalog-item {
  padding: 1rem;
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.catalog-item:hover {
  background-color: var(--color-background-mute);
}

.empty-message {
  text-align: center;
  color: var(--color-text);
}

.kubeconfig-box {
  margin-top: 2rem;
  padding: 1rem;
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  position: relative;
}

.copy-button {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  padding: 0.5rem;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  cursor: pointer;
}
</style>
