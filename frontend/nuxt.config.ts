// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  $development: undefined, $env: undefined, $meta: undefined, $production: undefined, $test: undefined,
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  ssr: false,
  css: ['~/assets/main.css'],

  colorMode: {
    preference: 'system', // default theme based on OS preference
    fallback: 'light', // fallback theme if no OS preference
    classSuffix: '', // removes '-mode' suffix from class names
    storageKey: 'nuxt-color-mode' // localStorage key
  },

  app: {
    head: {
      title: 'Calc',
      htmlAttrs: {
        lang: 'ru',
      },
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/icon.svg' },
      ]
    }
  },

  modules: ['@nuxt/eslint', '@nuxt/fonts', '@nuxt/icon', '@nuxt/image', '@nuxt/ui', // '@nuxtjs/tailwindcss'
  '@nuxtjs/color-mode']
})