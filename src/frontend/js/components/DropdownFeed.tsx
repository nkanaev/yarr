import type { Feed, Folder } from "../api-types";
import { useI18n } from "../i18n";
import Dropdown from "./Dropdown";
import Icon from "./Icon";

interface DropdownFeedProps {
  feed: Feed;
  folders: Folder[];
  onRename: (feed: Feed) => void;
  onUpdateLink: (feed: Feed) => void;
  onMove: (feed: Feed, folderId: number | null) => void;
  onMoveToNewFolder: (feed: Feed) => void;
  onDelete: (feed: Feed) => void;
}

export default function DropdownFeed({
  feed,
  folders,
  onRename,
  onUpdateLink,
  onMove,
  onMoveToNewFolder,
  onDelete,
}: DropdownFeedProps) {
  const { t } = useI18n();
  const selectedFolder = folders.find(f => f.id === feed.folder_id);

  return (
    <Dropdown
      toggleClass="c-button-link px-2"
      drop="end"
      title={t("settings_selected")}
      button={<Icon name="more-horizontal" />}>
      <div className="c-dropdown-header text-break line-clamp-1" role="heading" aria-level={2}>
        {feed.title}
      </div>
      {feed.link && (
        <a
          className="c-dropdown-item text-start text-decoration-none d-flex gap-1"
          href={feed.link}
          rel="noopener noreferrer"
          target="_blank"
          referrerPolicy="no-referrer">
          <Icon className="me-1" name="globe" />
          {t("website")}
        </a>
      )}
      {feed.feed_link && (
        <a
          className="c-dropdown-item text-start text-decoration-none d-flex gap-1"
          href={feed.feed_link}
          rel="noopener noreferrer"
          target="_blank"
          referrerPolicy="no-referrer">
          <Icon className="me-1" name="rss" />
          {t("feed_link")}
        </a>
      )}
      {(feed.link || feed.feed_link) && <div className="c-dropdown-divider" />}
      <button
        type="button"
        className="c-dropdown-item w-100 text-start d-flex gap-1"
        onClick={() => onRename(feed)}>
        <Icon className="me-1" name="edit" />
        {t("rename")}
      </button>
      {feed.feed_link && (
        <button
          type="button"
          className="c-dropdown-item w-100 text-start d-flex gap-1"
          onClick={() => onUpdateLink(feed)}>
          <Icon className="me-1" name="edit" />
          {t("change_link")}
        </button>
      )}
      <div className="c-dropdown-divider" />
      <div className="c-dropdown-header" role="heading" aria-level={2}>
        {t("move_to")}
      </div>
      <div className="c-dropdown-item c-dropdown-select d-flex align-items-center gap-1">
        <Icon className="me-1 flex-shrink-0" name="folder" />
        <span className="flex-grow-1 text-break line-clamp-1">
          {selectedFolder?.title ?? "──"}
        </span>
        <Icon className="flex-shrink-0" name="chevron-down" />
        <select
          value={feed.folder_id === null ? "" : feed.folder_id}
          onChange={e => {
            const v = (e.target as HTMLSelectElement).value;
            onMove(feed, v === "" ? null : Number(v));
          }}>
          <option value="">──</option>
          {folders.map(f => (
            <option key={f.id} value={f.id}>
              {f.title}
            </option>
          ))}
        </select>
      </div>
      <button
        type="button"
        className="c-dropdown-item w-100 text-start opacity-75 d-flex gap-1"
        onClick={() => onMoveToNewFolder(feed)}>
        <Icon className="me-1" name="folder-plus" />
        {t("new_folder")}
      </button>
      <div className="c-dropdown-divider" />
      <button
        type="button"
        className="c-dropdown-item w-100 text-start text-danger d-flex gap-1"
        onClick={() => onDelete(feed)}>
        <Icon className="me-1" name="trash" />
        {t("delete")}
      </button>
    </Dropdown>
  );
}
