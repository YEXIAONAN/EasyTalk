<script setup lang="ts"></script>

<template>
  <div class="welcome">
    <h1>EasyTalk</h1>
    <p>一个简单、轻量、自托管的 AI 对话工具。</p>
  </div>
</template>

<style scoped>
.welcome {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 8px;
}

h1 {
  font-size: 28px;
  font-weight: 650;
  margin: 0;
  color: var(--color-primary);
}

p {
  margin: 0;
  color: var(--color-text-secondary);
}
</style>