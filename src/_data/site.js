/**
 * Global site configuration data
 */
const site = /** @type {const} */ ({
  title: 'Ben Smith - Full Stack Web Developer',
  url: 'https://www.bensmith.sh',
  author: {
    name: 'Ben Smith',
    email: 'ben@bensmith.sh',
  },
  description:
    "I'm a web developer at Thuma with a passion for the command line, coffee, and crafting delightful user experiences for the web!",
  favicon: '/assets/favicon.ico',
});
site.descriptionNodes = site.description.split('Thuma');

export default site;
