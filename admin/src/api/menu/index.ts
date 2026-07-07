import request from '@/axios'

export const getMenuListApi = () => request.get({ url: '/admin/menus' })
export const getMenuFlatApi = () => request.get({ url: '/admin/menus/flat' })
export const createMenuApi = (data: any) => request.post({ url: '/admin/menus', data })
export const updateMenuApi = (id: number, data: any) => request.put({ url: `/admin/menus/${id}`, data })
export const deleteMenuApi = (id: number) => request.delete({ url: `/admin/menus/${id}` })
