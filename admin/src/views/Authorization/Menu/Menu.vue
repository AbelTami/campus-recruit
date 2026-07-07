<script setup lang="tsx">
import { reactive, ref, unref } from 'vue'
import { getMenuListApi, createMenuApi, updateMenuApi, deleteMenuApi } from '@/api/menu'
import { useTable } from '@/hooks/web/useTable'
import { Table, TableColumn } from '@/components/Table'
import { Search } from '@/components/Search'
import { FormSchema } from '@/components/Form'
import { ContentWrap } from '@/components/ContentWrap'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { Icon } from '@/components/Icon'
import { ElTag } from 'element-plus'
import Write from './components/Write.vue'
import Detail from './components/Detail.vue'

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getMenuListApi()
    return { list: res.data || [] }
  },
})

const { dataList, loading } = tableState
const { getList } = tableMethods

const tableColumns = reactive<TableColumn[]>([
  { field: 'index', label: '#', type: 'index', width: 60 },
  {
    field: 'meta.title', label: '菜单名称', minWidth: 200,
    slots: { default: (data: any) => <span class="font-500">{data.row.meta?.title}</span> },
  },
  {
    field: 'meta.icon', label: '图标', width: 80, align: 'center',
    slots: {
      default: (data: any) => {
        const icon = data.row.meta?.icon
        return icon ? <Icon icon={icon} size={20} /> : <span class="text-gray-400">-</span>
      },
    },
  },
  {
    field: 'path', label: '路径', minWidth: 180,
    slots: {
      default: (data: any) => <code class="text-13px bg-gray-100 px-6px py-2px rounded">{data.row.path || '-'}</code>,
    },
  },
  {
    field: 'component', label: '组件', width: 140,
    slots: {
      default: (data: any) => {
        const c = data.row.component
        if (c === '#') return <ElTag size="small" type="primary">Layout</ElTag>
        if (c === '##') return <ElTag size="small" type="info">Parent</ElTag>
        return <span class="text-13px text-gray-500">{c || '-'}</span>
      },
    },
  },
  {
    field: 'status', label: '状态', width: 90, align: 'center',
    slots: {
      default: (data: any) => data.row.status === 1
        ? <ElTag type="success" size="small">启用</ElTag>
        : <ElTag type="danger" size="small">禁用</ElTag>,
    },
  },
  {
    field: 'action', label: '操作', width: 200, fixed: 'right',
    slots: {
      default: (data: any) => {
        const row = data.row
        return (
          <div class="flex gap-6px">
            <BaseButton type="primary" size="small" onClick={() => action(row, 'edit')}>编辑</BaseButton>
            <BaseButton plain size="small" onClick={() => action(row, 'detail')}>详情</BaseButton>
            <BaseButton type="danger" size="small" plain onClick={() => delData(row)}>删除</BaseButton>
          </div>
        )
      },
    },
  },
])

const searchSchema = reactive<FormSchema[]>([
  { field: 'meta.title', label: '菜单名称', component: 'Input', componentProps: { placeholder: '搜索菜单' } },
])

const searchParams = ref({})
const setSearchParams = (data: any) => { searchParams.value = data; getList() }

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref()
const actionType = ref('')

const openAdd = () => { dialogTitle.value = '新增菜单'; currentRow.value = undefined; dialogVisible.value = true; actionType.value = '' }

const action = (row: any, type: string) => {
  dialogTitle.value = type === 'edit' ? '编辑菜单' : '菜单详情'
  actionType.value = type
  currentRow.value = row
  dialogVisible.value = true
}

const writeRef = ref<ComponentRef<typeof Write>>()
const saveLoading = ref(false)

const save = async () => {
  const formData = await unref(writeRef)?.submit()
  if (!formData) return
  saveLoading.value = true
  try {
    if (actionType.value === '') await createMenuApi(formData)
    else await updateMenuApi(currentRow.value.id, formData)
    dialogVisible.value = false
    getList()
  } finally { saveLoading.value = false }
}

const delLoading = ref(false)
const delData = async (row: any) => {
  delLoading.value = true
  try { await deleteMenuApi(row.id); getList() }
  finally { delLoading.value = false }
}
</script>

<template>
  <div class="p-16px h-full">
    <ContentWrap class="h-full">
      <div class="flex items-center justify-between mb-20px">
        <div>
          <span class="text-18px font-600">菜单管理</span>
          <span class="text-13px text-gray-400 ml-12px">共 {{ dataList.length }} 项</span>
        </div>
        <BaseButton type="primary" @click="openAdd">新增菜单</BaseButton>
      </div>

      <Search :schema="searchSchema" layout="inline" label-width="70px" @reset="setSearchParams" @search="setSearchParams" />

      <div class="mt-16px">
        <Table :columns="tableColumns" :data="dataList" :loading="loading"  default-expand-all row-key="id" @register="tableRegister" />
      </div>
    </ContentWrap>
  </div>

  <Dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
    <Write v-if="actionType !== 'detail'" ref="writeRef" :current-row="currentRow" />
    <Detail v-else :current-row="currentRow" />
    <template #footer>
      <div class="flex justify-end gap-2.5">
        <BaseButton size="large" @click="dialogVisible = false">取消</BaseButton>
        <BaseButton v-if="actionType !== 'detail'" size="large" type="primary" :loading="saveLoading" @click="save">{{ actionType === '' ? '创建菜单' : '保存修改' }}</BaseButton>
      </div>
    </template>
  </Dialog>
</template>
