export function getPostExcerpt(postContent: string) {
  let excerpt = postContent.substring(0, 140);

  // if we have an empty string just return
  if (excerpt.length === 0) {
    return excerpt;
  }
  // strip out any partial words by finding the last space in the string and appending a
  // '...' to it
  const lastSpaceIndex = excerpt.lastIndexOf(' ');
  excerpt = excerpt.substring(0, lastSpaceIndex);
  excerpt += '...';
  return excerpt;
}
