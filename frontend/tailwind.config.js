/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          light: '#333333',
          DEFAULT: '#222222', // Deep Black/Gray
          dark: '#000000',
        },
        accent: {
          DEFAULT: '#A08355', // Luxury Gold/Bronze
          light: '#C5A875',
        },
        muted: {
          DEFAULT: '#F9F9F9',
          dark: '#E5E5E5'
        }
      },
      fontFamily: {
        serif: ['"Playfair Display"', 'serif'],
        sans: ['Lato', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
