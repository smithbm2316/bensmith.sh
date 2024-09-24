export default function ({ eleventy }) {
  return {
    eleventyExcludeFromCollections: eleventy.env.runMode === 'build',
    draft: true,
  };
}
