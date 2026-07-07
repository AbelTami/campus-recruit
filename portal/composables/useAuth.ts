export function useAuth() {
  const { loggedIn, session, user, clear, fetch } = useUserSession()

  const token = computed(() => (session.value as any)?.token || '')
  const currentUser = computed(() => (session.value as any)?.user || null)

  async function login(username: string, password: string): Promise<boolean> {
    try {
      await $fetch('/api/auth/login', { method: 'POST', body: { username, password } })
      await fetch()
      return true
    } catch {
      return false
    }
  }

  async function register(data: { username: string; password: string; name: string; studentNo: string }): Promise<boolean> {
    try {
      await $fetch('/api/auth/register', { method: 'POST', body: data })
      return true
    } catch {
      return false
    }
  }

  async function logout() {
    await clear()
    await navigateTo('/')
  }

  return {
    token,
    user: currentUser,
    isLoggedIn: loggedIn,
    login,
    register,
    logout,
  }
}
