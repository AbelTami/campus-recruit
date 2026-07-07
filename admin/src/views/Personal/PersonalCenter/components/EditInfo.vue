<script lang="ts" setup>
import { FormSchema, Form } from '@/components/Form'
import { useForm } from '@/hooks/web/useForm'
import { useValidator } from '@/hooks/web/useValidator'
import { reactive, ref, watch, toRefs } from 'vue'
import { ElDivider, ElMessage, ElMessageBox } from 'element-plus'
import request from '@/axios'

const props = defineProps({
  userInfo: { type: Object, default: () => ({}) },
})
const emit = defineEmits(['update'])

const { required, phone, maxlength, email } = useValidator()

const formSchema = reactive<FormSchema[]>([
  { field: 'nickname', label: '昵称', component: 'Input', colProps: { span: 24 } },
  { field: 'phone', label: '手机号码', component: 'Input', colProps: { span: 24 } },
  { field: 'email', label: '邮箱', component: 'Input', colProps: { span: 24 } },
])

const rules = reactive({
  nickname: [required(), maxlength(50)],
  phone: [phone()],
  email: [email()],
})

const { formRegister, formMethods } = useForm()
const { setValues, getFormData, getElFormExpose } = formMethods

watch(() => props.userInfo, (value) => {
  if (!value) return
  setValues({
    nickname: value.nickname || value.username,
    phone: value.phoneNumber || '',
    email: value.email || '',
  })
}, { immediate: true, deep: true })

const saveLoading = ref(false)
const save = async () => {
  const elForm = await getElFormExpose()
  const valid = await elForm?.validate().catch(() => {})
  if (!valid) return
  ElMessageBox.confirm('是否确认修改?', '提示', {
    confirmButtonText: '确认', cancelButtonText: '取消', type: 'warning',
  }).then(async () => {
    try {
      saveLoading.value = true
      const data = await getFormData()
      await request.put({ url: '/admin/users/' + (props.userInfo.id || 0), data })
      emit('update')
      ElMessage.success('修改成功')
    } catch { ElMessage.error('修改失败') }
    finally { saveLoading.value = false }
  }).catch(() => {})
}
</script>

<template>
  <Form :rules="rules" @register="formRegister" :schema="formSchema" />
  <ElDivider />
  <BaseButton type="primary" @click="save">保存</BaseButton>
</template>
