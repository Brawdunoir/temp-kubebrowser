<script setup lang="ts">
import axios from 'axios'
import { ref, onMounted, watch } from 'vue'
import Home from './components/Home.vue'

const username = ref('')

const loading = ref(true)

onMounted(async () => {
  const response = await axios.get<string>('/api/me')
  username.value = response.data
  loading.value = false
})
</script>

<template>
  <header class="app-header">
    <h1>KubeBrowser</h1>
    <span v-if="!loading" class="username">{{ username }}</span>
  </header>
  <Home/>
</template>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  background-color: var(--color-background);
  border-bottom: 1px solid var(--color-border);
}

.logo {
  height: 50px;
}

.username {
  font-size: 1rem;
  color: var(--color-text);
}

header {
  line-height: 1.5;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
