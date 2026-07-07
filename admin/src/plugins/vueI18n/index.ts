import type { App } from 'vue'
import { i18n } from '@/i18n'

export { i18n }

export const setupI18n = (app: App<Element>) => {
  app.use(i18n)
}
