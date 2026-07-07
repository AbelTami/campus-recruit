import request from '@/axios'
import type { Enterprise } from './types'

export const getEnterpriseList = (params: Recordable): Promise<IResponse<{ list: Enterprise[]; total: number }>> =>
  request.get({ url: '/admin/enterprises', params })

export const getEnterpriseById = (id: number): Promise<IResponse<Enterprise>> =>
  request.get({ url: `/admin/enterprises/${id}` })

export const createEnterprise = (data: Partial<Enterprise>): Promise<IResponse<Enterprise>> =>
  request.post({ url: '/admin/enterprises', data })

export const updateEnterprise = (id: number, data: Partial<Enterprise>): Promise<IResponse<Enterprise>> =>
  request.put({ url: `/admin/enterprises/${id}`, data })

export const deleteEnterprise = (id: number): Promise<IResponse> =>
  request.delete({ url: `/admin/enterprises/${id}` })

export const batchDeleteEnterprises = (ids: number[]): Promise<IResponse> =>
  request.post({ url: '/admin/enterprises/delete', data: { ids } })
