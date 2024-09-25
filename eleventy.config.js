import pluginDirectoryOutput from '@11ty/eleventy-plugin-directory-output';
import pluginRSS from '@11ty/eleventy-plugin-rss';
import pluginSyntaxHighlight from '@11ty/eleventy-plugin-syntaxhighlight';

import config from '#config/index.js';
import baseConfig from '#config/base.js';

/** @param {import('@11ty/eleventy').UserConfig} eleventyConfig */
export default function configureEleventy(eleventyConfig) {
  eleventyConfig.addPlugin(baseConfig);

  eleventyConfig.setQuietMode(true);
  eleventyConfig.addPlugin(pluginDirectoryOutput);
  eleventyConfig.addPlugin(pluginSyntaxHighlight);
  eleventyConfig.addPlugin(pluginRSS);

  // configure the `src/assets` directory to be copied into our build without
  // Eleventy processing the files. This is where all our fonts, images,
  // styles, and other assets will go.
  eleventyConfig.setServerPassthroughCopyBehavior('passthrough');
  eleventyConfig.addPassthroughCopy(`${config.dir.input}/assets`);

  // ignore the `src/styles` directory so that tailwind manages it instead
  eleventyConfig.ignores.add(`${config.dir.input}/styles`);

  eleventyConfig.setServerOptions({
    // watch the compiled output of Tailwind
    watch: [`${config.dir.input}/assets/styles.css`],
    port: 2323,
  });

  return config;
}
