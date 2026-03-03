/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: 'var(--color-primary)',
        'primary-hover': 'var(--color-primary-hover)',
        secondary: 'var(--color-secondary)',
        bg: {
          base: 'var(--color-bg-base)',
          surface: 'var(--color-bg-surface)',
          accent: 'var(--color-bg-accent)',
        },
        border: {
          base: 'var(--color-border-base)',
        },
        text: {
          base: 'var(--color-text-base)',
          muted: 'var(--color-text-muted)',
        },
      }
    },
  },
  plugins: [
    require('@tailwindcss/line-clamp'),
  ],
}