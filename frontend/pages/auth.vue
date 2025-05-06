<script setup lang="ts">
import * as z from 'zod'
import {postLogin, postRegister} from "~/client";

const toast = useToast();
const auth = useAuth();

const schema = z.object({
  login: z.string().min(1),
  password: z.string().min(1)
});

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
  login: undefined,
  password: undefined
});

async function handleLogin() {
  if (state.login !== undefined && state.password !== undefined) {
    const resp = await postLogin({
      body: {
        login: state.login!,
        password: state.password!
      }
    });

    if (resp.error !== undefined) {
      toast.add({ title: 'Ошибка', description: resp.error.message, color: 'error' });
      return;
    }

    if (resp.data.token == undefined) {
      return;
    }

    auth.setToken(resp.data.token);
    toast.add({
      title: 'Добро пожаловать!',
      description: 'Вы успешно вошли!',
      color: 'success'
    });
    navigateTo('/');
  }
}

async function handleRegister() {
  if (state.login !== undefined && state.password !== undefined) {
    const resp = await postRegister({
      body: {
        login: state.login!,
        password: state.password!
      }
    });

    if (resp.error !== undefined) {
      toast.add({ title: 'Ошибка', description: resp.error.message, color: 'error' });
      return;
    }

    if (resp.data.token == undefined) {
      return;
    }

    auth.setToken(resp.data.token);
    toast.add({
      title: 'Добро пожаловать!',
      description: 'Вы успешно зарегистрировались!',
      color: 'success'
    });
    navigateTo('/');
  }
}
</script>

<template>
  <div class="w-full h-full flex justify-center p-20">
    <UForm :schema="schema" :state="state" class="space-y-4">
      <UFormField label="Логин" name="login">
        <UInput v-model="state.login" class="w-full"/>
      </UFormField>

      <UFormField label="Пароль" name="password">
        <UInput v-model="state.password" type="password" class="w-full"/>
      </UFormField>

      <UButton class="w-full text-center" @click="handleLogin">
        Войти
      </UButton>

      <UButton class="w-full text-center" @click="handleRegister">
        Зарегистрироваться
      </UButton>
    </UForm>
  </div>
</template>

<style scoped>

</style>