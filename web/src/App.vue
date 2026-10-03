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
        @navigate="currentView = $event"
        @select="select"
        @new-chat="onNewChat"
        @rename="rename"
        @delete="onDelete"
      />

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