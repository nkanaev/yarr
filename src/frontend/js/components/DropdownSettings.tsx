import { useRef } from "preact/hooks";
import type { Ref } from "preact";
import { useI18n } from "../i18n";
import type { AppSettings } from "../models/settings";
import type { DropdownHandle } from "./Dropdown";
import Dropdown from "./Dropdown";
import Icon from "./Icon";

const REFRESH_RATE_OPTIONS = [
  { title: "0", value: 0 },
  { title: "10m", value: 10 },
  { title: "30m", value: 30 },
  { title: "1h", value: 60 },
  { title: "2h", value: 120 },
  { title: "4h", value: 240 },
  { title: "12h", value: 720 },
  { title: "24h", value: 1440 },
];

interface DropdownSettingsProps {
  settings: AppSettings;
  refreshAvailable: boolean;
  onFetchAllFeeds: () => void;
  onImportOPML: (e: Event) => void;
  onShowNewFeed: () => void;
  dropdownRef?: Ref<DropdownHandle>;
}

export default function DropdownSettings({
  settings,
  refreshAvailable,
  onFetchAllFeeds,
  onImportOPML,
  onShowNewFeed,
  dropdownRef,
}: DropdownSettingsProps) {
  const { t } = useI18n();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const refreshRate = settings.values.refresh_rate.value;
  const refreshRateTitle =
    REFRESH_RATE_OPTIONS.find(o => o.value === refreshRate)?.title || "0";
  const maxRefreshRate =
    refreshRate === REFRESH_RATE_OPTIONS[REFRESH_RATE_OPTIONS.length - 1].value;

  const changeRefreshRate = (offset: number) => {
    const curIdx = REFRESH_RATE_OPTIONS.findIndex(o => o.value === refreshRate);
    if ((curIdx <= 0 && offset < 0) || (curIdx >= REFRESH_RATE_OPTIONS.length - 1 && offset > 0))
      return;
    settings.update({ refresh_rate: REFRESH_RATE_OPTIONS[curIdx + offset].value });
  };

  return (
    <Dropdown
      ref={dropdownRef}
      toggleClass="c-button-link px-2"
      drop="end"
      title={t("settings_feeds")}
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
                changeRefreshRate(-1);
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
                changeRefreshRate(1);
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
        <label
          className="c-dropdown-item mb-0 cursor-pointer w-100 d-flex gap-1"
          htmlFor="opml-import"
          onClick={e => e.stopPropagation()}>
          <Icon className="me-1" name="download" />
          {t("import")}
          <input
            ref={fileInputRef}
            type="file"
            id="opml-import"
            onChange={onImportOPML}
            name="opml"
            style={{ opacity: 0, width: "1px", height: 0, position: "absolute", zIndex: -1 }}
          />
        </label>
      </form>
      <a
        className="c-dropdown-item d-block text-start text-decoration-none d-flex gap-1"
        href="./opml/export">
        <Icon className="me-1" name="upload" />
        {t("export")}
      </a>
    </Dropdown>
  );
}
