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

  return (await res.json()) as T
}

export function listProviders(): Promise<Provider[]> {
  return request<Provider[]>('/providers')
}

export interface Info {
  name: string
  version: string
  commit: string
  repository: string
}

export function getInfo(): Promise<Info> {
  return request<Info>('/info')
}

export function reloadConfig(): Promise<{ success: boolean }> {
  return request<{ success: boolean }>('/config/reload', { method: 'POST' })
}