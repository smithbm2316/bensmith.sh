/**
 * @param {string} field
 * @param {string[]} validTypes
 * @param {string} inputPath
 * @throws
 */
function throwValidationError(field, validTypes, inputPath) {
  throw new Error(
    `Invalid frontmatter data for: ${inputPath}. ${field} is not ${validTypes.join(' or ')}`
  );
}

export default {
  layout: 'layouts/post.webc',
  tags: 'posts',
  eleventyDataSchema(data) {
    // ignore the root page
    if (data.page.url === '/words/') return;

    // https://github.com/11ty/eleventy/issues/879#issuecomment-2062585031
    if (data.draft !== undefined && typeof data.draft !== 'boolean') {
      throwValidationError(
        'draft',
        ['boolean', 'undefined'],
        data.page.inputPath
      );
    } else if (typeof data.layout !== 'string') {
      throwValidationError('layout', ['string'], data.page.inputPath);
    } else if (typeof data.title !== 'string') {
      throwValidationError('title', ['string'], data.page.inputPath);
    } else if (typeof data.date !== 'string') {
      throwValidationError('title', ['string'], data.page.inputPath);
    } else if (data.tags !== undefined && !Array.isArray(data.tags)) {
      throwValidationError(
        'title',
        ['string[]', 'undefined'],
        data.page.inputPath
      );
    }
  },
};
