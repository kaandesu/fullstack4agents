# Skill: Store Manager Pattern

How to create a new feature store (`*-manager.ts`) following the same pattern as `auth-manager.ts`.

The template already includes `auth-manager.ts`. Use this pattern for every new feature: `posts-manager.ts`, `workspace-manager.ts`, etc.

---

## Pattern

Every manager store:
1. Lives in `app/stores/<feature>-manager.ts`
2. Wraps all backend calls in `auth.API(...)` for automatic 401 handling
3. Uses `createToast` for user feedback
4. Exposes `loading` + data refs
5. Persists only what the user should see immediately on next visit

```ts
// app/stores/posts-manager.ts
// createToast is auto-imported from app/utils/create-toast.ts
// It delegates to useToastManagerStore — never import vue-sonner's toast directly

export type Post = {
  id: string
  title: string
  content: string
  created: string
}

export const usePostsManager = defineStore('PostsManager', () => {
  const auth = useAuthManager()
  const config = useRuntimeConfig()

  const loading = ref(false)
  const posts = ref<Post[]>([])
  const selectedPost = ref<Post | null>(null)

  // Wrap every Gin API call with auth.API() — handles 401 + token refresh automatically.
  // Use it at the call site (here), not inside helper functions.

  const fetchPosts = async () => {
    loading.value = true
    try {
      const res = await auth.API(() =>
        $fetch<{ data: Post[] }>('/api/posts', {
          baseURL: config.public.apiUrl as string,
          headers: { Authorization: `Bearer ${auth.token}` },
        })
      )
      posts.value = res.data ?? []
    } catch (e: any) {
      createToast({ message: 'Failed to load posts', type: 'error' })()
    } finally {
      loading.value = false
    }
  }

  const createPost = async (
    payload: { title: string; content: string },
    callbacks?: { onSuccess?: () => void }
  ) => {
    loading.value = true
    try {
      const res = await auth.API(() =>
        $fetch<{ data: Post }>('/api/posts', {
          method: 'POST',
          baseURL: config.public.apiUrl as string,
          headers: { Authorization: `Bearer ${auth.token}` },
          body: payload,
        })
      )
      posts.value.unshift(res.data!)
      createToast({ message: 'Post created', type: 'success' })()
      callbacks?.onSuccess?.()
    } catch (e: any) {
      createToast({ message: 'Failed to create post', toastOps: { description: e?.message }, type: 'error' })()
    } finally {
      loading.value = false
    }
  }

  const deletePost = async (id: string) => {
    try {
      await auth.API(() =>
        $fetch(`/api/posts/${id}`, {
          method: 'DELETE',
          baseURL: config.public.apiUrl as string,
          headers: { Authorization: `Bearer ${auth.token}` },
        })
      )
      posts.value = posts.value.filter((p) => p.id !== id)
      createToast({ message: 'Post deleted', type: 'success' })()
    } catch {
      createToast({ message: 'Failed to delete post', type: 'error' })()
    }
  }

  return { loading, posts, selectedPost, fetchPosts, createPost, deletePost }
}, {
  persist: {
    storage: piniaPluginPersistedstate.localStorage(),
    pick: ['posts'], // only persist what should survive a page reload
  },
})
```

---

## Key rules

- **Always use `auth.API(fn)`** for Gin API calls — never call `$fetch` directly from a store without it
- **Don't wrap auth operations themselves** (`signIn`, `signUp`) with `API()` — those are already handled inside `auth-manager.ts`
- **Keep `loading` granular** — one `loading` ref per store is usually enough; split if you have concurrent operations
- **Callbacks for navigation** — accept `onSuccess`/`onError` callbacks for actions that should trigger navigation, don't call `navigateTo` directly from the store

---

## PocketBase SDK calls (no Gin)

For direct PocketBase collection access (bypassing Gin), use `usePocketBase()` directly — no `API()` wrapper needed since the SDK handles its own retries:

```ts
const pb = usePocketBase()

const fetchPosts = async () => {
  const records = await pb.collection('posts').getFullList({ sort: '-created' })
  posts.value = records as unknown as Post[]
}
```
