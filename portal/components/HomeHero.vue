<script setup lang="ts">
const props = defineProps<{ searchQuery: string }>()

const emit = defineEmits<{
  (e: 'update:searchQuery', value: string): void
}>()
</script>

<template>
  <div class="relative -mt-16 min-h-[80vh] flex items-center justify-center overflow-hidden bg-gradient-to-br from-blue-600 via-blue-700 to-indigo-900">
    <div class="absolute inset-0 bg-grid-white opacity-10"></div>
    <div class="relative z-10 w-full max-w-4xl mx-auto px-4 text-center space-y-8 py-20">
      <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-white/20 text-white">
        🚀 实时数据统计
      </span>
      <h1 class="text-4xl md:text-6xl font-bold text-white leading-tight">
        发现你的
        <span class="text-blue-200">理想职业</span>
      </h1>
      <p class="text-xl text-blue-100 max-w-xl mx-auto">
        汇聚 30+ 优质企业，为大学生精准匹配就业机会
      </p>

      <div class="relative max-w-xl mx-auto">
        <div class="absolute inset-0 bg-blue-400 rounded-2xl blur-xl opacity-30"></div>
        <div class="relative bg-white rounded-2xl shadow-2xl p-2">
          <div class="flex gap-2">
            <div class="flex-1 relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-lg">🔍</span>
              <input
                :value="searchQuery"
                @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
                type="text"
                placeholder="搜索职位、企业..."
                class="w-full pl-12 pr-4 py-4 rounded-xl border-0 ring-1 ring-gray-100 focus:ring-2 focus:ring-blue-500 text-gray-900 placeholder:text-gray-400 outline-none"
              />
            </div>
            <button class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-4 rounded-xl font-semibold transition-all hover:shadow-lg">
              搜索
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
