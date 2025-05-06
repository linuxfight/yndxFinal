<script setup lang="ts">
defineProps<{
  link: string,
  icon: string,
  title: string,
  text: string,
}>()
</script>

<template>
  <a :href="link">
    <div class="bg-white dark:bg-dark p-6 rounded-xl shadow-sm hover:shadow-md transition-shadow">
      <UIcon :name="icon" class="text-4xl text-primary mb-4" />
      <h3 class="text-xl font-bold mb-2">{{title}}</h3>
      <p>
        {{text}}
      </p>
    </div>
  </a>
</template>

<style scoped>

</style>