import pluginWebc from '@11ty/eleventy-plugin-webc';
import config from './index.js';

/** @param {import('@11ty/eleventy').UserConfig} eleventyConfig */
export default function webcConfig(eleventyConfig) {
  // if a webc template's `ssr` property is truthy, ignore processing in prod
  // cli builds (not in our usage from a script with the programmatic api)
  eleventyConfig.addPreprocessor('ssr', 'webc', ({ ssr, eleventy }) => {
    if (
      ssr &&
      eleventy.env.runMode === 'build' &&
      eleventy.env.source === 'cli'
    ) {
      return false;
    }
  });

  eleventyConfig.addPlugin(pluginWebc, {
    // any `.webc` files found in the top-level of our `includes` directory
    // or in the `components` directory inside of our `includes` directory
    // will be processed as global webc components.
    components: `${config.dir.input}/${config.dir.includes}/components/*.webc`,
  });
}
