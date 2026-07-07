import { defineStore } from 'pinia'
import { store } from '../index'
import { useStorage } from '@/hooks/web/useStorage'

const { getStorage, setStorage } = useStorage('localStorage')

const localeMap = [
  { lang: 'zh-CN', name: '简体中文', flag: '🇨🇳' },
  { lang: 'en', name: 'English', flag: '🇺🇸' },
  { lang: 'ja', name: '日本語', flag: '🇯🇵' },
]

function loadLang(): string {
  try { const v = getStorage('lang'); if (v && localeMap.some(l => l.lang === v)) return v }
  catch { /* ignore corrupt value */ }
  return 'zh-CN'
}

export const useLocaleStore = defineStore('locales', {
  state: () => {
    const lang = loadLang()
    return {
      currentLocale: {
        lang,
        flag: (localeMap.find(l => l.lang === lang) || localeMap[0]).flag,
      },
      localeMap,
    }
  },
  getters: {
    getCurrentLocale: (state) => state.currentLocale,
    getLocaleMap: () => localeMap,
  },
  actions: {
    setCurrentLocale(locale: { lang: string }) {
      this.currentLocale.lang = locale?.lang
      const entry = localeMap.find(l => l.lang === locale?.lang)
      if (entry) this.currentLocale.flag = entry.flag
      setStorage('lang', locale?.lang)
    },
  },
})

export const useLocaleStoreWithOut = () => useLocaleStore(store)
