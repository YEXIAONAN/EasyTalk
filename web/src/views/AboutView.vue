<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { BRANDING } from '../config/branding'
import { useI18n } from '../i18n'
import { getInfo } from '../services/api'
import type { Info } from '../services/api'

const { t } = useI18n()

const info = ref<Info>({ name: BRANDING.name, version: '', commit: '', repository: '' })

onMounted(async () => {
  try {
    info.value = await getInfo()
  } catch {
    /* keep defaults */
  }
})

const repoLabel = computed(() =>
  info.value.repository.replace(/^https?:\/\/github\.com\//, ''),
)

const releasesUrl = computed(() => `${info.value.repository}/releases`)

const releaseTagUrl = computed(() =>
  /^v\d+\.\d+\.\d+$/.test(info.value.version)
    ? `${info.value.repository}/releases/tag/${info.value.version}`
    : '',
)

const licenseUrl = computed(() => `${info.value.repository}/blob/main/LICENSE`)
</script>

<template>
  <div class="about">
    <img class="mark" :src="BRANDING.logoMark" alt="EasyTalk" />
    <h1>{{ info.name || BRANDING.name }}</h1>
    <p class="version">{{ info.version || '—' }}</p>
    <p class="tagline">{{ t('about.tagline') }}</p>

    <div class="meta">
      <div class="meta-row">
        <span class="k">{{ t('about.version') }}</span>
        <span class="v">{{ info.version || '—' }}</span>
      </div>
      <div class="meta-row">
        <span class="k">{{ t('about.build') }}</span>
        <span class="v">{{ info.commit || '—' }}</span>
      </div>
      <div class="meta-row">
        <span class="k">{{ t('about.github') }}</span>
        <a class="v link" :href="info.repository" target="_blank" rel="noopener">{{ repoLabel }}</a>
      </div>
      <div class="meta-row">
        <span class="k">{{ t('about.releases') }}</span>
        <span class="v">
          <a class="link" :href="releasesUrl" target="_blank" rel="noopener">
            {{ t('about.viewReleases') }}
          </a>
          <template v-if="releaseTagUrl">
            <span class="sep">·</span>
            <a class="link" :href="releaseTagUrl" target="_blank" rel="noopener">
              {{ t('about.viewRelease') }}
            </a>
          </template>
        </span>
      </div>
    </div>

    <p class="desc">{{ t('about.description') }}</p>
    <p class="note">{{ t('about.apiKeyNote') }}</p>
    <p class="license">
      {{ t('about.license') }}:
      <a class="link" :href="licenseUrl" target="_blank" rel="noopener">MIT</a>
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
  overflow-y: auto;
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
  margin: 0 0 12px;
  color: var(--text-muted);
  font-size: 13px;
}

.tagline {
  margin: 0 0 24px;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.meta {
  width: 100%;
  max-width: 420px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 8px 16px;
  margin-bottom: 20px;
  text-align: left;
}

.meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 8px 0;
  border-bottom: 1px solid var(--border);
  font-size: 13px;
}

.meta-row:last-child {
  border-bottom: none;
}

.k {
  color: var(--text-muted);
}

.v {
  color: var(--text);
  font-weight: 550;
  text-align: right;
}

.sep {
  margin: 0 6px;
  color: var(--text-muted);
}

.link {
  color: var(--primary);
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}

.desc {
  max-width: 420px;
  margin: 0 0 8px;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.note {
  max-width: 420px;
  margin: 0 0 8px;
  font-size: 13px;
  color: var(--text-muted);
  line-height: 1.5;
}

.license {
  margin: 16px 0 0;
  font-size: 12.5px;
  color: var(--text-muted);
}
</style>