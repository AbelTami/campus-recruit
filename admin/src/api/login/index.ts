import request from '@/axios'
import type { UserLoginType, LoginResponseData } from './types'

export const loginApi = (data: UserLoginType): Promise<IResponse<LoginResponseData>> => {
  return request.post({ url: '/auth/login', data })
}

export const loginOutApi = (): Promise<IResponse> => {
  return request.post({ url: '/auth/logout' })
}

export const refreshTokenApi = (refreshToken: string): Promise<IResponse<LoginResponseData>> => {
  return request.post({ url: '/auth/refresh', data: { refresh_token: refreshToken } })
}
