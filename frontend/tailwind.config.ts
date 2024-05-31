import type { Config } from "tailwindcss";
import plugin from 'tailwindcss/plugin'
// @ts-expect-error
import flattenColorPalette from 'tailwindcss/lib/util/flattenColorPalette'

function addVariablesForColors({ addBase, theme }: any) {
  // const allColors = flattenColorPalette(theme('colors'))
  // const newVars = Object.fromEntries(
  //   Object.entries(allColors).map(([key, val]) => [`--${key}`, val])
  // )
  // addBase({
  //   ':root': newVars
  // })
}


const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {},
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    plugin(function ({ addUtilities }) {
      addUtilities({
        /**
         * Mimics the deprecated word-wrap: break-word; property (see: https://drafts.csswg.org/css-text-3/#word-break-property).
         *
         * Prefer Tailwinds `word-break` and only use this if soft wrap opportunities should be considered
         * (https://developer.mozilla.org/en-US/docs/Web/CSS/overflow-wrap).
         */
        '.hn-break-words': {
          'word-break': 'normal',
          'overflow-wrap': 'anywhere'
        }
      })
    }),
    addVariablesForColors
  ],
};
export default config;
