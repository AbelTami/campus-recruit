<script setup lang="ts">
const { loggedIn } = useUserSession()
if (loggedIn.value) { await navigateTo('/student') }

const { login } = useAuth()

const username = ref('')
const password = ref('')
const loading = ref(false)
const errors = ref<{ username?: string; password?: string; general?: string }>({})

function validate(): boolean {
  errors.value = {}
  if (!username.value.trim()) { errors.value.username = '请输入用户名'; return false }
  if (username.value.trim().length < 2) { errors.value.username = '用户名至少 2 个字符'; return false }
  if (!password.value) { errors.value.password = '请输入密码'; return false }
  if (password.value.length < 6) { errors.value.password = '密码至少 6 位'; return false }
  return true
}

function clearError(field: string) { delete (errors.value as any)[field] }

async function handleLogin() {
  if (!validate()) return
  loading.value = true; errors.value = {}
  const ok = await login(username.value.trim(), password.value)
  loading.value = false
  if (ok) await navigateTo('/student')
  else errors.value = { general: '用户名或密码错误，请重试' }
}
</script>

<template>
  <div class="max-w-md mx-auto px-4 py-12">
    <div class="text-center mb-8">
      <div class="w-16 h-16 rounded-2xl bg-blue-50 flex items-center justify-center mx-auto mb-4">
        <Icon name="heroicons:book-open" class="w-8 h-8 text-blue-600" />
      </div>
      <h1 class="text-2xl font-bold text-gray-900">学生登录</h1>
      <p class="text-gray-500 mt-2 text-sm">登录以查看投递记录和申请职位</p>
    </div>

    <div class="bg-white rounded-2xl border border-gray-100 p-6 space-y-5">
      <!-- Username -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1.5">用户名</label>
        <div class="relative">
          <Icon name="heroicons:user" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="username" type="text" placeholder="请输入用户名"
            class="w-full pl-10 pr-4 py-3 rounded-xl border text-sm outline-none transition-all"
            :class="errors.username ? 'border-red-300 focus:border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-100'"
            @input="clearError('username')"
          />
        </div>
        <p v-if="errors.username" class="mt-1.5 text-xs text-red-500 flex items-center gap-1">
          <Icon name="heroicons:exclamation-triangle" class="w-3.5 h-3.5" />{{ errors.username }}
        </p>
      </div>

      <!-- Password -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1.5">密码</label>
        <div class="relative">
          <Icon name="heroicons:lock-closed" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="password" type="password" placeholder="请输入密码"
            class="w-full pl-10 pr-4 py-3 rounded-xl border text-sm outline-none transition-all"
            :class="errors.password ? 'border-red-300 focus:border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-100'"
            @input="clearError('password')"
          />
        </div>
        <p v-if="errors.password" class="mt-1.5 text-xs text-red-500 flex items-center gap-1">
          <Icon name="heroicons:exclamation-triangle" class="w-3.5 h-3.5" />{{ errors.password }}
        </p>
      </div>

      <!-- General error -->
      <p v-if="errors.general" class="text-sm text-red-500 bg-red-50 rounded-lg px-4 py-2.5 flex items-center gap-2">
        <Icon name="heroicons:exclamation-triangle" class="w-4 h-4 shrink-0" />{{ errors.general }}
      </p>

      <button
        :disabled="loading"
        @click="handleLogin"
        class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white py-3 rounded-xl font-semibold text-sm transition-all hover:shadow-lg hover:shadow-blue-200"
      >
        {{ loading ? '登录中...' : '登录' }}
      </button>

      <p class="text-center text-sm text-gray-500">
        还没有账号？<NuxtLink to="/signup" class="text-blue-600 font-medium hover:underline">立即注册</NuxtLink>
      </p>

      <p class="text-center text-xs text-gray-400 border-t border-gray-100 pt-4">
        演示：管理员 admin / admin123 · 学生 student1 / 123456
      </p>
    </div>
  </div>
</template>
