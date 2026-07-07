<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { Table } from '@/components/Table'
import { Search } from '@/components/Search'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { ElInput, ElSelect, ElOption, ElTag, ElSwitch } from 'element-plus'
import { reactive, ref, unref, onMounted } from 'vue'
import { getEnterpriseList, createEnterprise, updateEnterprise, batchDeleteEnterprises } from '@/api/enterprise'
import { getIndustryList } from '@/api/enterprise/industry'
import type { Enterprise, Industry } from '@/api/enterprise/types'
import { useTable } from '@/hooks/web/useTable'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import request from '@/axios'

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { pageSize, currentPage } = tableState
    const res = await getEnterpriseList({ pageIndex: unref(currentPage), pageSize: unref(pageSize), ...unref(searchParams) })
    return { list: res.data.list || [], total: res.data.total || 0 }
  },
  fetchDelApi: async () => {
    const res = await batchDeleteEnterprises(unref(ids).map(Number))
    return !!res
  },
})

const { total, loading, dataList, pageSize, currentPage } = tableState
const { getList, getElTableExpose, delList } = tableMethods

const scaleMap: Record<string, string> = { micro: '1-50人', small: '50-200人', medium: '200-500人', large: '500-1000人', '1000+': '1000人以上', '500-1000': '500-1000人' }
const natureMap: Record<string, string> = { state_owned: '国有企业', private: '民营企业', foreign: '外资企业', joint_venture: '合资企业', public_institution: '事业单位', government: '政府机关' }
const scaleOptions = Object.entries(scaleMap).map(([k, v]) => ({ label: v, value: k }))
const natureOptions = Object.entries(natureMap).map(([k, v]) => ({ label: v, value: k }))

const crudSchemas = reactive<CrudSchema[]>([
  { field: 'selection', search: { hidden: true }, form: { hidden: true }, table: { type: 'selection' } },
  { field: 'index', label: '#', form: { hidden: true }, search: { hidden: true }, table: { type: 'index', width: 50 } },
  { field: 'keyword', label: '搜索', search: { component: 'Input', componentProps: { placeholder: '企业名称' } }, form: { hidden: true }, table: { hidden: true } },
  { field: 'name', label: '企业名称', minWidth: 160, search: { hidden: true } },
  { field: 'industry', label: '行业', minWidth: 100, search: { hidden: true }, table: { slots: { default: (d: any) => <ElTag size="small" type="info">{d.row.industry?.name || '-'}</ElTag> } } },
  { field: 'scale', label: '规模', minWidth: 90, search: { hidden: true }, table: { slots: { default: (d: any) => <span class="text-13px">{scaleMap[d.row.scale] || d.row.scale || '-'}</span> } } },
  { field: 'nature', label: '性质', minWidth: 80, search: { hidden: true }, table: { slots: { default: (d: any) => <span class="text-13px">{natureMap[d.row.nature] || d.row.nature || '-'}</span> } } },
  { field: 'city', label: '城市', minWidth: 75, search: { hidden: true } },
  { field: 'status', label: '状态', minWidth: 70, table: { slots: { default: (d: any) => d.row.status === 1 ? <ElTag type="success" size="small">正常</ElTag> : <ElTag type="danger" size="small">禁用</ElTag> } } },
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
const delData = async (row?: Enterprise) => {
  const el = await getElTableExpose()
  ids.value = row ? [String(row.id)] : el?.getSelectionRows().map((v: Enterprise) => String(v.id)) || []
  delLoading.value = true
  await delList(unref(ids).length).finally(() => { delLoading.value = false })
}

// ── detail dialog ──
const detailVisible = ref(false)
const currentEnterprise = ref<Enterprise>()
const viewDetail = (row: Enterprise) => { currentEnterprise.value = row; detailVisible.value = true }

// ── edit dialog ──
const editVisible = ref(false)
const editLoading = ref(false)
const editForm = reactive<Record<string, any>>({})
const editScale = ref('')
const editNature = ref('')
const editIndustryId = ref<number>()
const logoPreview = ref('')
const logoFile = ref<File | null>(null)
const isCreate = ref(false)
const industries = ref<Industry[]>([])

onMounted(async () => {
  const res = await getIndustryList()
  if (res.code === 0 && res.data) industries.value = res.data
})

function openCreate() {
  Object.assign(editForm, { name: '', shortName: '', city: '', address: '', website: '', description: '', contactName: '', contactPhone: '', contactEmail: '', status: 1 })
  editScale.value = 'medium'
  editNature.value = 'private'
  editIndustryId.value = undefined
  logoPreview.value = ''
  logoFile.value = null
  isCreate.value = true
  editVisible.value = true
}

function openEdit(row: Enterprise) {
  isCreate.value = false
  Object.assign(editForm, {
    name: row.name, shortName: row.shortName, city: row.city,
    address: row.address || '', website: row.website || '',
    description: row.description || '', contactName: row.contactName || '',
    contactPhone: row.contactPhone || '', contactEmail: row.contactEmail || '',
    status: row.status,
  })
  editScale.value = String(row.scale || '')
  editNature.value = String(row.nature || '')
  editIndustryId.value = row.industryId ? Number(row.industryId) : undefined
  logoPreview.value = row.logoUrl ? `/api/v1${row.logoUrl}` : ''
  logoFile.value = null
  currentEnterprise.value = row
  editVisible.value = true
}

function onLogoChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  logoFile.value = file
  logoPreview.value = URL.createObjectURL(file)
}

async function handleSave() {
  editLoading.value = true
  try {
    let logoUrl = currentEnterprise.value?.logoUrl
    if (logoFile.value) {
      const fd = new FormData()
      fd.append('file', logoFile.value)
      const res = await request.upload<{ url: string }>({ url: '/upload', data: fd })
      logoUrl = res.data.url
    }
    const data = { ...editForm, scale: editScale.value, nature: editNature.value, industryId: editIndustryId.value ?? null, logoUrl }
    if (isCreate.value) await createEnterprise(data)
    else await updateEnterprise(currentEnterprise.value!.id, data)
    editVisible.value = false
    getList()
  } finally { editLoading.value = false }
}

</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap class="h-full">
      <div class="flex items-center justify-between mb-20px">
        <div><span class="text-18px font-600">企业管理</span><span class="text-13px text-gray-400 ml-12px">共 {{ total }} 家</span></div>
        <div class="flex gap-2"><BaseButton type="primary" @click="openCreate">新增企业</BaseButton><BaseButton :loading="delLoading" type="danger" plain @click="delData()">批量删除</BaseButton></div>
      </div>

      <Search :schema="allSchemas.searchSchema" layout="inline" label-width="70px" @reset="setSearchParams" @search="setSearchParams" />
      <div class="mt-16px">
        <Table v-model:current-page="currentPage" v-model:page-size="pageSize" :columns="allSchemas.tableColumns" :data="dataList" :loading="loading" :show-overflow-tooltip="false" style="width: auto !important" @register="tableRegister" :pagination="{ total }" />
      </div>
    </ContentWrap>

    <!-- Detail dialog -->
    <Dialog v-model="detailVisible" title="企业详情" width="560px">
      <div v-if="currentEnterprise" class="flex flex-col gap-5">
        <!-- Logo + name -->
        <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
          <div class="w-[72px] h-[72px] rounded-xl bg-white flex items-center justify-center overflow-hidden shrink-0 ring-1 ring-black/5">
            <img v-if="currentEnterprise.logoUrl" :src="`/api/v1${currentEnterprise.logoUrl}`" class="w-full h-full object-cover" />
            <span v-else class="text-[28px]">🏢</span>
          </div>
          <div>
            <h3 class="text-[15px] font-bold text-gray-900">{{ currentEnterprise.name }}</h3>
            <p class="text-[13px] text-gray-500 mt-0.5">{{ currentEnterprise.shortName || currentEnterprise.name }}</p>
            <div class="flex items-center gap-2 mt-1.5">
              <span v-if="currentEnterprise.industry" class="text-[11px] bg-blue-50 text-blue-600 px-2 py-0.5 rounded-full font-medium">{{ currentEnterprise.industry.name }}</span>
              <span class="text-[11px]" :class="currentEnterprise.status === 1 ? 'text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full font-medium' : 'text-red-500 bg-red-50 px-2 py-0.5 rounded-full font-medium'">{{ currentEnterprise.status === 1 ? '正常' : '禁用' }}</span>
            </div>
          </div>
        </div>

        <!-- Basic info -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <span class="w-1 h-3.5 rounded-full bg-blue-500" />
            <h3 class="text-[13px] font-semibold text-gray-700">基本信息</h3>
          </div>
          <div class="grid grid-cols-2 gap-x-4 gap-y-2.5 text-[13px]">
            <div><span class="text-gray-400">城市</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.city || '-' }}</p></div>
            <div><span class="text-gray-400">规模</span><p class="text-gray-700 mt-0.5">{{ scaleMap[currentEnterprise.scale] || currentEnterprise.scale || '-' }}</p></div>
            <div><span class="text-gray-400">性质</span><p class="text-gray-700 mt-0.5">{{ natureMap[currentEnterprise.nature] || currentEnterprise.nature || '-' }}</p></div>
            <div><span class="text-gray-400">行业</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.industry?.name || '-' }}</p></div>
          </div>
        </div>

        <!-- Contact -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <span class="w-1 h-3.5 rounded-full bg-amber-500" />
            <h3 class="text-[13px] font-semibold text-gray-700">联系信息</h3>
          </div>
          <div class="grid grid-cols-2 gap-x-4 gap-y-2.5 text-[13px]">
            <div class="col-span-2"><span class="text-gray-400">地址</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.address || '-' }}</p></div>
            <div class="col-span-2"><span class="text-gray-400">官网</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.website || '-' }}</p></div>
            <div><span class="text-gray-400">联系人</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.contactName || '-' }}</p></div>
            <div><span class="text-gray-400">电话</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.contactPhone || '-' }}</p></div>
            <div class="col-span-2"><span class="text-gray-400">邮箱</span><p class="text-gray-700 mt-0.5">{{ currentEnterprise.contactEmail || '-' }}</p></div>
          </div>
        </div>

        <!-- Description -->
        <div v-if="currentEnterprise.description">
          <div class="flex items-center gap-2 mb-3">
            <span class="w-1 h-3.5 rounded-full bg-emerald-500" />
            <h3 class="text-[13px] font-semibold text-gray-700">企业简介</h3>
          </div>
          <p class="text-[13px] text-gray-600 leading-relaxed whitespace-pre-wrap">{{ currentEnterprise.description }}</p>
        </div>
      </div>
      <template #footer><BaseButton @click="detailVisible = false">关闭</BaseButton></template>
    </Dialog>

    <!-- Edit dialog -->
    <Dialog v-model="editVisible" :title="isCreate ? '新增企业' : '编辑企业'" width="680px" :close-on-click-modal="false">
      <div class="flex flex-col gap-6">
        <!-- Logo upload -->
        <div class="flex items-center gap-5 p-5 rounded-2xl bg-gradient-to-br from-gray-50 to-slate-50 border-2 border-dashed"
          :class="logoPreview ? 'border-blue-200' : 'border-gray-200 hover:border-gray-300'">
          <div class="relative w-24 h-24 rounded-2xl bg-white flex items-center justify-center overflow-hidden shrink-0 shadow-sm ring-1 ring-black/5"
            :class="logoPreview ? '' : 'border-2 border-dashed border-gray-200'">
            <img v-if="logoPreview" :src="logoPreview" class="w-full h-full object-cover" />
            <div v-else class="flex flex-col items-center gap-0.5">
              <span class="text-[28px] leading-none">🏢</span>
              <span class="text-[10px] text-gray-400">暂无</span>
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-[13px] font-semibold text-gray-800 mb-1">企业标识</p>
            <p class="text-[12px] text-gray-400 mb-3 leading-relaxed">上传企业 Logo 将展示在职位详情页和搜索结果中</p>
            <div class="flex items-center gap-2">
              <label class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-white border border-gray-200 text-[12px] font-medium text-gray-600 hover:border-blue-300 hover:text-blue-600 hover:shadow-sm cursor-pointer transition-all shadow-sm">
                {{ logoPreview ? '更换图片' : '上传 Logo' }}
                <input type="file" accept="image/*" @change="onLogoChange" class="hidden" />
              </label>
              <button v-if="logoPreview" @click="logoPreview=''; logoFile=null" class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-[12px] text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">移除</button>
            </div>
          </div>
        </div>

        <!-- Basic info -->
        <div>
          <div class="flex items-center gap-2 mb-4">
            <span class="w-1 h-3.5 rounded-full bg-blue-500" />
            <h3 class="text-[13px] font-semibold text-gray-700">基本信息</h3>
          </div>
          <div class="space-y-3">
            <div class="flex gap-3">
              <div class="flex-[2]"><label class="text-[12px] text-gray-500 mb-1.5 block">企业名称 <span class="text-red-400">*</span></label><ElInput v-model="editForm.name" placeholder="请输入企业全称" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">简称</label><ElInput v-model="editForm.shortName" placeholder="简称" size="large" /></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-[1.2]"><label class="text-[12px] text-gray-500 mb-1.5 block">所在城市</label><ElInput v-model="editForm.city" placeholder="例如：深圳" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">所属行业</label><ElSelect v-model="editIndustryId" placeholder="选择行业" size="large" popper-class="!z-[9999]" :key="'ind-' + editVisible"><ElOption v-for="ind in industries" :key="ind.id" :label="ind.name" :value="ind.id" /></ElSelect></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">企业规模</label><ElSelect v-model="editScale" placeholder="选择规模" size="large" popper-class="!z-[9999]" :key="'scale-' + editVisible"><ElOption v-for="s in scaleOptions" :key="s.value" :label="s.label" :value="s.value" /></ElSelect></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">企业性质</label><ElSelect v-model="editNature" placeholder="选择性质" size="large" popper-class="!z-[9999]" :key="'nature-' + editVisible"><ElOption v-for="n in natureOptions" :key="n.value" :label="n.label" :value="n.value" /></ElSelect></div>
            </div>
          </div>
        </div>

        <!-- Contact -->
        <div>
          <div class="flex items-center gap-2 mb-4">
            <span class="w-1 h-3.5 rounded-full bg-amber-500" />
            <h3 class="text-[13px] font-semibold text-gray-700">联系信息</h3>
          </div>
          <div class="space-y-3">
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">详细地址</label><ElInput v-model="editForm.address" placeholder="例如：深圳市南山区科技园路 1 号" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">官网地址</label><ElInput v-model="editForm.website" placeholder="例如：www.example.com" size="large" /></div>
            </div>
            <div class="flex gap-3">
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">联系人</label><ElInput v-model="editForm.contactName" placeholder="例如：张经理" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">联系电话</label><ElInput v-model="editForm.contactPhone" placeholder="例如：138xxxx8888" size="large" /></div>
              <div class="flex-1"><label class="text-[12px] text-gray-500 mb-1.5 block">电子邮箱</label><ElInput v-model="editForm.contactEmail" placeholder="例如：hr@example.com" size="large" /></div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div>
          <div class="flex items-center gap-2 mb-4">
            <span class="w-1 h-3.5 rounded-full bg-emerald-500" />
            <h3 class="text-[13px] font-semibold text-gray-700">企业简介</h3>
          </div>
          <label class="text-[12px] text-gray-500 mb-1.5 block">公司介绍</label>
          <ElInput v-model="editForm.description" type="textarea" :rows="4" placeholder="描述公司主营业务、企业文化、发展历程、核心优势等..." size="large" />
        </div>

        <!-- Status -->
        <div class="flex items-center justify-between py-3 px-4 rounded-xl bg-gray-50/70">
          <div class="flex items-center gap-3">
            <span class="text-[13px] text-gray-600 font-medium">企业状态</span>
            <span class="text-[12px] text-gray-400">禁用后前台将不展示该企业</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-[12px]" :class="editForm.status === 1 ? 'text-emerald-600' : 'text-gray-400'">{{ editForm.status === 1 ? '正常' : '禁用' }}</span>
            <ElSwitch v-model="editForm.status" :active-value="1" :inactive-value="0" size="large" inline-prompt active-text="开" inactive-text="关" style="--el-switch-on-color: #10b981" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2.5">
          <BaseButton size="large" @click="editVisible = false">取消</BaseButton>
          <BaseButton size="large" type="primary" :loading="editLoading" @click="handleSave">{{ isCreate ? '创建企业' : '保存修改' }}</BaseButton>
        </div>
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
:deep(.el-table__header-wrapper), :deep(.el-table__body-wrapper) { overflow-x: hidden; }
</style>
<style>
#app .el-table { width: auto !important; }
</style>
