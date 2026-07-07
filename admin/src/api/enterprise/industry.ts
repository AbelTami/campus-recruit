import request from '@/axios'
import type { Industry } from './types'

export const getIndustryList = (): Promise<IResponse<Industry[]>> =>
  request.get({ url: '/admin/industries' })
