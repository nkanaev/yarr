import { useEffect, useRef, useCallback } from "preact/hooks";
import type { Signal } from "@preact/signals";
import { useSignal, useComputed } from "@preact/signals";
import api from "../api";
import { scrollto, to } from "../utils";
import { useI18n } from "../i18n";
import type { AppSettings } from "../models/settings";
import type { AppRegistry } from "../models/registry";
import type { DropdownHandle } from "./Dropdown";
import type { Feed, Folder } from "../api-types";
import Drag from "./Drag";
import Icon from "./Icon";
import Tooltip from "./Tooltip";
import NewFeed from "./NewFeed";
import { showToast } from "./Toast";
import DropdownSettings from "./DropdownSettings";
import DropdownApp from "./DropdownApp";

interface TreeFeedNode {
  type: "feed";
  feed: Feed;
}

interface TreeFolderNode {
  type: "folder";
  folder: Folder;
  feeds: TreeFeedNode[];
}

type FeedTreeNode = TreeFolderNode | TreeFeedNode;

interface ColumnFeedListProps {
  settings: AppSettings;
  registry: AppRegistry;
  selectedFeedId: Signal<number | null>;
  requiresAuth: boolean;
  handleRef: { current: ColumnFeedListHandle | null };
  onShowShortcuts: () => void;
}

export interface ColumnFeedListHandle {
  navigate: (dir: number) => void;
}

export default function ColumnFeedList({
  settings,
  registry,
  selectedFeedId,
  requiresAuth,
  handleRef,
  onShowShortcuts,
}: ColumnFeedListProps) {
  const { t } = useI18n();
  const { feed, filter } = settings.values;
  const showNewFeed = useSignal(false);
  const menuDropdownRef = useRef<DropdownHandle>(null);
  const refRoot = useRef<HTMLDivElement>(null);
  const refFeed = useRef<HTMLDivElement>(null);

  const currentFeed = useComputed(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    return type === "feed" ? registry.feedsById.value[Number(guid)] || null : null;
  });
  const currentFolder = useComputed(() => {
    const [type, guid] = (settings.values.feed.value || "").split(":", 2);
    return type === "folder" ? registry.foldersById.value[Number(guid)] || null : null;
  });

  const resizeFeedList = (w: number) => {
    settings.update({ feed_list_width: Math.min(Math.max(200, w), 700) });
  };

  const mustHideFolder = (folder: Folder) =>
    Boolean(
      settings.values.filter.value &&
      currentFolder.value?.id !== folder.id &&
      currentFeed.value?.folder_id !== folder.id &&
      !registry.stats.value.folders[folder.id]?.[settings.values.filter.value] &&
      registry.feedsById.value[selectedFeedId.value ?? -1]?.folder_id !== folder.id,
    );

  const mustHideFeed = (feed: Feed) =>
    Boolean(
      settings.values.filter.value &&
      currentFeed.value?.id !== feed.id &&
      !registry.stats.value.feeds[feed.id]?.[settings.values.filter.value] &&
      selectedFeedId.value !== feed.id,
    );

  const feedTree = useComputed((): FeedTreeNode[] => {
    const [rootFeeds, folderFeeds] = registry.feeds.value.reduce(
      (acc, f) => (acc[f.folder_id === null ? 0 : 1].push(f), acc),
      [[] as Feed[], [] as Feed[]],
    );
    const byFolder: Record<number, Feed[]> = folderFeeds.reduce(
      (acc, f) => ((acc[f.folder_id as number] ||= []).push(f), acc),
      {} as Record<number, Feed[]>,
    );
    const feedNode = (feed: Feed): TreeFeedNode => ({ type: "feed", feed });

    return [
      ...registry.folders.value
        .filter(folder => !mustHideFolder(folder))
        .map(folder => ({
          type: "folder" as const,
          folder,
          feeds: (byFolder[folder.id] || []).filter(f => !mustHideFeed(f)).map(feedNode),
        })),
      ...rootFeeds.filter(f => !mustHideFeed(f)).map(feedNode),
    ];
  });

  const allTitle =
    filter.value === "unread"
      ? t("all_unread")
      : filter.value === "starred"
        ? t("all_starred")
        : t("all_feeds");

  const toggleFolderExpanded = async (folder: Folder) => {
    const is_expanded = !folder.is_expanded;
    registry.folders.value = registry.folders.value.map(f =>
      f.id === folder.id ? { ...f, is_expanded } : f,
    );
    const [err] = await to(api.folders.update(folder.id, { is_expanded }));
    if (err)
      showToast({ title: "fail_save_folder", description: err }, { level: "fail" });
  };

  // Selecting a feed reveals it inside its folder
  useEffect(() => {
    const [type, guid] = (feed.value || "").split(":", 2);
    if (type !== "feed") return;
    const folderId = registry.feedsById.value[Number(guid)]?.folder_id;
    if (folderId === null || folderId === undefined) return;
    const folder = registry.foldersById.value[folderId];
    if (!folder || folder.is_expanded) return;
    void toggleFolderExpanded(folder);
  }, [settings.values.feed.value]);

  // Scroll the selected node into view, retrying until it has been rendered
  const scrolledFor = useRef<string | null>(null);
  useEffect(() => {
    const key = feed.value || "";
    if (scrolledFor.current === key) return;

    const frame = requestAnimationFrame(() => {
      const target = refFeed.current;
      const scroll = refRoot.current;
      if (!target || !scroll) return;
      scrolledFor.current = key;
      scrollto(target, scroll);
    });
    return () => cancelAnimationFrame(frame);
  }, [feed.value, registry.feeds.value]);

  const fetchAllFeeds = async () => {
    if (registry.loadingFeeds.value > 0) return;
    const [err] = await to(api.feeds.refresh());
    if (err)
      showToast({ title: "fail_refresh", description: err }, { level: "fail" });
    else registry.refreshStats();
  };

  const importOPML = async (e: Event) => {
    const input = e.target as HTMLInputElement;
    const form = input.form;
    if (!form) return;
    menuDropdownRef.current?.hide();
    const [err] = await to(api.upload_opml(form));
    if (err)
      showToast({ title: "fail_import", description: err }, { level: "fail" });
    else {
      input.value = "";
      registry.refresh();
      registry.refreshStats();
    }
  };

  const logout = async () => {
    const [err] = await to(api.logout());
    if (err)
      showToast({ title: "fail_logout", description: err }, { level: "fail" });
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
          refreshAvailable={registry.refreshAvailable.value}
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
                registry.refresh();
                registry.refreshStats();
                showNewFeed.value = false;
                settings.values.feed.value = "feed:" + feed.id;
              }}
            />
          </div>
        )}

        <div id="feed-list-scroll" ref={refRoot} className="p-2 overflow-auto flex-shrink-1">
          <div className="d-flex flex-column gap-1">
            <div
              ref={feed.value === "" ? refFeed : undefined}
              className="c-listitem d-flex user-select-none gap-2"
              role="radio"
              aria-checked={feed.value === ""}
              onClick={() => settings.update({ feed: "" })}>
              <div className="flex-shrink-0 d-flex">
                <Icon name="layers" />
              </div>
              <div className="flex-grow-1 min-w-0 text-truncate">{allTitle}</div>
              {filter.value && (
                <div className="flex-shrink-0">
                  <span className="ps-2 text-end opacity-50">
                    {registry.stats.value.total[filter.value]}
                  </span>
                </div>
              )}
            </div>

            {feedTree.value.map(node => {
              if (node.type === "folder") {
                const folderKey = "folder:" + node.folder.id;
                const isSelected = feed.value === folderKey;
                return (
                  <div key={folderKey}>
                    <div
                      ref={isSelected ? refFeed : undefined}
                      className="c-listitem d-flex user-select-none gap-2"
                      role="radio"
                      aria-checked={isSelected}
                      onClick={() => settings.update({ feed: folderKey })}>
                      <div className="flex-shrink-0 d-flex">
                        <div
                          onClick={e => {
                            e.stopPropagation();
                            toggleFolderExpanded(node.folder);
                          }}
                          className="p-2 m-n2 d-inline-flex">
                          <Icon name={node.folder.is_expanded ? "chevron-down" : "chevron-right"} />
                        </div>
                      </div>
                      <div className="flex-grow-1 min-w-0 text-truncate">{node.folder.title}</div>
                      {filter.value && (
                        <div className="flex-shrink-0">
                          <span className="ps-2 text-end opacity-50">
                            {registry.stats.value.folders[node.folder.id]?.[filter.value]}
                          </span>
                        </div>
                      )}
                    </div>
                    {node.folder.is_expanded && node.feeds.length > 0 && (
                      <div className="d-flex flex-column gap-1 ps-3 pt-1">
                        {node.feeds.map(feedNode => {
                          const feedKey = "feed:" + feedNode.feed.id;
                          return (
                            <div
                              key={feedKey}
                              ref={feed.value === feedKey ? refFeed : undefined}
                              className="c-listitem d-flex user-select-none gap-2"
                              role="radio"
                              aria-checked={feed.value === feedKey}
                              onClick={() => settings.update({ feed: feedKey })}>
                              <div className="flex-shrink-0 d-flex align-items-center">
                                {!feedNode.feed.icon ? (
                                  <Icon className="flex-shrink-0" name="rss" />
                                ) : (
                                  <img
                                    className="c-icon flex-shrink-0"
                                    src={feedNode.feed.icon}
                                    alt=""
                                    loading="lazy"
                                  />
                                )}
                              </div>
                              <div className="flex-grow-1 min-w-0 text-truncate">
                                {feedNode.feed.title}
                              </div>
                              <div className="flex-shrink-0 d-flex">
                                {filter.value && (
                                  <span className="ps-2 text-end opacity-50">
                                    {registry.stats.value.feeds[feedNode.feed.id]?.[filter.value]}
                                  </span>
                                )}
                                {!filter.value && registry.feedErrors.value[feedNode.feed.id] && (
                                  <Tooltip label={registry.feedErrors.value[feedNode.feed.id]}>
                                    <Icon className="flex-shrink-0" name="alert-circle" />
                                  </Tooltip>
                                )}
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    )}
                  </div>
                );
              }

              const feedKey = "feed:" + node.feed.id;
              return (
                <div
                  key={feedKey}
                  ref={feed.value === feedKey ? refFeed : undefined}
                  className="c-listitem d-flex user-select-none gap-2"
                  role="radio"
                  aria-checked={feed.value === feedKey}
                  onClick={() => settings.update({ feed: feedKey })}>
                  <div className="flex-shrink-0 d-flex">
                    {!node.feed.icon ? (
                      <Icon className="flex-shrink-0" name="rss" />
                    ) : (
                      <img
                        className="c-icon flex-shrink-0"
                        src={node.feed.icon}
                        alt=""
                        loading="lazy"
                      />
                    )}
                  </div>
                  <div className="flex-grow-1 min-w-0 text-truncate">{node.feed.title}</div>
                  <div className="flex-shrink-0 d-flex">
                    {filter.value && (
                      <span className="ps-2 text-end opacity-50">
                        {registry.stats.value.feeds[node.feed.id]?.[filter.value]}
                      </span>
                    )}
                    {!filter.value && registry.feedErrors.value[node.feed.id] && (
                      <Tooltip label={registry.feedErrors.value[node.feed.id]}>
                        <Icon className="flex-shrink-0" name="alert-circle" />
                      </Tooltip>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {registry.loadingFeeds.value > 0 && (
        <div className="px-2 py-1 d-flex align-items-center border-top flex-shrink-0">
          <span className="c-spinner mx-2" />
          <span className="text-truncate cursor-default user-select-none">
            {t("refreshing_progress", { count: registry.loadingFeeds.value })}
          </span>
        </div>
      )}
    </div>
  );
}
