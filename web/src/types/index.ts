export interface Provider {
  name: string
  base_url: string
  api_key: string
  models: string[]
}

export interface ChatMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
}

export interface Conversation {
  id: string
  title: string
  createdAt: number
  updatedAt: number
  messages: ChatMessage[]
}