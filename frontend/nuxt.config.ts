// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: false },
  ssr: false,
  debug: false,
  css: ["~/assets/css/tailwind.css"],
  modules: [
    "@nuxtjs/tailwindcss",
    "shadcn-nuxt",
    "@nuxtjs/color-mode",
    "@formkit/auto-animate/nuxt",
    "@nuxt/icon",
    "@vueuse/motion/nuxt",
  ],
  piniaPluginPersistedstate: {
    storage: "localStorage",
  },
  devServer: {
    port: 3000,
    host: "0.0.0.0",
  },
  nitro: {
    preset: "static",
  },
  colorMode: {
    preference: "light",
    fallback: "light",
    classSuffix: "",
  },
  shadcn: {
    prefix: "",
    componentDir: "@/components/ui",
  },
  runtimeConfig: {
    public: {
      pbUrl: process.env.PB_URL || "http://localhost:8090",
      apiUrl: process.env.API_URL || "http://localhost:8313",
    },
  },
});

