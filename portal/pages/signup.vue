<script setup lang="ts">
import type { RegisterRequest } from '~/types'

const { register } = useAuth()
const form = ref<RegisterRequest>({ username: '', password: '', name: '', studentNo: '' })
const loading = ref(false)
const success = ref(false)
const errors = ref<Record<string, string>>({})

function validate(): boolean {
  errors.value = {}
  if (!form.value.username.trim()) { errors.value.username = '请输入用户名'; return false }
  if (form.value.username.trim().length < 3) { errors.value.username = '用户名至少 3 个字符'; return false }
  if (!/^[a-zA-Z0-9_]+$/.test(form.value.username.trim())) { errors.value.username = '用户名只能包含字母、数字和下划线'; return false }
  if (!form.value.password) { errors.value.password = '请输入密码'; return false }
  if (form.value.password.length < 6) { errors.value.password = '密码至少 6 位'; return false }
  if (form.value.password.length > 32) { errors.value.password = '密码不能超过 32 位'; return false }
  if (!form.value.name.trim()) { errors.value.name = '请输入真实姓名'; return false }
  return true
}

function clearError(field: string) { delete errors.value[field] }

const passwordStrength = computed(() => {
  const p = form.value.password
  if (!p) return 0
  let score = 0
  if (p.length >= 6) score++
  if (p.length >= 10) score++
  if (/[a-z]/.test(p) && /[A-Z]/.test(p)) score++
  if (/\d/.test(p)) score++
  if (/[^a-zA-Z0-9]/.test(p)) score++
  return Math.min(score, 4)
})

const strengthLabel = ['', '弱', '一般', '较强', '强']
const strengthColor = ['', 'bg-red-400', 'bg-orange-400', 'bg-blue-400', 'bg-green-400']

async function handleRegister() {
  if (!validate()) return
  loading.value = true; errors.value = {}
  const ok = await register({ ...form.value, username: form.value.username.trim(), name: form.value.name.trim() })
  loading.value = false
  if (ok) success.value = true
  else errors.value = { general: '注册失败，用户名可能已存在' }
}
</script>

<template>
  <div class="max-w-md mx-auto px-4 py-12">
    <div class="text-center mb-8">
      <div class="w-16 h-16 rounded-2xl bg-blue-50 flex items-center justify-center mx-auto mb-4">
        <Icon name="heroicons:user-group" class="w-8 h-8 text-blue-600" />
      </div>
      <h1 class="text-2xl font-bold text-gray-900">学生注册</h1>
      <p class="text-gray-500 mt-2 text-sm">加入大学生就业平台</p>
    </div>

    <!-- Success -->
    <div v-if="success" class="bg-white rounded-2xl border border-gray-100 p-8 text-center">
      <Icon name="heroicons:check-badge" class="w-16 h-16 text-green-500 mx-auto mb-4" />
      <h2 class="text-xl font-semibold mb-2">注册成功！</h2>
      <p class="text-gray-500 mb-6 text-sm">现在可以登录并开始求职之旅</p>
      <NuxtLink to="/login" class="inline-block bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-xl font-semibold text-sm transition-all">去登录</NuxtLink>
    </div>

    <!-- Form -->
    <div v-else class="bg-white rounded-2xl border border-gray-100 p-6 space-y-5">
      <!-- Username -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1.5">用户名 <span class="text-red-400">*</span></label>
        <div class="relative">
          <Icon name="heroicons:user" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input v-model="form.username" type="text" placeholder="字母、数字、下划线，至少3位"
            class="w-full pl-10 pr-4 py-3 rounded-xl border text-sm outline-none transition-all"
            :class="errors.username ? 'border-red-300 focus:border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-100'"
            @input="clearError('username')" />
        </div>
        <p v-if="errors.username" class="mt-1.5 text-xs text-red-500 flex items-center gap-1"><Icon name="heroicons:exclamation-triangle" class="w-3.5 h-3.5" />{{ errors.username }}</p>
      </div>

      <!-- Password -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1.5">密码 <span class="text-red-400">*</span></label>
        <div class="relative">
          <Icon name="heroicons:lock-closed" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input v-model="form.password" type="password" placeholder="6-32 位"
            class="w-full pl-10 pr-4 py-3 rounded-xl border text-sm outline-none transition-all"
            :class="errors.password ? 'border-red-300 focus:border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-100'"
            @input="clearError('password')" />
        </div>
        <div v-if="form.password" class="mt-2 flex items-center gap-2">
          <div class="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
            <div class="h-full rounded-full transition-all duration-300" :class="strengthColor[passwordStrength]" :style="{ width: (passwordStrength/4*100)+'%' }" />
          </div>
          <span class="text-xs text-gray-500 w-8">{{ strengthLabel[passwordStrength] }}</span>
        </div>
        <p v-if="errors.password" class="mt-1.5 text-xs text-red-500 flex items-center gap-1"><Icon name="heroicons:exclamation-triangle" class="w-3.5 h-3.5" />{{ errors.password }}</p>
      </div>

      <!-- Name -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1.5">真实姓名 <span class="text-red-400">*</span></label>
        <div class="relative">
          <Icon name="heroicons:identification" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input v-model="form.name" type="text" placeholder="填写真实姓名"
            class="w-full pl-10 pr-4 py-3 rounded-xl border text-sm outline-none transition-all"
            :class="errors.name ? 'border-red-300 focus:border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-100'"
            @input="clearError('name')" />
        </div>
        <p v-if="errors.name" class="mt-1.5 text-xs text-red-500 flex items-center gap-1"><Icon name="heroicons:exclamation-triangle" class="w-3.5 h-3.5" />{{ errors.name }}</p>
      </div>

      <!-- Student No -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1.5">学号 <span class="text-gray-400 font-normal">(选填)</span></label>
        <div class="relative">
          <Icon name="heroicons:academic-cap" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input v-model="form.studentNo" type="text" placeholder="选填"
            class="w-full pl-10 pr-4 py-3 rounded-xl border border-gray-200 text-sm outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100 transition-all" />
        </div>
      </div>

      <p v-if="errors.general" class="text-sm text-red-500 bg-red-50 rounded-lg px-4 py-2.5 flex items-center gap-2">
        <Icon name="heroicons:exclamation-triangle" class="w-4 h-4 shrink-0" />{{ errors.general }}
      </p>

      <button :disabled="loading" @click="handleRegister"
        class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white py-3 rounded-xl font-semibold text-sm transition-all hover:shadow-lg hover:shadow-blue-200">
        {{ loading ? '注册中...' : '注册' }}
      </button>

      <p class="text-center text-sm text-gray-500">已有账号？<NuxtLink to="/login" class="text-blue-600 font-medium hover:underline">去登录</NuxtLink></p>
    </div>
  </div>
</template>
