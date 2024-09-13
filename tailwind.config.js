/** @type {import('tailwindcss').Config} */
module.exports = {
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

