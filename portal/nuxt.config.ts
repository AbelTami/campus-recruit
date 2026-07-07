export default defineNuxtConfig({
  compatibilityDate: '2025-07-07',
  devtools: { enabled: true },

  srcDir: '.',

  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxt/icon',
    'nuxt-auth-utils',
  ],

  icon: {
    clientBundle: {
      scan: true,
      include: [
        { prefix: 'heroicons', limit: 50 }
      ]
    }
  },

  app: {
    pageTransition: { name: 'page', mode: 'out-in' },
    head: {
      title: '大学生就业平台',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '大学生就业需求分析平台' },
      ],
    },
  },

  css: ['@/assets/css/scroll.css'],

  routeRules: {
    '/positions': { swr: 600 },
    '/positions/**': { swr: 600 },
    '/enterprises': { swr: 3600 },
    '/enterprises/**': { swr: 3600 },
  },

  nitro: {
    devProxy: {
      '/api/v1': { target: 'http://127.0.0.1:8080', changeOrigin: true }
    }
  }
})
