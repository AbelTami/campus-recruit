<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { Table } from '@/components/Table'
import { Search } from '@/components/Search'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { ElInput, ElSelect, ElOption, ElTag } from 'element-plus'
import { reactive, ref, unref } from 'vue'
import { getPositionList, createPosition, updatePosition, batchDeletePositions } from '@/api/position'
import type { Position } from '@/api/position/types'
import { useTable } from '@/hooks/web/useTable'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { pageSize, currentPage } = tableState
    const res = await getPositionList({ pageIndex: unref(currentPage), pageSize: unref(pageSize), ...unref(searchParams) })
    return { list: res.data.list || [], total: res.data.total || 0 }
  },
  fetchDelApi: async () => { const res = await batchDeletePositions(unref(ids).map(Number)); return !!res },
})

const { total, loading, dataList, pageSize, currentPage } = tableState
const { getList, getElTableExpose, delList } = tableMethods

const eduOptions = ['不限', '专科', '本科', '硕士', '博士']
const statusOptions = [{ label: '招聘中', value: 1 }, { label: '已下架', value: 0 }]

const crudSchemas = reactive<CrudSchema[]>([
  { field: 'selection', search: { hidden: true }, form: { hidden: true }, table: { type: 'selection' } },
  { field: 'index', label: '#', form: { hidden: true }, search: { hidden: true }, table: { type: 'index', width: 50 } },
  { field: 'keyword', label: '搜索', search: { component: 'Input', componentProps: { placeholder: '职位名称' } }, form: { hidden: true }, table: { hidden: true } },
  { field: 'title', label: '职位名称', minWidth: 180, search: { hidden: true } },
  { field: 'enterprise', label: '企业', minWidth: 140, search: { hidden: true }, table: { slots: { default: (d: any) => <span class="text-13px">{d.row.enterprise?.shortName || d.row.enterprise?.name || '-'}</span> } } },
  { field: 'city', label: '城市', width: 70, search: { hidden: true } },
  { field: 'educationRequirement', label: '学历', width: 65, search: { hidden: true }, table: { slots: { default: (d: any) => <span class="text-13px">{d.row.educationRequirement || '不限'}</span> } } },
  { field: 'salary', label: '薪资', width: 130, search: { hidden: true }, table: { slots: { default: (d: any) => <span class="font-500 text-[var(--el-color-primary)]">{d.row.salaryMin && d.row.salaryMax ? `¥${(d.row.salaryMin/1000).toFixed(0)}-${(d.row.salaryMax/1000).toFixed(0)}K` : '面议'}</span> } } },
  { field: 'headcount', label: '人数', width: 55, align: 'center', search: { hidden: true } },
  { field: 'status', label: '状态', width: 80, search: { component: 'Select', componentProps: { options: statusOptions, placeholder: '状态' } }, table: { slots: { default: (d: any) => d.row.status === 1 ? <ElTag type="success" size="small">招聘中</ElTag> : <ElTag type="info" size="small">已下架</ElTag> } } },
  {
    field: 'action', label: '操作', width: 220, form: { hidden: true }, search: { hidden: true },
    table: { slots: { default: (d: any) => (
          <div class="flex gap-6px">
            <BaseButton size="small" plain onClick={() => viewDetail(d.row)}>详情</BaseButton>
            <BaseButton size="small" type="primary" onClick={() => openEdit(d.row)}>编辑</BaseButton>
            <BaseButton type="danger" size="small" plain onClick={() => delData(d.row)}>删除</BaseButton>
          </div>
        ) } },
  },
])

const { allSchemas } = useCrudSchemas(crudSchemas)
const searchParams = ref({})
const setSearchParams = (p: any) => { currentPage.value = 1; searchParams.value = p; getList() }

const ids = ref<string[]>([])
const delLoading = ref(false)
const delData = async (row?: Position) => {
  const el = await getElTableExpose()
  ids.value = row ? [String(row.id)] : el?.getSelectionRows().map((v: Position) => String(v.id)) || []
  delLoading.value = true
  await delList(unref(ids).length).finally(() => { delLoading.value = false })
}

// ── detail dialog ──
const detailVisible = ref(false)
const currentPosition = ref<Position>()
const viewDetail = (row: Position) => { currentPosition.value = row; detailVisible.value = true }

// ── edit dialog ──
const editVisible = ref(false)
const editLoading = ref(false)
const editForm = reactive<Record<string, any>>({})
const editEdu = ref('')
const editStatus = ref(1)
const isCreate = ref(false)

function openCreate() {
  Object.assign(editForm, { title: '', city: '', experienceRequirement: null, salaryMin: null, salaryMax: null, headcount: 1, description: '', requirement: '', welfare: '' })
  editEdu.value = '本科'
  editStatus.value = 1
  isCreate.value = true
  editVisible.value = true
}

function openEdit(row: Position) {
  isCreate.value = false
  Object.assign(editForm, {
    title: row.title, city: row.city || '',
    experienceRequirement: row.experienceRequirement, salaryMin: row.salaryMin, salaryMax: row.salaryMax,
    headcount: row.headcount, description: row.description || '',
    requirement: row.requirement || '', welfare: row.welfare || '',
  })
  editEdu.value = row.educationRequirement || '不限'
  editStatus.value = row.status
  currentPosition.value = row
  editVisible.value = true
}

async function handleSave() {
  editLoading.value = true
  try {
    const data = { ...editForm, educationRequirement: editEdu.value, status: editStatus.value }
    if (isCreate.value) await createPosition(data)
    else await updatePosition(currentPosition.value!.id, data)
    editVisible.value = false
    getList()
  } finally { editLoading.value = false }
}
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap class="h-full">
      <div class="flex items-center justify-between mb-20px">
        <div><span class="text-18px font-600">职位管理</span><span class="text-13px text-gray-400 ml-12px">共 {{ total }} 个</span></div>
        <div class="flex gap-2"><BaseButton type="primary" @click="openCreate">新增职位</BaseButton><BaseButton :loading="delLoading" type="danger" plain @click="delData()">批量删除</BaseButton></div>
      </div>
      <Search :schema="allSchemas.searchSchema" layout="inline" label-width="70px" @reset="setSearchParams" @search="setSearchParams" />
      <div class="mt-16px">
        <Table v-model:current-page="currentPage" v-model:page-size="pageSize" :columns="allSchemas.tableColumns" :data="dataList" :loading="loading" :show-overflow-tooltip="false" @register="tableRegister" :pagination="{ total }" />
      </div>
    </ContentWrap>

    <!-- Detail dialog -->
    <Dialog v-model="detailVisible" title="职位详情" width="600px">
      <div v-if="currentPosition" class="flex flex-col gap-5">
        <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
          <div class="w-[72px] h-[72px] rounded-xl bg-white flex items-center justify-center shrink-0 ring-1 ring-black/5">
            <span class="text-[28px]">💼</span>
          </div>
          <div>
            <h3 class="text-[15px] font-bold text-gray-900">{{ currentPosition.title }}</h3>
            <p class="text-[13px] text-gray-500 mt-0.5">{{ currentPosition.enterprise?.name || '-' }} · {{ currentPosition.city || '-' }}</p>
            <div class="flex items-center gap-2 mt-1.5">
              <span class="text-[11px]" :class="currentPosition.status === 1 ? 'text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full font-medium' : 'text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full font-medium'">{{ currentPosition.status === 1 ? '招聘中' : '已下架' }}</span>
              <span class="text-[11px] text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full font-medium">{{ currentPosition.salaryMin ? `¥${(currentPosition.salaryMin/1000).toFixed(0)}K-${(currentPosition.salaryMax!/1000).toFixed(0)}K` : '面议' }}</span>
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-blue-500" /><h3 class="text-[13px] font-semibold text-gray-700">基本信息</h3></div>
          <div class="grid grid-cols-2 gap-x-4 gap-y-2.5 text-[13px]">
            <div><span class="text-gray-400">企业</span><p class="text-gray-700 mt-0.5">{{ currentPosition.enterprise?.name || '-' }}</p></div>
            <div><span class="text-gray-400">城市</span><p class="text-gray-700 mt-0.5">{{ currentPosition.city || '-' }}</p></div>
            <div><span class="text-gray-400">学历要求</span><p class="text-gray-700 mt-0.5">{{ currentPosition.educationRequirement || '不限' }}</p></div>
            <div><span class="text-gray-400">经验要求</span><p class="text-gray-700 mt-0.5">{{ currentPosition.experienceRequirement ? currentPosition.experienceRequirement + '年' : '不限' }}</p></div>
            <div><span class="text-gray-400">招聘人数</span><p class="text-gray-700 mt-0.5">{{ currentPosition.headcount }} 人</p></div>
            <div><span class="text-gray-400">行业</span><p class="text-gray-700 mt-0.5">{{ currentPosition.industry?.name || '-' }}</p></div>
          </div>
        </div>

        <div v-if="currentPosition.description">
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-emerald-500" /><h3 class="text-[13px] font-semibold text-gray-700">职位描述</h3></div>
          <p class="text-[13px] text-gray-600 leading-relaxed whitespace-pre-wrap">{{ currentPosition.description }}</p>
        </div>
        <div v-if="currentPosition.requirement">
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-amber-500" /><h3 class="text-[13px] font-semibold text-gray-700">任职要求</h3></div>
          <p class="text-[13px] text-gray-600 leading-relaxed whitespace-pre-wrap">{{ currentPosition.requirement }}</p>
        </div>
        <div v-if="currentPosition.welfare">
          <div class="flex items-center gap-2 mb-3"><span class="w-1 h-3.5 rounded-full bg-purple-500" /><h3 class="text-[13px] font-semibold text-gray-700">福利待遇</h3></div>
          <p class="text-[13px] text-gray-600 leading-relaxed">{{ currentPosition.welfare }}</p>
        </div>
      </div>
      <template #footer><BaseButton @click="detailVisible = false">关闭</BaseButton></template>
    </Dialog>

    <!-- Edit dialog -->
    <Dialog v-model="editVisible" :title="isCreate ? '新增职位' : '编辑职位'" width="600px" :close-on-click-modal="false">
      <div class="flex flex-col gap-5">
        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-blue-500" /><h3 class="text-[13px] font-semibold text-gray-700">基本信息</h3></div>
          <div class="space-y-3">
            <div><label class="text-[12px] text-gray-500 mb-1.5 block">职位名称 <span class="text-red-400">*</span></label><ElInput v-model="editForm.title" placeholder="请输入职位名称" size="large" /></div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">城市</label><ElInput v-model="editForm.city" placeholder="例如：深圳" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">学历要求</label><ElSelect v-model="editEdu" placeholder="选择学历" size="large" :key="'edu-' + editVisible"><ElOption v-for="e in eduOptions" :key="e" :label="e" :value="e" /></ElSelect></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">经验要求（年）</label><ElInput v-model.number="editForm.experienceRequirement" placeholder="例如：2" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">招聘人数</label><ElInput v-model.number="editForm.headcount" placeholder="例如：3" size="large" /></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">薪资下限</label><ElInput v-model.number="editForm.salaryMin" placeholder="例如：12000" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">薪资上限</label><ElInput v-model.number="editForm.salaryMax" placeholder="例如：20000" size="large" /></div>
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-emerald-500" /><h3 class="text-[13px] font-semibold text-gray-700">职位描述</h3></div>
          <ElInput v-model="editForm.description" type="textarea" :rows="4" placeholder="描述岗位职责、工作内容..." size="large" />
        </div>
        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-amber-500" /><h3 class="text-[13px] font-semibold text-gray-700">任职要求</h3></div>
          <ElInput v-model="editForm.requirement" type="textarea" :rows="3" placeholder="学历、技能、经验等要求..." size="large" />
        </div>
        <div>
          <div class="flex items-center gap-2 mb-4"><span class="w-1 h-3.5 rounded-full bg-purple-500" /><h3 class="text-[13px] font-semibold text-gray-700">福利待遇</h3></div>
          <ElInput v-model="editForm.welfare" placeholder="例如：五险一金,弹性工作,年终奖（逗号分隔）" size="large" />
        </div>

        <div class="flex items-center justify-between pt-1 border-t border-gray-100">
          <span class="text-[12px] text-gray-500">职位状态</span>
          <span class="text-[12px] mr-2" :class="editStatus === 1 ? 'text-emerald-600' : 'text-gray-400'">{{ editStatus === 1 ? '招聘中' : '已下架' }}</span>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2.5">
          <BaseButton size="large" @click="editVisible = false">取消</BaseButton>
          <BaseButton size="large" type="primary" :loading="editLoading" @click="handleSave">{{ isCreate ? '创建职位' : '保存修改' }}</BaseButton>
        </div>
      </template>
    </Dialog>
  </div>
</template>
