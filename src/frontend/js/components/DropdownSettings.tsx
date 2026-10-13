import { useRef } from "preact/hooks";
import type { Ref } from "preact";
import { useI18n } from "../i18n";
import type { DropdownHandle } from "./Dropdown";
import Dropdown from "./Dropdown";
import Icon from "./Icon";

export type Theme = "system" | "light" | "sepia" | "night";

interface DropdownSettingsProps {
  refreshAvailable: boolean;
  refreshRate: number;
  refreshRateTitle: string;
  maxRefreshRate: boolean;
  onRefreshRateChange: (delta: number) => void;
  onFetchAllFeeds: () => void;
  onImportOPML: (e: Event) => void;
  onShowNewFeed: () => void;
  dropdownRef?: Ref<DropdownHandle>;
}

export default function DropdownSettings({
  refreshAvailable,
  refreshRate,
  refreshRateTitle,
  maxRefreshRate,
  onRefreshRateChange,
  onFetchAllFeeds,
  onImportOPML,
  onShowNewFeed,
  dropdownRef,
}: DropdownSettingsProps) {
  const { t } = useI18n();
  const fileInputRef = useRef<HTMLInputElement>(null);

  return (
    <Dropdown
      ref={dropdownRef}
      toggleClass="c-button-link px-2"
      drop="end"
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
    </Dropdown>
  );
}
