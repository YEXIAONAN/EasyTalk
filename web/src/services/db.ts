import type { Conversation } from '../types'

const DB_NAME = 'easytalk'
const STORE = 'conversations'
const VERSION = 1

let dbPromise: Promise<IDBDatabase> | null = null

function openDB(): Promise<IDBDatabase> {
  if (dbPromise) return dbPromise
  dbPromise = new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, VERSION)
    req.onupgradeneeded = () => {
      const db = req.result
      if (!db.objectStoreNames.contains(STORE)) {
        db.createObjectStore(STORE, { keyPath: 'id' })
      }
    }
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
  return dbPromise
}

function reqAsPromise<T>(req: IDBRequest<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
}

export async function getAllConversations(): Promise<Conversation[]> {
  const db = await openDB()
  const tx = db.transaction(STORE, 'readonly')
  return reqAsPromise(tx.objectStore(STORE).getAll() as IDBRequest<Conversation[]>)
}

export async function saveConversation(conv: Conversation): Promise<void> {
  const db = await openDB()
  const tx = db.transaction(STORE, 'readwrite')
  await reqAsPromise(tx.objectStore(STORE).put(conv))
}

export async function deleteConversation(id: string): Promise<void> {
  const db = await openDB()
  const tx = db.transaction(STORE, 'readwrite')
  await reqAsPromise(tx.objectStore(STORE).delete(id))
}