<script setup lang="ts">
import AppNavbar from "~/components/global/AppNavbar.vue";
</script>

<template>
  <UApp class="w-full h-full">
    <NuxtLayout>
      <AppNavbar/>
      <NuxtPage class="w-full h-full"/>
    </NuxtLayout>
  </UApp>
</template>

<style>
</style>