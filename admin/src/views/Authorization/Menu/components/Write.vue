<script setup lang="tsx">
import { Form, FormSchema } from '@/components/Form'
import { useForm } from '@/hooks/web/useForm'
import { PropType, reactive, watch, ref, unref, computed } from 'vue'
import { useValidator } from '@/hooks/web/useValidator'
import { getMenuListApi } from '@/api/menu'
import { cloneDeep } from 'lodash-es'

const { required } = useValidator()

const props = defineProps({
  currentRow: { type: Object as PropType<any>, default: () => null },
})

const name = computed(() => props.currentRow?.meta?.title || '新菜单')
const isEdit = computed(() => !!props.currentRow)

const cacheComponent = ref('')

const formSchema = reactive<FormSchema[]>([
  {
    field: 'type', label: '菜单类型', component: 'RadioButton', value: 1, colProps: { span: 24 },
    componentProps: {
      options: [{ label: '目录', value: 0 }, { label: '菜单', value: 1 }],
      on: { change: async (val: number) => {
        const formData = await getFormData()
        if (val === 1) { setSchema([{ field: 'component', path: 'componentProps.disabled', value: false }]); setValues({ component: unref(cacheComponent) }) }
        else { setSchema([{ field: 'component', path: 'componentProps.disabled', value: true }]); setValues({ component: formData.parentId === void 0 ? '#' : '##' }) }
      }},
    },
  },
  { field: 'parentId', label: '父级菜单', component: 'TreeSelect', componentProps: { nodeKey: 'id', clearable: true, expandOnClickNode: false, checkStrictly: true, props: { label: 'name', value: 'id', children: 'children' }, on: { change: async (val: number) => { const fd = await getFormData(); if (fd.type === 0) setValues({ component: val ? '##' : '#' }); else setValues({ component: unref(cacheComponent) ?? '' }) } } }, optionApi: async () => { const res = await getMenuListApi(); return res.data || [] } },
  { field: 'meta.title', label: '菜单名称', component: 'Input', componentProps: { placeholder: '侧边栏显示的名称' } },
  { field: 'path', label: '路由路径', component: 'Input', componentProps: { placeholder: '例: /employment 或 dashboard' } },
  { field: 'meta.icon', label: '图标', component: 'Input', componentProps: { placeholder: '例: ri:user-line' } },
  { field: 'component', label: '组件', component: 'Input', value: '#', componentProps: { disabled: true, placeholder: '# Layout / ## Parent / views/...', on: { change: (val: string) => { cacheComponent.value = val } } } },
  { field: 'name', label: '路由名称', component: 'Input', componentProps: { placeholder: '唯一标识，例: EmploymentDashboard' } },
  { field: 'status', label: '状态', component: 'Select', componentProps: { options: [{ label: '禁用', value: 0 }, { label: '启用', value: 1 }] } },
  { field: 'meta.hidden', label: '隐藏菜单', component: 'Switch', formItemProps: { labelWidth: '120px' } },
  { field: 'meta.alwaysShow', label: '始终显示', component: 'Switch', formItemProps: { labelWidth: '120px' } },
])

const rules = reactive({ component: [required()], path: [required()], 'meta.title': [required()] })
const { formRegister, formMethods } = useForm()
const { setValues, getFormData, getElFormExpose, setSchema } = formMethods

const submit = async () => {
  const elForm = await getElFormExpose()
  const valid = await elForm?.validate().catch(() => {})
  if (valid) return await getFormData()
}

watch(() => props.currentRow, (value) => {
  if (!value) return
  const currentRow = cloneDeep(value)
  cacheComponent.value = currentRow.type === 1 ? currentRow.component : ''
  setValues(currentRow)
}, { deep: true, immediate: true })

defineExpose({ submit })
</script>

<template>
  <div class="flex flex-col gap-5">
    <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
      <div class="w-[72px] h-[72px] rounded-xl bg-purple-500 flex items-center justify-center shrink-0">
        <span class="text-white text-2xl font-semibold">{{ isEdit ? name.charAt(0) : '+' }}</span>
      </div>
      <div>
        <h3 class="text-[15px] font-bold text-gray-900">{{ isEdit ? name : '创建新菜单' }}</h3>
        <p class="text-[13px] text-gray-500 mt-0.5">{{ isEdit ? '编辑菜单信息及路由配置' : '创建菜单项并配置路由和权限' }}</p>
      </div>
    </div>
    <Form :rules="rules" @register="formRegister" :schema="formSchema" label-width="120px" />
  </div>
</template>
