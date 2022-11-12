import rss from '@astrojs/rss';
import { PostsSchema } from '@models/Posts';
import sanitizeHtml from 'sanitize-html';

// import all markdown files with Vite's import.meta.glob synchronously and validate them
const posts = PostsSchema.parse(
  Object.values(
    import.meta.glob('./**/*.md', {
      eager: true,
    }),
  ),
);

// set up some custom XML tags to inject into the RSS feed
const customDataTags = [
  // enable Atom feed feature
  // prettier-ignore
  `<atom:link href="${import.meta.env.SITE}blog/feed.xml" rel="self" type="application/rss+xml" />`,
  // enable english language metadata
  `<language>en-us</language>`,
];

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
          pubDate: post.frontmatter.pubDate
            ? new Date(post.frontmatter.pubDate)
            : new Date(),
        };

        // check if we have a description or excerpt and use it for the RSS feed's description
        // property if we do
        const descriptionOrExcerpt =
          post.frontmatter.description ?? post.frontmatter.excerpt;
        if (descriptionOrExcerpt) {
          return {
            // include our common post data
            ...postData,
            description: descriptionOrExcerpt,
            // add a <content> section with the HTML output from our post, and encode it properly
            customData: `<content:encoded><![CDATA[${sanitizeHtml(
              post.compiledContent(),
            )}]]></content:encoded>`,
          };
        } else {
          // if we don't have a description, then don't generate a <content> tag either
          return postData;
        }
      }),
    // (optional) inject custom xml
    customData: customDataTags.join(''),
    // inject the `xmlns:content` attribute with the namespace that defines how the
    // <content:encoded> element should work (as it's not part of the RSS 2.0 spec by default)
    xmlns: {
      // enables Atom feed features
      // https://validator.w3.org/feed/docs/warning/MissingAtomSelfLink.html
      atom: 'http://www.w3.org/2005/Atom',
      content: 'http://purl.org/rss/1.0/modules/content/',
    },
  });
