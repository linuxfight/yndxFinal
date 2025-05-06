import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
    darkMode: 'class',
    theme: {
        extend: {
            fontFamily: {
                regular: ['JetBrains Regular', 'sans-serif'],
                medium: ['JetBrains Medium', 'sans-serif'],
                bold: ['JetBrains Bold', 'sans-serif'],
            },

            colors: {
                'app-dark': '#46444D',
                'app-light': '#D9D9D9',
            }
        },
    }
}