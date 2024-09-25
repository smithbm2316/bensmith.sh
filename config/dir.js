/**
 * the `dir` key in the object returned from the Eleventy config function
 *
 * @link {https://www.11ty.dev/docs/config/#configuration-options}
 * @link {https://stackoverflow.com/a/64687300/15089697}
 */
const dir = /** @type {const} */ ({
  input: 'src',
  output: '_site',
  includes: '_includes',
  data: '_data',
});

export default dir;
