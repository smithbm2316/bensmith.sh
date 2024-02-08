export default function ({ eleventy }) {
  return {
    title: 'Ben Smith - Full Stack Web Developer',
    url: eleventy.env.runMode === 'build' ? 'https://bensmith.sh' : '',
    author: {
      name: 'Ben Smith',
      email: 'bsmithdev@mailbox.org',
    },
    description:
      "I'm a web developer with a passion for the command line, coffee, and crafting delightful user experiences for the web!",
    // favicon: '/assets/favicon.svg',
  };
}
