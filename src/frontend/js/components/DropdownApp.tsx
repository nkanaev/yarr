import { useI18n } from "../i18n";
import type { Lang } from "../i18n";
import Dropdown from "./Dropdown";
import Icon from "./Icon";
import Tooltip from "./Tooltip";
import type { Theme } from "./DropdownSettings";

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

interface DropdownAppProps {
  theme: Theme;
  onThemeChange: (t: Theme) => void;
  onShowShortcuts: () => void;
  language: Lang;
  onLanguageChange: (lang: Lang) => void;
  requiresAuth: boolean;
  onLogout: () => void;
}

export default function DropdownApp({
  theme,
  onThemeChange,
  onShowShortcuts,
  language,
  onLanguageChange,
  requiresAuth,
  onLogout,
}: DropdownAppProps) {
  const { t } = useI18n();

  return (
    <Dropdown toggleClass="" drop="start" title={t("settings_main")} button={<Icon name="anchor" />}>
      <div className="c-dropdown-header" role="heading" aria-level={2}>
        {t("theme")}
      </div>
      <div className="row text-center m-0">
        {THEMES.map(th => (
          <Tooltip key={th} label={th} className="col-3 px-0">
            <button
              type="button"
              className={`c-button-link theme-swatch w-100 px-0 rounded-0 theme-${th}`}
              aria-label={th}
              aria-pressed={theme === th}
              onClick={e => {
                e.stopPropagation();
                onThemeChange(th);
              }}
            />
          </Tooltip>
        ))}
      </div>

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
            <Tooltip key={lang.code} label={lang.name} className="col-3 px-0">
              <button
                type="button"
                className="c-dropdown-item text-center w-100"
                aria-label={lang.name}
                aria-pressed={language === lang.code}
                onClick={e => {
                  e.stopPropagation();
                  onLanguageChange(lang.code);
                }}>
                {lang.code}
              </button>
            </Tooltip>
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
