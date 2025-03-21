<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import YAML from 'yaml'
import { copyToClipboard } from '../utils/clipboard'

interface KubeConfig {
  name: string
  kubeconfig: object
}

const kubeconfigs = ref<KubeConfig[]>([])
const searchQuery = ref<string>('')
const selectedKubeconfig = ref<string | null>(null)

const filteredKubeconfigs = computed(() => {
  if (!searchQuery.value) return kubeconfigs.value
  const query = searchQuery.value.toLowerCase()
  return kubeconfigs.value.filter(kubeconfig =>
    kubeconfig.name.toLowerCase().includes(query)
  )
})

onMounted(async () => {
  const response = await axios.get<KubeConfig[]>('/api/kubeconfigs')
  kubeconfigs.value = response.data
  console.log(kubeconfigs.value)
})

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
  <div class="home">
    <input
      v-model="searchQuery"
      type="text"
      placeholder="Search kubeconfigs..."
      class="search-bar"
    />
    <div v-if="filteredKubeconfigs.length" class="catalog">
      <div
        v-for="kubeconfig in filteredKubeconfigs"
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
.home {
  padding: 1rem;
}

.search-bar {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 1rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
}

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
