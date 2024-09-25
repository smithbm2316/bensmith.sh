import site from './site.js';

/**
 * Metadata for RSS/Atom/JSON feeds
 */
const metadata = /** @type {const} */ ({
  title: `Ben Smith's Blog`,
  subtitle: `Ben's writings and thoughts about tech`,
  language: 'en',
  url: site.url,
  author: site.author,
});

export default metadata;
