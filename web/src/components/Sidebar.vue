<script setup lang="ts">
import { ref } from 'vue'
import type { Conversation } from '../types'

defineProps<{
  active: string
  conversations: Conversation[]
  activeId: string | null
}>()

const emit = defineEmits<{
  navigate: [view: 'chat' | 'settings' | 'about']
  select: [id: string]
  newChat: []
  rename: [id: string, title: string]
  delete: [id: string]
}>()

const editingId = ref<string | null>(null)
const editTitle = ref('')

function startRename(c: Conversation) {
  editingId.value = c.id
  editTitle.value = c.title
}

function saveRename(id: string) {
  const title = editTitle.value.trim()
  if (title) emit('rename', id, title)
  editingId.value = null
}
</script>

<template>
  <aside class="sidebar">
    <div class="brand">
      <span class="logo">ET</span>
      <span class="name">EasyTalk</span>
    </div>

    <button class="new-chat" @click="emit('newChat')">
      <svg viewBox="0 0 16 16" width="14" height="14" fill="none">
        <path d="M8 3v10M3 8h10" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
      </svg>
      <span>New Chat</span>
    </button>

    <div class="scroll">
      <div class="section-title">History</div>
      <ul class="history">
        <li v-if="conversations.length === 0" class="history-empty">暂无对话</li>
        <li
          v-for="c in conversations"
          :key="c.id"
          class="history-item"
          :class="{ active: c.id === activeId }"
          @click="emit('select', c.id)"
        >
          <input
            v-if="editingId === c.id"
            v-model="editTitle"
            class="history-input"
            @click.stop
            @keydown.enter="saveRename(c.id)"
            @keydown.esc="editingId = null"
            @blur="saveRename(c.id)"
          />
          <span v-else class="history-title">{{ c.title }}</span>

          <span class="actions" @click.stop>
            <button class="mini-btn" title="重命名" @click="startRename(c)">
              <svg viewBox="0 0 16 16" width="13" height="13" fill="none">
                <path
                  d="M11.3 2.3l2.4 2.4-7.4 7.4-3 0.6 0.6-3 7.4-7.4zM9.8 3.8l2.4 2.4"
                  stroke="currentColor"
                  stroke-width="1.3"
                  stroke-linejoin="round"
                />
              </svg>
            </button>
            <button class="mini-btn danger" title="删除" @click="emit('delete', c.id)">
              <svg viewBox="0 0 16 16" width="13" height="13" fill="none">
                <path
                  d="M3 4.5h10M6.5 2.5h3M5 4.5l.5 9h5l.5-9"
                  stroke="currentColor"
                  stroke-width="1.3"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </button>
          </span>
        </li>
      </ul>
    </div>

    <nav class="bottom">
      <button
        class="nav-item"
        :class="{ active: active === 'settings' }"
        @click="emit('navigate', 'settings')"
      >
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <circle cx="8" cy="8" r="2.4" stroke="currentColor" stroke-width="1.4" />
          <path
            d="M8 1.8v1.4M8 12.8v1.4M1.8 8h1.4M12.8 8h1.4M3.6 3.6l1 1M11.4 11.4l1 1M12.4 3.6l-1 1M4.6 11.4l-1 1"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
          />
        </svg>
        <span>Settings</span>
      </button>
      <button
        class="nav-item"
        :class="{ active: active === 'about' }"
        @click="emit('navigate', 'about')"
      >
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <circle cx="8" cy="8" r="6.2" stroke="currentColor" stroke-width="1.4" />
          <path d="M8 7.2v4M8 5v.2" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
        </svg>
        <span>About</span>
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  height: 100%;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border);
  padding: 14px 12px;
  flex-shrink: 0;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 2px 6px 14px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 7px;
  background: var(--primary);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.3px;
}

.name {
  font-size: 15px;
  font-weight: 650;
}

.new-chat {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  height: 34px;
  background: var(--primary);
  color: #fff;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 550;
  margin-bottom: 16px;
  transition: background 0.12s ease;
}

.new-chat:hover {
  background: var(--primary-hover);
}

.scroll {
  flex: 1;
  overflow-y: auto;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.4px;
  padding: 0 6px 8px;
}

.history {
  list-style: none;
  margin: 0;
  padding: 0;
}

.history-empty {
  padding: 6px;
  font-size: 13px;
  color: var(--text-muted);
}

.history-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 8px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--text-secondary);
  font-size: 13px;
  transition: background 0.12s ease;
}

.history-item:hover {
  background: var(--bg-hover);
}

.history-item.active {
  background: var(--bg-hover);
  color: var(--text);
  font-weight: 550;
}

.history-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-input {
  flex: 1;
  min-width: 0;
  border: 1px solid var(--primary);
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 13px;
  outline: none;
}

.actions {
  display: none;
  align-items: center;
  gap: 2px;
}

.history-item:hover .actions {
  display: inline-flex;
}

.mini-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 5px;
  color: var(--text-muted);
}

.mini-btn:hover {
  background: var(--border);
  color: var(--text);
}

.mini-btn.danger:hover {
  color: #ef4444;
}

.bottom {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  color: var(--text-secondary);
  transition: background 0.12s ease, color 0.12s ease;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text);
}

.nav-item.active {
  background: var(--bg-hover);
  color: var(--text);
  font-weight: 550;
}
</style>