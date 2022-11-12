import getReadingTime from 'reading-time';
import { toString } from 'mdast-util-to-string';
import { getPostExcerpt } from '../src/utils';

export function injectFrontmatter() {
  return function (tree: any, { data }: any) {
    // layout for post
    data.astro.frontmatter.layout = '@layouts/BlogPost.astro';

    const postContent = toString(tree);
    // total time to read the post
    const readingTime = getReadingTime(postContent);
    data.astro.frontmatter.readingTime = readingTime.text;

    // excerpt for post
    const excerpt = getPostExcerpt(postContent);
    if (excerpt) {
      data.astro.frontmatter.excerpt = excerpt;
    }
  };
}
