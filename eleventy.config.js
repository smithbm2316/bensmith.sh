/**
 * The config object/instance for Eleventy
 * @typedef {import('@11ty/eleventy/src/UserConfig').default} EleventyConfig
 */

// 11ty plugins
import pluginDirectoryOutput from '@11ty/eleventy-plugin-directory-output';
import {
  EleventyRenderPlugin,
  InputPathToUrlTransformPlugin,
} from '@11ty/eleventy';
import pluginRSS from '@11ty/eleventy-plugin-rss';
import pluginSyntaxHighlight from '@11ty/eleventy-plugin-syntaxhighlight';
import pluginWebc from '@11ty/eleventy-plugin-webc';

// dependencies
import markdownItAnchor from 'markdown-it-anchor';
import slugify from '@sindresorhus/slugify';

/**
 * @param {EleventyConfig} eleventyConfig
 * @returns {Record<string, unknown>} Final configuration that we give to Eleventy
 */
export default function configureEleventy(eleventyConfig) {
  /** @see {@link https://www.11ty.dev/docs/config/#configuration-options Configuration Options} */
  const config = {
    dir: {
      input: 'src',
      output: '_site',
      includes: '_includes',
      data: '_data',
    },
    markdownTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    templateFormats: ['html', 'md', 'njk', '11ty.js', 'webc', 'css'],
  };

  // enable markdown headings getting turned into links automatically, using the same `slugify` function on the heading names as eleventy provides globally from `@sindresorhus/slugify`
  eleventyConfig.amendLibrary('md', (mdLib) => {
    mdLib.use(markdownItAnchor, {
      slugify,
      level: [2, 3, 4, 5, 6],
      permalink: markdownItAnchor.permalink.linkAfterHeader({
        style: 'aria-describedby',
        wrapper: [`<div class='post-heading-wrapper'>`, `</div>`],
      }),
    });
  });

  // create a collection that only contains the list of tags attached to posts
  eleventyConfig.addCollection('postTags', (collection) => {
    /** @type {Set<string>} */
    let tagsSet = new Set();
    let all = collection.getAll();
    let tagsToExclude = ['all', 'posts', 'postTags'];

    for (let item of all) {
      if (!item.data.tags) {
        continue;
      }
      for (let tag of item.data.tags) {
        if (tagsToExclude.includes(tag)) {
          continue;
        }
        tagsSet.add(tag);
      }
    }

    /** @type {string[]} */
    return Array.from(tagsSet).sort();
  });

  // datetime helpers
  // https://danabyerly.com/articles/time-is-on-your-side/
  eleventyConfig.addFilter('humanDate', (inputDate) => {
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'full',
      timeZone: 'UTC',
    }).format(new Date(inputDate));
  });
  eleventyConfig.addFilter('machineDate', (inputDate) => {
    return new Date(inputDate).toISOString();
  });
  eleventyConfig.addFilter('machineDatetime', (inputDate) => {
    return new Date(inputDate).toISOString();
  });

  // format page title helper
  eleventyConfig.addFilter('pageTitle', (title) => {
    if (!title || typeof title !== 'string') {
      return 'Ben Smith - Full Stack Web Developer';
    } else if (title.startsWith('custom:')) {
      return title.replace('custom:', '');
    } else {
      return `${title} - Ben Smith`;
    }
  });

  // configure the `src/assets` directory to be copied into our build without Eleventy processing the files. This is where all our fonts, images, styles, and other assets will go
  eleventyConfig.setServerPassthroughCopyBehavior('passthrough');
  eleventyConfig.addPassthroughCopy(`${config.dir.input}/assets`);

  /*
   * PLUGINS
   */
  eleventyConfig.addPlugin(EleventyRenderPlugin);
  eleventyConfig.addPlugin(InputPathToUrlTransformPlugin);

  eleventyConfig.setQuietMode(true);
  eleventyConfig.addPlugin(pluginDirectoryOutput);

  eleventyConfig.addPlugin(pluginWebc, {
    components: [
      // any `.webc` files found in the top-level of our `includes` directory or in the `components` directory inside of our `includes` directory will be processed as global webc components.
      `${config.dir.input}/${config.dir.includes}/components/*.webc`,
      // include <syntax-highlight> web component from 11ty plugin
      'npm:@11ty/eleventy-plugin-syntaxhighlight/*.webc',
    ],
  });

  eleventyConfig.addPlugin(pluginRSS);
  eleventyConfig.addPlugin(pluginSyntaxHighlight);

  /* END PLUGINS */

  eleventyConfig.setServerOptions({
    watch: [`${config.dir.input}/assets/styles.css`],
  });

  return config;
}
