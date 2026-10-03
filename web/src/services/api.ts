import type { Provider } from '../types'

const BASE = '/api'

export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  })

  if (!res.ok) {
    let message = res.statusText
    try {
      const body = await res.json()
      if (body && body.error) message = body.error
    } catch {
      /* ignore */
    }
    throw new ApiError(res.status, message)
  }

  if (res.status === 204) {
    return undefined as T
  }
  return (await res.json()) as T
}

export function listProviders(): Promise<Provider[]> {
  return request<Provider[]>('/providers')
}

export function createProvider(p: Provider): Promise<Provider> {
  return request<Provider>('/providers', { method: 'POST', body: JSON.stringify(p) })
}

export function updateProvider(name: string, p: Provider): Promise<Provider> {
  return request<Provider>(`/providers/${encodeURIComponent(name)}`, {
    method: 'PUT',
    body: JSON.stringify(p),
  })
}

export function deleteProvider(name: string): Promise<void> {
  return request<void>(`/providers/${encodeURIComponent(name)}`, { method: 'DELETE' })
}

export function testProvider(name: string): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>(`/providers/${encodeURIComponent(name)}/test`, {
    method: 'POST',
  })
}