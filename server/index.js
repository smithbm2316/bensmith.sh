import Eleventy from '@11ty/eleventy';
import polka from 'polka';
import sirv from 'sirv';
/**
 * The config object/instance for Eleventy
 * @typedef {import('@11ty/eleventy').UserConfig} EleventyConfig
 */

polka()
  .use(sirv('_site'))
  .get('/ping', async (req, res) => {
    let elev = new Eleventy('src/ping.webc', undefined, {
      // make sure to let eleventy know this is being run from a script with the
      // programmatic api, not from our normal dev/prod cli builds
      source: 'script',
      configPath: 'config/base.js',
      config(/** @type {EleventyConfig} */ eleventyConfig) {
        eleventyConfig.addGlobalData('props', {
          title: 'Pong',
          message: req.query?.message,
        });
      },
    });

    let json = await elev.toJSON();
    res.setHeader('content-type', 'text/html');
    res.statusCode = 200;
    res.end(json[0].content);
  })
  .listen(2323, (err) => {
    if (err) throw err;
    console.log('> Ready on localhost:2323~!');
  });
