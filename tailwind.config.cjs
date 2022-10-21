const plugin = require('tailwindcss/plugin');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        laserwave: {
          maximumBlue: '#40b4c4',
          hotPink: '#eb64B9',
          powderBlue: '#b4dce7',
          africanViolet: '#b381c5',
          pearlAqua: '#74dfc4',
          oldLavender: '#91889b',
          romanSilver: '#7b6995',
          mustard: '#ffe261',
          raisinBlack: '#27212e',
        },
        laserwaveContrast: {
          maximumBlue: '#1ed3ec',
          hotPink: '#ff52bf',
          powderBlue: '#acdfef',
          africanViolet: '#d887f5',
          pearlAqua: '#3feabf',
          oldLavender: '#b4abbe',
          romanSilver: '#b4a8c8',
          mustard: '#ffe261',
          raisinBlack: '#19151e',
        },
      },
    },
  },
  plugins: [],
  /* plugins: [
    plugin(function ({ addComponents, theme }) {
      addComponents({
        '.link': {
          boxShadow: 'inset 0 -0.25em 0 0 rgba(244, 114, 182, 0.45)',
          color: theme('colors.laserwave.hotPink'),
          transition: 'box-shadow 0.25s ease-in-out, color 0.25s linear',
          '&:hover': {
            boxShadow: 'inset 0 -1.5em 0 0 rgba(244, 114, 182, 0.45)',
            color: theme('colors.white'),
          },
          '&:focus-visible': {
            boxShadow: 'inset 0 -1.5em 0 0 rgba(244, 114, 182, 0.45)',
            color: theme('colors.white'),
          },
        },
      });
    }),
  ], */
};
