<script setup lang="ts">
import type { Enterprise, ApiResponse } from '~/types'

const route = useRoute()
const id = route.params.id as string

const { data: ent } = await useAsyncData(`ent-${id}`, () =>
  $fetch<ApiResponse<Enterprise>>('http://127.0.0.1:8080/api/v1/portal/enterprises/' + id).then(r => r.code === 0 ? r.data : null).catch(() => null)
)

const gradients = ['from-blue-500 to-blue-600','from-emerald-500 to-teal-600','from-violet-500 to-purple-600','from-amber-500 to-orange-600','from-rose-500 to-pink-600','from-cyan-500 to-sky-600','from-indigo-500 to-indigo-600','from-lime-500 to-green-600']
</script>

<template>
  <div v-if="ent" class="max-w-4xl mx-auto px-4 pt-8 pb-12 space-y-6">
    <div class="text-center">
      <div class="w-20 h-20 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-sm overflow-hidden bg-gray-100">
        <img v-if="ent.logoUrl" :src="ent.logoUrl" class="w-full h-full object-cover" />
        <span v-else class="text-white text-3xl font-bold bg-gradient-to-br w-full h-full flex items-center justify-center" :class="gradients[ent.name.charCodeAt(0) % gradients.length]">{{ ent.name.charAt(0) }}</span>
      </div>
      <h1 class="text-3xl font-bold text-gray-900">{{ ent.name }}</h1>
      <div class="flex flex-wrap justify-center gap-2 mt-3 text-sm text-gray-500">
        <span>{{ ent.industry?.name || '-' }}</span><span>·</span><span>{{ ent.city || '-' }}</span><span>·</span><span>{{ ent.scale || '-' }}</span>
      </div>
    </div>

    <div class="bg-white rounded-xl p-6 border border-gray-100">
      <h2 class="font-semibold text-lg mb-3">企业简介</h2>
      <p class="text-gray-600 text-sm whitespace-pre-wrap">{{ ent.description || '暂无' }}</p>
    </div>

    <div class="grid gap-4 sm:grid-cols-2">
      <div class="bg-white rounded-xl p-6 border border-gray-100">
        <h3 class="font-semibold mb-3">基本信息</h3>
        <div class="space-y-2 text-sm text-gray-600">
          <div>官网：{{ ent.website || '-' }}</div>
          <div>地址：{{ ent.address || '-' }}</div>
          <div>联系人：{{ ent.contactName || '-' }}</div>
          <div>电话：{{ ent.contactPhone || '-' }}</div>
          <div>邮箱：{{ ent.contactEmail || '-' }}</div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-6 border border-gray-100">
        <h3 class="font-semibold mb-3">企业规模</h3>
        <div class="space-y-2 text-sm text-gray-600">
          <div>规模：{{ ent.scale || '-' }}</div>
          <div>性质：{{ ent.nature || '-' }}</div>
          <div>行业：{{ ent.industry?.name || '-' }}</div>
          <div>城市：{{ ent.city || '-' }}</div>
          <div>状态：<span :class="ent.status === 1 ? 'text-green-600' : 'text-red-500'">{{ ent.status === 1 ? '正常' : '禁用' }}</span></div>
        </div>
      </div>
    </div>
  </div>
</template>
