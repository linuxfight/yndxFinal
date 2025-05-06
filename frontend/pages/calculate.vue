<script setup lang="ts">
import * as z from 'zod'
import {postCalculate} from "~/client";
import type {FormSubmitEvent} from "#ui/types";

const toast = useToast();
const auth = useAuth();

const schema = z.object({
  expression: z.string().min(1),
});

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
  expression: undefined,
});

async function onSubmit(_: FormSubmitEvent<Schema>) {
  const resp = await postCalculate({
    headers: {
      Authorization: `Bearer ${auth.getToken()}`
    },
    body: {
      expression: state.expression!
    }
  });

  if (resp.error) {
    toast.add({ title: 'Ошибка', description: resp.error.message, color: 'error' });
    return;
  }

  toast.add({ title: 'Успешно!', description: `Выражение ${resp.data.id} отправлено на решение`, color: 'success' });
  navigateTo('/expressions');
}
</script>

<template>
  <div class="w-full h-full flex justify-center p-20">
    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
      <UFormField label="Выражение" name="expression">
        <UInput v-model="state.expression" class="w-full"/>
      </UFormField>

      <UButton class="w-full text-center" type="submit">
        Решить
      </UButton>
    </UForm>
  </div>
</template>

<style scoped>

</style>