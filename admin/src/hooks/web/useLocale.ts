import { i18n } from '@/i18n'
import { useLocaleStoreWithOut } from '@/store/modules/locale'

export const useLocale = () => {
  const localeStore = useLocaleStoreWithOut()

  const changeLocale = (locale: string) => {
    i18n.global.locale.value = locale
    localeStore.setCurrentLocale({ lang: locale })
    document.documentElement.lang = locale
  }

  return { changeLocale }
}
