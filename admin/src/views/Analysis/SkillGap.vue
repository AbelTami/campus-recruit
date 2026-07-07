<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { Echart } from '@/components/Echart'
import { ref, onMounted } from 'vue'
import request from '@/axios'
import type { EChartsOption } from 'echarts'

interface GapRow { name: string; student: number; job: number }

const radarOption = ref<EChartsOption>({})

onMounted(async () => {
  const res = await request.get<GapRow[]>({ url: '/admin/analysis/skill-gap' })
  if (!res?.data) return
  const data = res.data.slice(0, 12)
  radarOption.value = {
    color: ['#409eff', '#e6a23c'],
    tooltip: {},
    legend: { data: ['学生拥有', '职位需求'], bottom: 0, textStyle: { color: '#909399' } },
    radar: {
      center: ['50%', '45%'],
      radius: '60%',
      indicator: data.map((r: GapRow) => ({ name: r.name, max: Math.max(r.student, r.job) * 1.3 })),
      axisName: { color: '#909399', fontSize: 11 },
    },
    series: [{
      type: 'radar',
      data: [
        { value: data.map((r: GapRow) => r.student), name: '学生拥有', symbol: 'circle', symbolSize: 4 },
        { value: data.map((r: GapRow) => r.job), name: '职位需求', symbol: 'circle', symbolSize: 4 },
      ],
    }],
  }
})
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap title="技能供需差距分析" class="h-full">
      <Echart :options="radarOption" height="460px" />
    </ContentWrap>
  </div>
</template>
