export default function ({ eleventy }) {
  if (eleventy.env.runMode === 'build') {
    return {
      permalink: false,
      eleventyExcludeFromCollections: true,
    };
  }
}
