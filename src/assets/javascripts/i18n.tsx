import { createContext } from "preact";
import { useContext, useState, useMemo, useCallback } from "preact/hooks";
import type { ComponentChildren } from "preact";
import _translations from "./i18n-translations.json" with { type: "json" };
import { FluentResource, FluentBundle, FluentVariable } from "@fluent/bundle";
import { Pattern } from "@fluent/bundle/esm/ast";

export type Lang = "en" | "de" | "fr" | "es" | "ja" | "pt" | "zh" | "ru";

const translations = _translations satisfies Record<string, Record<Lang, string>>;

export type TranslationKey = keyof typeof translations;

function createBundle(lang: Lang): FluentBundle {
  const ftl = Object.entries(translations)
    .map(([key, langs]) => `${key} = ${langs[lang]}`)
    .join("\n");
  const resource = new FluentResource(ftl);
  const bundle = new FluentBundle(lang);
  bundle.addResource(resource);
  return bundle;
}

export type TranslateFn = (code: TranslationKey, args?: Record<string, FluentVariable>) => string;

interface I18nContextType {
  t: TranslateFn;
  lang: Lang;
  setLang: (lang: Lang) => void;
}

const defaultBundle = createBundle("en");

const defaultT: TranslateFn = (code, args) => {
  const msg = defaultBundle.getMessage(code);
  if (msg?.value) {
    return defaultBundle.formatPattern(msg.value as Pattern, args);
  }
  return "";
};

const I18nContext = createContext<I18nContextType>({
  t: defaultT,
  lang: "en",
  setLang: () => {},
});

export function I18nProvider({
  initialLang = "en",
  children,
}: {
  initialLang?: Lang;
  children: ComponentChildren;
}) {
  const [lang, setLangState] = useState<Lang>(initialLang);
  const bundle = useMemo(() => createBundle(lang), [lang]);

  const t = useCallback<TranslateFn>(
    (code, args) => {
      const msg = bundle.getMessage(code);
      if (msg?.value) {
        return bundle.formatPattern(msg.value as Pattern, args);
      }
      return "";
    },
    [bundle],
  );

  const setLang = useCallback((l: Lang) => {
    document.documentElement.lang = l;
    setLangState(l);
  }, []);

  return <I18nContext.Provider value={{ t, lang, setLang }}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  return useContext(I18nContext);
}
