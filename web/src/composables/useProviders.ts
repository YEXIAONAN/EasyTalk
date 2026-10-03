import { ref } from 'vue'
import type { Provider } from '../types'
import * as api from '../services/api'

// Module-level state so the chat header and settings view share the same list.
const providers = ref<Provider[]>([])
const loading = ref(false)

async function fetchProviders() {
  loading.value = true
  try {
    providers.value = await api.listProviders()
  } finally {
    loading.value = false
  }
}

async function addProvider(p: Provider) {
  await api.createProvider(p)
  await fetchProviders()
}

async function editProvider(name: string, p: Provider) {
  await api.updateProvider(name, p)
  await fetchProviders()
}

async function removeProvider(name: string) {
  await api.deleteProvider(name)
  await fetchProviders()
}

export function useProviders() {
  return { providers, loading, fetchProviders, addProvider, editProvider, removeProvider }
}