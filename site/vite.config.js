import { globSync } from "glob";
import { browserslistToTargets } from "lightningcss";
import path from "node:path";
import { defineConfig } from "vite";

/** @type {Record<string, string>} */
let staticHTMLFiles = {};
for (let file of globSync("**/*.html")) {
  let pathnameRelativeToRoot = path.relative(
    ".",
    file.slice(0, file.length - path.extname(file).length),
  );
  staticHTMLFiles[pathnameRelativeToRoot] = file;
}

export default defineConfig({
  appType: "mpa",
  clearScreen: false,
  server: {
    port: 2323,
  },
  build: {
    cssMinify: "lightningcss",
    outDir: "../.site",
    rollupOptions: {
      input: staticHTMLFiles,
    },
  },
  preview: {
    port: 2323,
  },
  css: {
    transformer: "lightningcss",
    lightningcss: {
      targets: browserslistToTargets([
        "> 0.5%, last 2 versions, Firefox ESR, not dead",
      ]),
      drafts: {
        customMedia: true,
      },
    },
  },
});
