<script setup lang="tsx">
import { Form, FormSchema } from '@/components/Form'
import { useForm } from '@/hooks/web/useForm'
import { PropType, reactive, watch, ref, unref, nextTick, computed } from 'vue'
import { useValidator } from '@/hooks/web/useValidator'
import { useI18n } from '@/hooks/web/useI18n'
import { ElTree, ElCheckboxGroup, ElCheckbox } from 'element-plus'
import { getMenuListApi } from '@/api/menu'
import { filter, eachTree } from '@/utils/tree'
import { findIndex } from '@/utils'

const { t } = useI18n()
const { required } = useValidator()

const props = defineProps({
  currentRow: { type: Object as PropType<any>, default: () => null },
})

const name = computed(() => props.currentRow?.roleName || '新角色')
const avatar = computed(() => name.value.charAt(0))
const isEdit = computed(() => !!props.currentRow)

const treeRef = ref<typeof ElTree>()
const currentTreeData = ref()
const nodeClick = (treeData: any) => { currentTreeData.value = treeData }

const formSchema = ref<FormSchema[]>([
  { field: 'roleName', label: t('role.roleName'), component: 'Input' },
  { field: 'status', label: t('menu.status'), component: 'Select', componentProps: { options: [{ label: t('userDemo.disable'), value: 0 }, { label: t('userDemo.enable'), value: 1 }] } },
  {
    field: 'menu', label: t('role.menu'), colProps: { span: 24 },
    formItemProps: {
      slots: {
        default: () => (
          <div class="flex w-full">
            <div class="flex-1">
              <ElTree ref={treeRef} show-checkbox node-key="id" highlight-current check-strictly expand-on-click-node={false} data={treeData.value} onNode-click={nodeClick}>
                {{ default: (data: any) => <span>{data.data.meta.title}</span> }}
              </ElTree>
            </div>
            <div class="flex-1">
              {unref(currentTreeData) && unref(currentTreeData)?.permissionList ? (
                <ElCheckboxGroup v-model={unref(currentTreeData).meta.permission}>
                  {unref(currentTreeData)?.permissionList.map((v: any) => <ElCheckbox label={v.value}>{v.label}</ElCheckbox>)}
                </ElCheckboxGroup>
              ) : null}
            </div>
          </div>
        ),
      },
    },
  },
])

const rules = reactive({ roleName: [required()], role: [required()], status: [required()] })
const { formRegister, formMethods } = useForm()
const { setValues, getFormData, getElFormExpose } = formMethods

const treeData = ref<any[]>([])
const getMenuList = async () => {
  const res = await getMenuListApi()
  if (res) {
    treeData.value = res.data
    await nextTick()
    if (!props.currentRow?.id) return
    const assignedIds: number[] = props.currentRow.menuIds || []
    assignedIds.forEach((id: number) => unref(treeRef)?.setChecked(id, true, false))
  }
}
getMenuList()

const submit = async () => {
  const elForm = await getElFormExpose()
  const valid = await elForm?.validate().catch(() => {})
  if (valid) {
    const formData = await getFormData()
    const checkedKeys = unref(treeRef)?.getCheckedKeys() || []
    const data = filter(unref(treeData), (item: any) => checkedKeys.includes(item.id))
    formData.menu = data || []
    return formData
  }
}

watch(() => props.currentRow, (currentRow) => {
  if (!currentRow) return
  setValues(currentRow)
}, { deep: true, immediate: true })

defineExpose({ submit })
</script>

<template>
  <div class="flex flex-col gap-5">
    <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
      <div class="w-[72px] h-[72px] rounded-xl bg-amber-500 flex items-center justify-center shrink-0">
        <span class="text-white text-2xl font-semibold">{{ avatar }}</span>
      </div>
      <div>
        <h3 class="text-[15px] font-bold text-gray-900">{{ isEdit ? name : '创建新角色' }}</h3>
        <p class="text-[13px] text-gray-500 mt-0.5">{{ isEdit ? '编辑角色信息及菜单权限分配' : '创建角色并分配菜单和权限' }}</p>
      </div>
    </div>
    <Form :rules="rules" @register="formRegister" :schema="formSchema" label-width="100px" />
  </div>
</template>
