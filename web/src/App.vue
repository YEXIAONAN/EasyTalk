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
  conversations,
  activeId,
  messages,
  sending,
  load,
  select,
  newChat,
  rename,
  remove,
  clear,
  send,
  stop,
} = useChat()
const { settings, applyTheme } = useSettings()
const provider = ref('')
const model = ref('')

onMounted(() => {
  fetchProviders()
  load()
  applyTheme()
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', applyTheme)
})

watch(() => settings.theme, applyTheme)

// Keep the selected provider/model valid whenever the provider list changes.
watch(
  providers,
  (list) => {
    if (!list.some((p) => p.name === provider.value)) {
      const def = list.find((p) => p.name === settings.defaultProvider)
      provider.value = (def ?? list[0])?.name ?? ''
    }
    const current = list.find((p) => p.name === provider.value)
    const models = current?.models ?? []
    if (!models.includes(model.value)) {
      const defModel = settings.defaultModel && models.includes(settings.defaultModel)
      model.value = defModel ? settings.defaultModel : (models[0] ?? '')
    }
  },
  { immediate: true },
)

// When the provider changes, default to its first model.
watch(provider, (name) => {
  const current = providers.value.find((p) => p.name === name)
  const models = current?.models ?? []
  const defModel = settings.defaultModel && models.includes(settings.defaultModel)
  model.value = defModel ? settings.defaultModel : (models[0] ?? '')
})

function onSend(content: string) {
  send(provider.value, model.value, content)
}

function onSelect(id: string) {
  select(id)
  menuOpen.value = false
}

function onNavigate(view: View) {
  currentView.value = view
  menuOpen.value = false
}

function onNewChat() {
  newChat()
  currentView.value = 'chat'
}

function onDelete(id: string) {
  if (confirm('删除该对话？')) remove(id)
}
</script>

<template>
  <div class="app">
    <Sidebar
      :active="currentView"
      :conversations="conversations"
      :active-id="activeId"
      :mobile-open="menuOpen"
      @navigate="onNavigate"
      @select="onSelect"
      @new-chat="onNewChat"
      @rename="rename"
      @delete="onDelete"
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
        @clear="clear"
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