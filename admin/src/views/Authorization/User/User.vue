<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { Table } from '@/components/Table'
import { Search } from '@/components/Search'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import { ElTree, ElInput, ElCard, ElTag } from 'element-plus'
import { ref, unref, nextTick, reactive, watch } from 'vue'
import { getDepartmentApi, getUserByIdApi, saveUserApi, deleteUserByIdApi } from '@/api/department'
import { getRoleListApi } from '@/api/role'
import type { DepartmentItem, DepartmentUserItem } from '@/api/department/types'
import { useTable } from '@/hooks/web/useTable'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import Write from './components/Write.vue'
import Detail from './components/Detail.vue'

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { pageSize, currentPage } = tableState
    const res = await getUserByIdApi({
      id: unref(currentNodeKey),
      pageIndex: unref(currentPage),
      pageSize: unref(pageSize),
      ...unref(searchParams),
    })
    const list = (res.data.list || []).map((u: any) => ({ ...u, deptName: u.department?.label || '-' }))
    return { list, total: res.data.total || 0 }
  },
  fetchDelApi: async () => {
    const res = await deleteUserByIdApi(unref(ids))
    return !!res
  },
})

const { total, loading, dataList, pageSize, currentPage } = tableState
const { getList, getElTableExpose, delList } = tableMethods

const crudSchemas = reactive<CrudSchema[]>([
  { field: 'selection', search: { hidden: true }, form: { hidden: true }, detail: { hidden: true }, table: { type: 'selection' } },
  { field: 'index', label: '#', form: { hidden: true }, search: { hidden: true }, detail: { hidden: true }, table: { type: 'index', width: 50 } },
  {
    field: 'username', label: '用户名', width: 120,
    form: { component: 'Input', componentProps: { placeholder: '请输入用户名' } },
    search: { component: 'Input', componentProps: { placeholder: '按用户名搜索' } },
  },
  {
    field: 'nickname', label: '昵称', width: 120,
    form: { component: 'Input', componentProps: { placeholder: '请输入昵称' } },
    search: { hidden: true },
  },
  {
    field: 'role', label: '角色', search: { hidden: true },
    detail: {
      slots: {
        default: (data: any) => {
          const roles: string[] = data?.roles || []
          return roles.length
            ? <div class="flex flex-wrap gap-4px">{roles.map((r: string) => <ElTag size="small">{r}</ElTag>)}</div>
            : <span class="text-gray-400">未分配</span>
        },
      },
    },
    form: {
      component: 'Select', value: [],
      componentProps: { multiple: true, collapseTags: true, maxCollapseTags: 1, placeholder: '请选择角色' },
      optionApi: async () => {
        const res = await getRoleListApi()
        return res.data?.list?.map((v: any) => ({ label: v.roleName, value: String(v.code || v.id) })) || []
      },
    },
    table: {
      slots: {
        default: (data: any) => {
          const roles: string[] = data.row.roles || []
          return roles.length
            ? <div class="flex flex-wrap gap-4px">{roles.map((r: string) => <ElTag key={r} size="small">{r}</ElTag>)}</div>
            : <span class="text-gray-400">-</span>
        },
      },
    },
  },
  {
    field: 'deptName', label: '所属部门', search: { hidden: true }, form: { hidden: true },
  },
  {
    field: 'departmentId', label: '部门', search: { hidden: true }, table: { hidden: true }, detail: { hidden: true },
    form: {
      component: 'Select', value: undefined,
      componentProps: { placeholder: '请选择学院', clearable: true },
      optionApi: async () => {
        const res = await getDepartmentApi()
        const list: any[] = []
        function walk(nodes: any[]) {
          for (const n of nodes) {
            if (n.id > 0) list.push({ label: n.label, value: Number(n.id) })
            if (n.children) walk(n.children)
          }
        }
        const data = (res as any).data?.list || (res as any).list || []
        walk(data)
        return list
      },
    },
  },
  {
    field: 'status', label: '状态', width: 80, search: { hidden: true },
    form: {
      component: 'Select',
      componentProps: { options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }], placeholder: '请选择状态' },
    },
    table: {
      slots: {
        default: (data: any) => data.row.status === 1
          ? <ElTag type="success" size="small">启用</ElTag>
          : <ElTag type="danger" size="small">禁用</ElTag>,
      },
    },
  },
  {
    field: 'action', label: '操作', form: { hidden: true }, detail: { hidden: true }, search: { hidden: true },
    table: { width: 240, fixed: 'right',
      slots: {
        default: (data: any) => {
          const row = data.row as DepartmentUserItem
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
  },
])

const { allSchemas } = useCrudSchemas(crudSchemas)

const searchParams = ref({})
const setSearchParams = (params: any) => { currentPage.value = 1; searchParams.value = params; getList() }

const treeEl = ref<typeof ElTree>()
const currentNodeKey = ref('')
const departmentList = ref<DepartmentItem[]>([])

const fetchDepartment = async () => {
  const res = await getDepartmentApi()
  if (res?.data?.list) {
    departmentList.value = res.data.list
    const first = res.data.list[0]
    if (first?.children?.[0]) currentNodeKey.value = first.children[0].id
    await nextTick()
    unref(treeEl)?.setCurrentKey(currentNodeKey.value)
  }
}
fetchDepartment()

const currentDepartment = ref('')
watch(() => currentDepartment.value, (val) => { unref(treeEl)?.filter(val) })

const currentChange = (data: DepartmentItem) => {
  currentNodeKey.value = data.id
  currentPage.value = 1
  getList()
}

const filterNode = (value: string, data: DepartmentItem) => !value || data.label.includes(value)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<DepartmentUserItem>()
const actionType = ref('')

const openAdd = () => { dialogTitle.value = '新增用户'; currentRow.value = undefined; dialogVisible.value = true; actionType.value = '' }

const delLoading = ref(false)
const ids = ref<string[]>([])

const delData = async (row?: DepartmentUserItem) => {
  const elTableExpose = await getElTableExpose()
  ids.value = row ? [row.id] : elTableExpose?.getSelectionRows().map((v: DepartmentUserItem) => v.id) || []
  delLoading.value = true
  await delList(unref(ids).length).finally(() => { delLoading.value = false })
}

const action = (row: DepartmentUserItem, type: string) => {
  dialogTitle.value = type === 'edit' ? '编辑用户' : '用户详情'
  actionType.value = type
  currentRow.value = { ...row, role: (row as any).roles || [], departmentId: Number((row as any).department?.id) || undefined }
  dialogVisible.value = true
}

const writeRef = ref<ComponentRef<typeof Write>>()
const saveLoading = ref(false)

const save = async () => {
  const formData = await unref(writeRef)?.submit()
  if (!formData) return
  const payload: any = { username: formData.username, nickname: formData.nickname, status: formData.status, role: formData.role, departmentId: formData.departmentId ? Number(formData.departmentId) : null }
  if (currentRow.value) payload.id = currentRow.value.id
  saveLoading.value = true
  try {
    const res = await saveUserApi(payload)
    if (res) { currentPage.value = 1; getList(); dialogVisible.value = false }
  } finally { saveLoading.value = false }
}
</script>

<template>
  <div class="flex w-full h-full gap-16px p-16px">
    <el-card shadow="never" class="w-240px shrink-0 h-full">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-600 text-15px">部门结构</span>
          <span class="text-12px text-gray-400">{{ departmentList.length }} 个部门</span>
        </div>
      </template>
      <ElInput v-model="currentDepartment" placeholder="输入关键字过滤" clearable size="small" class="mb-12px" />
      <ElTree
        ref="treeEl" :data="departmentList" default-expand-all
        :expand-on-click-node="false" node-key="id"
        :current-node-key="currentNodeKey"
        :props="{ label: 'label' }"
        :filter-node-method="filterNode"
        highlight-current
        @current-change="currentChange"
      />
    </el-card>

    <ContentWrap class="flex-1 h-full">
      <div class="flex items-center justify-between mb-16px">
        <div>
          <span class="text-16px font-600">用户列表</span>
          <span class="text-13px text-gray-400 ml-8px">共 {{ total }} 人</span>
        </div>
        <div class="flex gap-8px">
          <BaseButton type="primary" @click="openAdd">新增用户</BaseButton>
          <BaseButton :loading="delLoading" type="danger" plain @click="delData()">批量删除</BaseButton>
        </div>
      </div>

      <Search :schema="allSchemas.searchSchema" layout="inline" label-width="60px" @reset="setSearchParams" @search="setSearchParams" />

      <Table
        v-model:current-page="currentPage" v-model:page-size="pageSize"
        :columns="allSchemas.tableColumns" :data="dataList" :loading="loading"
        @register="tableRegister" :pagination="{ total }"
        class="mt-10px"
      />
    </ContentWrap>

    <Dialog v-model="dialogVisible" :title="dialogTitle" width="600px" :close-on-click-modal="false">
      <Write v-if="actionType !== 'detail'" ref="writeRef" :form-schema="allSchemas.formSchema" :current-row="currentRow" />
      <Detail v-if="actionType === 'detail'" :detail-schema="allSchemas.detailSchema" :current-row="currentRow" />
      <template #footer>
        <div class="flex justify-end gap-2.5">
          <BaseButton size="large" @click="dialogVisible = false">取消</BaseButton>
          <BaseButton v-if="actionType !== 'detail'" size="large" type="primary" :loading="saveLoading" @click="save">{{ actionType === '' ? '创建用户' : '保存修改' }}</BaseButton>
        </div>
      </template>
    </Dialog>
  </div>
</template>
