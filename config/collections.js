/** @param {import('@11ty/eleventy').UserConfig} eleventyConfig */
export default function collectionsConfig(eleventyConfig) {
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
}
