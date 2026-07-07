<script setup lang="tsx">
import { PropType, ref, unref, nextTick } from 'vue'
import { Descriptions } from '@/components/Descriptions'
import { ElTag, ElTree } from 'element-plus'
import { findIndex } from '@/utils'
import { getMenuListApi } from '@/api/menu'

defineProps({
  currentRow: { type: Object as PropType<any>, default: () => undefined },
})

const filterPermissionName = (value: string) => {
  const index = findIndex(unref(currentTreeData)?.permissionList || [], (item) => item.value === value)
  return (unref(currentTreeData)?.permissionList || [])[index].label ?? ''
}

const renderTag = (enable?: boolean) => <ElTag type={!enable ? 'danger' : 'success'}>{enable ? '启用' : '禁用'}</ElTag>

const treeRef = ref<typeof ElTree>()
const currentTreeData = ref()
const nodeClick = (treeData: any) => { currentTreeData.value = treeData }

const treeData = ref<any[]>([])
const getMenuList = async () => {
  const res = await getMenuListApi()
  if (res) { treeData.value = res.data; await nextTick() }
}
getMenuList()

const detailSchema = ref<any[]>([
  { field: 'roleName', label: '角色名称' },
  { field: 'status', label: '状态', slots: { default: (data: any) => renderTag(data.status) } },
  { field: 'remark', label: '备注', span: 24 },
  {
    field: 'permissionList', label: '菜单分配', span: 24,
    slots: {
      default: () => (
        <div class="flex w-full">
          <div class="flex-1">
            <ElTree ref={treeRef} node-key="id" props={{ children: 'children', label: 'label' }} highlight-current expand-on-click-node={false} data={treeData.value} onNode-click={nodeClick}>
              {{ default: (data: any) => <span>{data?.data?.label || data?.data?.title}</span> }}
            </ElTree>
          </div>
          <div class="flex-1">{unref(currentTreeData) ? unref(currentTreeData)?.meta?.permission?.map((v: string) => <ElTag class="ml-2 mt-2">{filterPermissionName(v)}</ElTag>) : null}</div>
        </div>
      ),
    },
  },
])
</script>

<template>
  <div v-if="currentRow" class="flex flex-col gap-5">
    <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
      <div class="w-[72px] h-[72px] rounded-xl bg-amber-500 flex items-center justify-center shrink-0">
        <span class="text-white text-2xl font-semibold">{{ currentRow.roleName?.charAt(0) || 'R' }}</span>
      </div>
      <div>
        <h3 class="text-[15px] font-bold text-gray-900">{{ currentRow.roleName }}</h3>
        <p class="text-[13px] text-gray-500 mt-0.5">{{ currentRow.status ? '启用' : '禁用' }}</p>
      </div>
    </div>
    <Descriptions :schema="detailSchema" :data="currentRow" :column="2" />
  </div>
</template>
