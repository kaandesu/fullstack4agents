import { defineStore, storeToRefs } from 'pinia'
import createToast from '~/utils/create-toast'

export type UserAccount = {
  id: string
  email: string
  name?: string
  avatar?: string
}

type Callbacks = {
  onSuccess?: (data?: any) => any
  onError?: (error?: any) => any
}

export const useAuthManager = defineStore(
  'AuthManager',
  () => {
    const pb = usePocketBase()
    const loading = ref(false)
    const initialized = ref(false)

    // PocketBase SDK owns the token — we derive from it
    const token = computed(() => pb.authStore.token)
    const isAuthenticated = computed(() => pb.authStore.isValid)

    // Cached user data for fast rendering — synced on every auth change
    const userAccount = ref<UserAccount | null>(pb.authStore.model as UserAccount | null)

    pb.authStore.onChange((_, model) => {
      userAccount.value = model as UserAccount | null
    })

    /**
     * Wrap any Gin API call with this.
     * On 401, it refreshes the PocketBase session and retries once.
     *
     * Usage — always wrap at the call site, not in the function definition:
     *   const data = await auth.API(() => $fetch('/api/items', { ... }))
     */
    const API = async (fn: () => Promise<any>) => {
      try {
        return await fn()
      } catch (e: any) {
        if (e?.statusCode === 401 || e?.status === 401) {
          const refreshed = await refreshToken()
          if (refreshed) return await fn()
          throw e
        }
        throw e
      }
    }

    const signIn = ({
      email,
      password,
      onSuccess,
      onError,
    }: { email: string; password: string } & Callbacks) => {
      loading.value = true
      return pb
        .collection('users')
        .authWithPassword(email, password)
        .then((auth) => {
          userAccount.value = auth.record as unknown as UserAccount
          createToast({ message: 'Login successful', type: 'success' })()
          onSuccess?.(auth)
        })
        .catch((e) => {
          createToast({
            message: 'Login failed',
            toastOps: { description: e?.message ?? 'Invalid credentials' },
            type: 'error',
          })()
          onError?.(e)
        })
        .finally(() => {
          loading.value = false
        })
    }

    const signUp = ({
      email,
      password,
      name,
      onSuccess,
      onError,
    }: { email: string; password: string; name?: string } & Callbacks) => {
      loading.value = true
      return pb
        .collection('users')
        .create({ email, password, passwordConfirm: password, name })
        .then(async () => {
          await pb.collection('users').authWithPassword(email, password)
          userAccount.value = pb.authStore.model as UserAccount
          createToast({
            message: 'Account created!',
            toastOps: { description: 'Welcome aboard' },
            type: 'success',
          })()
          onSuccess?.(pb.authStore.model)
        })
        .catch((e) => {
          createToast({
            message: 'Registration failed',
            toastOps: { description: e?.message ?? 'Could not create account' },
            type: 'error',
          })()
          onError?.(e)
        })
        .finally(() => {
          loading.value = false
        })
    }

    const refreshToken = async ({ onSuccess, onError }: Callbacks = {}): Promise<boolean> => {
      if (!pb.authStore.token) {
        onError?.({ message: 'No token available' })
        return false
      }
      try {
        await pb.collection('users').authRefresh()
        userAccount.value = pb.authStore.model as UserAccount
        onSuccess?.(pb.authStore.model)
        initialized.value = true
        return true
      } catch {
        logout()
        onError?.({ message: 'Session expired' })
        return false
      }
    }

    const logout = async () => {
      pb.authStore.clear()
      userAccount.value = null
      initialized.value = false
      await navigateTo('/auth/sign-in')
    }

    return {
      loading,
      initialized,
      isAuthenticated,
      token,
      userAccount,
      API,
      signIn,
      signUp,
      refreshToken,
      logout,
    }
  },
  {
    persist: {
      storage: piniaPluginPersistedstate.localStorage(),
      // PocketBase SDK persists the token itself — we only cache userAccount
      pick: ['userAccount'],
    },
  },
)
