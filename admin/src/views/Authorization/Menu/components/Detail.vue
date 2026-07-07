<script setup lang="tsx">
import { PropType, ref, computed } from 'vue'
import { Descriptions, DescriptionsSchema } from '@/components/Descriptions'
import { Icon } from '@/components/Icon'
import { ElTag } from 'element-plus'

const props = defineProps({
  currentRow: { type: Object as PropType<any>, default: () => undefined },
})

const name = computed(() => props.currentRow?.meta?.title || props.currentRow?.name || '菜单')
const icon = computed(() => props.currentRow?.meta?.icon || 'ri:menu-line')

const renderTag = (enable?: boolean) => <ElTag type={!enable ? 'danger' : 'success'}>{enable ? '启用' : '禁用'}</ElTag>

const detailSchema = ref<DescriptionsSchema[]>([
  { field: 'type', label: '菜单类型', span: 24, slots: { default: (data: any) => <>{data.type === 1 ? '菜单' : '目录'}</> } },
  { field: 'parentName', label: '父级菜单' },
  { field: 'meta.title', label: '菜单名称' },
  { field: 'component', label: '组件', slots: { default: (data: any) => <>{data.component === '#' ? '顶级目录' : data.component === '##' ? '子目录' : data.component}</> } },
  { field: 'name', label: '组件名称' },
  { field: 'meta.icon', label: '图标', slots: { default: (data: any) => data.icon ? <Icon icon={data.icon} /> : null } },
  { field: 'path', label: '路径' },
  { field: 'meta.activeMenu', label: '高亮菜单' },
  { field: 'permissionList', label: '按钮权限', span: 24, slots: { default: (data: any) => <>{data?.permissionList?.map((v: any) => <ElTag class="mr-1" key={v.value}>{v.label}</ElTag>)}</> } },
  { field: 'menuState', label: '菜单状态', slots: { default: (data: any) => renderTag(data.menuState) } },
  { field: 'meta.hidden', label: '是否隐藏', slots: { default: (data: any) => renderTag(data.enableHidden) } },
  { field: 'meta.alwaysShow', label: '是否一直显示', slots: { default: (data: any) => renderTag(data.enableDisplay) } },
  { field: 'meta.noCache', label: '是否清除缓存', slots: { default: (data: any) => renderTag(data.enableCleanCache) } },
  { field: 'meta.breadcrumb', label: '是否显示面包屑', slots: { default: (data: any) => renderTag(data.enableShowCrumb) } },
  { field: 'meta.affix', label: '是否固定标签页', slots: { default: (data: any) => renderTag(data.enablePinnedTab) } },
  { field: 'meta.noTagsView', label: '是否隐藏标签页', slots: { default: (data: any) => renderTag(data.enableHiddenTab) } },
  { field: 'meta.canTo', label: '是否可跳转', slots: { default: (data: any) => renderTag(data.enableSkip) } },
])
</script>

<template>
  <div v-if="currentRow" class="flex flex-col gap-5">
    <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
      <div class="w-[72px] h-[72px] rounded-xl bg-purple-500 flex items-center justify-center shrink-0">
        <Icon v-if="icon" :icon="icon" :size="28" color="#fff" />
        <span v-else class="text-white text-2xl font-semibold">M</span>
      </div>
      <div>
        <h3 class="text-[15px] font-bold text-gray-900">{{ name }}</h3>
        <p class="text-[13px] text-gray-500 mt-0.5">{{ currentRow.path || '-' }}</p>
      </div>
    </div>
    <Descriptions :schema="detailSchema" :data="currentRow" :column="2" />
  </div>
</template>
