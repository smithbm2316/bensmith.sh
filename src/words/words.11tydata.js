import { z } from 'zod';

export default {
  layout: 'layouts/post.webc',
  tags: 'posts',
  /** @param {string} content - pass the page's parsed content to this directly */
  parseHeadings(content) {
    // pool all the article headings from the page content into an array so we can create a TOC easily
    let headings = [];
    let matches = content.matchAll(
      /^<h[23456].*>(?<heading>.*)<\/h[23456]>$/gm
    );
    for (let match of matches) {
      if (match.groups?.heading) {
        headings.push(match.groups.heading);
      }
    }
    return headings;
  },
  eleventyDataSchema(data) {
    // ignore the root page
    if (data.page.url === '/words/') return;

    // https://github.com/11ty/eleventy/issues/879#issuecomment-2062585031
    let result = z
      .object({
        draft: z.boolean().optional(),
        layout: z.string().min(1),
        title: z.string().min(1),
        date: z.string().min(1),
        tags: z.array(z.string()),
      })
      .safeParse(data);
    if (!result.success) {
      throw new Error(
        JSON.stringify(
          {
            inputPath: data.page.inputPath,
            fieldErrors: result.error.flatten().fieldErrors,
          },
          null,
          2
        )
      );
    }
  },
};
