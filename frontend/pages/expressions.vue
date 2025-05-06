<script setup lang="ts">
import {type DtoApiError, type DtoListAllExpressionsResponse, getExpressions} from "~/client";
import { ref, onMounted, onBeforeUnmount } from 'vue';

const auth = useAuth();
const toast = useToast();
const response = ref<DtoListAllExpressionsResponse>();
const error = ref<DtoApiError>();
let updateInterval: NodeJS.Timeout;

if (!auth.checkTokenExists()) {
  navigateTo('/');
}

async function loadExpressions() {
  try {
    const { data, error: fetchError } = await getExpressions({
      headers: {
        Authorization: `Bearer ${auth.getToken()}`
      }
    });

    if (fetchError) throw fetchError;
    response.value = data;
    error.value = undefined;
  } catch (err) {
    toast.add({ title: 'Ошибка', description: `Ошибка: ${err}`, color: 'error' })
  }
}

onMounted(async () => {
  // Initial load
  await loadExpressions();

  // Set up periodic updates
  updateInterval = setInterval(async () => {
    await loadExpressions();
  }, 5000);
});

onBeforeUnmount(() => {
  // Cleanup interval when component is destroyed
  clearInterval(updateInterval);
});

function getColor(status: string): string {
  switch (status) {
    case 'DONE': return 'text-green-500 dark:text-green-400';
    case 'FAILED': return 'text-red-500 dark:text-red-400';
    case 'PROCESSING': return 'text-yellow-500 dark:text-yellow-400';
    default: return 'text-gray-500 dark:text-gray-400';
  }
}
</script>

<template>
  <UContainer class="p-20 dark:bg-dark-2">
    <!-- Loading state -->
    <div v-if="!response && !error" class="text-center">
      <USkeleton class="h-8 w-[300px] mb-4" />
      <USkeleton class="h-8 w-[200px]" />
    </div>

    <!-- Error state -->
    <div v-if="error" class="text-red-500">
      Error loading expressions: {{ error.message }}
    </div>

    <!-- Data loaded -->
    <ul v-if="response?.expressions">
      <li
          v-for="expression in response.expressions"
          :key="expression.id"
          class="p-4"
      >
        <div
            :class="[getColor(expression.status ?? 'PROCESSING'), 'border rounded-xl p-5 text-2xl transition-colors']"
        >
          <h1>{{ expression.id }} = {{ expression.result }}</h1>
        </div>
      </li>
    </ul>
  </UContainer>
</template>

<style scoped>
</style>