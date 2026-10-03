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

  return (
    <Dropdown
      toggleClass="c-button-link px-2"
      drop="right"
      title={t("feed_settings")}
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
      {folders.map(
        f =>
          f.id !== feed.folder_id && (
            <button
              key={f.id}
              type="button"
              className="c-dropdown-item w-100 text-start d-flex gap-1"
              onClick={() => onMove(feed, f.id)}>
              <Icon className="me-1" name="folder" />
              <span className="text-break line-clamp-1">{f.title}</span>
            </button>
          ),
      )}
      {feed.folder_id && (
        <button
          type="button"
          className="c-dropdown-item w-100 text-start opacity-75 d-flex gap-1"
          onClick={() => onMove(feed, null)}>
          <Icon className="me-1" name="folder-minus" />
          ──
        </button>
      )}
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
