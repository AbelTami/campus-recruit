<script setup lang="tsx">
import { reactive, ref, unref } from 'vue'
import { getRoleListApi, createRoleApi, updateRoleApi, deleteRoleApi } from '@/api/role'
import { useTable } from '@/hooks/web/useTable'
import { Table, TableColumn } from '@/components/Table'
import { Search } from '@/components/Search'
import { FormSchema } from '@/components/Form'
import { ContentWrap } from '@/components/ContentWrap'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { ElTag } from 'element-plus'
import Write from './components/Write.vue'
import Detail from './components/Detail.vue'

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getRoleListApi()
    return { list: res.data.list || [], total: res.data.list?.length || 0 }
  },
})

const { dataList, loading, total } = tableState
const { getList } = tableMethods

const tableColumns = reactive<TableColumn[]>([
  { field: 'index', label: '#', type: 'index', width: 50 },
  { field: 'roleName', label: '角色名称', width: 150 },
  { field: 'code', label: '角色编码', width: 150 },
  {
    field: 'status', label: '状态', width: 80,
    slots: {
      default: (data: any) => data.row.status === 1
        ? <ElTag type="success" size="small">启用</ElTag>
        : <ElTag type="danger" size="small">禁用</ElTag>,
    },
  },
  { field: 'remark', label: '备注', minWidth: 200 },
  {
    field: 'action', label: '操作', width: 240, fixed: 'right',
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
  { field: 'roleName', label: '角色名称', component: 'Input', componentProps: { placeholder: '搜索角色' } },
])

const searchParams = ref({})
const setSearchParams = (data: any) => { searchParams.value = data; getList() }

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref()
const actionType = ref('')

const openAdd = () => { dialogTitle.value = '新增角色'; currentRow.value = undefined; dialogVisible.value = true; actionType.value = '' }

const action = (row: any, type: string) => {
  dialogTitle.value = type === 'edit' ? '编辑角色' : '角色详情'
  actionType.value = type
  currentRow.value = row
  dialogVisible.value = true
}

const writeRef = ref<ComponentRef<typeof Write>>()
const saveLoading = ref(false)

const save = async () => {
  const write = unref(writeRef)
  const formData = await write?.submit()
  if (!formData) return
  saveLoading.value = true
  try {
    if (actionType.value === '') await createRoleApi(formData)
    else await updateRoleApi(currentRow.value.id, formData)
    dialogVisible.value = false
    getList()
  } finally { saveLoading.value = false }
}

const delLoading = ref(false)
const delData = async (row: any) => {
  delLoading.value = true
  try {
    await deleteRoleApi(row.id)
    getList()
  } finally { delLoading.value = false }
}
</script>

<template>
  <ContentWrap class="h-full">
    <div class="flex items-center justify-between mb-16px">
      <div>
        <span class="text-16px font-600">角色列表</span>
        <span class="text-13px text-gray-400 ml-8px">共 {{ total }} 个</span>
      </div>
      <BaseButton type="primary" @click="openAdd">新增角色</BaseButton>
    </div>

    <Search :schema="searchSchema" layout="inline" label-width="70px" @reset="setSearchParams" @search="setSearchParams" />

    <Table
      :columns="tableColumns" :data="dataList" :loading="loading"
      :pagination="{ total }" @register="tableRegister"
      class="mt-10px"
    />
  </ContentWrap>

  <Dialog v-model="dialogVisible" :title="dialogTitle" width="700px" :close-on-click-modal="false">
    <Write v-if="actionType !== 'detail'" ref="writeRef" :current-row="currentRow" />
    <Detail v-else :current-row="currentRow" />
    <template #footer>
      <div class="flex justify-end gap-2.5">
        <BaseButton size="large" @click="dialogVisible = false">取消</BaseButton>
        <BaseButton v-if="actionType !== 'detail'" size="large" type="primary" :loading="saveLoading" @click="save">{{ actionType === '' ? '创建角色' : '保存修改' }}</BaseButton>
      </div>
    </template>
  </Dialog>
</template>
