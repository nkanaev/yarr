import type { Signal } from "@preact/signals";
import { useComputed } from "@preact/signals";
import { useEffect, useRef } from "preact/hooks";
import api from "../api";
import { scrollto, to } from "../utils";
import { useI18n } from "../i18n";
import type { AppSettings } from "../models/settings";
import type { AppRegistry } from "../models/registry";
import type { AppItemList } from "../models/itemlist";
import type { Feed, Folder } from "../api-types";
import Drag from "./Drag";
import Icon from "./Icon";
import Tooltip from "./Tooltip";
import RelativeTime from "./RelativeTime";
import DropdownFeed from "./DropdownFeed";
import DropdownFolder from "./DropdownFolder";
import { showToast } from "./Toast";

interface ColumnItemListProps {
  settings: AppSettings;
  registry: AppRegistry;
  itemList: AppItemList;
  itemSelected: Signal<number | null>;
  onMarkAllRead: () => void;
}

export default function ColumnItemList({
  settings,
  registry,
  itemList,
  itemSelected,
  onMarkAllRead,
}: ColumnItemListProps) {
  const { t } = useI18n();

  const currentFeed = useComputed(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    return type === "feed" ? registry.feedsById.value[Number(guid)] || null : null;
  });
  const currentFolder = useComputed(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    return type === "folder" ? registry.foldersById.value[Number(guid)] || null : null;
  });
  const searchScope = useComputed(() => {
    return (
      currentFeed.value?.title ||
      currentFolder.value?.title ||
      (settings.values.filter.value === "unread"
        ? t("all_unread")
        : settings.values.filter.value === "starred"
          ? t("all_starred")
          : t("all_feeds"))
    );
  });

  const resizeItemList = (w: number) => {
    settings.update({ item_list_width: Math.min(Math.max(200, w), 700) });
  };

  const itemListRef = useRef<HTMLDivElement>(null);
  const spinnerRef = useRef<HTMLDivElement>(null);

  // Load the next page when the loading indicator enters the viewport.
  // Re-runs after each page so short lists keep filling the column.
  useEffect(() => {
    const scroll = itemListRef.current;
    const target = spinnerRef.current;
    if (!scroll || !target) return;
    const io = new IntersectionObserver(
      entries => entries.some(e => e.isIntersecting) && itemList.loadMore(),
      { root: scroll, rootMargin: "0px 0px 100px 0px" },
    );
    io.observe(target);
    return () => io.disconnect();
  }, [itemList, itemList.items.value.length, itemList.loading.value]);

  // Reveal the selected item
  useEffect(() => {
    if (itemSelected.value === null) return;
    const timer = setTimeout(() => {
      const el = itemListRef.current;
      if (!el) return;
      const selected = el.querySelector('[aria-checked="true"]');
      if (selected) scrollto(selected, el);
    }, 0);
    return () => clearTimeout(timer);
  }, [itemSelected.value]);

  // Reset scroll when the item query changes
  useEffect(() => {
    if (itemListRef.current) itemListRef.current.scrollTop = 0;
  }, [
    settings.values.feed.value,
    settings.values.filter.value,
    settings.values.sort_newest_first.value,
  ]);

  // Feed & folder actions
  const renameFeed = async (feed: Feed) => {
    const title = prompt(t("prompt_new_title"), feed.title);
    if (!title) return;
    const [err] = await to(api.feeds.update(feed.id, { title }));
    if (err) showToast({ title: "fail_save_feed", description: err }, { level: "fail" });
    else
      registry.feeds.value = registry.feeds.value.map(f =>
        f.id === feed.id ? { ...f, title } : f,
      );
  };

  const updateFeedLink = async (feed: Feed) => {
    const feed_link = prompt(t("prompt_feed_link"), feed.feed_link);
    if (feed_link === null) return;
    const [err] = await to(api.feeds.update(feed.id, { feed_link }));
    if (err) showToast({ title: "fail_save_feed", description: err }, { level: "fail" });
    else
      registry.feeds.value = registry.feeds.value.map(f =>
        f.id === feed.id ? { ...f, feed_link } : f,
      );
  };

  const moveFeed = async (feed: Feed, folder_id: number | null) => {
    const [err] = await to(api.feeds.update(feed.id, { folder_id }));
    if (err) showToast({ title: "fail_save_feed", description: err }, { level: "fail" });
    else {
      registry.feeds.value = registry.feeds.value.map(f =>
        f.id === feed.id ? { ...f, folder_id } : f,
      );
      registry.refreshStats();
    }
  };

  const moveFeedToNewFolder = async (feed: Feed) => {
    const title = prompt(t("prompt_folder_name"));
    if (!title) return;
    const [folderErr, folder] = await to(api.folders.create({ title }));
    if (folderErr) {
      showToast({ title: "fail_save_folder", description: folderErr }, { level: "fail" });
      return;
    }
    const [updateErr] = await to(api.feeds.update(feed.id, { folder_id: folder.id }));
    if (updateErr) {
      showToast({ title: "fail_save_feed", description: updateErr }, { level: "fail" });
      return;
    }
    await registry.refresh();
    registry.refreshStats();
  };

  const deleteFeed = async (feed: Feed) => {
    if (!confirm(t("confirm_delete", { name: feed.title }))) return;
    const [err] = await to(api.feeds.delete(feed.id));
    if (err) showToast({ title: "fail_save_feed", description: err }, { level: "fail" });
    else {
      settings.values.feed.value = null;
      registry.refreshStats();
      registry.refresh();
    }
  };

  const renameFolder = async (folder: Folder) => {
    const title = prompt(t("prompt_new_title"), folder.title);
    if (!title) return;
    const [err] = await to(api.folders.update(folder.id, { title }));
    if (err) showToast({ title: "fail_save_folder", description: err }, { level: "fail" });
    else {
      registry.folders.value = registry.folders.value
        .map(f => (f.id === folder.id ? { ...f, title } : f))
        .sort((a, b) => a.title.localeCompare(b.title));
    }
  };

  const deleteFolder = async (folder: Folder) => {
    if (!confirm(t("confirm_delete", { name: folder.title }))) return;
    const [err] = await to(api.folders.delete(folder.id));
    if (err) showToast({ title: "fail_save_folder", description: err }, { level: "fail" });
    else {
      settings.values.feed.value = null;
      registry.refreshStats();
      registry.refresh();
    }
  };

  return (
    <div
      id="col-item-list"
      className="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
      style={{ width: `${settings.values.item_list_width.value}px` }}>
      <Drag width={settings.values.item_list_width.value} onResize={resizeItemList} />
      <div className="px-2 py-1 d-flex gap-1 align-items-center">
        <Tooltip label={t("show_feeds")}>
          <button
            type="button"
            className="c-button-pill d-md-none"
            aria-label={t("show_feeds")}
            onClick={() => (settings.values.feed.value = null)}>
            <Icon name="chevron-left" />
          </button>
        </Tooltip>
        <Tooltip
          label={`${t("show_first")}: ${
            settings.values.sort_newest_first.value ? t("new") : t("old")
          }`}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("show_first")}
            onClick={() => {
              const next = !settings.values.sort_newest_first.value;
              settings.update({ sort_newest_first: next });
            }}>
            <Icon
              name={settings.values.sort_newest_first.value ? "sort-new-first" : "sort-old-first"}
            />
          </button>
        </Tooltip>

        <div className="c-search flex-grow-1">
          <Icon name="search" />
          <input
            id="searchbar"
            type="search"
            className="d-block"
            value={itemList.search}
            placeholder={t("search_placeholder", { scope: searchScope.value })}
            onInput={e => (itemList.search.value = (e.target as HTMLInputElement).value)}
            onKeyDown={e => {
              e.stopPropagation();
              if (e.key === "Enter") (e.target as HTMLInputElement).blur();
              if (e.key === "Escape") {
                itemList.search.value = "";
                (e.target as HTMLInputElement).blur();
              }
            }}
          />
        </div>

        {settings.values.filter.value === "unread" && (
          <Tooltip label={t("mark_all_read")}>
            <button
              type="button"
              className="c-button-pill"
              aria-label={t("mark_all_read")}
              onClick={onMarkAllRead}>
              <Icon name="check" />
            </button>
          </Tooltip>
        )}

        {currentFeed.value ? (
          <DropdownFeed
            feed={currentFeed.value}
            folders={registry.folders.value}
            onRename={renameFeed}
            onUpdateLink={updateFeedLink}
            onMove={moveFeed}
            onMoveToNewFolder={moveFeedToNewFolder}
            onDelete={deleteFeed}
          />
        ) : currentFolder.value ? (
          <DropdownFolder
            folder={currentFolder.value}
            onRename={renameFolder}
            onDelete={deleteFolder}
          />
        ) : (
          <Tooltip label={t("settings_selected")}>
            <button type="button" className="c-button-link c-button-pill px-2" disabled>
              <Icon name="more-horizontal" />
            </button>
          </Tooltip>
        )}
      </div>

      <div
        id="item-list-scroll"
        ref={itemListRef}
        className="d-flex flex-column p-2 overflow-auto border-top flex-grow-1 gap-1">
        {itemList.items.value.map(item => (
          <div
            key={item.id}
            className="c-listitem d-flex flex-column user-select-none"
            role="radio"
            aria-checked={itemSelected.value === item.id}
            onClick={() => (itemSelected.value = item.id)}>
            <div
              style={{ lineHeight: "100%", opacity: 0.7, marginBottom: "0.1rem" }}
              className="d-flex align-items-center">
              <Icon
                small
                name={item.status === "starred" ? "star-full" : "circle-full"}
                className={`indicator me-1 ${item.status === "read" ? "is-hidden" : ""}`}
              />
              <small className="flex-fill text-truncate me-1">
                {(registry.feedsById.value[item.feed_id] || {}).title}
              </small>
              <small className="flex-shrink-0">
                <RelativeTime val={item.date} locale={settings.values.language.value} />
              </small>
            </div>
            <div className="text-break line-clamp-3">{item.title || t("untitled")}</div>
          </div>
        ))}
        {(itemList.hasMore.value || itemList.loading.value) && (
          <div
            ref={spinnerRef}
            className={`text-center my-3 ${!itemList.items.value.length ? "show-delayed" : ""}`}>
            <span className="c-spinner" />
          </div>
        )}
      </div>

      {currentFeed.value && registry.feedErrors.value[currentFeed.value.id] && (
        <div className="px-3 py-2 border-top text-danger text-break">
          {registry.feedErrors.value[currentFeed.value.id]}
        </div>
      )}
    </div>
  );
}
