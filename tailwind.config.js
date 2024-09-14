/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'selector',
  content: [
    "./components/**/*.{go,templ}",
    "./handlers/**/*.go",
    "./static/**/*.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

