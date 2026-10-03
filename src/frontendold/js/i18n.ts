import { reactive, type Plugin } from "vue";
import _translations from "./i18n-translations.json" with { type: "json" };
import { FluentResource, FluentBundle, FluentVariable, Message } from "@fluent/bundle";
import { Pattern } from "@fluent/bundle/esm/ast";

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
    const i18nstate = reactive({
      bundle: new FluentBundle("en"),
    });

    const t = (code: TranslationKey, args?: Record<string, FluentVariable>): string => {
      const msg = i18nstate.bundle.getMessage(code);
      if (msg?.value) {
        return i18nstate.bundle.formatPattern(msg.value as Pattern, args);
      }
      return "";
    };
    t.set = (lang: Lang) => {
      const ftl = ftlFrom(lang);
      const resource = new FluentResource(ftl);
      let bundle = new FluentBundle(lang);
      bundle.addResource(resource);
      i18nstate.bundle = bundle;
    };
    app.config.globalProperties.$t = t;
  },
};
