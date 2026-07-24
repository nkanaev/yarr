import esbuild from "esbuild";
import vuePlugin from "esbuild-plugin-vue3";

const args = process.argv.slice(2);
const isWatch = args.includes("--watch");
const isMinify = args.includes("--minify") || process.env.NODE_ENV === "production";

const options = {
  entryPoints: [
    "src/frontend/js/main.ts",
    "src/frontend/css/app.css",
  ],
  bundle: true,
  outdir: "src/assets/static",
  entryNames: "bundle",
  alias: { vue: "vue/dist/vue.esm-bundler.js" },
  plugins: [vuePlugin()],
  supported: { nesting: false },
  target: "es2020",
  sourcemap: !isMinify,
  minify: isMinify,
};

if (isWatch) {
  const ctx = await esbuild.context(options);
  await ctx.watch();
  console.log("Watching for changes...");
} else {
  await esbuild.build(options);
}
