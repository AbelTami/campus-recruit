import request from '@/axios'
import type { Application } from './types'

export const getApplicationList = (params: Recordable): Promise<IResponse<{ list: Application[]; total: number }>> =>
  request.get({ url: '/admin/applications', params })

export const updateApplicationStatus = (id: number, status: string, note?: string): Promise<IResponse<Application>> =>
  request.put({ url: `/admin/applications/${id}/status`, data: { status, note: note || '' } })

export const deleteApplication = (id: number): Promise<IResponse> =>
  request.delete({ url: `/admin/applications/${id}` })
