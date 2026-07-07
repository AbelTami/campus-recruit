<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { BaseButton } from '@/components/Button'
import { ElCard, ElRow, ElCol } from 'element-plus'
import { ref, onMounted } from 'vue'
import request from '@/axios'

interface Report { reportTitle: string; totalStudents: number; employed: number; employRate: number; avgSalary: number | null; totalEnterprises: number; totalPositions: number; totalApplications: number; topIndustry: string; topCity: string; topSkill: string }

const report = ref<Report | null>(null)
const today = new Date().toLocaleDateString('zh-CN')

onMounted(async () => {
  const res = await request.get<Report>({ url: '/admin/analysis/report' })
  if (res?.data) report.value = res.data
})

function handlePrint() { window.print() }
</script>

<template>
  <div class="report-page p-16px h-full">
    <ContentWrap class="h-full">
      <div class="flex items-center justify-between mb-20px no-print">
        <span class="text-18px font-600">{{ report?.reportTitle || '就业质量报告' }}</span>
        <BaseButton type="primary" @click="handlePrint">打印报告</BaseButton>
      </div>

      <template v-if="report">
        <!-- 打印头部 -->
        <div class="print-only text-center mb-16px">
          <h1 class="text-26px font-700 mb-6px" style="color:#1a1a1a;">{{ report.reportTitle }}</h1>
          <p class="text-14px" style="color:#666;">生成日期：{{ today }}&nbsp;&nbsp;|&nbsp;&nbsp;数据来源：大学生就业需求分析系统&nbsp;&nbsp;|&nbsp;&nbsp;机密等级：内部</p>
          <hr style="border:1px solid #409eff;margin:14px 0;">
        </div>

        <!-- 一、核心指标 -->
        <h2 class="section-heading no-print">一、核心指标</h2>
        <h2 class="section-heading print-only">一、核心指标</h2>
        <el-row :gutter="16" class="mb-20px">
          <el-col v-for="c in [
            { l:'在校学生总数', v:report.totalStudents.toLocaleString()+' 人', clr:'#409eff', sub:'涵盖24个学院' },
            { l:'已就业人数', v:report.employed.toLocaleString()+' 人', clr:'#67c23a', sub:'就业率 '+report.employRate.toFixed(1)+'%' },
            { l:'平均薪资', v:report.avgSalary?'¥'+(report.avgSalary/1000).toFixed(1)+'K/月':'N/A', clr:'#e6a23c', sub:'毕业生起薪水平' },
            { l:'在招职位', v:report.totalPositions.toLocaleString()+' 个', clr:'#f56c6c', sub:'来自 '+report.totalEnterprises+' 家合作企业' },
          ]" :key="c.l" :xs="12" :sm="6">
            <el-card shadow="hover" class="report-card">
              <div class="text-13px text-gray-400 mb-4px">{{ c.l }}</div>
              <div class="text-26px font-700 mb-2px" :style="{ color: c.clr }">{{ c.v }}</div>
              <div class="text-12px text-gray-400">{{ c.sub }}</div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 二、就业市场分析 -->
        <h2 class="section-heading">二、就业市场分析</h2>
        <el-row :gutter="16" class="mb-20px">
          <el-col :xs="24" :md="8" v-for="item in [
            { l:'热门需求行业', v:report.topIndustry, icon:'🏭', desc:'该行业在当前招聘市场中职位数量最多，是毕业生就业的主要方向。' },
            { l:'热门就业城市', v:report.topCity, icon:'🏙️', desc:'该城市提供的就业岗位最为集中，薪资水平与生活成本需综合考量。' },
            { l:'企业最需技能', v:report.topSkill, icon:'💡', desc:'掌握该技能的学生在求职中具有显著竞争优势，建议纳入重点培养计划。' },
          ]" :key="item.l" class="mb-16px">
            <el-card shadow="never" class="report-card">
              <div class="flex items-center gap-8px mb-8px">
                <span class="text-20px">{{ item.icon }}</span>
                <span class="text-14px text-gray-500">{{ item.l }}</span>
              </div>
              <span class="report-highlight">{{ item.v }}</span>
              <p class="text-13px text-gray-500 leading-20px">{{ item.desc }}</p>
            </el-card>
          </el-col>
        </el-row>

        <!-- 三、业务运营数据 -->
        <h2 class="section-heading">三、业务运营数据</h2>
        <div class="stats-table mb-20px">
          <div class="stats-row">
            <div class="stats-cell label">合作企业总数</div>
            <div class="stats-cell value">{{ report.totalEnterprises }} 家</div>
            <div class="stats-cell label">在招职位总数</div>
            <div class="stats-cell value">{{ report.totalPositions }} 个</div>
          </div>
          <div class="stats-row">
            <div class="stats-cell label">累计投递次数</div>
            <div class="stats-cell value">{{ report.totalApplications }} 次</div>
            <div class="stats-cell label">毕业生就业率</div>
            <div class="stats-cell value">{{ report.employRate.toFixed(1) }}%</div>
          </div>
          <div class="stats-row">
            <div class="stats-cell label">覆盖学院数</div>
            <div class="stats-cell value">24 个</div>
            <div class="stats-cell label">人均投递次数</div>
            <div class="stats-cell value">{{ report.totalStudents > 0 ? (report.totalApplications / report.totalStudents).toFixed(1) : '0' }} 次</div>
          </div>
        </div>

        <!-- 四、小结 -->
        <h2 class="section-heading">四、总结与建议</h2>
        <div class="summary-box p-16px rounded-8px mb-20px" style="background:#f5f7fa;line-height:1.8;">
          <p class="text-14px mb-8px">本报告基于系统当前数据自动生成。截至 {{ today }}，系统共管理 <strong>{{ report.totalStudents }}</strong> 名在校学生，覆盖 <strong>24</strong> 个学院。</p>
          <p class="text-14px mb-8px">当前整体就业率为 <strong>{{ report.employRate.toFixed(1) }}%</strong>，热门就业方向集中在 <strong>{{ report.topIndustry }}</strong> 行业，<strong>{{ report.topCity }}</strong> 是毕业生首选就业城市。</p>
          <p class="text-14px mb-8px">企业端最紧缺的技能为 <strong>{{ report.topSkill }}</strong>，建议就业指导中心在后续课程和培训中加强该技能的培养，提升学生就业竞争力。</p>
          <p class="text-14px">系统目前与 <strong>{{ report.totalEnterprises }}</strong> 家企业建立合作关系，在招职位 <strong>{{ report.totalPositions }}</strong> 个，累计完成 <strong>{{ report.totalApplications }}</strong> 次简历投递。</p>
        </div>

        <!-- 打印尾部 -->
        <div class="print-only text-center text-12px pt-12px" style="color:#999;border-top:1px solid #e5e5e5;">
          本报告由大学生就业需求分析系统自动生成&nbsp;&nbsp;|&nbsp;&nbsp;{{ today }}&nbsp;&nbsp;|&nbsp;&nbsp;第 1 页
        </div>
      </template>
    </ContentWrap>
  </div>
</template>

<style scoped>
.section-heading {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  border-left: 4px solid #409eff;
  padding-left: 10px;
  margin-bottom: 14px;
  margin-top: 8px;
}
.report-card :deep(.el-card__body) { padding: 16px 20px; }

.stats-table { border: 1px solid #ebeef5; border-radius: 6px; overflow: hidden; }
.stats-row { display: flex; border-bottom: 1px solid #ebeef5; }
.stats-row:last-child { border-bottom: none; }
.stats-cell { padding: 12px 16px; font-size: 14px; }
.stats-cell.label { flex: 1; background: #fafafa; color: #666; }
.stats-cell.value { flex: 1; font-weight: 600; color: #333; }
.report-highlight { display: inline-block; padding: 5px 12px; background: #ecf5ff; color: #409eff; border-radius: 4px; font-size: 13px; font-weight: 600; white-space: nowrap; margin-bottom: 8px; }

.print-only { display: none; }
.no-print { }

@media print {
  .no-print { display: none !important; }
  .print-only { display: block !important; }
  .report-page { padding: 0 !important; }
  .section-heading { border-left-color: #333 !important; color: #111 !important; break-after: avoid; }
  .summary-box { background: #fff !important; border: 1px solid #ddd !important; break-inside: avoid; }
  .stats-table { border-color: #ddd !important; break-inside: avoid; }
  .stats-row { border-color: #ddd !important; }
  .stats-cell.label { background: #f9f9f9 !important; }
  .report-card { break-inside: avoid; }
  .el-row { break-inside: avoid; }
  h2 { break-after: avoid; }
}
</style>
