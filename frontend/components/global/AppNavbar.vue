<script setup lang="ts">
const colorMode = useColorMode();
const auth = useAuth();
const toast = useToast();

const isDark = computed({
  get() {
    return colorMode.value === 'dark'
  },
  set(_isDark) {
    colorMode.preference = _isDark ? 'dark' : 'light'
  }
});

function logOut() {
  toast.add({
    title: 'Ещё увидимся!',
    description: 'Вы вышли из аккаунта',
    color: 'warning'
  });
  auth.deleteToken();
}
</script>

<template>
  <header class="flex items-center justify-between p-6 dark:bg-dark text-3xl">
    <!-- Left-aligned item -->
    <div class="flex gap-5">
      <NuxtLink to="/">
        <UButton
            icon="i-lucide-house"
            color="neutral"
            variant="outline"
            class="text-2xl p-3 rounded-2xl">
          <h1>Главная</h1>
        </UButton>
      </NuxtLink>

      <NuxtLink v-if="auth.checkTokenExists()" to="/expressions">
        <UButton
            icon="i-lucide-list"
            color="neutral"
            variant="outline"
            class="text-2xl p-3 rounded-2xl">
          <h1>Все вычисления</h1>
        </UButton>
      </NuxtLink>
    </div>

    <!-- Right-aligned group -->
    <div class="ml-auto flex gap-5">
      <!-- Auth/Logout button -->

      <NuxtLink v-if="auth.checkTokenExists()" to="/">
        <UButton
            icon="i-lucide-log-out"
            color="neutral"
            variant="outline"
            class="text-2xl p-3 rounded-2xl"
            @click="logOut">
          <h1>Выйти</h1>
        </UButton>
      </NuxtLink>

      <!-- Theme toggle -->
      <ClientOnly v-if="!colorMode?.forced">
        <UButton
            :icon="!isDark ? 'i-lucide-moon' : 'i-lucide-sun'"
            color="neutral"
            variant="outline"
            class="text-2xl p-3 rounded-2xl"
            @click="isDark = !isDark"
        />
      </ClientOnly>
    </div>
  </header>
</template>

<style scoped>
</style>
