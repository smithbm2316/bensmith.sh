export default {
  pagination: {
    data: 'collections',
    size: 1,
    alias: 'tag',
    // https://www.11ty.dev/docs/pagination/#filtering-values
    filter: ['all', 'posts', 'postTags'],
    // https://www.11ty.dev/docs/pagination/#add-all-pagination-pages-to-collections
    addAllPagesToCollections: true,
  },
  // https://github.com/11ty/eleventy-plugin-webc/issues/42#issuecomment-1374152043
  permalink({ tag }) {
    return `/tags/${this.slugify(tag)}/`;
  },
};
