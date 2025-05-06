<template>
  <div class="h-full w-full">
    <!-- Hero Section -->
    <section class="p-20 bg-white dark:bg-dark">
      <InfoHero/>
    </section>

    <!-- Features Section -->
    <section id="features" class="p-20 bg-light dark:bg-dark-2">
      <div class="text-center mb-16">
        <h2 class="text-3xl font-bold mb-4">Дополнительный функционал</h2>
        <p class="max-w-2xl mx-auto">
          Что реализовано?
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">
        <InfoFeature
            link="https://calc-backend.lxft.tech/docs"
            title="Swagger (тык)"
            text="Документация API"
            icon="i-lucide-scroll-text"
        />

        <InfoFeature
            link="https://calc.lxft.tech"
            title="Деплой (тык)"
            text="Доступно на тут"
            icon="i-lucide-cloud"
        />

        <InfoFeature
            link="https://github.com/linuxfight/yndxFinal/actions"
            title="CI (тык)"
            icon="i-lucide-repeat"
            text="Тесты и билд"
        />
      </div>
    </section>


    <section id="demo" class="py-20 bg-light dark:bg-dark-2">
      <InfoDemo/>
    </section>

    <InfoFooter/>
  </div>
</template>

<script setup lang="ts">
import InfoDemo from "~/components/index/InfoDemo.vue";
import InfoFeature from "~/components/index/InfoFeature.vue";
import InfoHero from "~/components/index/InfoHero.vue";
import InfoFooter from "~/components/index/InfoFooter.vue";
</script>

<style scoped>
</style>