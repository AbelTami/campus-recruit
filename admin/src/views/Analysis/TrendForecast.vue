<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { Echart } from '@/components/Echart'
import { ref, onMounted } from 'vue'
import request from '@/axios'
import type { EChartsOption } from 'echarts'

interface TrendRow { year: string; rate: number }

const lineOption = ref<EChartsOption>({})

onMounted(async () => {
  const res = await request.get<TrendRow[]>({ url: '/admin/analysis/trend-forecast' })
  if (!res?.data) return
  const data = res.data
  const isForecast = data.map((_: TrendRow, i: number) => i >= data.length - 2)

  lineOption.value = {
    tooltip: { trigger: 'axis', valueFormatter: (v: any) => `${v}%` },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map((r: TrendRow) => r.year), axisLabel: { color: '#909399' } },
    yAxis: { type: 'value', name: '就业率(%)', axisLabel: { color: '#909399' }, splitLine: { lineStyle: { type: 'dashed', color: '#ebeef5' } } },
    series: [
      {
        type: 'line', name: '历史数据',
        data: data.filter((_: TrendRow, i: number) => !isForecast[i]).map((r: TrendRow) => r.rate),
        smooth: true, symbol: 'circle', symbolSize: 8, color: '#409eff',
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(64,158,255,0.3)' }, { offset: 1, color: 'rgba(64,158,255,0.02)' }] } },
      },
      {
        type: 'line', name: '预测值',
        data: [...Array(data.length - 2).fill(null), ...data.filter((_: TrendRow, i: number) => isForecast[i]).map((r: TrendRow) => r.rate)],
        smooth: true, symbol: 'diamond', symbolSize: 10, color: '#e6a23c',
        lineStyle: { type: 'dashed', width: 2 },
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(230,162,60,0.2)' }, { offset: 1, color: 'rgba(230,162,60,0.02)' }] } },
      },
    ],
  }
})
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap title="就业率趋势预测" class="h-full">
      <Echart :options="lineOption" height="420px" />
    </ContentWrap>
  </div>
</template>
