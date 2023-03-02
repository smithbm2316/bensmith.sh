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
/*
 * MARKDOWN PLUGINS
 */
import { injectFrontmatter } from './plugins/remarkPlugins';

// https://astro.build/config
export default defineConfig({
  integrations: [
    prefetch(),
    mdx({ extendPlugins: false, remarkPlugins: [] }),
    image({
      serviceEntryPoint: '@astrojs/image/sharp',
    }),
    solidJs(),
    sitemap(),
  ],
  markdown: {
    extendDefaultPlugins: true,
    remarkPlugins: [injectFrontmatter],
    syntaxHighlight: 'shiki',
    shikiConfig: {
      theme: 'css-variables',
      wrap: false,
    },
  },
  output: 'static',
  server: {
    port: 2323,
  },
  site: 'https://bensmith.sh',
  /* site:
    process.env.ASTRO_ENV === 'production'
      ? 'https://bensmith.sh'
      : 'http://localhost:2323', */
});
