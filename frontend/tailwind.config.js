/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        calc: {
          dark: '#0f172a',
          surface: '#1e293b',
          display: '#020617',
          key: '#334155',
          'key-hover': '#475569',
          accent: '#3b82f6',
          'accent-hover': '#2563eb',
          operator: '#f59e0b',
          'operator-hover': '#d97706',
          danger: '#ef4444',
          'danger-hover': '#dc2626',
        }
      }
    },
  },
  plugins: [],
}
