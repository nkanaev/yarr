import { useEffect, useRef, useMemo, useCallback } from "preact/hooks";
import { useSignal, useComputed, useModel } from "@preact/signals";
import { useI18n } from "../i18n";
import api, { NetworkError, HTTPError } from "../api";
import { scrollto, debounce, to } from "../utils";
import type { DropdownHandle } from "../components/Dropdown";
import Drag from "../components/Drag";
import Modal from "../components/Modal";
import Shortcuts from "../components/Shortcuts";
import RelativeTime from "../components/RelativeTime";
import Icon from "../components/Icon";
import Tooltip from "../components/Tooltip";
import FeedTree from "../components/FeedTree";
import NewFeed from "../components/NewFeed";
import ToastContainer, { showToast } from "../components/Toast";
import DropdownSettings, { Theme } from "../components/DropdownSettings";
import DropdownApp from "../components/DropdownApp";
import DropdownFeed from "../components/DropdownFeed";
import DropdownFolder from "../components/DropdownFolder";
import ArticleView from "../components/ArticleView";
import type { FeedTreeNode, TreeFeedNode } from "../components/FeedTree";
import { setupKeybindings, KeyActions } from "../key";
import { SettingsModel } from "../models/settings";
import type { Feed, Folder, Item, FeedStat, ItemStatus, ItemListQuery } from "../api-types";

type Stats = { unread: number; starred: number };

const THEME_COLORS: Record<string, string> = {
  night: "#0e0e0e",
  sepia: "#f4f0e5",
  light: "#fff",
};

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

const TITLE = document.title;

export default function App() {
  const { t, setLang } = useI18n();
  const appSettings = window.app?.settings || {};

  // Feeds & folders
  const folders = useSignal<Folder[]>([]);
  const feeds = useSignal<Feed[]>([]);
  const feedErrors = useSignal<Record<number, string>>({});
  const feedStats = useSignal<Record<number, FeedStat>>({});

  // Items
  const items = useSignal<Item[]>([]);
  const itemsHasMore = useSignal(true);
  const itemSelected = useSignal<number | null>(null);
  const itemSelectedDetails = useSignal<Item | null>(null);
  const itemSelectedReadability = useSignal("");
  const itemFullscreen = useSignal(false);
  const itemSearch = useSignal("");

  // Status & modal
  const showModal = useSignal<"" | "shortcuts">("");
  const showNewFeed = useSignal(false);
  const refreshAvailable = useSignal(true);
  const loadingFeeds = useSignal(0);
  const loadingItems = useSignal(false);
  const loadingReadability = useSignal(false);

  // Element refs
  const itemListRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const menuDropdownRef = useRef<DropdownHandle>(null);
  const activeReqIdRef = useRef(0);

  // Error helper
  const errDescription = useCallback(
    (err: unknown): string | undefined => {
      if (err instanceof HTTPError)
        return t("error_server", { code: err.status, text: err.statusText });
      if (err instanceof NetworkError) return t("error_network");
      return undefined;
    },
    [t],
  );

  const settings = useModel(
    () =>
      new SettingsModel(appSettings, err =>
        showToast(
          { title: t("fail_save_settings"), description: errDescription(err) },
          { level: "fail" },
        ),
      ),
  );

  // Theme meta tag
  const updateMetaTheme = useCallback((themeName: Theme) => {
    let active = themeName;
    if (active === "system") {
      const dark = window.matchMedia?.("(prefers-color-scheme: dark)")?.matches;
      active = dark ? "night" : "light";
    }
    const metaTag: HTMLMetaElement | null = document.querySelector("meta[name='theme-color']");
    if (metaTag) metaTag.content = THEME_COLORS[active] || "";
    document.documentElement.dataset.theme = themeName;
  }, []);

  useEffect(() => {
    updateMetaTheme(settings.values.theme_name.value);
    const mql = window.matchMedia?.("(prefers-color-scheme: dark)");
    const handler = () => updateMetaTheme(settings.values.theme_name.value);
    mql?.addEventListener("change", handler);
    return () => mql?.removeEventListener("change", handler);
  }, [settings.values.theme_name.value, updateMetaTheme]);

  const resizeFeedList = (w: number) => {
    settings.update({ feed_list_width: Math.min(Math.max(200, w), 700) });
  };

  const resizeItemList = (w: number) => {
    settings.update({ item_list_width: Math.min(Math.max(200, w), 700) });
  };

  // Computed signals
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
  const stats = useComputed(() => {
    const statsFeeds: Record<number, Stats> = {};
    const statsFolders: Record<number, Stats> = {};
    const statsTotal: Stats = { unread: 0, starred: 0 };

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

  const refreshRateTitle = useComputed(() => {
    return (
      REFRESH_RATE_OPTIONS.find(o => o.value === settings.values.refresh_rate.value)?.title || "0"
    );
  });

  // Keep document title in sync with unread count
  useEffect(() => {
    const unread = stats.value.total.unread;
    document.title = TITLE + (unread ? ` (${unread})` : "");
  }, [stats.value.total.unread]);

  // Refresh feeds & stats
  const refreshFeeds = useCallback(async () => {
    const [err, values] = await to(Promise.all([api.folders.list(), api.feeds.list()]));
    if (err) {
      showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" });
      return;
    }
    folders.value = values[0];
    feeds.value = values[1];
  }, [t, errDescription]);

  const refreshStats = useCallback(
    async (loopMode?: boolean) => {
      const [err, data] = await to(api.status());
      if (err) {
        showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" });
        return;
      }

      loadingFeeds.value = data.running;
      refreshAvailable.value = data.refresh;
      if (data.running) setTimeout(() => refreshStats(true), 500);

      feedStats.value = data.stats.reduce(
        (acc, stat) => ({ ...acc, [stat.feed_id]: stat }),
        {} as Record<number, FeedStat>,
      );

      const [feedErr, errors] = await to(api.feeds.list_errors());
      if (feedErr) {
        showToast(
          { title: t("fail_load"), description: errDescription(feedErr) },
          { level: "fail" },
        );
        return;
      }
      feedErrors.value = errors;
    },
    [t, errDescription],
  );

  useEffect(() => {
    refreshStats();
    refreshFeeds();
  }, [refreshStats, refreshFeeds]);

  // Feed tree
  const mustHideFolder = (folder: Folder) =>
    Boolean(
      settings.values.filter.value &&
      currentFolder.value?.id !== folder.id &&
      currentFeed.value?.folder_id !== folder.id &&
      !stats.value.folders[folder.id]?.[settings.values.filter.value] &&
      feedsById.value[itemSelectedDetails.value?.feed_id ?? -1]?.folder_id !== folder.id,
    );

  const mustHideFeed = (feed: Feed) =>
    Boolean(
      settings.values.filter.value &&
      currentFeed.value?.id !== feed.id &&
      !stats.value.feeds[feed.id]?.[settings.values.filter.value] &&
      itemSelectedDetails.value?.feed_id !== feed.id,
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

  // Items query & refresh
  const getItemsQuery = useCallback((): ItemListQuery => {
    const query: ItemListQuery = {};
    if (settings.values.feed.value) {
      const [type, guid] = settings.values.feed.value.split(":", 2);
      if (type === "feed") query.feed_id = guid;
      else if (type === "folder") query.folder_id = guid;
    }
    if (settings.values.filter.value) query.status = settings.values.filter.value;
    if (itemSearch.value) query.search = itemSearch.value;
    if (!settings.values.sort_newest_first.value) query.oldest_first = true;
    return query;
  }, []);

  const refreshItems = useCallback(
    async (loadMore = false) => {
      if (settings.values.feed.value === null) {
        items.value = [];
        itemsHasMore.value = false;
        loadingItems.value = false;
        return;
      }

      const reqId = ++activeReqIdRef.current;

      if (!loadMore) {
        items.value = [];
        itemsHasMore.value = true;
        if (itemListRef.current) itemListRef.current.scrollTop = 0;
      }

      const query = getItemsQuery();
      if (loadMore && items.value.length > 0) {
        query.after = items.value[items.value.length - 1].id;
      }

      loadingItems.value = true;
      const [err, data] = await to(api.items.list(query));

      if (reqId !== activeReqIdRef.current) return;
      loadingItems.value = false;

      if (err) {
        showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" });
        return;
      }

      items.value = loadMore ? items.value.concat(data.list) : data.list;
      itemsHasMore.value = data.has_more;

      if (data.has_more && !loadMore) {
        setTimeout(() => {
          const el = itemListRef.current;
          if (!el || el.scrollHeight === 0) return;
          const scale =
            (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
          if (el.scrollHeight - el.scrollTop - el.offsetHeight < 70 * scale) {
            refreshItems(true);
          }
        }, 0);
      }
    },
    [getItemsQuery, t, errDescription],
  );

  useEffect(() => {
    itemSelected.value = null;
    refreshItems(false);
  }, [
    settings.values.feed.value,
    settings.values.filter.value,
    settings.values.sort_newest_first.value,
  ]);

  useEffect(() => {
    const timer = setTimeout(() => refreshItems(false), 500);
    return () => clearTimeout(timer);
  }, [itemSearch.value]);

  const loadMoreItems = useCallback(() => {
    if (!itemsHasMore.value || loadingItems.value) return;
    const el = itemListRef.current;
    if (!el || el.scrollHeight === 0) return;
    const scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
    const closeToBottom = el.scrollHeight - el.scrollTop - el.offsetHeight < 70 * scale;
    if (
      closeToBottom ||
      (itemSelected.value &&
        items.value.length &&
        itemSelected.value === items.value[items.value.length - 1].id)
    ) {
      refreshItems(true);
    }
  }, [refreshItems]);

  const onItemListScroll = useMemo(() => debounce(loadMoreItems, 200), [loadMoreItems]);

  // Item details & auto-read
  useEffect(() => {
    itemSelectedReadability.value = "";
    if (itemSelected.value === null) {
      itemSelectedDetails.value = null;
      itemFullscreen.value = false;
      return;
    }

    if (contentRef.current) contentRef.current.scrollTop = 0;

    to(api.items.get(itemSelected.value)).then(([itemErr, item]) => {
      if (itemErr) {
        showToast(
          { title: t("fail_load"), description: errDescription(itemErr) },
          { level: "fail" },
        );
        return;
      }
      itemSelectedDetails.value = item;

      if (item.status === "unread") {
        to(api.items.update(item.id, { status: "read" })).then(([updateErr]) => {
          if (updateErr) {
            showToast(
              { title: t("fail_update_article"), description: errDescription(updateErr) },
              { level: "fail" },
            );
            return;
          }
          const cur = feedStats.value[item.feed_id];
          if (cur) {
            feedStats.value = {
              ...feedStats.value,
              [item.feed_id]: { ...cur, unread: Math.max(0, (cur.unread || 1) - 1) },
            };
          }
          items.value = items.value.map(i => (i.id === item.id ? { ...i, status: "read" } : i));
          item.status = "read";
        });
      }
    });
  }, [itemSelected.value, t, errDescription]);

  // Status toggle
  const toggleItemStatus = async (item: Item, targetstatus: ItemStatus) => {
    const oldstatus = item.status;
    const newstatus = item.status !== targetstatus ? targetstatus : "read";

    const [err] = await to(api.items.update(item.id, { status: newstatus }));
    if (err) {
      showToast(
        { title: t("fail_update_article"), description: errDescription(err) },
        { level: "fail" },
      );
      return;
    }

    const cur = feedStats.value[item.feed_id];
    if (cur) {
      let unread = cur.unread || 0;
      let starred = cur.starred || 0;
      if (oldstatus === "unread") unread -= 1;
      if (oldstatus === "starred") starred -= 1;
      if (newstatus === "unread") unread += 1;
      if (newstatus === "starred") starred += 1;
      feedStats.value = { ...feedStats.value, [item.feed_id]: { ...cur, unread, starred } };
    }

    items.value = items.value.map(i => (i.id === item.id ? { ...i, status: newstatus } : i));
    item.status = newstatus;
    if (itemSelectedDetails.value?.id === item.id) {
      itemSelectedDetails.value = { ...itemSelectedDetails.value, status: newstatus };
    }
  };

  const markItemsRead = async () => {
    const markQuery = getItemsQuery();
    const [err] = await to(
      api.items.mark_read({ folder_id: markQuery.folder_id, feed_id: markQuery.feed_id }),
    );
    if (err) {
      showToast(
        { title: t("fail_update_article"), description: errDescription(err) },
        { level: "fail" },
      );
      return;
    }
    items.value = [];
    itemSelected.value = null;
    itemsHasMore.value = false;
    refreshStats();
  };

  // Feed & folder actions
  const renameFeed = async (feed: Feed) => {
    const title = prompt(t("prompt_new_title"), feed.title);
    if (!title) return;
    const [err] = await to(api.feeds.update(feed.id, { title }));
    if (err)
      showToast(
        { title: t("fail_save_feed"), description: errDescription(err) },
        { level: "fail" },
      );
    else feeds.value = feeds.value.map(f => (f.id === feed.id ? { ...f, title } : f));
  };

  const updateFeedLink = async (feed: Feed) => {
    const feed_link = prompt(t("prompt_feed_link"), feed.feed_link);
    if (feed_link === null) return;
    const [err] = await to(api.feeds.update(feed.id, { feed_link }));
    if (err)
      showToast(
        { title: t("fail_save_feed"), description: errDescription(err) },
        { level: "fail" },
      );
    else feeds.value = feeds.value.map(f => (f.id === feed.id ? { ...f, feed_link } : f));
  };

  const moveFeed = async (feed: Feed, folder_id: number | null) => {
    const [err] = await to(api.feeds.update(feed.id, { folder_id }));
    if (err)
      showToast(
        { title: t("fail_save_feed"), description: errDescription(err) },
        { level: "fail" },
      );
    else {
      feeds.value = feeds.value.map(f => (f.id === feed.id ? { ...f, folder_id } : f));
      refreshStats();
    }
  };

  const moveFeedToNewFolder = async (feed: Feed) => {
    const title = prompt(t("prompt_folder_name"));
    if (!title) return;
    const [folderErr, folder] = await to(api.folders.create({ title }));
    if (folderErr) {
      showToast(
        { title: t("fail_save_folder"), description: errDescription(folderErr) },
        { level: "fail" },
      );
      return;
    }
    const [updateErr] = await to(api.feeds.update(feed.id, { folder_id: folder.id }));
    if (updateErr) {
      showToast(
        { title: t("fail_save_feed"), description: errDescription(updateErr) },
        { level: "fail" },
      );
      return;
    }
    await refreshFeeds();
    refreshStats();
  };

  const deleteFeed = async (feed: Feed) => {
    if (!confirm(t("confirm_delete", { name: feed.title }))) return;
    const [err] = await to(api.feeds.delete(feed.id));
    if (err)
      showToast(
        { title: t("fail_save_feed"), description: errDescription(err) },
        { level: "fail" },
      );
    else {
      settings.values.feed.value = null;
      refreshStats();
      refreshFeeds();
    }
  };

  const renameFolder = async (folder: Folder) => {
    const title = prompt(t("prompt_new_title"), folder.title);
    if (!title) return;
    const [err] = await to(api.folders.update(folder.id, { title }));
    if (err)
      showToast(
        { title: t("fail_save_folder"), description: errDescription(err) },
        { level: "fail" },
      );
    else {
      folders.value = folders.value
        .map(f => (f.id === folder.id ? { ...f, title } : f))
        .sort((a, b) => a.title.localeCompare(b.title));
    }
  };

  const deleteFolder = async (folder: Folder) => {
    if (!confirm(t("confirm_delete", { name: folder.title }))) return;
    const [err] = await to(api.folders.delete(folder.id));
    if (err)
      showToast(
        { title: t("fail_save_folder"), description: errDescription(err) },
        { level: "fail" },
      );
    else {
      settings.values.feed.value = null;
      refreshStats();
      refreshFeeds();
    }
  };

  const toggleFolderExpanded = async (folder: Folder) => {
    const is_expanded = !folder.is_expanded;
    folders.value = folders.value.map(f => (f.id === folder.id ? { ...f, is_expanded } : f));
    const [err] = await to(api.folders.update(folder.id, { is_expanded }));
    if (err)
      showToast(
        { title: t("fail_save_folder"), description: errDescription(err) },
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
  }, [settings.values.feed.value, feedsById.value]);

  const toggleReadability = async () => {
    if (itemSelectedReadability.value) {
      itemSelectedReadability.value = "";
      return;
    }
    if (!itemSelectedDetails.value?.link) return;
    loadingReadability.value = true;
    const [err, data] = await to(api.crawl(itemSelectedDetails.value.link));
    loadingReadability.value = false;
    if (err)
      showToast(
        { title: t("fail_readability"), description: errDescription(err) },
        { level: "fail" },
      );
    else itemSelectedReadability.value = data?.content || "";
  };

  const fetchAllFeeds = async () => {
    if (loadingFeeds.value) return;
    const [err] = await to(api.feeds.refresh());
    if (err)
      showToast({ title: t("fail_refresh"), description: errDescription(err) }, { level: "fail" });
    else refreshStats();
  };

  const importOPML = async (e: Event) => {
    const input = e.target as HTMLInputElement;
    const form = input.form;
    if (!form) return;
    menuDropdownRef.current?.hide();
    const [err] = await to(api.upload_opml(form));
    if (err)
      showToast({ title: t("fail_import"), description: errDescription(err) }, { level: "fail" });
    else {
      input.value = "";
      refreshFeeds();
      refreshStats();
    }
  };

  const changeRefreshRate = (offset: number) => {
    const curIdx = REFRESH_RATE_OPTIONS.findIndex(
      o => o.value === settings.values.refresh_rate.value,
    );
    if ((curIdx <= 0 && offset < 0) || (curIdx >= REFRESH_RATE_OPTIONS.length - 1 && offset > 0))
      return;
    const newRate = REFRESH_RATE_OPTIONS[curIdx + offset].value;
    settings.update({ refresh_rate: newRate });
  };

  // Keyboard navigation
  const navigateToItem = useCallback(
    (dir: number) => {
      if (itemSelected.value === null) {
        if (items.value.length) itemSelected.value = items.value[0].id;
        return;
      }
      const curIdx = items.value.findIndex(x => x.id === itemSelected.value);
      if (curIdx === -1) {
        if (items.value.length) itemSelected.value = items.value[0].id;
        return;
      }
      const newIdx = curIdx + dir;
      if (newIdx < 0 || newIdx >= items.value.length) return;
      itemSelected.value = items.value[newIdx].id;

      setTimeout(() => {
        const scroll = document.querySelector("#item-list-scroll");
        const handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);
        loadMoreItems();
      }, 0);
    },
    [loadMoreItems],
  );

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

    const nextVal = list[newIdx];
    settings.update({ feed: nextVal });

    setTimeout(() => {
      const scroll = document.querySelector("#feed-list-scroll");
      const handle = scroll?.querySelector('[aria-checked="true"]');
      if (handle && scroll) scrollto(handle, scroll);
    }, 0);
  }, []);

  // Keybindings
  const keyActionsRef = useRef<KeyActions>({} as KeyActions);
  keyActionsRef.current = {
    enabled: () => showModal.value === "",
    openItemLink: () => {
      if (itemSelectedDetails.value?.link)
        window.open(itemSelectedDetails.value.link, "_blank", "noopener,noreferrer");
    },
    toggleReadability,
    toggleItemRead: () =>
      itemSelectedDetails.value && toggleItemStatus(itemSelectedDetails.value, "unread"),
    markAllRead: () => settings.values.filter.value === "unread" && markItemsRead(),
    toggleItemStarred: () =>
      itemSelectedDetails.value && toggleItemStatus(itemSelectedDetails.value, "starred"),
    nextItem: () => navigateToItem(+1),
    prevItem: () => navigateToItem(-1),
    nextFeed: () => navigateToFeed(+1),
    prevFeed: () => navigateToFeed(-1),
    closeItem: () => (itemSelected.value = null),
    showShortcuts: () => (showModal.value = "shortcuts"),
    showAll: () => settings.update({ filter: "" }),
    showUnread: () => settings.update({ filter: "unread" }),
    showStarred: () => settings.update({ filter: "starred" }),
  };

  useEffect(() => setupKeybindings(keyActionsRef), []);

  return (
    <div
      className={`d-flex ${settings.values.feed.value !== null ? "feed-selected" : ""} ${
        itemSelected.value !== null ? "item-selected" : ""
      } ${itemFullscreen.value ? "item-fullscreen" : ""}`}>
      {/* Feed list column */}
      <div
        id="col-feed-list"
        className="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
        style={{ width: `${settings.values.feed_list_width.value}px` }}>
        <Drag width={settings.values.feed_list_width.value} onResize={resizeFeedList} />
        <div className="px-2 py-1 d-flex justify-content-between gap-2">
          <DropdownApp
            theme={settings.values.theme_name.value}
            onThemeChange={th => settings.update({ theme_name: th })}
            onShowShortcuts={() => (showModal.value = "shortcuts")}
            language={settings.values.language.value}
            onLanguageChange={code => {
              setLang(code);
              settings.update({ language: code });
            }}
            requiresAuth={window.app?.requiresAuth}
            onLogout={async () => {
              const [err] = await to(api.logout());
              if (err)
                showToast(
                  { title: t("fail_logout"), description: errDescription(err) },
                  { level: "fail" },
                );
              else document.location.reload();
            }}
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
            refreshAvailable={refreshAvailable.value}
            refreshRate={settings.values.refresh_rate.value}
            refreshRateTitle={refreshRateTitle.value}
            maxRefreshRate={
              settings.values.refresh_rate.value ===
              REFRESH_RATE_OPTIONS[REFRESH_RATE_OPTIONS.length - 1].value
            }
            onRefreshRateChange={changeRefreshRate}
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
              value={settings.values.feed.value}
              filterSelected={settings.values.filter.value}
              stats={stats.value}
              feedErrors={feedErrors.value}
              onSelect={val => settings.update({ feed: val })}
              onToggleFolder={toggleFolderExpanded}
            />
          </div>
        </div>

        {loadingFeeds.value > 0 && (
          <div className="px-2 py-1 d-flex align-items-center border-top flex-shrink-0">
            <span className="c-spinner mx-2" />
            <span className="text-truncate cursor-default user-select-none">
              {t("refreshing_progress", { count: loadingFeeds.value })}
            </span>
          </div>
        )}
      </div>

      {/* Item list column */}
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
              value={itemSearch}
              placeholder={t("search_placeholder", { scope: searchScope.value })}
              onInput={e => (itemSearch.value = (e.target as HTMLInputElement).value)}
              onKeyDown={e => {
                e.stopPropagation();
                if (e.key === "Enter") (e.target as HTMLInputElement).blur();
                if (e.key === "Escape") {
                  itemSearch.value = "";
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
                onClick={markItemsRead}>
                <Icon name="check" />
              </button>
            </Tooltip>
          )}

          {currentFeed.value ? (
            <DropdownFeed
              feed={currentFeed.value}
              folders={folders.value}
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
          className="d-flex flex-column p-2 overflow-auto border-top flex-grow-1 gap-1"
          onScroll={onItemListScroll}>
          {items.value.map(item => (
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
                  name="circle-full"
                  className={`indicator me-1 ${item.status === "unread" ? "" : "is-hidden"}`}
                />
                <Icon
                  small
                  name="star-full"
                  className={`indicator me-1 ${item.status === "starred" ? "" : "is-hidden"}`}
                />
                <small className="flex-fill text-truncate me-1">
                  {(feedsById.value[item.feed_id] || {}).title}
                </small>
                <small className="flex-shrink-0">
                  <RelativeTime val={item.date} locale={settings.values.language.value} />
                </small>
              </div>
              <div className="text-break line-clamp-3">{item.title || t("untitled")}</div>
            </div>
          ))}
          {(itemsHasMore.value || loadingItems.value) && (
            <div className={`text-center my-3 ${items.value.length === 0 ? "show-delayed" : ""}`}>
              <span className="c-spinner" />
            </div>
          )}
        </div>

        {currentFeed.value && feedErrors.value[currentFeed.value.id] && (
          <div className="px-3 py-2 border-top text-danger text-break">
            {feedErrors.value[currentFeed.value.id]}
          </div>
        )}
      </div>

      {/* Item reader column */}
      <div id="col-item" className="vh-100 d-flex flex-column w-100" style={{ minWidth: 0 }}>
        <ArticleView
          item={itemSelectedDetails.value}
          feedTitle={(feedsById.value[itemSelectedDetails.value?.feed_id ?? -1] || {}).title || ""}
          feedId={(feedsById.value[itemSelectedDetails.value?.feed_id ?? -1] || {}).id}
          theme={{ font: settings.values.theme_font.value, size: settings.values.theme_size.value }}
          onFontChange={font => settings.update({ theme_font: font })}
          onFontSizeChange={delta => {
            const size = +(settings.values.theme_size.value + 0.1 * delta).toFixed(1);
            settings.update({ theme_size: size });
          }}
          readabilityContent={itemSelectedReadability.value}
          loadingReadability={loadingReadability.value}
          onToggleReadability={toggleReadability}
          onToggleStarred={item => toggleItemStatus(item, "starred")}
          onToggleRead={item => toggleItemStatus(item, "unread")}
          onSelectFeed={id => (settings.values.feed.value = `feed:${id}`)}
          onNavigate={navigateToItem}
          canPrev={Boolean(items.value.length && itemSelected.value !== items.value[0].id)}
          canNext={Boolean(
            items.value.length && itemSelected.value !== items.value[items.value.length - 1].id,
          )}
          isFullscreen={itemFullscreen.value}
          onToggleFullscreen={() => (itemFullscreen.value = !itemFullscreen.value)}
          onClose={() => (itemSelected.value = null)}
          contentRef={contentRef}
          language={settings.values.language.value}
        />
      </div>

      {/* Dialog modal */}
      <Modal open={showModal.value !== ""} onHide={() => (showModal.value = "")}>
        <button
          type="button"
          className="c-button-link outline-none position-absolute top-0 end-0 p-2 m-2"
          style={{ lineHeight: 1 }}
          onClick={() => (showModal.value = "")}>
          <Icon name="x" />
        </button>
        {showModal.value === "shortcuts" && <Shortcuts />}
      </Modal>

      <ToastContainer />
    </div>
  );
}
