import { createI18n } from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import ja from 'element-plus/es/locale/lang/ja'
const elementLocales: Record<string, any> = { 'zh-CN': zhCn, en, ja }

const messages: Record<string, any> = {}

const langFiles = import.meta.glob('./locales/*/**/*.ts', { eager: true }) as Record<string, { default: any }>

for (const path in langFiles) {
  const match = path.match(/locales\/([\w-]+)\/(.+)\.ts$/)
  if (!match) continue
  const [, locale] = match
  if (!messages[locale]) messages[locale] = {}
  const mod = langFiles[path].default
  Object.assign(messages[locale], mod)
  // merge element-plus locale
  if (elementLocales[locale]) {
    messages[locale].el = elementLocales[locale]
  }
}

export const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages,
  silentTranslationWarn: true,
  silentFallbackWarn: true,
  missingWarn: false,
}) as any

export const availableLocales = Object.keys(messages).map(lang => ({
  lang,
  name: ({ 'zh-CN': '简体中文', en: 'English', ja: '日本語' } as Record<string, string>)[lang] || lang,
}))
