import getReadingTime from 'reading-time';
import { toString } from 'mdast-util-to-string';

export function injectFrontmatter() {
  return function (tree: any, { data }: any) {
    // layout for post
    data.astro.frontmatter.layout = '@layouts/BlogPost.astro';

    const textOnPage = toString(tree);
    const readingTime = getReadingTime(textOnPage);
    data.astro.frontmatter.readingTime = readingTime.text;
  };
}
