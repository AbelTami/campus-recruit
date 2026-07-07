<script setup lang="ts">
import type { Enterprise, ApiResponse, PaginatedData } from '~/types'

const keyword = ref('')

const { data } = await useAsyncData('enterprises', () =>
  $fetch<ApiResponse<PaginatedData<Enterprise>>>('http://127.0.0.1:8080/api/v1/portal/enterprises', {
    query: { pageSize: 50, keyword: keyword.value }
  }).then(r => r.code === 0 ? r.data : { list: [], total: 0 }).catch(() => ({ list: [], total: 0 }))
)

const gradients = ['from-blue-500 to-blue-600', 'from-emerald-500 to-teal-600', 'from-violet-500 to-purple-600', 'from-amber-500 to-orange-600', 'from-rose-500 to-pink-600', 'from-cyan-500 to-sky-600', 'from-indigo-500 to-indigo-600', 'from-lime-500 to-green-600']

function avatarColor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash)
  return gradients[Math.abs(hash) % gradients.length]
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 pt-8 pb-12">
    <div class="text-center mb-10">
      <h1 class="text-3xl font-bold text-gray-900">企业名录</h1>
      <p class="text-gray-500 mt-2">了解合作企业，找到心仪的公司</p>
      <div class="max-w-md mx-auto mt-6">
        <input v-model="keyword" type="text" placeholder="搜索企业..." class="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none" />
      </div>
    </div>

    <div v-if="!data?.list?.length" class="text-center py-16 text-gray-400">暂无企业</div>
    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <NuxtLink v-for="ent in data.list" :key="ent.id" :to="`/enterprises/${ent.id}`" class="group bg-white hover:shadow-xl transition-all duration-300 p-6 border border-gray-100 hover:border-blue-200 rounded-2xl text-center">
        <div class="w-14 h-14 rounded-xl flex items-center justify-center mx-auto mb-3 shadow-sm group-hover:shadow-md group-hover:scale-105 transition-all duration-300 overflow-hidden bg-gray-100">
          <img v-if="ent.logoUrl" :src="ent.logoUrl" class="w-full h-full object-cover" />
          <span v-else class="text-white text-xl font-bold bg-gradient-to-br w-full h-full flex items-center justify-center" :class="avatarColor(ent.name)">{{ ent.name.charAt(0) }}</span>
        </div>
        <h3 class="font-semibold text-gray-900">{{ ent.name }}</h3>
        <div class="flex flex-wrap justify-center gap-1.5 mt-2">
          <span v-if="ent.industry" class="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full">{{ ent.industry.name }}</span>
          <span v-if="ent.city" class="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full">{{ ent.city }}</span>
        </div>
        <p class="text-xs text-gray-400 mt-2 line-clamp-2">{{ ent.description || '-' }}</p>
      </NuxtLink>
    </div>
  </div>
</template>
