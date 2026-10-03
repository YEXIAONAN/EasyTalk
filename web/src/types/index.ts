export interface Provider {
  name: string
  models: string[]
}

export interface ChatMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
  createdAt: number
}

export type Theme = 'light' | 'dark' | 'system'

export interface Settings {
  defaultProvider: string
  defaultModel: string
  streamResponse: boolean
  theme: Theme
  temperature: number | null
  maxTokens: number | null
  topP: number | null
}