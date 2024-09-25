import { InputPathToUrlTransformPlugin as pluginInputPathToUrl } from '@11ty/eleventy';

import dir from './dir.js';
import collectionsConfig from './collections.js';
import filtersConfig from './filters.js';
import markdownConfig from './markdown.js';
import webcConfig from './webc.js';

/**
 * @param {import('@11ty/eleventy').UserConfig} eleventyConfig
 * @returns {Record<string, unknown>} Final configuration that we give to Eleventy
 */
export default function baseConfig(eleventyConfig) {
  // set defualt frontmatter data template format to Javascript instead of YAML
  eleventyConfig.setFrontMatterParsingOptions({ language: 'javascript' });

  eleventyConfig.addPlugin(collectionsConfig);
  eleventyConfig.addPlugin(filtersConfig);
  eleventyConfig.addPlugin(markdownConfig);
  eleventyConfig.addPlugin(webcConfig);

  eleventyConfig.addPlugin(pluginInputPathToUrl);

  return {
    dir,
    markdownTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    templateFormats: ['html', 'md', 'njk', '11ty.js', 'webc'],
  };
}
