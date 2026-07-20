const esbuild = require("esbuild");
const vuePlugin = require("esbuild-plugin-vue3");

esbuild
  .build({
    entryPoints: [
      "src/frontend/js/main.ts",
      "src/frontend/css/app.css",
    ],
    bundle: true,
    outdir: "src/assets/static",
    entryNames: "bundle",
    alias: { vue: "vue/dist/vue.esm-bundler.js" },
    plugins: [vuePlugin()],
  })
  .catch(() => process.exit(1));
