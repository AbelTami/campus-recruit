import request from '@/axios'
import type { Position } from './types'

export const getPositionList = (params: Recordable): Promise<IResponse<{ list: Position[]; total: number }>> =>
  request.get({ url: '/admin/positions', params })

export const getPositionById = (id: number): Promise<IResponse<Position>> =>
  request.get({ url: `/admin/positions/${id}` })

export const createPosition = (data: Partial<Position>): Promise<IResponse<Position>> =>
  request.post({ url: '/admin/positions', data })

export const updatePosition = (id: number, data: Partial<Position>): Promise<IResponse<Position>> =>
  request.put({ url: `/admin/positions/${id}`, data })

export const deletePosition = (id: number): Promise<IResponse> =>
  request.delete({ url: `/admin/positions/${id}` })

export const batchDeletePositions = (ids: number[]): Promise<IResponse> =>
  request.post({ url: '/admin/positions/delete', data: { ids } })
