<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { BaseButton } from '@/components/Button'
import { ElTable, ElTableColumn, ElTag, ElSelect, ElOption } from 'element-plus'
import { ref, onMounted } from 'vue'
import request from '@/axios'

interface Student { id: number; name: string; studentNo: string }
interface MatchResult { id: number; title: string; enterprise: string; city: string; salary: string; matchRate: number; matchSkill: number }

const students = ref<Student[]>([])
const selectedStudent = ref<number | null>(null)
const results = ref<MatchResult[]>([])
const loading = ref(false)

onMounted(async () => {
  const res = await request.get<{ list: Student[] }>({ url: '/admin/students', params: { pageSize: 100 } })
  if (res?.data) students.value = res.data.list || []
})

async function doMatch() {
  if (!selectedStudent.value) return
  loading.value = true
  const res = await request.get<MatchResult[]>({ url: '/admin/analysis/match-recommend', params: { studentId: selectedStudent.value } })
  if (res?.data) results.value = res.data
  loading.value = false
}
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap class="h-full">
      <div class="text-18px font-600 mb-16px">人岗智能匹配</div>
      <div class="flex items-center gap-12px mb-20px">
        <span class="text-14px">选择学生：</span>
        <ElSelect v-model="selectedStudent" placeholder="请选择学生" filterable style="width:280px" @change="doMatch">
          <ElOption v-for="s in students" :key="s.id" :label="`${s.name} (${s.studentNo})`" :value="s.id" />
        </ElSelect>
      </div>

      <ElTable v-if="results.length" :data="results" v-loading="loading" border stripe>
        <ElTableColumn type="index" label="#" width="50" />
        <ElTableColumn prop="title" label="推荐职位" min-width="180" />
        <ElTableColumn prop="enterprise" label="企业" width="140" />
        <ElTableColumn prop="city" label="城市" width="80" />
        <ElTableColumn prop="salary" label="薪资" width="120" />
        <ElTableColumn label="匹配度" width="120">
          <template #default="{ row }">
            <div class="flex items-center gap-8px">
              <div class="flex-1 h-8px bg-gray-200 rounded-full overflow-hidden">
                <div class="h-full rounded-full transition-all" :style="{ width: row.matchRate+'%', background: row.matchRate>=60?'#67c23a':row.matchRate>=40?'#e6a23c':'#f56c6c' }" />
              </div>
              <span class="text-13px font-600">{{ row.matchRate }}%</span>
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn label="匹配技能" width="90">
          <template #default="{ row }"><ElTag size="small" type="success">{{ row.matchSkill }}</ElTag></template>
        </ElTableColumn>
      </ElTable>
      <div v-else-if="selectedStudent && !loading" class="text-center text-gray-400 py-40px">未找到匹配职位</div>
    </ContentWrap>
  </div>
</template>
