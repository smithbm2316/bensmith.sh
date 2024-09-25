import site from '#src/_data/site.js';

/** @param {import('@11ty/eleventy').UserConfig} eleventyConfig */
export default function filtersConfig(eleventyConfig) {
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
      return site.title;
    } else if (title.startsWith('custom:')) {
      return title.replace('custom:', '');
    } else {
      return `${title} - Ben Smith`;
    }
  });
}
