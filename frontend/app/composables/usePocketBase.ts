import PocketBase from 'pocketbase'

// Singleton shared across the app — PocketBase manages its own auth persistence
let _pb: PocketBase | null = null

export function usePocketBase(): PocketBase {
  if (!_pb) {
    const config = useRuntimeConfig()
    _pb = new PocketBase((config.public.pbUrl as string) || 'http://localhost:8090')
  }
  return _pb
}
