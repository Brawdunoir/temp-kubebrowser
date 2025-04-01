<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import type { Kubeconfig } from '../types/Kubeconfig'
import KubeconfigCatalog from './KubeconfigCatalog.vue'
import KubeconfigDisplay from './KubeconfigDisplay.vue'
import HelloComponent from './HelloComponent.vue'
import SearchBox from './SearchBox.vue'
import { BsEmojiSurpriseFill } from '@kalimahapps/vue-icons';

const kubeconfigs = ref<Kubeconfig[]>([])
const searchQuery = ref('')
const selectedKubeconfig = ref<string | null>(null)
const indexSelected = ref<number | null>(null)
const emptyKubeconfigs = ref(false)

const filteredKubeconfigs = computed(() => {
  if (!searchQuery.value) return kubeconfigs.value
  selectedKubeconfig.value = null
  indexSelected.value = null
  const query = searchQuery.value.toLowerCase()
  const filtered = kubeconfigs.value.filter((kubeconfig) => kubeconfig.name.toLowerCase().includes(query))
  return filtered
})

function updateSelectedKubeconfig(kubeconfig: string, index: number) {
  indexSelected.value = index
  selectedKubeconfig.value = kubeconfig
}

onMounted(async () => {
  if (import.meta.env.DEV) {
    // Mock response for development
    kubeconfigs.value = [
      { name: 'Cluster number 1', kubeconfig: { apiVersion: 'v1', kind: 'Config' } },
      { name: 'Cluster number 2', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 3', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 4', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 5', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 6', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 7', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Another cluster', kubeconfig: { apiVersion: 'v1', kind: 'Another' } },
      { name: 'Cluster number 4', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 5', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 6', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Cluster number 7', kubeconfig: { apiVersion: 'v1', kind: 'Config2' } },
      { name: 'Another cluster', kubeconfig: { apiVersion: 'v1', kind: 'Another' } },
    ]
  } else {
    const response = await axios.get<Kubeconfig[]>('/api/kubeconfigs')
    if (!response.data.length) {
      emptyKubeconfigs.value = true
    }
    kubeconfigs.value = response.data
  }
})
</script>

<template>
  <HelloComponent />
  <div v-if="emptyKubeconfigs" class="flex flex-col gap-4 items-center justify-center">
    <BsEmojiSurpriseFill class="w-10 h-10 text-gray-600"/>
    <p class="text-gray-300">oops, it seems like you don't have acces to any clusters</p>
  </div>
  <div v-else class="my-10">
    <div class="flex space-x-8 my-8">
      <div class="space-y-4">
        <SearchBox v-model="searchQuery" placeholder="Search clusters..." />
        <KubeconfigCatalog
          :kubeconfigs="filteredKubeconfigs"
          :index-selected="indexSelected"
          @kubeconfig-selected="updateSelectedKubeconfig"
        />
      </div>
      <KubeconfigDisplay class="w-full" :kubeconfig="selectedKubeconfig" :catalog-length="filteredKubeconfigs.length" />
    </div>
  </div>
</template>
