import { z } from 'zod'

const bodySchema = z.object({
  username: z.string().min(1),
  password: z.string().min(1),
})

export default defineEventHandler(async (event) => {
  const { username, password } = await readValidatedBody(event, bodySchema.parse)

  try {
    const res = await $fetch<{ code: number; data: { access_token: string; userInfo: any }; message?: string }>(
      'http://127.0.0.1:8080/api/v1/auth/login',
      { method: 'POST', body: { username, password } }
    )

    if (res.code === 0 && res.data?.access_token) {
      await setUserSession(event, {
        token: res.data.access_token,
        user: res.data.userInfo,
      })
      return { success: true }
    }

    throw createError({ status: 401, message: res.message || '用户名或密码错误' })
  } catch (err: any) {
    if (err.statusCode) throw err // re-throw H3 errors
    throw createError({ status: 401, message: '用户名或密码错误' })
  }
})
