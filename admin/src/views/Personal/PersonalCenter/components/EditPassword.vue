<script setup lang="ts">
import { Form, FormSchema } from '@/components/Form'
import { useForm } from '@/hooks/web/useForm'
import { reactive, ref } from 'vue'
import { useValidator } from '@/hooks/web/useValidator'
import { ElMessage, ElMessageBox, ElDivider } from 'element-plus'
import request from '@/axios'

const { required } = useValidator()

const formSchema = reactive<FormSchema[]>([
  { field: 'oldPassword', label: '旧密码', component: 'InputPassword', colProps: { span: 24 } },
  { field: 'newPassword', label: '新密码', component: 'InputPassword', colProps: { span: 24 }, componentProps: { strength: true } },
  { field: 'newPassword2', label: '确认新密码', component: 'InputPassword', colProps: { span: 24 }, componentProps: { strength: true } },
])

const rules = reactive({
  oldPassword: [required()],
  newPassword: [required(), { asyncValidator: async (_: any, val: string, callback: any) => {
    const fd = await getFormData()
    if (val !== fd.newPassword2) callback(new Error('新密码与确认新密码不一致'))
    else callback()
  }}],
  newPassword2: [required(), { asyncValidator: async (_: any, val: string, callback: any) => {
    const fd = await getFormData()
    if (val !== fd.newPassword) callback(new Error('确认新密码与新密码不一致'))
    else callback()
  }}],
})

const { formRegister, formMethods } = useForm()
const { getFormData, getElFormExpose } = formMethods

const saveLoading = ref(false)
const save = async () => {
  const elForm = await getElFormExpose()
  const valid = await elForm?.validate().catch(() => {})
  if (!valid) return
  ElMessageBox.confirm('是否确认修改密码?', '提示', {
    confirmButtonText: '确认', cancelButtonText: '取消', type: 'warning',
  }).then(async () => {
    try {
      saveLoading.value = true
      const fd = await getFormData()
      await request.put({ url: '/admin/users/password', data: { oldPassword: fd.oldPassword, newPassword: fd.newPassword } })
      ElMessage.success('密码修改成功')
    } catch { ElMessage.error('修改失败') }
    finally { saveLoading.value = false }
  }).catch(() => {})
}
</script>

<template>
  <Form :rules="rules" @register="formRegister" :schema="formSchema" />
  <ElDivider />
  <BaseButton type="primary" @click="save">确认修改</BaseButton>
</template>
