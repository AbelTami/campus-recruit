<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { Echart } from '@/components/Echart'
import { Icon } from '@/components/Icon'
import { ElRow, ElCol, ElCard } from 'element-plus'
import { ref, onMounted } from 'vue'
import request from '@/axios'
import type { EChartsOption } from 'echarts'

interface Overview { totalStudents: number; employedCount: number; employmentRate: number; avgSalary: number | null; totalEnterprises: number; totalPositions: number }
interface IndustryRow { name: string; value: number }
interface TrendRow { year: string; rate: number }

const cards = ref([
  { label: '在校学生', value: '-', icon: 'ri:user-line', color: '#409eff' },
  { label: '已就业', value: '-', icon: 'ri:briefcase-line', color: '#67c23a' },
  { label: '就业率', value: '-', icon: 'ri:line-chart-line', color: '#e6a23c' },
  { label: '在招职位', value: '-', icon: 'ri:money-cny-circle-line', color: '#f56c6c' },
])

const barOption = ref<EChartsOption>({})
const pieOption = ref<EChartsOption>({})
const lineOption = ref<EChartsOption>({})

async function fetchAll() {
  const [ov, ind, trend] = await Promise.all([
    request.get<Overview>({ url: '/admin/dashboard/overview' }),
    request.get<IndustryRow[]>({ url: '/admin/dashboard/industry-dist' }),
    request.get<TrendRow[]>({ url: '/admin/dashboard/employment-trend' }),
  ])

  if (ov?.data) {
    const d = ov.data
    cards.value = [
      { label: '在校学生', value: d.totalStudents.toLocaleString(), icon: 'ri:user-line', color: '#409eff' },
      { label: '已就业', value: d.employedCount.toLocaleString(), icon: 'ri:briefcase-line', color: '#67c23a' },
      { label: '就业率', value: `${d.employmentRate.toFixed(1)}%`, icon: 'ri:line-chart-line', color: '#e6a23c' },
      { label: '在招职位', value: d.totalPositions.toLocaleString(), icon: 'ri:money-cny-circle-line', color: '#f56c6c' },
    ]
  }

  if (ind?.data) {
    barOption.value = {
      color: ['#409eff'],
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: ind.data.map((r: IndustryRow) => r.name), axisLabel: { color: '#909399', rotate: 30 } },
      yAxis: { type: 'value', name: '职位数', axisLabel: { color: '#909399' }, splitLine: { lineStyle: { type: 'dashed', color: '#ebeef5' } } },
      series: [{ type: 'bar', data: ind.data.map((r: IndustryRow) => r.value), barWidth: '50%', itemStyle: { borderRadius: [6, 6, 0, 0] } }],
    }
  }

  if (trend?.data) {
    lineOption.value = {
      color: ['#409eff'],
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: trend.data.map((r: TrendRow) => r.year), boundaryGap: false, axisLabel: { color: '#909399' } },
      yAxis: { type: 'value', name: '就业率(%)', axisLabel: { color: '#909399' }, splitLine: { lineStyle: { type: 'dashed', color: '#ebeef5' } } },
      series: [{
        type: 'line', data: trend.data.map((r: TrendRow) => r.rate), smooth: true, symbol: 'circle', symbolSize: 8,
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(64,158,255,0.35)' }, { offset: 1, color: 'rgba(64,158,255,0.02)' }] } },
      }],
    }
  }
}

onMounted(() => fetchAll())
</script>

<template>
  <div class="p-16px">
    <el-row :gutter="16" class="mb-16px">
      <el-col v-for="c in cards" :key="c.label" :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-13px text-gray-400 mb-6px">{{ c.label }}</div>
              <div class="text-26px font-600" :style="{ color: c.color }">{{ c.value }}</div>
            </div>
            <div class="w-48px h-48px rounded-xl flex items-center justify-center" :style="{ backgroundColor: c.color + '14' }">
              <Icon :icon="c.icon" :color="c.color" :size="24" />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="mb-16px">
      <el-col :xs="24" :lg="15">
        <ContentWrap title="各行业职位分布">
          <Echart :options="barOption" height="300px" />
        </ContentWrap>
      </el-col>
      <el-col :xs="24" :lg="9">
        <ContentWrap title="就业率趋势">
          <Echart :options="lineOption" height="300px" />
        </ContentWrap>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.stat-card :deep(.el-card__body) { padding: 18px 20px; }
</style>
