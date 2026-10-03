import { useRef } from "preact/hooks";
import type { Ref } from "preact";
import type { Lang } from "../i18n";
import { useI18n } from "../i18n";
import type { DropdownHandle } from "./Dropdown";
import Dropdown from "./Dropdown";
import Icon from "./Icon";

export type Theme = "system" | "light" | "sepia" | "night";

const THEMES: Theme[] = ["light", "sepia", "night", "system"];

const LANGUAGES: { code: Lang; name: string }[] = [
  { code: "en", name: "English" },
  { code: "de", name: "Deutsch" },
  { code: "es", name: "Español" },
  { code: "fr", name: "Français" },
  { code: "ja", name: "日本語" },
  { code: "pt", name: "Português" },
  { code: "ru", name: "Русский" },
  { code: "zh", name: "简体中文" },
];

interface DropdownSettingsProps {
  theme: Theme;
  onThemeChange: (t: Theme) => void;
  refreshAvailable: boolean;
  refreshRate: number;
  refreshRateTitle: string;
  maxRefreshRate: boolean;
  onRefreshRateChange: (delta: number) => void;
  onFetchAllFeeds: () => void;
  onImportOPML: (e: Event) => void;
  onShowShortcuts: () => void;
  onShowNewFeed: () => void;
  language: Lang;
  onLanguageChange: (lang: Lang) => void;
  requiresAuth: boolean;
  onLogout: () => void;
  dropdownRef?: Ref<DropdownHandle>;
}

export default function DropdownSettings({
  theme,
  onThemeChange,
  refreshAvailable,
  refreshRate,
  refreshRateTitle,
  maxRefreshRate,
  onRefreshRateChange,
  onFetchAllFeeds,
  onImportOPML,
  onShowShortcuts,
  onShowNewFeed,
  language,
  onLanguageChange,
  requiresAuth,
  onLogout,
  dropdownRef,
}: DropdownSettingsProps) {
  const { t } = useI18n();
  const fileInputRef = useRef<HTMLInputElement>(null);

  return (
    <Dropdown
      ref={dropdownRef}
      toggleClass="c-button-link px-2"
      drop="right"
      title={t("settings")}
      button={<Icon name="more-horizontal" />}>
      <button
        type="button"
        className="c-dropdown-item w-100 text-start d-flex gap-1"
        onClick={onShowNewFeed}>
        <Icon className="me-1" name="plus" />
        {t("new_feed")}
      </button>

      {refreshAvailable && (
        <>
          <div className="c-dropdown-divider" />
          <button
            type="button"
            className="c-dropdown-item w-100 text-start d-flex gap-1"
            onClick={onFetchAllFeeds}>
            <Icon className="me-1" name="rotate-cw" />
            {t("refresh_feeds")}
          </button>
        </>
      )}

      <div className="c-dropdown-divider" />
      <div className="c-dropdown-header" role="heading" aria-level={2}>
        {t("theme")}
      </div>
      <div className="row text-center m-0">
        {THEMES.map(th => (
          <button
            key={th}
            type="button"
            className={`c-button-link theme-swatch col-3 px-0 rounded-0 theme-${th}`}
            title={th}
            aria-label={th}
            aria-pressed={theme === th}
            onClick={e => {
              e.stopPropagation();
              onThemeChange(th);
            }}
          />
        ))}
      </div>

      {refreshAvailable && (
        <>
          <div className="c-dropdown-divider" />
          <div className="c-dropdown-header" role="heading" aria-level={2}>
            {t("auto_refresh")}
          </div>
          <div className="row text-center m-0">
            <button
              type="button"
              className="c-dropdown-item col-4 px-0 d-flex gap-1 justify-content-center"
              onClick={e => {
                e.stopPropagation();
                onRefreshRateChange(-1);
              }}
              disabled={!refreshRate}>
              <Icon name="chevron-down" />
            </button>
            <div className="col-4 d-flex align-items-center justify-content-center user-select-none">
              {refreshRateTitle}
            </div>
            <button
              type="button"
              className="c-dropdown-item col-4 px-0 d-flex gap-1 justify-content-center"
              onClick={e => {
                e.stopPropagation();
                onRefreshRateChange(1);
              }}
              disabled={maxRefreshRate}>
              <Icon name="chevron-up" />
            </button>
          </div>
        </>
      )}

      <div className="c-dropdown-divider" />
      <div className="c-dropdown-header" role="heading" aria-level={2}>
        {t("subscriptions")}
      </div>
      <form onSubmit={e => e.preventDefault()} tabIndex={-1}>
        <input
          ref={fileInputRef}
          type="file"
          id="opml-import"
          onChange={onImportOPML}
          name="opml"
          style={{ opacity: 0, width: "1px", height: 0, position: "absolute", zIndex: -1 }}
        />
        <label
          className="c-dropdown-item mb-0 cursor-pointer w-100 d-flex gap-1"
          htmlFor="opml-import"
          onClick={e => e.stopPropagation()}>
          <Icon className="me-1" name="download" />
          {t("import")}
        </label>
      </form>
      <a
        className="c-dropdown-item d-block text-start text-decoration-none d-flex gap-1"
        href="./opml/export">
        <Icon className="me-1" name="upload" />
        {t("export")}
      </a>

      <div className="c-dropdown-divider" />
      <button
        type="button"
        className="c-dropdown-item w-100 text-start d-flex gap-1"
        onClick={onShowShortcuts}>
        <Icon className="me-1" name="help-circle" />
        {t("shortcuts")}
      </button>

      <div className="c-dropdown-divider" />
      <div className="c-dropdown-header" role="heading" aria-level={2}>
        A / あ / 文
      </div>
      <div className="container">
        <div className="row">
          {LANGUAGES.map(lang => (
            <button
              key={lang.code}
              type="button"
              className="c-dropdown-item text-center col-3 px-0"
              aria-label={lang.name}
              aria-pressed={language === lang.code}
              title={lang.name}
              onClick={e => {
                e.stopPropagation();
                onLanguageChange(lang.code);
              }}>
              {lang.code}
            </button>
          ))}
        </div>
      </div>

      {requiresAuth && (
        <>
          <div className="c-dropdown-divider" />
          <button type="button" className="c-dropdown-item w-100 text-start" onClick={onLogout}>
            <Icon className="me-1" name="log-out" />
            {t("log_out")}
          </button>
        </>
      )}
    </Dropdown>
  );
}
