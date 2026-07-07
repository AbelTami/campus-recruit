import request from '@/axios'

export const getRoleListApi = () => request.get({ url: '/admin/roles/list' })
export const createRoleApi = (data: any) => request.post({ url: '/admin/roles', data })
export const updateRoleApi = (id: number, data: any) => request.put({ url: `/admin/roles/${id}`, data })
export const deleteRoleApi = (id: number) => request.delete({ url: `/admin/roles/${id}` })
