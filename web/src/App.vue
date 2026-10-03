<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import Sidebar from './components/Sidebar.vue'
import ChatView from './components/ChatView.vue'
import SettingsView from './views/SettingsView.vue'
import AboutView from './views/AboutView.vue'
import { useProviders } from './composables/useProviders'

type View = 'chat' | 'settings' | 'about'

const currentView = ref<View>('chat')

const { providers, fetchProviders } = useProviders()
const provider = ref('')
const model = ref('')
const sending = ref(false)

onMounted(fetchProviders)

// Keep the selected provider/model valid whenever the provider list changes.
watch(
  providers,
  (list) => {
    if (!list.some((p) => p.name === provider.value)) {
      provider.value = list[0]?.name ?? ''
      model.value = list[0]?.models[0] ?? ''
    } else if (!model.value) {
      const current = list.find((p) => p.name === provider.value)
      model.value = current?.models[0] ?? ''
    }
  },
  { immediate: true },
)

// When the provider changes, default to its first model.
watch(provider, (name) => {
  const current = providers.value.find((p) => p.name === name)
  model.value = current?.models[0] ?? ''
})
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