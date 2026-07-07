<script setup lang="ts">
import { Form, FormSchema } from '@/components/Form'
import { useForm } from '@/hooks/web/useForm'
import { PropType, reactive, watch, computed } from 'vue'
import { DepartmentUserItem } from '@/api/department/types'
import { useValidator } from '@/hooks/web/useValidator'

const { required } = useValidator()

const props = defineProps({
  currentRow: { type: Object as PropType<DepartmentUserItem>, default: () => undefined },
  formSchema: { type: Array as PropType<FormSchema[]>, default: () => [] },
})

const name = computed(() => props.currentRow?.username || '新用户')
const avatar = computed(() => name.value.charAt(0))
const isEdit = computed(() => !!props.currentRow)

const rules = reactive({
  username: [required()],
  account: [required()],
  'department.id': [required()],
})

const { formRegister, formMethods } = useForm()
const { setValues, getFormData, getElFormExpose } = formMethods

const submit = async () => {
  const elForm = await getElFormExpose()
  const valid = await elForm?.validate().catch(() => {})
  if (valid) return await getFormData()
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
      <div class="w-[72px] h-[72px] rounded-full bg-blue-500 flex items-center justify-center shrink-0">
        <span class="text-white text-2xl font-semibold">{{ avatar }}</span>
      </div>
      <div>
        <h3 class="text-[15px] font-bold text-gray-900">{{ isEdit ? name : '创建新用户' }}</h3>
        <p class="text-[13px] text-gray-500 mt-0.5">{{ isEdit ? '编辑用户信息及部门角色分配' : '填写信息创建新的系统用户' }}</p>
      </div>
    </div>
    <Form :rules="rules" @register="formRegister" :schema="formSchema" label-width="100px" />
  </div>
</template>
