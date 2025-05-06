<template>
  <NuxtLayout>
    <AppNavbar/>
    <NuxtPage class="w-full h-full"/>
    <AppFooter/>
  </NuxtLayout>
</template>


<script setup lang="ts">
import AppNavbar from "~/components/global/AppNavbar.vue";
import AppFooter from "~/components/global/AppFooter.vue";
</script>

<style>
</style>