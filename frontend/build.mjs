import * as esbuild from "esbuild";
import { copyFileSync, mkdirSync } from "fs";

// Create dist directory
mkdirSync("./dist", { recursive: true });

// Copy index.html to dist
copyFileSync("./index.html", "./dist/index.html");

// Build bundle.js
await esbuild.build({
  entryPoints: ["./src/main.ts"],
  bundle: true,
  outfile: "./dist/bundle.js",
  format: "esm",
  platform: "browser",
  target: "es2020",
});
