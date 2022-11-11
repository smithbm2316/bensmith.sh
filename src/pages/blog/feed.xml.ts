import rss from '@astrojs/rss';
import { z } from 'zod';

// set up a zod schema for what content we will receive from the import.meta.glob execution below
const PostsSchema = z.array(
  z.object({
    url: z.string(),
    frontmatter: z.object({
      title: z.string(),
      pubDate: z.string(),
      draft: z.boolean().optional(),
    }),
    compiledContent: z.function().returns(z.string()),
    rawContent: z.function().returns(z.string()),
  }),
);

// import all markdown files with Vite's import.meta.glob synchronously and validate them
const posts = PostsSchema.parse(
  Object.values(
    import.meta.glob('./**/*.md', {
      eager: true,
    }),
  ),
);

export const get = () =>
  rss({
    title: 'Ben Smith\u2019s Blog',
    description: 'Ben\u2019s writings and thoughts about tech',
    // SITE will use "site" from your project's astro.config.
    site: import.meta.env.SITE,
    // filter out any draft posts, and then loop through each post to create our RSS feed
    items: posts
      .filter((post) => !post?.frontmatter?.draft)
      .map((post) => {
        // we will use the values below always, even if we don't have a description for the post
        const postData = {
          title: post.frontmatter.title,
          link: post.url,
          pubDate: new Date(post.frontmatter.pubDate),
        };

        // Parse the raw content of our post to get the first 140 characters for our description
        let description = post.rawContent().substring(0, 140);
        if (description) {
          // strip out any partial words by finding the last space in the string and appending a
          // '...' to it
          const lastSpaceIndex = description.lastIndexOf(' ');
          description = description.substring(0, lastSpaceIndex);
          description += '...';

          return {
            // include our common post data
            ...postData,
            description,
            // add a <content> section with the HTML output from our post, and encode it properly
            customData: `<content:encoded><![CDATA[${post.compiledContent()}]]></content:encoded>`,
          };
        } else {
          // if we don't have a description, then don't generate a <content> tag either
          return postData;
        }
      }),
    // (optional) inject custom xml
    customData: `<language>en-us</language>`,
    // inject the `xmlns:content` attribute with the namespace that defines how the
    // <content:encoded> element should work (as it's not part of the RSS 2.0 spec by default)
    xmlns: {
      content: 'http://purl.org/rss/1.0/modules/content/',
    },
  });
