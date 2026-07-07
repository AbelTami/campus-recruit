<script setup lang="ts">
import type { Position, PortalStats, ApiResponse, PaginatedData } from '~/types'

const { loggedIn } = useUserSession()

const searchQuery = ref('')
const locationFilter = ref('')

const { data: stats } = await useLazyAsyncData('stats', () =>
  $fetch<ApiResponse<PortalStats>>('http://127.0.0.1:8080/api/v1/portal/stats').then(r => r.code === 0 ? r.data : null).catch(() => null)
)

const { data: positions, status } = await useLazyAsyncData('home-positions', () =>
  $fetch<ApiResponse<PaginatedData<Position>>>('http://127.0.0.1:8080/api/v1/portal/positions', { query: { pageIndex: 1, pageSize: 50 } })
    .then(r => r.code === 0 ? r.data.list : []).catch(() => [])
)

const featuredJobs = computed(() => {
  if (!positions.value) return []
  return positions.value.filter(p => p.status === 1).slice(0, 9)
})

const filteredJobs = computed(() => {
  if (!positions.value) return []
  return positions.value.filter((p: Position) => {
    const match = !searchQuery.value || p.title.includes(searchQuery.value) || (p.enterprise?.name || '').includes(searchQuery.value)
    const loc = !locationFilter.value || p.city === locationFilter.value
    return match && loc
  }).slice(0, 9)
})

const cities = ['深圳', '北京', '上海', '杭州', '广州', '成都', '武汉', '南京']
const salaryFmt = (min: number | null, max: number | null) => min ? `¥${(min/1000).toFixed(0)}-${(max!/1000).toFixed(0)}K` : '面议'
</script>

<template>
  <div>
    <!-- Hero -->
    <HomeHero v-model:searchQuery="searchQuery" />

    <!-- Stats Bar -->
    <ClientOnly>
      <div v-if="stats" class="max-w-6xl mx-auto px-4 -mt-12 relative z-20">
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div v-for="s in [
          { label: '在校学生', value: stats.totalStudents.toLocaleString(), icon: 'heroicons:users', color: 'from-blue-500 to-blue-600' },
          { label: '在招职位', value: stats.totalPositions.toLocaleString(), icon: 'heroicons:briefcase', color: 'from-emerald-500 to-emerald-600' },
          { label: '合作企业', value: stats.totalEnterprises.toLocaleString(), icon: 'heroicons:building-office-2', color: 'from-violet-500 to-violet-600' },
          { label: '就业率', value: (stats.employmentRate || 0).toFixed(1) + '%', icon: 'heroicons:chart-bar', color: 'from-amber-500 to-amber-600' },
        ]" :key="s.label" class="bg-white rounded-2xl shadow-lg p-5 hover:shadow-xl transition-shadow duration-300">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-xl bg-gradient-to-br flex items-center justify-center shrink-0" :class="s.color">
              <Icon :name="s.icon" class="w-6 h-6 text-white" />
            </div>
            <div>
              <div class="text-2xl font-bold text-gray-900">{{ s.value }}</div>
              <div class="text-xs text-gray-500">{{ s.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    </ClientOnly>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-4 py-16">
      <!-- City Filter -->
      <div class="mb-10">
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-sm font-semibold text-gray-700 mr-2">热门城市：</span>
          <button
            v-for="city in cities" :key="city"
            @click="locationFilter = locationFilter === city ? '' : city"
            class="px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200"
            :class="locationFilter === city
              ? 'bg-blue-600 text-white shadow-md scale-105'
              : 'bg-white text-gray-600 shadow-sm ring-1 ring-gray-200 hover:ring-blue-400 hover:text-blue-600'"
          >{{ city }}</button>
        </div>
      </div>

      <!-- Section Title -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">{{ searchQuery ? '搜索结果' : '热门职位' }}</h2>
          <p class="text-gray-500 text-sm mt-1" v-if="searchQuery">包含"{{ searchQuery }}"的职位</p>
        </div>
        <NuxtLink to="/positions" class="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center gap-1">
          查看全部 <Icon name="heroicons:arrow-right" class="w-4 h-4" />
        </NuxtLink>
      </div>

      <!-- Skeleton -->
      <div v-if="status === 'pending'" class="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="n in 6" :key="n" class="bg-white rounded-2xl border border-gray-100 p-6 animate-pulse">
          <div class="h-5 bg-gray-200 rounded w-3/4 mb-3" />
          <div class="h-4 bg-gray-100 rounded w-1/2 mb-4" />
          <div class="h-4 bg-gray-100 rounded w-2/3 mb-2" />
          <div class="h-4 bg-gray-100 rounded w-1/3" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="!filteredJobs.length" class="text-center py-20">
        <Icon name="heroicons:magnifying-glass" class="w-16 h-16 mx-auto mb-4 text-gray-300" />
        <p class="text-lg text-gray-500">暂无匹配职位</p>
        <button @click="searchQuery='';locationFilter=''" class="text-blue-600 text-sm mt-2 hover:underline">清除筛选条件</button>
      </div>

      <!-- Job Grid -->
      <div v-else class="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="(job, idx) in filteredJobs" :key="job.id"
          class="animate-fadeIn"
          :style="{ animationDelay: `${idx * 80}ms` }"
        >
          <NuxtLink
            :to="`/positions/${job.id}`"
            class="block group bg-white hover:shadow-xl transition-all duration-300 p-6 border border-gray-100 hover:border-blue-200 rounded-2xl h-full"
          >
            <div class="flex justify-between items-start gap-3">
              <div class="min-w-0">
                <h3 class="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors line-clamp-1">{{ job.title }}</h3>
                <p class="text-gray-500 text-sm mt-1">{{ job.enterprise?.name }}</p>
              </div>
              <span class="bg-blue-50 text-blue-700 text-xs px-3 py-1.5 rounded-full font-semibold whitespace-nowrap shrink-0">
                {{ salaryFmt(job.salaryMin, job.salaryMax) }}
              </span>
            </div>

            <div class="mt-4 flex flex-wrap items-center gap-3 text-sm text-gray-500">
              <span class="flex items-center gap-1"><Icon name="heroicons:map-pin" class="w-4 h-4 text-blue-400" />{{ job.city || '全国' }}</span>
              <span class="flex items-center gap-1"><Icon name="heroicons:book-open" class="w-4 h-4 text-blue-400" />{{ job.educationRequirement || '不限' }}</span>
            </div>

            <div class="mt-4 flex flex-wrap gap-2">
              <span v-if="job.industry" class="bg-gray-50 text-xs text-gray-500 px-3 py-1 rounded-full border border-gray-100">{{ job.industry.name }}</span>
              <span v-if="job.experienceRequirement" class="bg-gray-50 text-xs text-gray-500 px-3 py-1 rounded-full border border-gray-100">{{ job.experienceRequirement }}年经验</span>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- CTA -->
      <div v-if="!loggedIn" class="mt-20 text-center bg-gradient-to-r from-blue-600 to-indigo-700 rounded-3xl p-12 text-white">
        <h2 class="text-3xl font-bold mb-3">准备好开启求职之旅了吗？</h2>
        <p class="text-blue-100 text-lg mb-8">免费注册，浏览最新职位，一键投递简历</p>
        <div class="flex justify-center gap-4">
          <NuxtLink to="/signup" class="bg-white text-blue-600 px-8 py-3 rounded-xl font-semibold hover:bg-blue-50 transition-all shadow-lg hover:shadow-xl">立即注册</NuxtLink>
          <NuxtLink to="/positions" class="bg-blue-500 text-white px-8 py-3 rounded-xl font-semibold hover:bg-blue-400 transition-all border border-blue-400">浏览职位</NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.animate-fadeIn {
  animation: fadeIn 0.5s ease-out both;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
