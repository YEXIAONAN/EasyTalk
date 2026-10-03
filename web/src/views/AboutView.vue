<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { BRANDING } from '../config/branding'

const version = ref('v0.1.0')

onMounted(async () => {
  try {
    const res = await fetch('/api/health')
    const data = await res.json()
    if (data?.version) version.value = data.version
  } catch {
    /* keep default */
  }
})
</script>

<template>
  <div class="about">
    <img class="mark" :src="BRANDING.logoMark" alt="EasyTalk" />
    <h1>{{ BRANDING.name }}</h1>
    <p class="version">{{ version }}</p>
    <p class="desc">
      一个简单、轻量、自托管的 AI 对话工具，通过统一界面连接你自己的 AI API。
    </p>
    <p class="desc sub">
      无需注册、无需数据库、无需云服务。所有对话都在本地浏览器中完成，API Key
      仅保存在你的设备上。
    </p>
  </div>
</template>

<style scoped>
.about {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  padding: 32px;
}

.mark {
  width: 60px;
  height: 60px;
  border-radius: 14px;
  margin-bottom: 16px;
}

h1 {
  font-size: 24px;
  margin: 0 0 4px;
}

.version {
  margin: 0 0 16px;
  color: var(--text-muted);
  font-size: 13px;
}

.desc {
  max-width: 420px;
  margin: 0 0 8px;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.desc.sub {
  font-size: 13px;
  color: var(--text-muted);
}
</style>