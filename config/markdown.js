import markdownItAnchor from 'markdown-it-anchor';
import slugify from '@sindresorhus/slugify';

/** @param {import('@11ty/eleventy').UserConfig} eleventyConfig */
export default function markdownConfiguration(eleventyConfig) {
  // if a post's `draft` property is truthy, ignore processing it in prod
  eleventyConfig.addPreprocessor('drafts', 'md', ({ draft, eleventy }) => {
    if (draft && eleventy.env.runMode === 'build') {
      return false;
    }
  });

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
}
