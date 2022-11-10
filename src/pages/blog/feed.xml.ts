import rss from '@astrojs/rss';

interface Post {
  url: string;
  frontmatter: {
    title: string;
    pubDate?: Date;
    draft?: boolean;
  };
}

const postImportResult = import.meta.glob('./**/*.{md, mdx}', { eager: true });
const posts = Object.values(postImportResult) as Post[];

export const get = () =>
  rss({
    // `<title>` field in output xml
    title: 'Ben Smith\u2019s Blog',
    // `<description>` field in output xml
    description: 'Ben\u2019s writings and thoughts about tech',
    // base URL for RSS <item> links
    // SITE will use "site" from your project's astro.config.
    site: import.meta.env.SITE,
    // list of `<item>`s in output xml
    items: posts
      .filter((post) => !post.frontmatter.draft)
      .map((post) => ({
        link: post.url,
        title: post.frontmatter.title,
        pubDate: post.frontmatter.pubDate ?? new Date(),
      })),
    // (optional) inject custom xml
    customData: `<language>en-us</language>`,
    stylesheet: '/pretty-feed-v3.xsl',
  });
