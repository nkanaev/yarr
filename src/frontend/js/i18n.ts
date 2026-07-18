import type { Plugin } from "vue";
import _translations from "./i18n-translations.json" with {type: "json"}
import { FluentResource, FluentBundle, FluentVariable } from "@fluent/bundle";

export type Lang = "en" | "de" | "fr" | "es" | "ja" | "pt" | "zh" | "ru";

const translations = _translations satisfies Record<string, Record<Lang, string>>;

export type TranslationKey = keyof typeof translations;

function ftlFrom(lang: Lang) {
  return Object.entries(translations)
    .map(([key, langs]) => `${key} = ${langs[lang]}`)
    .join("\n");
}
export default {
  install(app: any) {
    let bundle = undefined as FluentBundle | undefined;
    app.config.globalProperties.$setLang = function (lang: Lang) {
      const ftl = ftlFrom(lang);
      const resource = new FluentResource(ftl);
      bundle = new FluentBundle(lang);
      bundle.addResource(resource);
    };
    app.config.globalProperties.$t = function (
      code: TranslationKey,
      args?: Record<string, FluentVariable>,
    ): string | undefined {
      if (!bundle) return;
      const msg = bundle.getMessage(code);
      if (!msg || !msg.value) return;
      return bundle.formatPattern(msg.value, args);
    };
  },
};
