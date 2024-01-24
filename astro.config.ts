// https://astro.build/config
import { defineConfig } from 'astro/config';
/*
 * ASTRO INTEGRATIONS
 */
import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';
import icon from 'astro-icon';
/*
 * MARKDOWN PLUGINS
 */
import { injectFrontmatter } from './plugins/remarkPlugins';

// https://astro.build/config
export default defineConfig({
  integrations: [mdx({ remarkPlugins: [] }), sitemap(), icon()],
  markdown: {
    remarkPlugins: [injectFrontmatter],
    syntaxHighlight: 'shiki',
    shikiConfig: {
      theme: 'css-variables',
      wrap: false,
    },
  },
  output: 'static',
  prefetch: true,
  server: {
    port: 2323,
  },
  site: 'https://bensmith.sh',
  /* site:
    process.env.ASTRO_ENV === 'production'
      ? 'https://bensmith.sh'
      : 'http://localhost:2323', */
});
