import "vue";
import type { Lang, TranslationKey } from "./i18n";
import { FluentVariable } from "@fluent/bundle";

declare module "vue" {
  interface ComponentCustomProperties {
    $setLang: (lang: Lang) => void;
    $t: (
      code: TranslationKey,
      args?: Record<string, FluentVariable>,
    ) => string | undefined;
  }
}

export {};
