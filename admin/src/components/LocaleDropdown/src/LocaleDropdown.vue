<script setup lang="ts">
import { computed, unref } from 'vue'
import { ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'
import { useLocaleStore } from '@/store/modules/locale'
import { useLocale } from '@/hooks/web/useLocale'
import { propTypes } from '@/utils/propTypes'
import { useDesign } from '@/hooks/web/useDesign'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('locale-dropdown')

defineProps({
  color: propTypes.string.def('')
})

const localeStore = useLocaleStore()

const langMap = computed(() => localeStore.getLocaleMap)

const currentLang = computed(() => localeStore.getCurrentLocale)

const setLang = (lang: LocaleType) => {
  if (lang === unref(currentLang).lang) return
  localeStore.setCurrentLocale({ lang })
  useLocale().changeLocale(lang)
}

const labelMap: Record<string, string> = { 'zh-CN': '中文', en: 'English', ja: '日本語' }
</script>

<template>
  <ElDropdown :class="prefixCls" trigger="click" @command="setLang" :hide-on-click="true">
    <button class="inline-flex items-center gap-1.5 cursor-pointer px-3 py-1.5 rounded-full border border-gray-200 bg-white hover:border-gray-300 hover:shadow-sm transition-all text-gray-600 select-none text-[13px] font-medium">
      <span class="text-[15px] leading-none">{{ currentLang.flag }}</span>
      <span>{{ labelMap[currentLang.lang] || currentLang.lang }}</span>
      <Icon icon="ri:arrow-down-s-line" :size="14" class="text-gray-400 transition-transform duration-200" />
    </button>
    <template #dropdown>
      <ElDropdownMenu class="!rounded-xl !shadow-lg !shadow-gray-200/50 !border !border-gray-100 !min-w-[170px] !py-1.5 !px-1">
        <ElDropdownItem
          v-for="item in langMap"
          :key="item.lang"
          :command="item.lang"
          class="!rounded-lg !mb-0.5 last:!mb-0"
          :class="currentLang.lang === item.lang
            ? '!text-blue-600 !bg-blue-50/80 hover:!bg-blue-50'
            : 'hover:!bg-gray-50'"
        >
          <span class="inline-flex items-center gap-3 w-full py-0.5">
            <span class="text-base leading-none">{{ item.flag }}</span>
            <span class="text-[13px] flex-1">{{ item.name }}</span>
            <Icon v-if="currentLang.lang === item.lang" icon="ri:check-line" :size="18" class="shrink-0 text-blue-500" />
          </span>
        </ElDropdownItem>
      </ElDropdownMenu>
    </template>
  </ElDropdown>
</template>
