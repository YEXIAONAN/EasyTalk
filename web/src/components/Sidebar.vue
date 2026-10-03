<script setup lang="ts">
defineProps<{
  active: string
  mobileOpen: boolean
}>()

const emit = defineEmits<{
  navigate: [view: 'chat' | 'settings' | 'about']
  clear: []
}>()
</script>

<template>
  <aside class="sidebar" :class="{ open: mobileOpen }">
    <div class="brand">
      <span class="logo">ET</span>
      <span class="name">EasyTalk</span>
    </div>

    <div class="scroll"></div>

    <nav class="bottom">
      <button class="nav-item" @click="emit('clear')">
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <path
            d="M13.5 8A5.5 5.5 0 1 1 8 2.5M13.5 8V3.5M13.5 8h-4.5"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        <span>Clear Session</span>
      </button>
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

.scroll {
  flex: 1;
  overflow-y: auto;
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

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 40;
    transform: translateX(-100%);
    transition: transform 0.2s ease;
    box-shadow: none;
  }

  .sidebar.open {
    transform: translateX(0);
    box-shadow: 0 0 32px rgba(0, 0, 0, 0.24);
  }
}
</style>