export default defineNuxtConfig({
  compatibilityDate: '2026-07-12',
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    backendBaseUrl: 'http://api',
  },
  typescript: {
    strict: true,
  },
})
