<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import YAML from 'yaml'
import { copyToClipboard } from '../utils/clipboard'
import type { Kubeconfig } from '../types/Kubeconfig'
import KubeconfigCatalog from './KubeconfigCatalog.vue'
import Hello from './Hello.vue'

const kubeconfigs = ref<Kubeconfig[]>([])
const searchQuery = ref<string>('')

const filteredKubeconfigs = computed(() => {
  if (!searchQuery.value) return kubeconfigs.value
  const query = searchQuery.value.toLowerCase()
  return kubeconfigs.value.filter((kubeconfig) => kubeconfig.name.toLowerCase().includes(query))
})

onMounted(async () => {
  if (import.meta.env.DEV) {
    // Mock response for development
    kubeconfigs.value = [
      { name: 'Cluster 1', kubeconfig: { apiVersion: 'v1', kind: 'Config' } },
      { name: 'Cluster 2', kubeconfig: { apiVersion: 'v1', kind: 'Config' } },
    ]
  } else {
    const response = await axios.get<Kubeconfig[]>('/api/kubeconfigs')
    kubeconfigs.value = response.data
  }
})
</script>

<template>
  <Hello />
  <div class="my-10">
    <input
      v-model="searchQuery"
      type="text"
      placeholder="Search kubeconfigs..."
      class="p-3 rounded-md bg-gray-800 w-full border-2 border-gray-600"
    />
    <KubeconfigCatalog :kubeconfigs="filteredKubeconfigs" />
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
