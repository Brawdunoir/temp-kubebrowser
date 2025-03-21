<script setup lang="ts">
import axios from 'axios'
import { ref, onMounted } from 'vue'
import Home from './components/Home.vue'
import Header from './components/Header.vue'
import Hello from './components/Hello.vue'

const username = ref('')

const loading = ref(true)

onMounted(async () => {
  if (import.meta.env.DEV) {
    // Mock response for development
    username.value = 'Firstname Lastname'
    loading.value = false
  } else {
    const response = await axios.get<string>('/api/me')
    username.value = response.data
    loading.value = false
  }
})
</script>

<template>
  <Header/>
  <Hello/>
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
