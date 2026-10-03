<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import Sidebar from './components/Sidebar.vue'
import ChatView from './components/ChatView.vue'
import SettingsView from './views/SettingsView.vue'
import AboutView from './views/AboutView.vue'
import { useProviders } from './composables/useProviders'
import { useChat } from './composables/useChat'
import { useSettings } from './composables/useSettings'

type View = 'chat' | 'settings' | 'about'

const currentView = ref<View>('chat')
const menuOpen = ref(false)

const { providers, fetchProviders } = useProviders()
const {
  messages,
  sending,
  requestCount,
  inputTokens,
  outputTokens,
  totalTokens,
  cachedTokens,
  lastResponseMs,
  firstTokenMs,
  clear,
  send,
  stop,
} = useChat()
const { settings, applyTheme } = useSettings()
const provider = ref('')
const model = ref('')

onMounted(() => {
  fetchProviders()
  applyTheme()
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', applyTheme)
})

watch(() => settings.theme, applyTheme)

// Keep the selected provider/model valid whenever the provider list changes.
watch(
  providers,
  (list) => {
    if (!list.some((p) => p.name === provider.value)) {
      provider.value = list[0]?.name ?? ''
    }
    const current = list.find((p) => p.name === provider.value)
    const models = current?.models ?? []
    if (!models.includes(model.value)) {
      model.value = models[0] ?? ''
    }
  },
  { immediate: true },
)

// When the provider changes, default to its first model.
watch(provider, (name) => {
  const current = providers.value.find((p) => p.name === name)
  model.value = current?.models[0] ?? ''
})

function onSend(content: string) {
  send(provider.value, model.value, content)
}

function onClearSession() {
  clear()
  currentView.value = 'chat'
  menuOpen.value = false
}

function onNavigate(view: View) {
  currentView.value = view
  menuOpen.value = false
}
</script>

<template>
  <div class="app">
    <Sidebar
      :active="currentView"
      :mobile-open="menuOpen"
      :provider="provider"
      :model="model"
      :request-count="requestCount"
      :input-tokens="inputTokens"
      :output-tokens="outputTokens"
      :cached-tokens="cachedTokens"
      :total-tokens="totalTokens"
      :last-response-ms="lastResponseMs"
      :first-token-ms="firstTokenMs"
      @navigate="onNavigate"
      @clear="onClearSession"
    />

    <div v-if="menuOpen" class="backdrop" @click="menuOpen = false"></div>

    <main class="main">
      <ChatView
        v-if="currentView === 'chat'"
        v-model:provider="provider"
        v-model:model="model"
        :providers="providers"
        :messages="messages"
        :sending="sending"
        @send="onSend"
        @stop="stop"
        @toggle-menu="menuOpen = !menuOpen"
        @settings="onNavigate('settings')"
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

.backdrop {
  display: none;
}

@media (max-width: 768px) {
  .backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.35);
    z-index: 30;
  }
}
</style>