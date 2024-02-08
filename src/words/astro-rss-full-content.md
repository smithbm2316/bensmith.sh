---
layout: layouts/post.webc
title: Add the content of your blog posts to your Astro RSS feed
date: 2022-12-09
showMeTheCode: true
tags:
  - astro
  - rss
  - open-web
---

*Note: the content of this post is likely out of date as Astro has had multiple major version releases since the time that this post was written.*

If you haven't seen it yet, [Astro](https://astro.build) is a *delightful* tool for building faster websites. It supports static-site generation (SSG), server-side rendering (SSR), as well as islands architecture for adding small pockets of interactivity to your websites with your frontend Javascript framework of choice.

If you're sitting here wondering *what in the world did any of that just mean?*, don't worry! Astro is approachable to web developers of all skill levels. Whether you've barely started learning HTML and CSS, are a veteran of the web development scene, or are somewhere inbetween like myself, Astro will welcome you with open arms. I recommend you start with [their tutorial here](https://docs.astro.build/en/tutorial/0-introduction/). Their documentation team sets the bar *incredibly high* for what the docs of a modern Javascript framework should look like.

I will preface this post with the following: the content of this post is *not* targeted at beginners. This post is directed at existing Astro users or existing web developers new to Astro that want to add an RSS feed to their shiny new Astro site with the `@astrojs/rss` package. I recently re-wrote my personal site (the one you're currently on) with Astro and wanted a slightly different implementation to the RSS feed than the one that their integration *used* to supply out-of-the-box. As of [`@astrojs/rss` v1.1.0](https://github.com/withastro/astro/releases/tag/%40astrojs%2Frss%401.1.0) I was able to upstream the main feature of my custom solution directly into the `@astrojs/rss` package itself ([the docs for that update are here](https://docs.astro.build/en/guides/rss/#including-full-post-content)). We'll cover both how to build out my original solution as well as how to implement the new addition to the `@astrojs/rss` package. If that sounds interesting to you, read on!

## The @astrojs/rss package

Before we cover the implementation though, let's take a moment to talk about [Astro‘s RSS package](https://docs.astro.build/en/guides/rss/). It will generate a RSS feed for you based upon a glob of markdown files that you supply to it with [Vite's `import.meta.glob()`](https://docs.astro.build/en/guides/rss/#generating-items) function. Depending on exactly how you configure it, Astro can generate a title, description, published date, and additional custom XML tags that you supply to the `rss()` helper function. It will not include the content of each of your posts into the feed itself, since that is not supported by the latest RSS specification (v2.0). While that spec is now well over a decade old, likje HTML it was created with the intention of being extended in a forgiving manner.

Now you might be asking yourself *why would you even want* the content of your post to be included in your feed? For many users simply having the title and a short description of your post is enough. Since RSS is my favorite way to consume blog posts and other forms of content, I prefer feeds that include the content of the post so that I can read it from the comfort of my favorite RSS reader. I particularly love using [Feedbin](https://feedbin.com), which has a nice minimal styling to it as well as syntax highlighting for HTML `<code>` blocks.

### A brief love letter to RSS

RSS is a widely celebrated tenant of the [Open Web](https://www.w3.org/wiki/Open_Web_Platform) and [Indie Web](https://indieweb.org/). The ability to modify the content of any webpage that you visit inside your browser with custom Javascript functionality or CSS styles is one of my favorite parts of being a web developer and building applications for the web. If a website does not support dark mode, I can choose to inject custom CSS onto the page in order to make whatever site I’m reading at the moment more pleasing to my eye. The same is true for the blog posts that I subscribe to in my RSS reader. The ability to *choose exactly what feeds I subscribe to* is the other main feature of RSS that brings me so much joy. Instead of using a product like Twitter or Facebook that have algorithms *designed* to recommend and curate content for me, I'm able to control *exactly* what feeds, content creators, and other media I subscribe to in my RSS reader.

### Extending the RSS v2.0 specification

Now let me circle back to the actual topic of this post. Some of you more familiar with the RSS spec might be wondering why I don’t just inject the full post content into the `<description>` tag, since that is part of the RSS v2.0 standard and has wide support by many RSS readers. From the reading that I’ve done from [the people who maintain the RSS spec](https://www.rssboard.org/rss-specification#hrelementsOfLtitemgt), they recommend to avoid this. The descriptions are supposed to be used for a short piece of text that is either an excerpt from the post or description of the post's content. In addition to that, I still would like to be able to have the ability to add my own handwritten descriptions for longer posts that would be well-served by a summary. Luckily for us the [RSS spec maintainers encourage extending the spec](https://www.rssboard.org/rss-specification#extendingRss) to fit the use cases that different developers and users might have. The extension that I ended up using seems to be well-supported by many feed readers.

The way that we can extend the specifications is by adding some metadata for a “namespace”. This allows a RSS feed to give feed readers extra instructions on how to parse a non-standard tag. In the code snippets that I will show below, you’ll see the namespace that is attached to my RSS feed in order to support the [new `<content>` tag](https://www.rssboard.org/rss-profile#namespace-elements-content). The `<content:encoded>` tag expects to contain valid HTML, not just raw text. This gives feed readers and the end-user/consumer of the content more flexibility to style their content how they see fit.

### Show me the code!

Let's take a look at the full example that I wrote for my own personal site before I
upstreamed the implementation to the `@astrojs/rss` package. First, we need to get a list of our posts and their content using the [`import.meta.glob`](https://docs.astro.build/en/reference/api-reference/#importmeta) function. 

```typescript
import rss from '@astrojs/rss';
import sanitizeHtml from 'sanitize-html';

const posts = Object.values(
  import.meta.glob('./posts/**/*.md', {
    eager: true, 
  }),
);
```

Next, we'll set up a couple of custom config options to pass to the config object that Astro's `rss` function accepts. First, we'll set the `site` variable to be equal to our base URL. Then the `customData` variable will be a string which we pass some extra XML tags that are necessary to enable this for Atom feed readers and set the language to `en-us` (change this to whatever language you are specifying yourself). Lastly, we'll set up the `xmlns` variable which will be an object that allows us to inject `<xmlns:content>` tags into the feed so that we can tell feed readers to use special tags that are not necessarily a part of the RSS 2.0 spec. We will set up our feed to enable Atom and the `<content:encoding>` tags by default.

```typescript
const site = import.meta.env.SITE;
const customData = [
  `<atom:link href="${site}feed.xml" rel="self" type="application/rss+xml" />`,
  `<language>en-us</language>`,
].join('');
const xmlns = {
  atom: 'http://www.w3.org/2005/Atom',
  content: 'http://purl.org/rss/1.0/modules/content/',
};
```

Once we've taken care of all that setup, we can export a default function called `get` which Astro will treat as a `GET` [API endpoint](https://docs.astro.build/en/guides/endpoints/) from which we can return the call to the `rss` function provided by Astro. We can pass along our `site`, `customData`, and `xmlns` to the object accepted by the `rss` function as well as our `title` and `description` for the site. Note: the `\u2019` that you see below is the unicode code for an apostrophe so that it's properly escaped in our feed.

The last piece of the puzzle that we need to pass to the `rss` function is our `posts`, where I've used the `sanitize-html` package to sanitize our post's content before injecting it into our feed. Once we've sanitized the content, we can pass it to the `customData` property where we can tell Astro to inject a tag using the `<content:encoded>` tag that contains our post's sanitized HTML.

```typescript
export const get = () =>
  rss({
    title: 'Ben Smith\u2019s Blog',
    description: 'Ben\u2019s writings and thoughts about tech',
    site,
    items: posts.map((post) => {
      const postHTML = sanitizeHtml(post.compiledContent());
      return {
        title: post.frontmatter.title,
        link: post.url,
        pubDate: new Date(post.frontmatter.pubDate),
        description: descriptionOrExcerpt,
        customData: `<content:encoded><![CDATA[${postHTML}]]></content:encoded>`,
      };
    }),
    customData,
    xmlns,
  });
```

### Final code

Now let's drastically simplify that with the new updates to the `@astrojs/rss` post below with the [contribution that I made to the Astro project](https://github.com/withastro/astro/pull/5366) which removes a lot of the tedious configuration that I covered above. Now when you pass the list of our `posts` to the `items` property, you don't have to add those extra tags yourself as you can pass your sanitized post content directly to the new `content` property! The releveant meta tags, namespace definitions, Atom feed-specific tags, and other convoluted configuration is in Astro's `rss` package itself now.

```typescript
import rss from '@astrojs/rss';
import sanitizeHtml from 'sanitize-html';

const posts = Object.values(
  import.meta.glob('./posts/**/*.md', {
    eager: true, 
  }),
);

export const get = () => rss({
  title: 'Ben Smith\u2019s Blog',
  description: 'Ben\u2019s writings and thoughts about tech',
  site: import.meta.env.SITE,
  items: posts.map((post) => ({
    link: post.url,
    title: post.frontmatter.title,
    pubDate: post.frontmatter.pubDate,
    content: sanitizeHtml(post.compiledContent()),
  }))
});
```

Thanks for reading! If this post was helpful, please let me know, I'd love to hear any feedback you might have.
