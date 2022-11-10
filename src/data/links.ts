export interface LinkWithIcon {
  title: string;
  url: string;
  pack?: string;
  icon: string;
}

export const internalSiteLinks: LinkWithIcon[] = [
  {
    title: 'Home',
    url: '/',
    icon: 'home',
  },
  {
    title: 'About',
    url: '/about',
    icon: 'user-circle',
  },
  {
    title: 'Projects',
    url: '/projects',
    icon: 'briefcase',
  },
  {
    title: 'Blog',
    url: '/blog',
    icon: 'newspaper',
  },
  {
    title: 'Uses',
    url: '/uses',
    icon: 'chip',
  },
  {
    title: 'Linktree',
    url: '/me',
    icon: 'link',
  },
];

export const externalLinks: LinkWithIcon[] = [
  {
    title: 'GitHub',
    url: 'https://github.com/smithbm2316',
    pack: 'mdi',
    icon: 'github',
  },
  {
    title: 'YouTube',
    url: 'https://youtube.com/@smithbm2316',
    pack: 'mdi',
    icon: 'youtube',
  },
  {
    title: 'Mastodon',
    url: 'https://fosstodon.org/@smithbm2316',
    pack: 'mdi',
    icon: 'mastodon',
  },
  {
    title: 'Twitter',
    url: 'https://twitter.com/smithbm2316',
    pack: 'mdi',
    icon: 'twitter',
  },
  {
    title: 'Email',
    url: 'mailto:bsmithdev@mailbox.org',
    icon: 'mail-open',
  },
  {
    title: 'LinkedIn',
    url: 'https://linkedin.com/in/smithbm2316',
    pack: 'mdi',
    icon: 'linkedin',
  },
];
