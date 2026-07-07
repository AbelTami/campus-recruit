<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { Table } from '@/components/Table'
import { Search } from '@/components/Search'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { ElTag, ElSelect, ElOption, ElInput } from 'element-plus'
import { reactive, ref, unref } from 'vue'
import { getApplicationList, updateApplicationStatus, deleteApplication } from '@/api/application'
import type { Application } from '@/api/application/types'
import { useTable } from '@/hooks/web/useTable'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'

const statusMap: Record<string, string> = { pending: '待处理', viewed: '已查看', interview: '面试中', offer: '已发Offer', accepted: '已接受', rejected: '已拒绝' }
const statusColors: Record<string, string> = { pending: 'info', viewed: '', interview: 'warning', offer: 'primary', accepted: 'success', rejected: 'danger' }

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { pageSize, currentPage } = tableState
    const res = await getApplicationList({ pageIndex: unref(currentPage), pageSize: unref(pageSize), ...unref(searchParams) })
    return { list: res.data.list || [], total: res.data.total || 0 }
  },
  fetchDelApi: async () => { const r = await deleteApplication(Number(unref(ids)[0])); return !!r },
})

const { total, loading, dataList, pageSize, currentPage } = tableState
const { getList, getElTableExpose, delList } = tableMethods

const crudSchemas = reactive<CrudSchema[]>([
  { field: 'index', label: '#', form: { hidden: true }, search: { hidden: true }, table: { type: 'index', width: 50 } },
  { field: 'keyword', label: '搜索', search: { component: 'Input', componentProps: { placeholder: '学生姓名/职位' } }, form: { hidden: true }, table: { hidden: true } },
  {
    field: 'student', label: '学生', width: 100, search: { hidden: true },
    table: { slots: { default: (d: any) => <span class="font-500">{d.row.student?.name || '-'}</span> } },
  },
  {
    field: 'position', label: '职位', minWidth: 180, search: { hidden: true },
    table: { slots: { default: (d: any) => <span>{d.row.position?.title || '-'}</span> } },
  },
  {
    field: 'enterprise', label: '企业', width: 140, search: { hidden: true },
    table: { slots: { default: (d: any) => <span class="text-13px">{d.row.enterprise?.name || '-'}</span> } },
  },
  {
    field: 'status', label: '状态', width: 100,
    search: { component: 'Select', componentProps: { options: Object.entries(statusMap).map(([k, v]) => ({ label: v, value: k })), placeholder: '状态' } },
    table: {
      slots: {
        default: (d: any) => <ElTag size="small" type={statusColors[d.row.status] || 'info'}>{statusMap[d.row.status] || d.row.status}</ElTag>,
      },
    },
  },
  {
    field: 'createdAt', label: '投递时间', width: 160, search: { hidden: true },
    table: { slots: { default: (d: any) => <span class="text-13px text-gray-500">{d.row.createdAt?.slice(0, 10)}</span> } },
  },
  {
    field: 'action', label: '操作', width: 220, fixed: 'right', form: { hidden: true }, search: { hidden: true },
    table: {
      slots: {
        default: (d: any) => {
          const row = d.row as Application
          return (
            <div class="flex gap-6px items-center">
              <ElSelect modelValue={row.status} size="small" style="width:100px" onChange={(v: string) => changeStatus(row, v)}>
                {Object.entries(statusMap).map(([k, v]) => <ElOption key={k} label={v} value={k} />)}
              </ElSelect>
              <BaseButton type="danger" size="small" plain onClick={() => delData(row)}>删除</BaseButton>
            </div>
          )
        },
      },
    },
  },
])

const { allSchemas } = useCrudSchemas(crudSchemas)
const searchParams = ref({})
const setSearchParams = (p: any) => { currentPage.value = 1; searchParams.value = p; getList() }

const ids = ref<string[]>([])
const delLoading = ref(false)
const delData = async (row?: Application) => {
  ids.value = row ? [String(row.id)] : []
  delLoading.value = true
  await delList(unref(ids).length).finally(() => { delLoading.value = false })
}

const changeStatus = async (row: Application, newStatus: string) => {
  await updateApplicationStatus(row.id, newStatus)
  row.status = newStatus
  getList()
}
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap class="h-full">
      <div class="flex items-center justify-between mb-20px">
        <div><span class="text-18px font-600">投递管理</span><span class="text-13px text-gray-400 ml-12px">共 {{ total }} 条</span></div>
      </div>
      <Search :schema="allSchemas.searchSchema" layout="inline" label-width="70px" @reset="setSearchParams" @search="setSearchParams" />
      <div class="mt-16px">
        <Table v-model:current-page="currentPage" v-model:page-size="pageSize" :columns="allSchemas.tableColumns" :data="dataList" :loading="loading" @register="tableRegister" :pagination="{ total }" />
      </div>
    </ContentWrap>
  </div>
</template>
