<script setup lang="ts">
const auth = useAuth();
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

      <NuxtLink v-if="!auth.checkTokenExists()" to="/auth">
        <UButton
            icon="i-lucide-log-in"
            color="neutral"
            variant="outline"
            class="text-2xl p-5 rounded-xl">
          <h1>Авторизация</h1>
        </UButton>
      </NuxtLink>

      <NuxtLink v-if="auth.checkTokenExists()" to="/calculate">
        <UButton
            icon="i-lucide-calculator"
            color="neutral"
            variant="outline"
            class="text-2xl p-5 rounded-xl">
          <h1>Вычислить</h1>
        </UButton>
      </NuxtLink>
    </div>
    <div class="relative">
      <div class="bg-white p-4 rounded-xl shadow-xl">
        <img
            src="/assets/demo.png"
            alt="Calculator App Screenshot"
            class="rounded-lg w-full"
        >
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>