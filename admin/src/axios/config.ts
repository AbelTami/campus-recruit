import { AxiosResponse, InternalAxiosRequestConfig } from './types'
import { ElMessage, ElMessageBox } from 'element-plus'
import qs from 'qs'
import { SUCCESS_CODE, TRANSFORM_REQUEST_DATA } from '@/constants'
import { useUserStoreWithOut } from '@/store/modules/user'
import { objToFormData } from '@/utils'

const defaultRequestInterceptors = (config: InternalAxiosRequestConfig) => {
  if (
    config.method === 'post' &&
    config.headers['Content-Type'] === 'application/x-www-form-urlencoded'
  ) {
    config.data = qs.stringify(config.data)
  } else if (
    TRANSFORM_REQUEST_DATA &&
    config.method === 'post' &&
    config.headers['Content-Type'] === 'multipart/form-data' &&
    !(config.data instanceof FormData)
  ) {
    config.data = objToFormData(config.data)
  }
  if (config.method === 'get' && config.params) {
    let url = config.url as string
    url += '?'
    const keys = Object.keys(config.params)
    for (const key of keys) {
      if (config.params[key] !== void 0 && config.params[key] !== null) {
        url += `${key}=${encodeURIComponent(config.params[key])}&`
      }
    }
    url = url.substring(0, url.length - 1)
    config.params = {}
    config.url = url
  }
  return config
}

let tokenExpiredShowing = false

const defaultResponseInterceptors = (response: AxiosResponse) => {
  if (response?.config?.responseType === 'blob') {
    return response
  } else if (response.data.code === SUCCESS_CODE) {
    return response.data
  } else {
    if (response?.data?.code === 40100) {
      if (tokenExpiredShowing) return Promise.reject(response.data)
      tokenExpiredShowing = true
      const userStore = useUserStoreWithOut()
      ElMessageBox.alert('登录已过期，请重新登录', '提示', { confirmButtonText: '确定', type: 'warning' })
        .then(() => { userStore.logout() })
        .catch(() => {})
        .finally(() => { tokenExpiredShowing = false })
    } else {
      ElMessage.error(response?.data?.message)
    }
    return Promise.reject(response.data)
  }
}

export { defaultResponseInterceptors, defaultRequestInterceptors }
