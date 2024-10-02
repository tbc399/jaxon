/** @type {import('tailwindcss').Config} */

const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: [
    './internal/**/*.{templ,go,html}',
    'node_modules/preline/dist/*.js'
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          '"Inter var", sans-serif', {fontFeatureSettings: '"ss01"'}
        ],
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('preline/plugin'),
    require('daisyui')
  ],
  daisyui: {
    darkTheme: "sunset",
    themes: [
        "sunset",
    ]
  }
}

