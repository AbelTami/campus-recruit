<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { Echart } from '@/components/Echart'
import { ElRow, ElCol } from 'element-plus'
import { ref, onMounted } from 'vue'
import request from '@/axios'
import type { EChartsOption } from 'echarts'

interface DemandRow { name: string; value: number }

const barOption = ref<EChartsOption>({})
const pieOption = ref<EChartsOption>({})

onMounted(async () => {
  const [cityRes, indRes] = await Promise.all([
    request.get<DemandRow[]>({ url: '/admin/analysis/city-demand' }),
    request.get<DemandRow[]>({ url: '/admin/dashboard/industry-dist' }),
  ])

  if (cityRes?.data) {
    barOption.value = {
      color: ['#409eff'],
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: cityRes.data.map((r: DemandRow) => r.name), axisLabel: { color: '#909399' } },
      yAxis: { type: 'value', name: '职位数', axisLabel: { color: '#909399' }, splitLine: { lineStyle: { type: 'dashed', color: '#ebeef5' } } },
      series: [{ type: 'bar', data: cityRes.data.map((r: DemandRow) => r.value), barWidth: '50%', itemStyle: { borderRadius: [6, 6, 0, 0] } }],
    }
  }

  if (indRes?.data) {
    pieOption.value = {
      color: ['#409eff','#67c23a','#e6a23c','#f56c6c','#909399','#ff9800','#8bc34a','#00bcd4','#9c27b0','#3f51b5'],
      tooltip: { trigger: 'item' },
      legend: { bottom: 0, textStyle: { color: '#909399', fontSize: 11 } },
      series: [{ type: 'pie', radius: ['45%', '72%'], center: ['50%', '43%'], label: { show: false }, emphasis: { label: { show: true } },
        data: indRes.data.slice(0, 10).map((r: DemandRow) => ({ name: r.name, value: r.value })) }],
    }
  }
})
</script>

<template>
  <div class="p-16px h-full">
    <el-row :gutter="16" class="h-full">
      <el-col :xs="24" :lg="14">
        <ContentWrap title="城市职位需求排行" class="h-full">
          <Echart :options="barOption" height="400px" />
        </ContentWrap>
      </el-col>
      <el-col :xs="24" :lg="10">
        <ContentWrap title="行业需求占比" class="h-full">
          <Echart :options="pieOption" height="400px" />
        </ContentWrap>
      </el-col>
    </el-row>
  </div>
</template>
