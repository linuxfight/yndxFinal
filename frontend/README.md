# Веб интерфейс

# Стек
Сделан на Nuxt/VueJS + NuxtUI + TailwindCSS + Zod.

Деплой выполняется на Cloudflare Pages через команду:
```shell
# pnpm
pnpm run deploy

# или npm
npm run deploy
```

Фронтэнд делает запросы к бэкэнду в облаке - https://calc-backend.lxft.tech.

Использование бэкэнда на localhost:8080 - не предусмотрено, т.к. при локальной разработке используется SwaggerUI.

# Разработка
1. Установите NodeJS (https://nodejs.org/)
```shell
# pnpm
pnpm install
pnpm run dev

# или npm
npm install
npm run dev
```

# Структура проекта
```shell
.
├── README.md (этот файл)
├── app.vue (главный файл, из него подключется layout)
├── assets (шрифты, стили, демо фото)
│ ├── JetBrainsMono-Bold.woff2
│ ├── JetBrainsMono-Medium.woff2
│ ├── JetBrainsMono-Regular.woff2
│ ├── demo.png
│ └── main.css
├── client (openapi клиент для взаимодействия с оркестратором)
│ ├── client.gen.ts
│ ├── index.ts
│ ├── sdk.gen.ts
│ └── types.gen.ts
├── components (здесь находятся компоненты для разных страниц)
│ ├── global
│ │ └── AppNavbar.vue
│ └── index
│     ├── InfoDemo.vue
│     ├── InfoFeature.vue
│     ├── InfoFooter.vue
│     └── InfoHero.vue
├── eslint.config.mjs (конфиг линтера)
├── nuxt.config.ts (конфиг фреймворка)
├── package.json (конфиг проекта)
├── pages (layout)
│ ├── auth.vue
│ ├── calculate.vue
│ ├── expressions.vue
│ └── index.vue
├── pnpm-lock.yaml (lockfile, как go.sum)
├── public (файлы, которые должны быть доступны напрямую по /)
│ ├── icon.svg
│ └── robots.txt
├── server
│ └── tsconfig.json
├── tailwind.config.ts (конфиг стилей)
├── tsconfig.json (конфиг TypeScript)
└── utils (утилиты для авторизации, взаимодействия с Cookie)
    └── auth.ts
```
