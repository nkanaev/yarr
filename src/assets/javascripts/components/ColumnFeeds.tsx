import { useEffect, useRef, useCallback } from "preact/hooks";
import type { Signal } from "@preact/signals";
import { useSignal, useComputed } from "@preact/signals";
import api from "../api";
import { to } from "../utils";
import { useI18n } from "../i18n";
import type { AppSettings } from "../models/settings";
import type { DropdownHandle } from "./Dropdown";
import type { Feed, Folder, FeedStat } from "../api-types";
import Drag from "./Drag";
import Icon from "./Icon";
import Tooltip from "./Tooltip";
import FeedTree from "./FeedTree";
import type { FeedTreeNode, TreeFeedNode, FeedTreeStats } from "./FeedTree";
import NewFeed from "./NewFeed";
import { showToast } from "./Toast";
import DropdownSettings from "./DropdownSettings";
import DropdownApp from "./DropdownApp";

interface ColumnFeedsProps {
  settings: AppSettings;
  feeds: Signal<Feed[]>;
  folders: Signal<Folder[]>;
  feedStats: Signal<Record<number, FeedStat>>;
  feedErrors: Record<number, string>;
  loadingFeeds: number;
  refreshAvailable: boolean;
  selectedFeedId: Signal<number | null>;
  requiresAuth: boolean;
  handleRef: { current: ColumnFeedsHandle | null };
  refreshFeeds: () => Promise<void>;
  refreshStats: () => Promise<void> | void;
  describeError: (err: unknown) => string | undefined;
  onShowShortcuts: () => void;
}

export interface ColumnFeedsHandle {
  navigate: (dir: number) => void;
}

export default function ColumnFeeds({
  settings,
  feeds,
  folders,
  feedStats,
  feedErrors,
  loadingFeeds,
  refreshAvailable,
  selectedFeedId,
  requiresAuth,
  handleRef,
  refreshFeeds,
  refreshStats,
  describeError,
  onShowShortcuts,
}: ColumnFeedsProps) {
  const { t } = useI18n();
  const showNewFeed = useSignal(false);
  const menuDropdownRef = useRef<DropdownHandle>(null);

  const feedsById = useComputed(() => {
    return feeds.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Feed>);
  });
  const foldersById = useComputed(() => {
    return folders.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Folder>);
  });
  const currentFeed = useComputed(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    return type === "feed" ? feedsById.value[Number(guid)] || null : null;
  });
  const currentFolder = useComputed(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    return type === "folder" ? foldersById.value[Number(guid)] || null : null;
  });
  const stats = useComputed((): FeedTreeStats => {
    const statsFeeds: Record<number, { unread: number; starred: number }> = {};
    const statsFolders: Record<number, { unread: number; starred: number }> = {};
    const statsTotal = { unread: 0, starred: 0 };

    for (const feed of feeds.value) {
      const n = feedStats.value[feed.id];
      if (!n) continue;
      const fStat = { unread: n.unread || 0, starred: n.starred || 0 };
      statsFeeds[feed.id] = fStat;

      if (feed.folder_id !== null) {
        if (!statsFolders[feed.folder_id]) statsFolders[feed.folder_id] = { unread: 0, starred: 0 };
        statsFolders[feed.folder_id].unread += fStat.unread;
        statsFolders[feed.folder_id].starred += fStat.starred;
      }
      statsTotal.unread += fStat.unread;
      statsTotal.starred += fStat.starred;
    }
    return { feeds: statsFeeds, folders: statsFolders, total: statsTotal };
  });

  const resizeFeedList = (w: number) => {
    settings.update({ feed_list_width: Math.min(Math.max(200, w), 700) });
  };

  const mustHideFolder = (folder: Folder) =>
    Boolean(
      settings.values.filter.value &&
      currentFolder.value?.id !== folder.id &&
      currentFeed.value?.folder_id !== folder.id &&
      !stats.value.folders[folder.id]?.[settings.values.filter.value] &&
      feedsById.value[selectedFeedId.value ?? -1]?.folder_id !== folder.id,
    );

  const mustHideFeed = (feed: Feed) =>
    Boolean(
      settings.values.filter.value &&
      currentFeed.value?.id !== feed.id &&
      !stats.value.feeds[feed.id]?.[settings.values.filter.value] &&
      selectedFeedId.value !== feed.id,
    );

  const feedTree = useComputed((): FeedTreeNode[] => {
    const [rootFeeds, folderFeeds] = feeds.value.reduce(
      (acc, f) => (acc[f.folder_id === null ? 0 : 1].push(f), acc),
      [[] as Feed[], [] as Feed[]],
    );
    const byFolder: Record<number, Feed[]> = folderFeeds.reduce(
      (acc, f) => ((acc[f.folder_id as number] ||= []).push(f), acc),
      {} as Record<number, Feed[]>,
    );
    const feedNode = (feed: Feed): TreeFeedNode => ({ type: "feed", feed });

    return [
      ...folders.value
        .filter(folder => !mustHideFolder(folder))
        .map(folder => ({
          type: "folder" as const,
          folder,
          feeds: (byFolder[folder.id] || []).filter(f => !mustHideFeed(f)).map(feedNode),
        })),
      ...rootFeeds.filter(f => !mustHideFeed(f)).map(feedNode),
    ];
  });

  const toggleFolderExpanded = async (folder: Folder) => {
    const is_expanded = !folder.is_expanded;
    folders.value = folders.value.map(f => (f.id === folder.id ? { ...f, is_expanded } : f));
    const [err] = await to(api.folders.update(folder.id, { is_expanded }));
    if (err)
      showToast(
        { title: t("fail_save_folder"), description: describeError(err) },
        { level: "fail" },
      );
  };

  // Selecting a feed reveals it inside its folder
  useEffect(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    if (type !== "feed") return;
    const folderId = feedsById.value[Number(guid)]?.folder_id;
    if (folderId === null || folderId === undefined) return;
    const folder = foldersById.value[folderId];
    if (!folder || folder.is_expanded) return;
    void toggleFolderExpanded(folder);
  }, [settings.values.feed.value, feedsById.value, foldersById.value]);

  const fetchAllFeeds = async () => {
    if (loadingFeeds > 0) return;
    const [err] = await to(api.feeds.refresh());
    if (err)
      showToast(
        { title: t("fail_refresh"), description: describeError(err) },
        { level: "fail" },
      );
    else refreshStats();
  };

  const importOPML = async (e: Event) => {
    const input = e.target as HTMLInputElement;
    const form = input.form;
    if (!form) return;
    menuDropdownRef.current?.hide();
    const [err] = await to(api.upload_opml(form));
    if (err)
      showToast(
        { title: t("fail_import"), description: describeError(err) },
        { level: "fail" },
      );
    else {
      input.value = "";
      refreshFeeds();
      refreshStats();
    }
  };

  const logout = async () => {
    const [err] = await to(api.logout());
    if (err)
      showToast(
        { title: t("fail_logout"), description: describeError(err) },
        { level: "fail" },
      );
    else document.location.reload();
  };

  const navigateToFeed = useCallback((dir: number) => {
    const list: string[] = [""];
    for (const node of feedTree.value) {
      if (node.type === "folder") {
        list.push("folder:" + node.folder.id);
        if (node.folder.is_expanded) {
          for (const f of node.feeds) list.push("feed:" + f.feed.id);
        }
      } else {
        list.push("feed:" + node.feed.id);
      }
    }

    const curIdx = list.indexOf(settings.values.feed.value || "");
    if (curIdx === -1) {
      settings.values.feed.value = "";
      return;
    }
    const newIdx = curIdx + dir;
    if (newIdx < 0 || newIdx >= list.length) return;

    settings.update({ feed: list[newIdx] });
  }, []);

  useEffect(() => {
    handleRef.current = { navigate: navigateToFeed };
    return () => {
      handleRef.current = null;
    };
  }, [handleRef, navigateToFeed]);

  return (
    <div
      id="col-feed-list"
      className="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
      style={{ width: `${settings.values.feed_list_width.value}px` }}>
      <Drag width={settings.values.feed_list_width.value} onResize={resizeFeedList} />
      <div className="px-2 py-1 d-flex justify-content-between gap-2">
        <DropdownApp
          settings={settings}
          onShowShortcuts={onShowShortcuts}
          requiresAuth={requiresAuth}
          onLogout={logout}
        />
        <div class="d-flex flex-row gap-1">
          <Tooltip label={t("unread")}>
            <button
              type="button"
              className="c-button-pill"
              aria-label={t("unread")}
              aria-pressed={settings.values.filter.value === "unread"}
              onClick={() => settings.update({ filter: "unread" })}>
              <Icon name="circle-full" />
            </button>
          </Tooltip>
          <Tooltip label={t("starred")}>
            <button
              type="button"
              className="c-button-pill"
              aria-label={t("starred")}
              aria-pressed={settings.values.filter.value === "starred"}
              onClick={() => settings.update({ filter: "starred" })}>
              <Icon name="star-full" />
            </button>
          </Tooltip>
          <Tooltip label={t("all")}>
            <button
              type="button"
              className="c-button-pill"
              aria-label={t("all")}
              aria-pressed={settings.values.filter.value === ""}
              onClick={() => settings.update({ filter: "" })}>
              <Icon name="assorted" />
            </button>
          </Tooltip>
        </div>
        <DropdownSettings
          dropdownRef={menuDropdownRef}
          settings={settings}
          refreshAvailable={refreshAvailable}
          onFetchAllFeeds={fetchAllFeeds}
          onImportOPML={importOPML}
          onShowNewFeed={() => (showNewFeed.value = true)}
        />
      </div>
      <div class="d-flex flex-column overflow-auto border-top flex-grow-1">
        {showNewFeed.value && (
          <div className="px-2 py-3 overflow-auto flex-shrink-0 border-bottom">
            <NewFeed
              folderId={currentFeed.value?.folder_id ?? currentFolder.value?.id ?? null}
              onCancel={() => (showNewFeed.value = false)}
              onCreated={feed => {
                refreshFeeds();
                refreshStats();
                showNewFeed.value = false;
                settings.values.feed.value = "feed:" + feed.id;
              }}
            />
          </div>
        )}

        <div id="feed-list-scroll" className="p-2 overflow-auto flex-shrink-1">
          <FeedTree
            tree={feedTree.value}
            feed={settings.values.feed}
            filter={settings.values.filter}
            stats={stats.value}
            feedErrors={feedErrors}
            onSelect={val => settings.update({ feed: val })}
            onToggleFolder={toggleFolderExpanded}
          />
        </div>
      </div>

      {loadingFeeds > 0 && (
        <div className="px-2 py-1 d-flex align-items-center border-top flex-shrink-0">
          <span className="c-spinner mx-2" />
          <span className="text-truncate cursor-default user-select-none">
            {t("refreshing_progress", { count: loadingFeeds })}
          </span>
        </div>
      )}
    </div>
  );
}
