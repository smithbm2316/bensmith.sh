/**
 * The object returned from the Eleventy config function
 *
 * @link {https://www.11ty.dev/docs/config/#configuration-options}
 * @link {https://stackoverflow.com/a/64687300/15089697}
 */
const config = /** @type {const} */ ({
  dir: {
    input: 'src',
    output: '_site',
    includes: '_includes',
    data: '_data',
  },
  markdownTemplateEngine: 'njk',
  htmlTemplateEngine: 'njk',
  templateFormats: [
    'html',
    'md',
    'njk',
    '11ty.js',
    'webc',
    'vue',
    // copy over these files as plain text
    'css',
    'txt',
    'webmanifest',
  ],
});

export default config;
