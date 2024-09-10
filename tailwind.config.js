/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./components/**/*.{go,templ}",
    "./handlers/**/*.go",
    "./cmd/**/*.go",],
  theme: {
    extend: {},
  },
  plugins: [],
}

