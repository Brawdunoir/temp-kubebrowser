<script setup lang="ts">
import axios from 'axios'
import { ref, onMounted } from 'vue'
import Home from './components/Home.vue'

const userName = ref<string>('')

onMounted(async () => {
  const response = await axios.get<string>('http://localhost:8080/api/me')
  userName.value = response.data
})
</script>

<template>
  <header class="app-header">
    <h1>KubeBrowser</h1>
    <span class="username">{{ userName }}</span>
  </header>
  <Home/>
</template>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 3rem;
  background-color: var(--color-background-soft);
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
