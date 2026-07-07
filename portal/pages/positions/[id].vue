<script setup lang="ts">
import type { Position, Application, ApiResponse, PaginatedData } from '~/types'

const route = useRoute()
const id = route.params.id as string
const { token, isLoggedIn } = useAuth()

const { data: pos } = await useAsyncData(`pos-${id}`, () =>
  $fetch<ApiResponse<Position>>('http://127.0.0.1:8080/api/v1/portal/positions/' + id).then(r => r.code === 0 ? r.data : null).catch(() => null)
)

const applying = ref(false)
const applied = ref(false)
const failed = ref(false)
const showDuplicateModal = ref(false)

// Pre-check: has the current user already applied to this position?
if (import.meta.client && isLoggedIn.value) {
  $fetch<ApiResponse<PaginatedData<Application>>>('http://127.0.0.1:8080/api/v1/portal/student/applications', {
    query: { positionId: id, pageSize: 1 },
    headers: { Authorization: `Bearer ${token.value}` }
  }).then(r => {
    if (r.code === 0 && r.data?.list?.length) applied.value = true
  }).catch(() => {})
}

async function handleApply() {
  if (!isLoggedIn.value) { await navigateTo('/login'); return }
  applying.value = true; failed.value = false
  try {
    const res = await $fetch<ApiResponse<null>>('http://127.0.0.1:8080/api/v1/portal/student/apply', {
      method: 'POST', body: { positionId: Number(id) },
      headers: { Authorization: `Bearer ${token.value}` }
    })
    if (res.code === 0) { applied.value = true }
    else if (res.code === 40900) { showDuplicateModal.value = true }
    else { failed.value = true }
  } catch { failed.value = true }
  applying.value = false
}
</script>

<template>
  <div v-if="pos" class="max-w-5xl mx-auto px-4 py-6">
    <!-- Breadcrumb -->
    <div class="flex items-center gap-2 text-sm text-gray-400 mb-6">
      <NuxtLink to="/positions" class="hover:text-blue-600 transition-colors">职位搜索</NuxtLink>
      <span>/</span>
      <span class="text-gray-600">{{ pos.title }}</span>
    </div>

    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Main content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Header -->
        <div class="bg-white rounded-2xl border border-gray-100 p-6">
          <div class="flex items-start justify-between gap-4 mb-4">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">{{ pos.title }}</h1>
              <div class="flex flex-wrap items-center gap-x-4 gap-y-1 mt-2 text-sm text-gray-500">
                <span class="flex items-center gap-1"><Icon name="heroicons:building-office" class="w-4 h-4" />{{ pos.enterprise?.name }}</span>
                <span class="flex items-center gap-1"><Icon name="heroicons:map-pin" class="w-4 h-4" />{{ pos.city }}</span>
                <span v-if="pos.status === 1" class="bg-green-50 text-green-700 text-xs px-2 py-0.5 rounded-full font-medium">招聘中</span>
                <span v-else class="bg-gray-100 text-gray-500 text-xs px-2 py-0.5 rounded-full">已下架</span>
              </div>
            </div>
            <span class="text-2xl font-bold text-blue-600 whitespace-nowrap shrink-0">
              {{ pos.salaryMin ? `¥${(pos.salaryMin/1000).toFixed(0)}K-${(pos.salaryMax/1000).toFixed(0)}K` : '薪资面议' }}
            </span>
          </div>

          <div class="flex flex-wrap gap-3">
            <div v-if="pos.educationRequirement" class="flex items-center gap-1.5 text-xs bg-gray-50 text-gray-600 px-3 py-1.5 rounded-full">
              <Icon name="heroicons:book-open" class="w-3.5 h-3.5" />{{ pos.educationRequirement }}
            </div>
            <div v-if="pos.experienceRequirement" class="flex items-center gap-1.5 text-xs bg-gray-50 text-gray-600 px-3 py-1.5 rounded-full">
              <Icon name="heroicons:clock" class="w-3.5 h-3.5" />{{ pos.experienceRequirement }}年经验
            </div>
            <div class="flex items-center gap-1.5 text-xs bg-gray-50 text-gray-600 px-3 py-1.5 rounded-full">
              <Icon name="heroicons:user-group" class="w-3.5 h-3.5" />招{{ pos.headcount }}人
            </div>
            <div v-if="pos.industry" class="flex items-center gap-1.5 text-xs bg-blue-50 text-blue-600 px-3 py-1.5 rounded-full">
              {{ pos.industry.name }}
            </div>
          </div>
        </div>

        <!-- Description -->
        <div v-if="pos.description" class="bg-white rounded-2xl border border-gray-100 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <span class="w-1 h-5 bg-blue-600 rounded-full"></span>职位描述
          </h2>
          <div class="text-gray-600 text-sm leading-relaxed space-y-3 whitespace-pre-wrap">{{ pos.description }}</div>
        </div>

        <!-- Requirements -->
        <div v-if="pos.requirement" class="bg-white rounded-2xl border border-gray-100 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <span class="w-1 h-5 bg-blue-600 rounded-full"></span>任职要求
          </h2>
          <div class="text-gray-600 text-sm leading-relaxed whitespace-pre-wrap">{{ pos.requirement }}</div>
        </div>

        <!-- Welfare -->
        <div v-if="pos.welfare" class="bg-white rounded-2xl border border-gray-100 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <span class="w-1 h-5 bg-blue-600 rounded-full"></span>福利待遇
          </h2>
          <div class="flex flex-wrap gap-2">
            <span v-for="w in pos.welfare.split(/[,，、]/)" :key="w" class="bg-orange-50 text-orange-700 text-xs px-3 py-1.5 rounded-full">{{ w.trim() }}</span>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-4">
        <!-- Company card -->
        <div class="bg-white rounded-2xl border border-gray-100 p-5">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center shadow-sm shrink-0 overflow-hidden bg-gray-100">
              <img v-if="pos.enterprise?.logoUrl" :src="pos.enterprise?.logoUrl" class="w-full h-full object-cover" />
              <span v-else class="text-white text-lg font-bold bg-gradient-to-br from-blue-500 to-blue-600 w-full h-full flex items-center justify-center">{{ (pos.enterprise?.name || '?').charAt(0) }}</span>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">{{ pos.enterprise?.name }}</h3>
              <p class="text-xs text-gray-500">{{ pos.enterprise?.industry?.name || '' }}</p>
            </div>
          </div>
          <NuxtLink :to="`/enterprises/${pos.enterpriseId}`" class="block text-center text-sm text-blue-600 hover:text-blue-700 border border-blue-200 rounded-lg py-2 hover:bg-blue-50 transition-colors">
            查看企业详情
          </NuxtLink>
        </div>

        <!-- Salary card -->
        <div class="bg-white rounded-2xl border border-gray-100 p-5">
          <h3 class="text-sm font-semibold text-gray-700 mb-3">薪资信息</h3>
          <div class="text-2xl font-bold text-blue-600 mb-1">{{ pos.salaryMin ? `¥${(pos.salaryMin/1000).toFixed(0)}K-${(pos.salaryMax/1000).toFixed(0)}K` : '面议' }}</div>
          <p class="text-xs text-gray-400">{{ pos.salaryType === 'yearly' ? '年薪' : '月薪' }}</p>
        </div>

        <!-- Job info card -->
        <div class="bg-white rounded-2xl border border-gray-100 p-5">
          <h3 class="text-sm font-semibold text-gray-700 mb-3">职位信息</h3>
          <div class="space-y-2.5 text-sm text-gray-600">
            <div class="flex justify-between"><span class="text-gray-400">学历要求</span><span>{{ pos.educationRequirement || '不限' }}</span></div>
            <div class="flex justify-between"><span class="text-gray-400">工作经验</span><span>{{ pos.experienceRequirement ? pos.experienceRequirement+'年' : '不限' }}</span></div>
            <div class="flex justify-between"><span class="text-gray-400">招聘人数</span><span>{{ pos.headcount }}人</span></div>
            <div class="flex justify-between"><span class="text-gray-400">工作城市</span><span>{{ pos.city }}</span></div>
            <div class="flex justify-between"><span class="text-gray-400">行业领域</span><span>{{ pos.industry?.name || '-' }}</span></div>
          </div>
        </div>

        <!-- Apply button -->
        <button
          @click="handleApply"
          :disabled="applying || applied || pos.status !== 1"
          class="w-full py-3 rounded-xl font-semibold text-base transition-all duration-200"
          :class="applied
            ? 'bg-green-500 text-white cursor-default'
            : failed
              ? 'bg-red-500 text-white'
              : pos.status !== 1
                ? 'bg-gray-200 text-gray-400 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700 text-white shadow-lg shadow-blue-200 hover:shadow-xl hover:-translate-y-0.5 active:scale-[0.98]'"
        >
          <span v-if="applying">投递中...</span>
          <span v-else-if="applied" class="flex items-center justify-center gap-1.5">已投递 <Icon name="heroicons:check" class="w-5 h-5" /></span>
          <span v-else-if="failed" class="flex items-center justify-center gap-1.5">投递失败 <Icon name="heroicons:x-mark" class="w-5 h-5" /></span>
          <span v-else-if="pos.status !== 1">已下架</span>
          <span v-else>{{ isLoggedIn ? '立即投递简历' : '登录后投递' }}</span>
        </button>
      </div>
    </div>
    <!-- Duplicate modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showDuplicateModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="showDuplicateModal = false" />
          <div class="relative bg-white rounded-2xl shadow-2xl p-8 max-w-sm w-full text-center">
            <div class="w-16 h-16 rounded-full bg-blue-50 flex items-center justify-center mx-auto mb-5">
              <Icon name="heroicons:information-circle" class="w-8 h-8 text-blue-500" />
            </div>
            <h2 class="text-xl font-bold text-gray-900 mb-2">已投递过该职位</h2>
            <p class="text-gray-500 text-sm mb-6">你之前已经向这个职位投递过简历了，可以在个人中心查看投递进度</p>
            <div class="flex gap-3">
              <button @click="showDuplicateModal = false" class="flex-1 py-2.5 rounded-lg border border-gray-200 text-gray-600 text-sm font-medium hover:bg-gray-50 transition-colors">知道了</button>
              <NuxtLink to="/student" class="flex-1 py-2.5 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 transition-colors text-center">查看投递</NuxtLink>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>

</template>

<style scoped>
.modal-enter-active { transition: all 0.25s ease-out; }
.modal-leave-active { transition: all 0.2s ease-in; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.modal-enter-from .relative { transform: scale(0.95); }
</style>
