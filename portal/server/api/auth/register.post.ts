import { z } from 'zod'

const bodySchema = z.object({
  username: z.string().min(3),
  password: z.string().min(6).max(32),
  name: z.string().min(1),
  studentNo: z.string().optional(),
})

export default defineEventHandler(async (event) => {
  const body = await readValidatedBody(event, bodySchema.parse)

  try {
    const res = await $fetch<{ code: number; message?: string }>(
      'http://127.0.0.1:8080/api/v1/auth/register',
      { method: 'POST', body }
    )

    if (res.code === 0) {
      return { success: true }
    }

    throw createError({ status: 400, message: res.message || '注册失败' })
  } catch (err: any) {
    if (err.statusCode) throw err
    throw createError({ status: 400, message: '注册失败，用户名可能已存在' })
  }
})
