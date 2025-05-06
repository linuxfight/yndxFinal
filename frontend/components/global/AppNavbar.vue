<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import ThemeSwitch from "~/components/global/ThemeSwitch.vue";

const items = ref<NavigationMenuItem[][]>([
  [
    {
      label: 'Calc',
      icon: 'i-lucide-calculator',
      to: '/'
    }
  ]
])
</script>

<template>
  <header class="flex justify-between items-center p-6 dark:bg-dark text-3xl">

    <UNavigationMenu color="neutral" :items="items" class="w-full text-primary font-bold" />

    <ThemeSwitch/>
  </header>
</template>
