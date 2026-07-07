import { i18n } from '@/i18n'

export const useI18n = () => {
  if (!i18n) return { t: (key: string) => key }
  const { t } = i18n.global
  return { t }
}

export const t = (key: string) => key
