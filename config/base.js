/*
 * Defines a base configuration to be consumed by my node.js server that runs
 * Eleventy with the programmatic API to generate server-side rendered routes
 */

// import any 11ty dependencies or plugins needed for the base config
import { InputPathToUrlTransformPlugin as pluginInputPathToUrl } from '@11ty/eleventy';
// import all modular configurations to add as eleventy plugins
import config from './index.js';
import collectionsConfig from './collections.js';
import filtersConfig from './filters.js';
import markdownConfig from './markdown.js';
import webcConfig from './webc.js';

/** @param {import('@11ty/eleventy').UserConfig} eleventyConfig */
export default function baseConfig(eleventyConfig) {
  eleventyConfig.setFrontMatterParsingOptions({ language: 'javascript' });

  eleventyConfig.addPlugin(collectionsConfig);
  eleventyConfig.addPlugin(filtersConfig);
  eleventyConfig.addPlugin(markdownConfig);
  eleventyConfig.addPlugin(webcConfig);

  eleventyConfig.addPlugin(pluginInputPathToUrl);

  return config;
}
