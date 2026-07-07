<script setup lang="ts">
import { shallowRef, computed } from '#imports'
import type { Application, ApiResponse, PaginatedData } from '~/types'

definePageMeta({ middleware: ['auth'] })

const { token, user, logout } = useAuth()

const { data: apps, status, refresh, error: fetchError } = useAsyncData('my-apps', () =>
  $fetch<ApiResponse<PaginatedData<Application>>>('http://127.0.0.1:8080/api/v1/portal/student/applications', {
    headers: { Authorization: `Bearer ${token.value}` }
  }).then(r => r.code === 0 ? r.data : { list: [], total: 0 }).catch(() => ({ list: [], total: 0 })),
{ server: false }
)

// ── filter state ──
const activeFilter = shallowRef<'all' | 'pending' | 'interview' | 'offer'>('all')

// ── status definitions ──
const statusLabel: Record<string, string> = { pending: '待处理', viewed: '已查看', interview: '面试中', offer: '已发Offer', accepted: '已接受', rejected: '已拒绝' }

const statusBorder: Record<string, string> = { pending: 'border-l-blue-300', viewed: 'border-l-sky-300', interview: 'border-l-amber-300', offer: 'border-l-purple-300', accepted: 'border-l-emerald-300', rejected: 'border-l-gray-200' }
const statusDot: Record<string, string> = { pending: 'bg-blue-400 ring-blue-400/20', viewed: 'bg-sky-400 ring-sky-400/20', interview: 'bg-amber-400 ring-amber-400/20', offer: 'bg-purple-400 ring-purple-400/20', accepted: 'bg-emerald-400 ring-emerald-400/20', rejected: 'bg-gray-300 ring-gray-300/20' }
const statusTag: Record<string, string> = { pending: 'text-blue-500 bg-blue-50/60', viewed: 'text-sky-500 bg-sky-50/60', interview: 'text-amber-500 bg-amber-50/60', offer: 'text-purple-500 bg-purple-50/60', accepted: 'text-emerald-500 bg-emerald-50/60', rejected: 'text-gray-400 bg-gray-50/80' }

const filterTabs = [
  { key: 'all' as const, label: '全部' },
  { key: 'pending' as const, label: '待处理' },
  { key: 'interview' as const, label: '面试中' },
  { key: 'offer' as const, label: '已获Offer' },
]

// ── derived data ──
const list = computed<Application[]>(() => apps.value?.list || [])

const filteredList = computed(() => {
  if (activeFilter.value === 'all') return list.value
  if (activeFilter.value === 'pending') return list.value.filter(a => a.status === 'pending' || a.status === 'viewed')
  if (activeFilter.value === 'interview') return list.value.filter(a => a.status === 'interview')
  if (activeFilter.value === 'offer') return list.value.filter(a => a.status === 'offer' || a.status === 'accepted')
  return list.value
})

const stats = computed(() => {
  const t = list.value.length
  const pending = list.value.filter(a => a.status === 'pending').length
  const viewed = t - pending
  const weekAgo = Date.now() - 7 * 86400000
  const thisWeek = list.value.filter(a => new Date(a.createdAt).getTime() > weekAgo).length
  return {
    total: t,
    pending: pending + list.value.filter(a => a.status === 'viewed').length,
    interview: list.value.filter(a => a.status === 'interview').length,
    offer: list.value.filter(a => a.status === 'offer' || a.status === 'accepted').length,
    viewRate: t > 0 ? Math.round((viewed / t) * 100) : 0,
    thisWeek,
  }
})

const employmentStatus = computed(() => {
  if (stats.value.offer > 0) return 'offer' as const
  if (stats.value.interview > 0) return 'interview' as const
  if (stats.value.total > 0) return 'active' as const
  return 'seeking' as const
})

// ── relative time ──
function relativeTime(dateStr: string): string {
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diff = now - then
  const mins = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

// ── navigation ──
function goToFilter(key: typeof activeFilter.value) {
  activeFilter.value = key
}

function goToApp(app: Application) {
  if (app.position?.id) navigateTo(`/positions/${app.position.id}`)
}

function retry() {
  refresh()
}
</script>

<template>
  <div class="max-w-5xl mx-auto px-4 py-6">
    <!-- ═══ Profile Header ═══ -->
    <div class="bg-white rounded-xl border border-gray-100/60 p-5 mb-6">
      <div class="flex items-center gap-3">
        <div class="w-11 h-11 rounded-full bg-blue-500 flex items-center justify-center shrink-0">
          <span class="text-white text-sm font-semibold">{{ (user?.nickname || user?.username || '?').charAt(0) }}</span>
        </div>
        <h1 class="text-[15px] font-semibold text-gray-800">{{ user?.nickname || user?.username }}</h1>
        <span class="text-[11px] rounded-full px-2 py-0.5 font-medium"
          :class="employmentStatus === 'offer'
            ? 'text-emerald-600 bg-emerald-50'
            : employmentStatus === 'interview'
              ? 'text-amber-600 bg-amber-50'
              : employmentStatus === 'active'
                ? 'text-blue-500 bg-blue-50'
                : 'text-gray-400 bg-gray-100'"
        >{{ employmentStatus === 'offer' ? '已获Offer' : employmentStatus === 'interview' ? '面试中' : employmentStatus === 'active' ? '投递中' : '求职中' }}</span>
        <div class="flex-1" />
        <NuxtLink to="/student/profile" class="text-[12px] text-gray-400 hover:text-gray-600 transition-colors">编辑资料</NuxtLink>
        <NuxtLink to="/positions" class="text-[12px] text-gray-400 hover:text-gray-600 transition-colors">浏览职位</NuxtLink>
        <button @click="logout" class="text-[12px] text-gray-300 hover:text-gray-500 transition-colors">退出</button>
      </div>

      <!-- Stats -->
      <div class="grid grid-cols-6 mt-4 pt-4 border-t border-gray-100/60">
        <button @click="goToFilter('all')" class="text-center group">
          <div class="text-lg font-semibold tabular-nums" :class="activeFilter === 'all' ? 'text-gray-700' : 'text-gray-300 group-hover:text-gray-400 transition-colors'">{{ stats.total }}</div>
          <div class="text-[11px] mt-0.5" :class="activeFilter === 'all' ? 'text-gray-500 font-medium' : 'text-gray-400'">全部投递</div>
        </button>
        <button @click="goToFilter('pending')" class="text-center group">
          <div class="text-lg font-semibold tabular-nums" :class="stats.pending > 0 ? 'text-blue-500' : activeFilter === 'pending' ? 'text-gray-500' : 'text-gray-300 group-hover:text-gray-400 transition-colors'">{{ stats.pending }}</div>
          <div class="text-[11px] mt-0.5" :class="activeFilter === 'pending' ? 'text-blue-500 font-medium' : 'text-gray-400'">待处理</div>
        </button>
        <button disabled class="text-center">
          <div class="text-lg font-semibold tabular-nums" :class="stats.viewRate >= 80 ? 'text-emerald-500' : stats.viewRate >= 40 ? 'text-amber-500' : stats.total > 0 ? 'text-red-400' : 'text-gray-300'">{{ stats.total > 0 ? stats.viewRate + '%' : '-' }}</div>
          <div class="text-[11px] text-gray-400 mt-0.5">简历查看率</div>
        </button>
        <button disabled class="text-center">
          <div class="text-lg font-semibold tabular-nums" :class="stats.thisWeek > 0 ? 'text-blue-500' : 'text-gray-300'">{{ stats.thisWeek }}</div>
          <div class="text-[11px] text-gray-400 mt-0.5">本周投递</div>
        </button>
        <button @click="goToFilter('interview')" class="text-center group">
          <div class="text-lg font-semibold tabular-nums" :class="stats.interview > 0 ? 'text-amber-500' : activeFilter === 'interview' ? 'text-gray-500' : 'text-gray-300 group-hover:text-gray-400 transition-colors'">{{ stats.interview }}</div>
          <div class="text-[11px] mt-0.5" :class="activeFilter === 'interview' ? 'text-amber-500 font-medium' : 'text-gray-400'">面试中</div>
        </button>
        <button @click="goToFilter('offer')" class="text-center group">
          <div class="text-lg font-semibold tabular-nums" :class="stats.offer > 0 ? 'text-emerald-500' : activeFilter === 'offer' ? 'text-gray-500' : 'text-gray-300 group-hover:text-gray-400 transition-colors'">{{ stats.offer }}</div>
          <div class="text-[11px] mt-0.5" :class="activeFilter === 'offer' ? 'text-emerald-500 font-medium' : 'text-gray-400'">已获Offer</div>
        </button>
      </div>
    </div>

    <!-- ═══ Filter Tabs + Refresh ═══ -->
    <div class="flex items-center justify-between mb-3">
      <div class="flex items-center gap-0.5 bg-gray-100/60 rounded-lg p-0.5">
        <button
          v-for="tab in filterTabs" :key="tab.key"
          @click="goToFilter(tab.key)"
          class="text-[12px] px-3 py-1 rounded-md transition-all duration-200 font-medium"
          :class="activeFilter === tab.key ? 'bg-white text-gray-700 shadow-sm' : 'text-gray-400 hover:text-gray-500'"
        >
          {{ tab.label }}
          <span v-if="tab.key !== 'all'" class="ml-1 tabular-nums opacity-60">{{ stats[tab.key === 'pending' ? 'pending' : tab.key] }}</span>
        </button>
      </div>
      <button @click="retry" class="text-[12px] text-gray-300 hover:text-gray-500 transition-colors inline-flex items-center gap-1" :class="{ 'animate-spin': status === 'pending' }">
        <Icon name="heroicons:arrow-path" class="w-3.5 h-3.5" />
      </button>
    </div>

    <!-- ═══ Content ═══ -->

    <!-- Loading skeleton -->
    <div v-if="status === 'pending'" class="space-y-2">
      <div v-for="n in 4" :key="n" class="bg-white rounded-lg border border-gray-100/60 px-4 py-3.5">
        <div class="flex items-start gap-3">
          <div class="mt-1.5 w-2 h-2 rounded-full bg-gray-200 shrink-0" />
          <div class="flex-1 space-y-2.5 animate-pulse">
            <div class="h-4 bg-gray-100 rounded w-3/5" />
            <div class="h-3 bg-gray-50 rounded w-2/5" />
            <div class="h-3 bg-gray-50 rounded w-1/4" />
          </div>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="fetchError" class="bg-white rounded-xl border border-red-100 py-16 text-center">
      <div class="w-12 h-12 rounded-full bg-red-50 flex items-center justify-center mx-auto mb-3">
        <Icon name="heroicons:exclamation-triangle" class="w-6 h-6 text-red-300" />
      </div>
      <p class="text-sm text-gray-500 mb-4">加载失败，请重试</p>
      <button @click="retry" class="text-[13px] text-blue-500 hover:text-blue-600 transition-colors font-medium">重新加载</button>
    </div>

    <!-- Empty -->
    <div v-else-if="!list.length" class="bg-white rounded-xl border border-dashed border-gray-150 py-16 text-center">
      <div class="w-12 h-12 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-3">
        <Icon name="heroicons:inbox" class="w-6 h-6 text-gray-300" />
      </div>
      <p class="text-sm text-gray-400 mb-4">还没有投递记录</p>
      <NuxtLink to="/positions" class="inline-flex items-center gap-1.5 text-[13px] text-gray-500 hover:text-blue-600 transition-colors">
        <Icon name="heroicons:magnifying-glass" class="w-3.5 h-3.5 opacity-50" />浏览职位
      </NuxtLink>
    </div>

    <!-- Empty for current filter -->
    <div v-else-if="!filteredList.length" class="bg-white rounded-xl border border-gray-100/60 py-16 text-center">
      <p class="text-sm text-gray-400 mb-4">该分类下暂无投递</p>
      <button @click="goToFilter('all')" class="text-[13px] text-gray-500 hover:text-blue-600 transition-colors">查看全部</button>
    </div>

    <!-- Application list -->
    <div v-else class="space-y-2">
      <div
        v-for="(app, idx) in filteredList" :key="app.id"
        class="bg-white rounded-lg border border-gray-100/60 border-l-[3px] px-4 py-3.5 hover:shadow-sm hover:-translate-y-px transition-all duration-200 cursor-pointer animate-fadeIn"
        :class="statusBorder[app.status] || 'border-l-gray-200'"
        :style="{ animationDelay: `${idx * 40}ms` }"
        @click="goToApp(app)"
      >
        <div class="flex items-start gap-3">
          <div class="mt-1.5 w-2 h-2 rounded-full ring-2 shrink-0" :class="statusDot[app.status] || 'bg-gray-300 ring-gray-300/20'" />

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <h3 class="text-[14px] font-medium text-gray-700 truncate">{{ app.position?.title || '职位已删除' }}</h3>
                <p class="text-[12px] text-gray-400 mt-0.5">
                  {{ app.enterprise?.name || '-' }}
                  <span class="mx-1.5 text-gray-200/80">·</span>
                  {{ app.position?.city || '-' }}
                </p>
              </div>
              <span class="text-[11px] px-2 py-0.5 rounded-full font-medium shrink-0 leading-relaxed" :class="statusTag[app.status] || 'text-gray-400 bg-gray-50'">
                {{ statusLabel[app.status] || app.status }}
              </span>
            </div>

            <div class="flex items-center mt-2.5 text-[11px] text-gray-300/80">
              <Icon name="heroicons:clock" class="w-3 h-3 mr-1 opacity-30" />
              {{ relativeTime(app.createdAt) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fadeIn {
  animation: fadeIn 0.35s ease-out both;
}
</style>
