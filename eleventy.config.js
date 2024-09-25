import pluginDirectoryOutput from '@11ty/eleventy-plugin-directory-output';
import pluginRSS from '@11ty/eleventy-plugin-rss';
import pluginSyntaxHighlight from '@11ty/eleventy-plugin-syntaxhighlight';

import dir from '#config/dir.js';
import baseConfig from '#config/base.js';

/**
 * @param {import('@11ty/eleventy').UserConfig} eleventyConfig
 * @returns {Record<string, unknown>} Final configuration that we give to Eleventy
 */
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
  eleventyConfig.addPassthroughCopy(`${dir.input}/assets`);

  // ignore the `src/styles` directory so that tailwind manages it instead
  eleventyConfig.ignores.add(`${dir.input}/styles`);

  eleventyConfig.setServerOptions({
    // watch the compiled output of Tailwind
    watch: [`${dir.input}/assets/styles.css`],
  });

  return {
    dir,
    markdownTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    templateFormats: [
      'html',
      'md',
      'njk',
      '11ty.js',
      'webc',
      // copy over these files as plain text
      'css',
      'txt',
      'webmanifest',
    ],
  };
}
