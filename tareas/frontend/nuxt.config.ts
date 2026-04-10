export default defineNuxtConfig({
  devtools: { enabled: true },
  css: ['~/assets/css/bulma-overrides.scss'],

  modules: ['@pinia/nuxt'],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api'
    }
  },

  nitro: {
    routeRules: {
      '/api/**': {
        proxy: 'http://localhost:8080/**'
      }
    }
  },

  experimental: {
    appManifest: false
  },

  compatibilityDate: '2025-04-01'
})