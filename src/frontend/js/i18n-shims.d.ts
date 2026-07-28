import "vue";
import type { Lang, TranslationKey } from "./i18n";
import { FluentVariable } from "@fluent/bundle";

declare module "vue" {
  interface ComponentCustomProperties {
    $t: {
      (code: TranslationKey, args?: Record<string, FluentVariable>): string;
      set: (lang: Lang) => void;
    };
  }
}

export {};
