import { z } from 'zod';

export const PostSchema = z.object({
  url: z.string(),
  frontmatter: z.object({
    title: z.string(),
    layout: z.string(),
    pubDate: z.string(),
    tags: z.string().optional(),
    description: z.string().optional(),
    excerpt: z.string().optional(),
    draft: z.boolean().optional(),
  }),
  compiledContent: z.function().returns(z.string()),
  rawContent: z.function().returns(z.string()),
});
export type Post = z.infer<typeof PostSchema>;

export const PostsSchema = z.array(PostSchema);
export type Posts = z.infer<typeof PostsSchema>;
