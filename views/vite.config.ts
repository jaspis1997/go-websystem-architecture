import { defineConfig } from "vite";
import solid from "vite-plugin-solid";

import fs from "node:fs";

import inputPlugin from "@macropygia/vite-plugin-glob-input";
import suidPlugin from "@suid/vite-plugin";
import path, { resolve } from "node:path";
// @ts-expect-error
import handlebars from "vite-plugin-handlebars";

function pluginServeAdditionalStaticFiles(){
  return {
    name: "serve-additional-static-files",
    // biome-ignore lint/suspicious/noExplicitAny: <explanation>
    configureServer(server: any) {
      // biome-ignore lint/suspicious/noExplicitAny: <explanation>
      server.middlewares.use((req:any, res:any, next:any) => {
        if (req.url.startsWith("/static", req.url)) {
          // const filePath = path.join("./public", req.url.replace(routePrefix, ''));
          const filePath = path.join("./public", req.url.replace("/static", ''));
          console.debug(req.url,filePath);
          // ファイルが存在する場合はレスポンスを返す
          if (fs.existsSync(filePath) && fs.statSync(filePath).isFile()) {
            res.setHeader('Content-Type', getMimeType(filePath));
            res.statusCode = 200;
            fs.createReadStream(filePath).pipe(res);
            return;
          }
        }
        next();
      });
    },
  }
}

function getMimeType(filePath: string) {
  const ext = path.extname(filePath).toLowerCase();
  const mimeTypes: Record<string, string> = {
    '.html': 'text/html',
    '.css': 'text/css',
    '.js': 'application/javascript',
    '.json': 'application/json',
    '.png': 'image/png',
    '.jpg': 'image/jpeg',
    '.gif': 'image/gif',
    '.svg': 'image/svg+xml',
    '.txt': 'text/plain',
  };
  return mimeTypes[ext] || 'application/octet-stream';
}

export default defineConfig({
  root: "./",
  publicDir: "",
  assetsInclude: ["public/**/*"],
  build: {
    outDir: "dist",
    // sourcemap: true,
  },
  plugins: [
    pluginServeAdditionalStaticFiles(),
    solid(),
    suidPlugin(),
    inputPlugin({
      patterns: ["./templates/**/*.html","./partials/**/*.html"],
    }),
    handlebars({
      partialDirectory: resolve(__dirname, "./src/partials"),
    }),
  ],
  css: {
    modules: {
      localsConvention: "dashes",
    },
    preprocessorOptions: {
      scss: {
        api: "modern-compiler",
      },
    },
  },
  resolve: {
    alias: {
      "@src": resolve(__dirname, "./src"),
      "@partial": resolve(__dirname, "./src/partials"),
    },
  },
});
