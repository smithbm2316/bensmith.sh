// https://astro.build/config
import { defineConfig } from 'astro/config';
/*
 * ASTRO INTEGRATIONS
 */
import image from '@astrojs/image';
import mdx from '@astrojs/mdx';
import prefetch from '@astrojs/prefetch';
import sitemap from '@astrojs/sitemap';
import solidJs from '@astrojs/solid-js';
import tailwind from '@astrojs/tailwind';
/*
 * MARKDOWN PLUGINS
 */
// transform external links to use rel='nofollow' and open in new tab with target='_blank'
// import rehypeExternalLinks from './plugins/rehypeExternalLinks';

// https://astro.build/config
export default defineConfig({
  integrations: [
    tailwind({
      config: {
        applyBaseStyles: false,
      },
    }),
    prefetch(),
    mdx(),
    image(),
    solidJs(),
    sitemap(),
  ],
  output: 'static',
  markdown: {
    syntaxHighlight: 'shiki',
    shikiConfig: {
      theme: 'css-variables',
      wrap: true,
    },
  },
  server: {
    port: 2323,
  },
  vite: {
    ssr: {
      external: ['svgo'],
    },
  },
});
