<script setup lang="ts">
import { shallowRef, computed } from '#imports'
import type { ApiResponse } from '~/types'

definePageMeta({ middleware: ['auth'] })

const { token } = useAuth()

interface StudentProfile {
  id: number
  name: string
  studentNo: string
  gender?: number
  educationLevel?: string
  grade?: string
  graduationDate?: string
  phone?: string
  email?: string
  wechat?: string
  qq?: string
  hometownCity?: string
  expectedCity?: string
  expectedIndustry?: string
  expectedSalaryMin?: number
  expectedSalaryMax?: number
  resumeUrl?: string
}

const { data: profile, status, refresh } = useAsyncData('student-profile', () =>
  $fetch<ApiResponse<StudentProfile>>('http://127.0.0.1:8080/api/v1/portal/student/profile', {
    headers: { Authorization: `Bearer ${token.value}` }
  }).then(r => r.code === 0 ? r.data : null).catch(() => null),
{ server: false }
)

const form = shallowRef({
  gender: undefined as number | undefined,
  educationLevel: '',
  grade: '',
  graduationDate: '',
  phone: '',
  email: '',
  wechat: '',
  qq: '',
  hometownCity: '',
  expectedCity: '',
  expectedIndustry: '',
  expectedSalaryMin: undefined as number | undefined,
  expectedSalaryMax: undefined as number | undefined,
  resumeUrl: '',
})

const saving = shallowRef(false)
const saved = shallowRef(false)

// Seed form from loaded profile
watch(profile, (p) => {
  if (!p) return
  form.value = {
    gender: p.gender,
    educationLevel: p.educationLevel || '',
    grade: p.grade || '',
    graduationDate: p.graduationDate || '',
    phone: p.phone || '',
    email: p.email || '',
    wechat: p.wechat || '',
    qq: p.qq || '',
    hometownCity: p.hometownCity || '',
    expectedCity: p.expectedCity || '',
    expectedIndustry: p.expectedIndustry || '',
    expectedSalaryMin: p.expectedSalaryMin,
    expectedSalaryMax: p.expectedSalaryMax,
    resumeUrl: p.resumeUrl || '',
  }
}, { immediate: true })

const completions = computed(() => {
  const f = form.value
  let c = 0
  if (f.gender !== undefined) c++
  if (f.educationLevel) c++
  if (f.grade) c++
  if (f.phone || f.email) c++
  if (f.expectedCity) c++
  if (f.expectedIndustry) c++
  if (f.expectedSalaryMin || f.expectedSalaryMax) c++
  return c
})

const genderOptions = [
  { label: '男', value: 1 },
  { label: '女', value: 2 },
]

const eduOptions = ['高中', '专科', '本科', '硕士', '博士']

async function save() {
  saving.value = true
  try {
    await $fetch('http://127.0.0.1:8080/api/v1/portal/student/profile', {
      method: 'PUT',
      headers: { Authorization: `Bearer ${token.value}` },
      body: JSON.parse(JSON.stringify(form.value)),
    })
    saved.value = true
    setTimeout(() => { saved.value = false }, 2000)
  } catch { /* ignored */ }
  saving.value = false
}
</script>

<template>
  <div class="max-w-2xl mx-auto px-4 py-6">
    <!-- Breadcrumb -->
    <div class="flex items-center gap-1.5 text-[13px] text-gray-300 mb-6">
      <NuxtLink to="/student" class="hover:text-gray-500 transition-colors">个人中心</NuxtLink>
      <Icon name="heroicons:chevron-right" class="w-3 h-3 opacity-40" />
      <span class="text-gray-500">编辑资料</span>
    </div>

    <!-- Loading -->
    <div v-if="status === 'pending'" class="space-y-4">
      <div class="bg-white rounded-xl border border-gray-100/60 p-6 space-y-4 animate-pulse">
        <div class="h-3.5 bg-gray-100 rounded w-20" />
        <div class="grid grid-cols-2 gap-3">
          <div class="h-10 bg-gray-50 rounded-lg" />
          <div class="h-10 bg-gray-50 rounded-lg" />
        </div>
      </div>
      <div class="bg-white rounded-xl border border-gray-100/60 p-6 space-y-4 animate-pulse">
        <div class="h-3.5 bg-gray-100 rounded w-20" />
        <div class="grid grid-cols-2 gap-3">
          <div class="h-10 bg-gray-50 rounded-lg" v-for="n in 6" :key="n" />
        </div>
      </div>
    </div>

    <div v-else-if="profile" class="space-y-4">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center shrink-0">
            <span class="text-white text-sm font-semibold">{{ profile.name?.charAt(0) }}</span>
          </div>
          <div>
            <h1 class="text-base font-semibold text-gray-800">{{ profile.name }}</h1>
            <p class="text-[12px] text-gray-400">完善资料，提高求职匹配度</p>
          </div>
        </div>
        <div class="text-right">
          <div class="text-lg font-semibold tabular-nums" :class="completions === 7 ? 'text-emerald-500' : completions >= 4 ? 'text-blue-500' : 'text-gray-400'">{{ completions }}</div>
          <div class="text-[10px] text-gray-400">/ 7 项</div>
        </div>
      </div>

      <!-- Basic info (read-only) -->
      <div class="bg-white rounded-xl border border-gray-100/60 overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-100/60 flex items-center gap-2">
          <Icon name="heroicons:identification" class="w-4 h-4 text-gray-400" />
          <h2 class="text-[13px] font-medium text-gray-500">基本信息</h2>
        </div>
        <div class="p-5 grid grid-cols-2 gap-3">
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">姓名</label>
            <div class="px-3 py-2.5 rounded-lg bg-gray-50/80 text-sm text-gray-500 border border-transparent">{{ profile.name }}</div>
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">学号</label>
            <div class="px-3 py-2.5 rounded-lg bg-gray-50/80 text-sm text-gray-500 border border-transparent">{{ profile.studentNo }}</div>
          </div>
        </div>
      </div>

      <!-- Personal details -->
      <div class="bg-white rounded-xl border border-gray-100/60 overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-100/60 flex items-center gap-2">
          <Icon name="heroicons:user" class="w-4 h-4 text-gray-400" />
          <h2 class="text-[13px] font-medium text-gray-500">个人详情</h2>
        </div>
        <div class="p-5 grid grid-cols-2 gap-3">
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">性别</label>
            <select v-model.number="form.gender" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all bg-white">
              <option :value="undefined" disabled>请选择</option>
              <option v-for="g in genderOptions" :key="g.value" :value="g.value">{{ g.label }}</option>
            </select>
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">学历</label>
            <select v-model="form.educationLevel" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all bg-white">
              <option value="">请选择</option>
              <option v-for="e in eduOptions" :key="e" :value="e">{{ e }}</option>
            </select>
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">年级</label>
            <input v-model="form.grade" placeholder="2024 级" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">毕业日期</label>
            <input v-model="form.graduationDate" type="date" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">手机</label>
            <input v-model="form.phone" placeholder="手机号" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">邮箱</label>
            <input v-model="form.email" type="email" placeholder="example@mail.com" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">微信</label>
            <input v-model="form.wechat" placeholder="微信号" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">QQ</label>
            <input v-model="form.qq" placeholder="QQ 号" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div class="col-span-2">
            <label class="text-[11px] text-gray-400 mb-1 block">籍贯</label>
            <input v-model="form.hometownCity" placeholder="例如：广东深圳" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
        </div>
      </div>

      <!-- Career -->
      <div class="bg-white rounded-xl border border-gray-100/60 overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-100/60 flex items-center gap-2">
          <Icon name="heroicons:briefcase" class="w-4 h-4 text-gray-400" />
          <h2 class="text-[13px] font-medium text-gray-500">求职意向</h2>
        </div>
        <div class="p-5 grid grid-cols-2 gap-3">
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">期望城市</label>
            <input v-model="form.expectedCity" placeholder="深圳" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">期望行业</label>
            <input v-model="form.expectedIndustry" placeholder="互联网 / IT" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">期望薪资下限</label>
            <div class="relative">
              <input v-model.number="form.expectedSalaryMin" type="number" placeholder="8,000" class="w-full pl-7 pr-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-[11px] text-gray-400">¥</span>
            </div>
          </div>
          <div>
            <label class="text-[11px] text-gray-400 mb-1 block">期望薪资上限</label>
            <div class="relative">
              <input v-model.number="form.expectedSalaryMax" type="number" placeholder="15,000" class="w-full pl-7 pr-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-[11px] text-gray-400">¥</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Contact -->
      <div class="bg-white rounded-xl border border-gray-100/60 overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-100/60 flex items-center gap-2">
          <Icon name="heroicons:envelope" class="w-4 h-4 text-gray-400" />
          <h2 class="text-[13px] font-medium text-gray-500">简历附件</h2>
        </div>
        <div class="p-5">
          <input v-model="form.resumeUrl" placeholder="简历链接地址" class="w-full px-3 py-2.5 rounded-lg border border-gray-150 text-sm text-gray-700 placeholder:text-gray-300 focus:border-blue-400 focus:ring-1 focus:ring-blue-100 outline-none transition-all" />
        </div>
      </div>

      <!-- Save -->
      <button
        @click="save"
        :disabled="saving"
        class="w-full py-2.5 rounded-lg text-[13px] font-medium transition-all duration-200"
        :class="saved
          ? 'bg-emerald-500 text-white shadow-sm shadow-emerald-200'
          : 'bg-gray-900 hover:bg-gray-800 text-white shadow-sm'"
      >
        <span v-if="saved" class="inline-flex items-center gap-1.5"><Icon name="heroicons:check" class="w-4 h-4" />已保存</span>
        <span v-else-if="saving" class="inline-flex items-center gap-1.5"><Icon name="heroicons:arrow-path" class="w-4 h-4 animate-spin" />保存中...</span>
        <span v-else>保存资料</span>
      </button>
    </div>
  </div>
</template>
