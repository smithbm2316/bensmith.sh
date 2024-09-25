import site from '#src/_data/site.js';

export default {
  layout: null,
  metadata: {
    title: `Ben Smith's Blog`,
    subtitle: `Ben's writings and thoughts about tech`,
    language: 'en',
    url: site.url,
    author: site.author,
  },
  eleventyImport: {
    collections: ['posts'],
  },
};
