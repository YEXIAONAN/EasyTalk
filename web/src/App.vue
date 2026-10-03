<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './components/Sidebar.vue'
import ChatView from './components/ChatView.vue'
import SettingsView from './views/SettingsView.vue'
import AboutView from './views/AboutView.vue'
import type { Provider } from './types'

type View = 'chat' | 'settings' | 'about'

const currentView = ref<View>('chat')

// Providers are loaded from the backend in a later phase.
const providers = ref<Provider[]>([])
const provider = ref('')
const model = ref('')
const sending = ref(false)
</script>

<template>
  <div class="app">
    <Sidebar :active="currentView" @navigate="currentView = $event" />

    <main class="main">
      <ChatView
        v-if="currentView === 'chat'"
        v-model:provider="provider"
        v-model:model="model"
        :providers="providers"
        :sending="sending"
        @settings="currentView = 'settings'"
      />
      <SettingsView v-else-if="currentView === 'settings'" />
      <AboutView v-else />
    </main>
  </div>
</template>

<style scoped>
.app {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.main {
  flex: 1;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>