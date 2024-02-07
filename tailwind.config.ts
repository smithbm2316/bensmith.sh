import { type Config } from 'tailwindcss';
import plugin from 'tailwindcss/plugin';

const config: Config = {
  content: ['./src/**/*.{astro,md,mdx}'],
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
  plugins: [
    // overwrite your post draft with the draft inside of obsidian
    // cp -iv ~/documents/notes/blog-drafts/tailwind-is-good-enough.md src/pages/blog/tailwind-is-good-enough.md
    // add variants to allow us to apply styles only to the kiosk or web app
    plugin(function ({ addVariant }) {
      addVariant('kiosk', ':is([data-rendering-on="kiosk"] &)');
      addVariant('web', ':is([data-rendering-on="web"] &)');
    }),
  ],
};

export default config;
