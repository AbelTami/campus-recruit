<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { Echart } from '@/components/Echart'
import { ref, onMounted } from 'vue'
import request from '@/axios'
import type { EChartsOption } from 'echarts'

interface CollegeRow { name: string; rate: number; total: number }

const barOption = ref<EChartsOption>({})

onMounted(async () => {
  const res = await request.get<CollegeRow[]>({ url: '/admin/dashboard/college-employment' })
  if (!res?.data) return
  const data = res.data
  const colors = data.map((_: CollegeRow) => _.rate >= 50 ? '#67c23a' : _.rate >= 30 ? '#e6a23c' : '#f56c6c')
  barOption.value = {
    color: colors,
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, valueFormatter: (v: any) => `${v}% (${data[v.dataIndex]?.total || 0}人)` },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map((r: CollegeRow) => r.name), axisLabel: { color: '#909399', rotate: 30 } },
    yAxis: { type: 'value', name: '就业率(%)', max: 100, axisLabel: { color: '#909399' }, splitLine: { lineStyle: { type: 'dashed', color: '#ebeef5' } } },
    series: [{ type: 'bar', data: data.map((r: CollegeRow) => r.rate), barWidth: '50%', itemStyle: { borderRadius: [6, 6, 0, 0] },
      label: { show: true, position: 'top', formatter: '{c}%', color: '#606266' } }],
  }
})
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap title="各学院就业率" class="h-full">
      <Echart :options="barOption" height="420px" />
    </ContentWrap>
  </div>
</template>
