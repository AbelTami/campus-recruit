<script setup lang="ts">
import type { Position, ApiResponse, PaginatedData } from '~/types'

const route = useRoute()
const page = ref(1)
const keyword = ref('')
const city = ref('')
const education = ref('')
const fadeKey = ref(0)

const cities = ['深圳', '北京', '上海', '杭州', '广州', '合肥', '南京', '成都', '武汉', '西安']
const educations = ['本科', '硕士', '博士', '专科']

// SSR initial fetch — enables SWR caching
const { data: initialData, status } = await useAsyncData('positions', () =>
  $fetch<ApiResponse<PaginatedData<Position>>>('http://127.0.0.1:8080/api/v1/portal/positions', {
    query: { pageIndex: 1, pageSize: 12 }
  }).then(r => r.code === 0 ? r.data : { list: [], total: 0 }).catch(() => ({ list: [], total: 0 }))
)

const list = ref<Position[]>(initialData.value?.list || [])
const total = ref(initialData.value?.total || 0)
const loading = ref(false)

// Client-side fetch for filters / pagination
async function fetchPositions() {
  loading.value = true
  try {
    const q = new URLSearchParams({ pageIndex: String(page.value), pageSize: '12' })
    if (keyword.value) q.set('keyword', keyword.value)
    if (city.value) q.set('city', city.value)
    if (education.value) q.set('educationRequirement', education.value)
    const res = await $fetch<ApiResponse<PaginatedData<Position>>>(`http://127.0.0.1:8080/api/v1/portal/positions?${q}`)
    if (res.code === 0) { list.value = res.data.list; total.value = res.data.total }
  } catch { /* ignore */ }
  loading.value = false
}

let debounceTimer: ReturnType<typeof setTimeout>
function onKeywordInput() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; fetchPositions() }, 300)
}

function goPage(p: number) {
  if (p < 1 || p > Math.ceil(total.value / 12) || p === page.value) return
  page.value = p
  fadeKey.value++
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function clearFilters() { city.value = ''; education.value = ''; keyword.value = ''; onFilterChange() }
function onFilterChange() { page.value = 1; fetchPositions() }

watch([page], () => fetchPositions())
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 pt-8 pb-12">
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">职位搜索</h1>
      <p class="text-gray-500 mt-2">找到最适合你的工作机会</p>
    </div>

    <!-- Filters -->
    <!-- Search bar -->
    <div class="flex items-center gap-2 bg-white rounded-lg border border-gray-200 hover:border-gray-300 px-3 py-2 mb-5 transition-all focus-within:ring-0">
      <Icon name="heroicons:magnifying-glass" class="w-4 h-4 text-gray-400 shrink-0" />
      <input v-model="keyword" type="text" placeholder="搜索职位、企业..." class="flex-1 border-0 outline-none ring-0 focus:ring-0 text-sm text-gray-700 placeholder:text-gray-400" @input="onKeywordInput" />
      <span v-if="keyword" @click="keyword=''; onFilterChange()" class="cursor-pointer text-gray-300 hover:text-gray-500 shrink-0">
        <Icon name="heroicons:x-mark" class="w-4 h-4" />
      </span>
      <span v-if="total && !keyword" class="text-xs text-gray-400 shrink-0">{{ total }} 个职位</span>
    </div>

    <!-- City filter pills -->
    <div class="flex flex-wrap items-center gap-2 mb-3">
      <span class="text-xs font-medium text-gray-500 mr-1">城市：</span>
      <button v-for="c in ['']" :key="''" @click="city=''; onFilterChange()" class="px-3 py-1.5 rounded-lg text-xs transition-all" :class="!city ? 'bg-blue-600 text-white shadow' : 'bg-white border border-gray-200 text-gray-600 hover:border-blue-300'">全部</button>
      <button v-for="c in cities" :key="c" @click="city = city === c ? '' : c; onFilterChange()" class="px-3 py-1.5 rounded-lg text-xs transition-all" :class="city === c ? 'bg-blue-600 text-white shadow' : 'bg-white border border-gray-200 text-gray-600 hover:border-blue-300'">{{ c }}</button>
    </div>

    <!-- Education filter pills -->
    <div class="flex flex-wrap items-center gap-2 mb-3">
      <span class="text-xs font-medium text-gray-500 mr-1">学历：</span>
      <button @click="education=''; onFilterChange()" class="px-3 py-1.5 rounded-lg text-xs transition-all" :class="!education ? 'bg-blue-600 text-white shadow' : 'bg-white border border-gray-200 text-gray-600 hover:border-blue-300'">全部</button>
      <button v-for="e in educations" :key="e" @click="education = education === e ? '' : e; onFilterChange()" class="px-3 py-1.5 rounded-lg text-xs transition-all" :class="education === e ? 'bg-blue-600 text-white shadow' : 'bg-white border border-gray-200 text-gray-600 hover:border-blue-300'">{{ e }}</button>
      <button v-if="city || education || keyword" @click="clearFilters" class="ml-auto text-xs text-blue-600 hover:text-blue-700 font-medium">清除全部</button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="n in 6" :key="n" class="bg-white rounded-2xl border border-gray-100 p-6 animate-pulse">
        <div class="h-5 bg-gray-200 rounded w-3/4 mb-3" />
        <div class="h-4 bg-gray-100 rounded w-1/2 mb-4" />
        <div class="h-4 bg-gray-100 rounded w-2/3" />
      </div>
    </div>

    <!-- Empty -->
    <div v-else-if="!list.length" class="text-center py-20">
      <Icon name="heroicons:magnifying-glass" class="w-16 h-16 mx-auto mb-4 text-gray-300" />
      <p class="text-lg text-gray-500">暂无匹配职位</p>
      <p class="text-sm text-gray-400 mt-1">尝试调整筛选条件</p>
    </div>

    <!-- Results -->
    <div v-else :key="fadeKey" class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3 animate-fadeIn">
      <NuxtLink v-for="pos in list" :key="pos.id" :to="`/positions/${pos.id}`" class="group bg-white hover:shadow-xl transition-all duration-300 p-6 border border-gray-100 hover:border-blue-200 rounded-2xl">
        <div class="flex justify-between items-start gap-3">
          <div class="min-w-0">
            <h3 class="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors line-clamp-1">{{ pos.title }}</h3>
            <p class="text-gray-500 text-sm mt-1">{{ pos.enterprise?.name }}</p>
          </div>
          <span class="bg-blue-50 text-blue-700 text-xs px-3 py-1.5 rounded-full font-semibold whitespace-nowrap shrink-0">
            {{ pos.salaryMin ? `¥${(pos.salaryMin/1000).toFixed(0)}-${(pos.salaryMax/1000).toFixed(0)}K` : '面议' }}
          </span>
        </div>
        <div class="mt-4 flex flex-wrap items-center gap-3 text-sm text-gray-500">
          <span class="flex items-center gap-1"><Icon name="heroicons:map-pin" class="w-4 h-4 text-blue-400" />{{ pos.city || '全国' }}</span>
          <span class="flex items-center gap-1"><Icon name="heroicons:book-open" class="w-4 h-4 text-blue-400" />{{ pos.educationRequirement || '不限' }}</span>
        </div>
      </NuxtLink>
    </div>

    <!-- Pagination -->
    <div v-if="total > 12" class="flex justify-center items-center mt-10 gap-1">
      <button @click="goPage(page-1)" :disabled="page===1" class="px-3 py-2 rounded-lg text-sm border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all flex items-center gap-1">
        <Icon name="heroicons:chevron-left" class="w-4 h-4" />上一页
      </button>
      <template v-for="p in Math.ceil(total/12)" :key="p">
        <button v-if="p === 1 || p === Math.ceil(total/12) || Math.abs(p - page) <= 1"
          @click="goPage(p)"
          class="w-10 h-10 rounded-lg text-sm font-medium transition-all duration-200"
          :class="page === p ? 'bg-blue-600 text-white shadow-md scale-105' : 'text-gray-600 hover:bg-gray-100'"
        >{{ p }}</button>
        <span v-else-if="Math.abs(p - page) === 2" class="text-gray-400">...</span>
      </template>
      <button @click="goPage(page+1)" :disabled="page===Math.ceil(total/12)" class="px-3 py-2 rounded-lg text-sm border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all flex items-center gap-1">
        下一页<Icon name="heroicons:chevron-right" class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
