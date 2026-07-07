import request from '@/axios'
import { DepartmentListResponse, DepartmentUserParams, DepartmentUserResponse } from './types'

export const getDepartmentApi = () => {
  return request.get<DepartmentListResponse>({ url: '/admin/colleges' })
}

export const getUserByIdApi = (params: DepartmentUserParams) => {
  return request.get<DepartmentUserResponse>({ url: '/admin/users', params })
}

export const deleteUserByIdApi = (ids: string[] | number[]) => {
  return request.post({ url: '/admin/users/delete', data: { ids } })
}

export const saveUserApi = (data: any) => {
  if (data.id) {
    return request.put({ url: `/admin/users/${data.id}`, data })
  }
  return request.post({ url: '/admin/users', data })
}

export const saveDepartmentApi = (data: any) => {
  if (data.id) {
    return request.put({ url: `/admin/colleges/${data.id}`, data })
  }
  return request.post({ url: '/admin/colleges', data })
}

export const deleteDepartmentApi = (ids: string[] | number[]) => {
  return request.post({ url: '/admin/colleges/delete', data: { ids } })
}

export const getDepartmentTableApi = (params: any) => {
  return request.get({ url: '/admin/colleges/table/list', params })
}
