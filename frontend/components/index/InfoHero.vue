<script setup lang="ts">

</script>

<template>
  <div class="grid md:grid-cols-2 gap-12 items-center">
    <div>
      <h1 class="text-4xl md:text-5xl font-bold mb-4">
        Калькулятор
      </h1>
      <p class="text-lg mb-8">
        API и веб-интерфейс для решения математических выражений
      </p>
      <div class="flex flex-wrap gap-4">
        <UButton variant="outline" size="lg" to="#demo">
          Демо API
          <template #trailing>
            <UIcon name="i-lucide-play" />
          </template>
        </UButton>
        <UButton>
          FUNCTIONS
        </UButton>
      </div>
    </div>
    <div class="relative">
      <div class="bg-white p-4 rounded-xl shadow-xl">
        <img
            src="https://placehold.co/600x400/e2e8f0/475569?text=Calculator+App+Screenshot"
            alt="Calculator App Screenshot"
            class="rounded-lg w-full"
        >
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>