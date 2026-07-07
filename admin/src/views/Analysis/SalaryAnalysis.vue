<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { Echart } from '@/components/Echart'
import { ref, onMounted } from 'vue'
import request from '@/axios'
import type { EChartsOption } from 'echarts'

interface SalaryRow { name: string; avg: number }

const barOption = ref<EChartsOption>({})

onMounted(async () => {
  const res = await request.get<SalaryRow[]>({ url: '/admin/dashboard/salary-analysis' })
  if (!res?.data) return
  const data = res.data
  barOption.value = {
    color: ['#67c23a'],
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, valueFormatter: (v: any) => `¥${v}K` },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map((r: SalaryRow) => r.name), axisLabel: { color: '#909399', rotate: 30 } },
    yAxis: { type: 'value', name: '平均月薪(K)', axisLabel: { color: '#909399' }, splitLine: { lineStyle: { type: 'dashed', color: '#ebeef5' } } },
    series: [{ type: 'bar', data: data.map((r: SalaryRow) => r.avg), barWidth: '50%', itemStyle: { borderRadius: [6, 6, 0, 0] },
      label: { show: true, position: 'top', formatter: '¥{c}K', color: '#606266' } }],
  }
})
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap title="各行业平均薪资" class="h-full">
      <Echart :options="barOption" height="420px" />
    </ContentWrap>
  </div>
</template>
