<script setup lang="ts">
import { PropType, computed } from 'vue'
import { DepartmentUserItem } from '@/api/department/types'
import { Descriptions } from '@/components/Descriptions'

const props = defineProps({
  currentRow: { type: Object as PropType<DepartmentUserItem>, default: () => undefined },
  detailSchema: { type: Array as PropType<any[]>, default: () => [] },
})

const name = computed(() => props.currentRow?.username || '-')
const avatar = computed(() => name.value.charAt(0))
const roleName = computed(() => props.currentRow?.role || '-')
const dept = computed(() => props.currentRow?.department?.label || '-')
</script>

<template>
  <div v-if="currentRow" class="flex flex-col gap-5">
    <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
      <div class="w-[72px] h-[72px] rounded-full bg-blue-500 flex items-center justify-center shrink-0">
        <span class="text-white text-2xl font-semibold">{{ avatar }}</span>
      </div>
      <div>
        <h3 class="text-[15px] font-bold text-gray-900">{{ name }}</h3>
        <p class="text-[13px] text-gray-500 mt-0.5">{{ currentRow.username }}</p>
        <div class="flex items-center gap-2 mt-1.5">
          <span class="text-[11px] bg-blue-50 text-blue-600 px-2 py-0.5 rounded-full font-medium">{{ dept }}</span>
          <span class="text-[11px] bg-emerald-50 text-emerald-600 px-2 py-0.5 rounded-full font-medium">{{ roleName }}</span>
        </div>
      </div>
    </div>
    <Descriptions :schema="detailSchema" :data="currentRow" :column="2" />
  </div>
</template>
