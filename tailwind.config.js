/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/views/**/*.{templ,html,js}"],
  theme: {
    colors: {
      primary: "#3B71CA",
      secondary: "#9FA6B2",

      success: "#14A44D",

      danger: "#DC4C64",
      warning: "#E4A11B",

      info: "#54B4D3",

      light: "#F0FAFB",
      white: "#FAFAFA",

      dark: "#1F2937",

    },
    extend: {},
  },
  /* ?
      require('@tailwindcss/forms'),
      require('@tailwindcss/typography'),
  */
  plugins: [],
}