---
title: Add the content of your blog posts to your Astro RSS feed
tags: astro, rss, open-web
pubDate: November 11, 2022
showMeTheCode: true
draft: true
---
If you haven't seen it yet, [Astro](https://astro.build) is a *delightful* tool for building faster websites. It supports static-site generation (SSG), server-side rendering (SSR), as well as islands architecture for adding small pockets of interactivity to your websites with your favorite frontend Javascript framework of choice.

If you're sitting here wondering *what in the world did any of that just mean?*, don't worry! Astro is very approachable to web developers of all skill levels, whether you've barely started learning HTML and CSS, have been working in the industry for a little while (like myself), or are a greybeard of the web development scene. I recommend you start with [their tutorial here](https://docs.astro.build/en/tutorial/0-introduction/). 

However, the topic of this post is not targeted at beginners; rather this post is directed at existing Astro users or existing web developers new to Astro that want to use the RSS integration on their site. I recently re-wrote my personal site (the one you're currently on!) with Astro and wanted a slightly different implementation to the RSS feed than the one that their integration currently supplies out-of-the-box. If that sounds interesting to you, read on!

## The @astrojs/rss package

[Astro‘s RSS package](https://docs.astro.build/en/guides/rss) will generate an RSS feed for you, based upon the glob of markdown files that you supply to it. Depending on exactly how you configure it, Astro can generate a title, description, published date, and additional custom XML tags that you supply to Astro in the form of a raw XML string. It will not include the content of each of your posts into the feed itself, since that is not supported by the latest RSS specification (v2.0). The RSS v2.0 spec is now well over a decade old, but thankfully just like HTML it was created with the intention of being extended in a forgiving manner.

Now you might be asking yourself why would you even want the content of your post to be included in your feed? For many users simply having the title and a short description of your post will be enough. Since RSS is my favorite way to consume blog posts and other media nowadays, I prefer feeds that include the full content of the post. This enables me to read the full content of the post from the comfort of my favorite RSS reader. It also allows more tech-savvy users to build their own RSS client with whatever styling they might want.

### A short love-letter to RSS

RSS is a widely-celebrated tenant of the Open Web and Indie Web. The ability to modify the content of any webpage that you visit inside your browser with custom Javascript functionality or CSS styles is one of my favorite parts of being a web developer and building applications for the web. If a website does not support dark mode, I can choose to inject custom CSS onto the page in order to make whatever site I’m reading at the moment more pleasing to my eye. the same is true for the blog posts that I subscribe to in my RSS reader. That feature of RSS is just one of many that I love, on top of the incredible feature of being able to subscribe to *exactly* what feeds you want to instead of having an algorithm designed to keep you doomscrolling supplies you.


### Extending the RSS v2.0 specification

Some of you might be wondering why I don’t just inject the full post content into the `<description>`, since that is part of the RSS v2.0 standard and has wide support by many RSS readers. From the reading that I’ve done from [the people who maintain the RSS spec](https://www.rssboard.org/rss-specification#hrelementsOfLtitemgt), they recommend to avoid this. The descriptions are supposed to be used for a short piece of text that is either an excerpt from the post or describes what is contained in the post briefly. In addition to that, I still would like to be able to have the ability to add my own handwritten description for longer posts that would be well-served by a summary. Thankfully, the [RSS spec maintainers encourage extending the spec](https://www.rssboard.org/rss-specification#extendingRss) to fit the use case that different developers and users might have, and the extension that I used seems to be well-supported by many feed readers.

The way that we can extend the specifications is by adding some metadata for a “namespace”. This allows a RSS feed to give feed readers extra instructions on how to parse a non-standard tag. In the code snippet that I will show below, you’ll see the name space that is attached to my RSS feed in order to support the [new \<content\> tag](https://www.rssboard.org/rss-profile#namespace-elements-content). The `<content:encoded>` tag expects to contain valid HTML, not just raw text. This gives feed readers and the end-user/consumer of the content more flexibility to style their content how they see fit.

## Show me the code!

```typescript
// src/pages/feed.xml.ts
import rss from '@astrojs/rss';
// We'll assume that our feed is located at src/pages/feed.xml and our blog posts are all located in the src/pages/blog folder

// import all markdown files with Vite's import.meta.glob synchronously from our blog folder
const posts = Object.values(
  import.meta.glob('./blog/**/*.md', {
    eager: true,
  }),
);

// SITE will use "site" from your project's astro.config.
const SITE = import.meta.env.SITE;

// set up some custom XML tags to inject into the RSS feed
const customDataTags = [
  // enable Atom feed, as some RSS readers use that format
  // https://www.fpds.gov/wiki/index.php/FAADC_Atom_Feed_Specifications_V_1.0
  `<atom:link href="${SITE}feed.xml" rel="self" type="application/rss+xml" />`,
  // enable language metadata
  `<language>en-us</language>`,
];

export const get = () =>
  rss({
    // \u2019 is the unicode code for an apostrophe
    title: 'Ben Smith\u2019s Blog',
    description: 'Ben\u2019s writings and thoughts about tech',
    site: SITE,
    items: posts
      // remove any draft posts
      .filter((post) => !post?.frontmatter?.draft)
      // loop through each post to create our RSS feed
      .map((post) => {
        // we will use the values below always, even if we don't have a description for the post
        const postData = {
          title: post.frontmatter.title,
          link: post.url,
          pubDate: new Date(post.frontmatter.pubDate),
        };

        // check if we have a description or excerpt property in our posts's frontmatter, and use it
        // for the RSS feed's description property if we do
        const descriptionOrExcerpt =
          post.frontmatter.description ?? post.frontmatter.excerpt;
        if (descriptionOrExcerpt) {
          return {
            // include our common post data
            ...postData,
            description: descriptionOrExcerpt,
            // add a <content> section with the HTML output from our post, and encode it properly
            // https://docs.astro.build/en/guides/rss/#1-importmetaglob-result
            customData: `<content:encoded><![CDATA[${post.compiledContent()}]]></content:encoded>`,
          };
        } else {
          // if we don't have a description, then don't generate a <content> tag or description
          return postData;
        }
      }),
    // inject custom tags defined above as a string so that we have support for the Atom feed
    // standard and give RSS readers information about what language our posts are in
    customData: customDataTags.join(''),
    // inject the `xmlns:content` attribute with the namespace that defines how the
    // <content:encoded> element should work (as it's not part of the RSS 2.0 spec by default)
    xmlns: {
      // the namespace that enables Atom feed
      atom: 'http://www.w3.org/2005/Atom',
      // the namespace that enables the <content:encoding> tag
      content: 'http://purl.org/rss/1.0/modules/content/',
    },
  });
```

Thanks for reading, I hope that this post was helpful! If it was, please let me know on [Twitter](https://twitter.com/smithbm2316) or [Mastodon](https://fosstodon.org/@smithbm2316)! I'd love to hear any feedback and if this post helped you get up and running with content in your RSS feed on your Astro site.
