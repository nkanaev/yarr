import { useEffect, useRef, useMemo, useCallback } from "preact/hooks";
import { useSignal, useComputed, signal } from "@preact/signals";
import type { Lang } from "../i18n";
import { useI18n } from "../i18n";
import api, { NetworkError, HTTPError } from "../api";
import { scrollto, debounce, to } from "../utils";
import type { DropdownHandle } from "../components/Dropdown";
import Drag from "../components/Drag";
import Modal from "../components/Modal";
import Shortcuts from "../components/Shortcuts";
import RelativeTime from "../components/RelativeTime";
import Icon from "../components/Icon";
import FeedTree from "../components/FeedTree";
import NewFeed from "../components/NewFeed";
import ToastContainer, { showToast } from "../components/Toast";
import DropdownSettings, { Theme } from "../components/DropdownSettings";
import DropdownFeed from "../components/DropdownFeed";
import DropdownFolder from "../components/DropdownFolder";
import ArticleView, { ThemeFont } from "../components/ArticleView";
import type { FeedTreeNode, TreeFeedNode } from "../components/FeedTree";
import { setupKeybindings, KeyActions } from "../key";
import type {
  Feed,
  Folder,
  Item,
  FeedStat,
  ItemStatus,
  ItemListQuery,
  ItemMarkQuery,
  Settings,
} from "../api-types";

type Filter = "" | "starred" | "unread";
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

  // Signals for state
  const filterSelected = useSignal<Filter>((appSettings.filter as Filter) || "");
  const theme = useSignal<{ name: Theme; font: ThemeFont; size: number }>({
    name: (appSettings.theme_name as Theme) || "system",
    font: (appSettings.theme_font as ThemeFont) || "",
    size: (appSettings.theme_size as number) || 1,
  });
  const refreshRate = useSignal<number>(appSettings.refresh_rate || 0);
  const language = useSignal<Lang>(appSettings.language || "en");

  // Feeds & folders
  const folders = useSignal<Folder[]>([]);
  const feeds = useSignal<Feed[]>([]);
  const feedSelected = useSignal<string | null>(appSettings.feed ?? null);
  const feedListWidth = useSignal<number>(appSettings.feed_list_width || 300);
  const feedErrors = useSignal<Record<number, string>>({});
  const feedStats = useSignal<Record<number, FeedStat>>({});

  // Items
  const itemListWidth = useSignal<number>(appSettings.item_list_width || 300);
  const items = useSignal<Item[]>([]);
  const itemsHasMore = useSignal(true);
  const itemSelected = useSignal<number | null>(null);
  const itemSelectedDetails = useSignal<Item | null>(null);
  const itemSelectedReadability = useSignal("");
  const itemFullscreen = useSignal(false);
  const itemSearch = useSignal("");
  const itemSortNewestFirst = useSignal<boolean>(appSettings.sort_newest_first ?? true);

  // Status & modal
  const showModal = useSignal<"" | "shortcuts" | "newfeed">("");
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

  // Consolidated settings update
  const updateSettings = useCallback(
    (partial: Partial<Settings>) => {
      api.settings.update(partial).catch(err => {
        showToast(
          { title: t("fail_save_settings"), description: errDescription(err) },
          { level: "fail" },
        );
      });
    },
    [t, errDescription],
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
    updateMetaTheme(theme.value.name);
    const mql = window.matchMedia?.("(prefers-color-scheme: dark)");
    const handler = () => updateMetaTheme(theme.value.name);
    mql?.addEventListener("change", handler);
    return () => mql?.removeEventListener("change", handler);
  }, [theme.value.name, updateMetaTheme]);

  // Debounced width saving
  const saveFeedListWidth = useMemo(
    () => debounce((w: number) => updateSettings({ feed_list_width: w }), 1000),
    [updateSettings],
  );
  const saveItemListWidth = useMemo(
    () => debounce((w: number) => updateSettings({ item_list_width: w }), 1000),
    [updateSettings],
  );

  const resizeFeedList = (w: number) => {
    const clamped = Math.min(Math.max(200, w), 700);
    feedListWidth.value = clamped;
    saveFeedListWidth(clamped);
  };

  const resizeItemList = (w: number) => {
    const clamped = Math.min(Math.max(200, w), 700);
    itemListWidth.value = clamped;
    saveItemListWidth(clamped);
  };

  // Computed signals
  const feedsById = useComputed(() => {
    return feeds.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Feed>);
  });
  const foldersById = useComputed(() => {
    return folders.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Folder>);
  });
  const currentFeed = useComputed(() => {
    const [type, guid] = (feedSelected.value || "").split(":", 2);
    return type === "feed" ? feedsById.value[Number(guid)] || null : null;
  });
  const currentFolder = useComputed(() => {
    const [type, guid] = (feedSelected.value || "").split(":", 2);
    return type === "folder" ? foldersById.value[Number(guid)] || null : null;
  });
  const searchScope = useComputed(() => {
    return (
      currentFeed.value?.title ||
      currentFolder.value?.title ||
      (filterSelected.value === "unread"
        ? t("all_unread")
        : filterSelected.value === "starred"
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
    return REFRESH_RATE_OPTIONS.find(o => o.value === refreshRate.value)?.title || "0";
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
      filterSelected.value &&
      currentFolder.value?.id !== folder.id &&
      currentFeed.value?.folder_id !== folder.id &&
      !stats.value.folders[folder.id]?.[filterSelected.value] &&
      feedsById.value[itemSelectedDetails.value?.feed_id ?? -1]?.folder_id !== folder.id,
    );

  const mustHideFeed = (feed: Feed) =>
    Boolean(
      filterSelected.value &&
      currentFeed.value?.id !== feed.id &&
      !stats.value.feeds[feed.id]?.[filterSelected.value] &&
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
    if (feedSelected.value) {
      const [type, guid] = feedSelected.value.split(":", 2);
      if (type === "feed") query.feed_id = guid;
      else if (type === "folder") query.folder_id = guid;
    }
    if (filterSelected.value) query.status = filterSelected.value;
    if (itemSearch.value) query.search = itemSearch.value;
    if (!itemSortNewestFirst.value) query.oldest_first = true;
    return query;
  }, []);

  const refreshItems = useCallback(
    async (loadMore = false) => {
      if (feedSelected.value === null) {
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
  }, [feedSelected.value, filterSelected.value, itemSortNewestFirst.value]);

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
      feedSelected.value = null;
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
      folders.value = [...folders.value.map(f => (f.id === folder.id ? { ...f, title } : f))].sort(
        (a, b) => a.title.localeCompare(b.title),
      );
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
      feedSelected.value = null;
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
    const curIdx = REFRESH_RATE_OPTIONS.findIndex(o => o.value === refreshRate.value);
    if ((curIdx <= 0 && offset < 0) || (curIdx >= REFRESH_RATE_OPTIONS.length - 1 && offset > 0))
      return;
    const newRate = REFRESH_RATE_OPTIONS[curIdx + offset].value;
    refreshRate.value = newRate;
    updateSettings({ refresh_rate: newRate });
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

  const navigateToFeed = useCallback(
    (dir: number) => {
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

      const curIdx = list.indexOf(feedSelected.value || "");
      if (curIdx === -1) {
        feedSelected.value = "";
        return;
      }
      const newIdx = curIdx + dir;
      if (newIdx < 0 || newIdx >= list.length) return;

      const nextVal = list[newIdx];
      feedSelected.value = nextVal;
      updateSettings({ feed: nextVal });

      setTimeout(() => {
        const scroll = document.querySelector("#feed-list-scroll");
        const handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);
      }, 0);
    },
    [updateSettings],
  );

  // Keybindings
  const keyActionsRef = useRef<KeyActions>({} as KeyActions);
  keyActionsRef.current = {
    openItemLink: () => {
      if (itemSelectedDetails.value?.link)
        window.open(itemSelectedDetails.value.link, "_blank", "noopener,noreferrer");
    },
    toggleReadability,
    toggleItemRead: () =>
      itemSelectedDetails.value && toggleItemStatus(itemSelectedDetails.value, "unread"),
    markAllRead: () => filterSelected.value === "unread" && markItemsRead(),
    toggleItemStarred: () =>
      itemSelectedDetails.value && toggleItemStatus(itemSelectedDetails.value, "starred"),
    nextItem: () => navigateToItem(1),
    previousItem: () => navigateToItem(-1),
    nextFeed: () => navigateToFeed(1),
    previousFeed: () => navigateToFeed(-1),
    closeItem: () => (itemSelected.value = null),
    showShortcuts: () => (showModal.value = "shortcuts"),
    showAll: () => {
      filterSelected.value = "";
      updateSettings({ filter: "" });
    },
    showUnread: () => {
      filterSelected.value = "unread";
      updateSettings({ filter: "unread" });
    },
    showStarred: () => {
      filterSelected.value = "starred";
      updateSettings({ filter: "starred" });
    },
  };

  useEffect(() => setupKeybindings(keyActionsRef), []);

  return (
    <div
      className={`d-flex ${feedSelected.value !== null ? "feed-selected" : ""} ${
        itemSelected.value !== null ? "item-selected" : ""
      } ${itemFullscreen.value ? "item-fullscreen" : ""}`}>
      {/* Feed list column */}
      <div
        id="col-feed-list"
        className="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
        style={{ width: `${feedListWidth.value}px` }}>
        <Drag width={feedListWidth.value} onResize={resizeFeedList} />
        <div className="px-2 py-1 d-flex align-items-center">
          <Icon className="mx-2" name="anchor" />
          <div className="flex-grow-1" />
          <button
            type="button"
            className="c-button-pill ms-1"
            aria-pressed={filterSelected.value === "unread"}
            title={t("unread")}
            onClick={() => {
              filterSelected.value = "unread";
              updateSettings({ filter: "unread" });
            }}>
            <Icon name="circle-full" />
          </button>
          <button
            type="button"
            className="c-button-pill mx-1"
            aria-pressed={filterSelected.value === "starred"}
            title={t("starred")}
            onClick={() => {
              filterSelected.value = "starred";
              updateSettings({ filter: "starred" });
            }}>
            <Icon name="star-full" />
          </button>
          <button
            type="button"
            className="c-button-pill me-1"
            aria-pressed={filterSelected.value === ""}
            title={t("all")}
            onClick={() => {
              filterSelected.value = "";
              updateSettings({ filter: "" });
            }}>
            <Icon name="assorted" />
          </button>
          <div className="flex-grow-1" />

          <DropdownSettings
            dropdownRef={menuDropdownRef}
            theme={theme.value.name}
            onThemeChange={th => {
              theme.value = { ...theme.value, name: th };
              updateSettings({ theme_name: th });
            }}
            refreshAvailable={refreshAvailable.value}
            refreshRate={refreshRate.value}
            refreshRateTitle={refreshRateTitle.value}
            maxRefreshRate={
              refreshRate.value === REFRESH_RATE_OPTIONS[REFRESH_RATE_OPTIONS.length - 1].value
            }
            onRefreshRateChange={changeRefreshRate}
            onFetchAllFeeds={fetchAllFeeds}
            onImportOPML={importOPML}
            onShowShortcuts={() => (showModal.value = "shortcuts")}
            onShowNewFeed={() => (showModal.value = "newfeed")}
            language={language.value}
            onLanguageChange={code => {
              setLang(code);
              language.value = code;
              updateSettings({ language: code });
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
        </div>

        <div id="feed-list-scroll" className="p-2 overflow-auto border-top flex-grow-1">
          <FeedTree
            tree={feedTree.value}
            value={feedSelected.value}
            filterSelected={filterSelected.value}
            stats={stats.value}
            feedErrors={feedErrors.value}
            onSelect={val => {
              feedSelected.value = val;
              updateSettings({ feed: val });
            }}
            onToggleFolder={toggleFolderExpanded}
          />
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
        style={{ width: `${itemListWidth.value}px` }}>
        <Drag width={itemListWidth.value} onResize={resizeItemList} />
        <div className="px-2 py-1 d-flex gap-1 align-items-center">
          <button
            type="button"
            className="c-button-pill d-md-none"
            onClick={() => (feedSelected.value = null)}
            title={t("show_feeds")}>
            <Icon name="chevron-left" />
          </button>
          <button
            type="button"
            className="c-button-pill"
            onClick={() => {
              const next = !itemSortNewestFirst.value;
              itemSortNewestFirst.value = next;
              updateSettings({ sort_newest_first: next });
            }}
            title={`${t("show_first")}: ${itemSortNewestFirst.value ? t("new") : t("old")}`}>
            <Icon name={itemSortNewestFirst.value ? "sort-new-first" : "sort-old-first"} />
          </button>

          <div className="c-search flex-grow-1">
            <Icon name="search" />
            <input
              id="searchbar"
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

          {filterSelected.value === "unread" && (
            <button
              type="button"
              className="c-button-pill"
              onClick={markItemsRead}
              title={t("mark_all_read")}>
              <Icon name="check" />
            </button>
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
            <button type="button" className="c-button-link c-button-pill px-2" disabled>
              <Icon name="more-horizontal" />
            </button>
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
                  <RelativeTime val={item.date} locale={language.value} />
                </small>
              </div>
              <div className="text-break line-clamp-3">{item.title || t("untitled")}</div>
            </div>
          ))}
          {(itemsHasMore.value || loadingItems.value) && (
            <div className="text-center my-3">
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
          theme={theme.value}
          onFontChange={font => {
            theme.value = { ...theme.value, font };
            updateSettings({ theme_font: font });
          }}
          onFontSizeChange={delta => {
            const size = +(theme.value.size + 0.1 * delta).toFixed(1);
            theme.value = { ...theme.value, size };
            updateSettings({ theme_size: size });
          }}
          readabilityContent={itemSelectedReadability.value}
          loadingReadability={loadingReadability.value}
          onToggleReadability={toggleReadability}
          onToggleStarred={item => toggleItemStatus(item, "starred")}
          onToggleRead={item => toggleItemStatus(item, "unread")}
          onSelectFeed={id => (feedSelected.value = `feed:${id}`)}
          onNavigate={navigateToItem}
          canPrev={Boolean(items.value.length && itemSelected.value !== items.value[0].id)}
          canNext={Boolean(
            items.value.length && itemSelected.value !== items.value[items.value.length - 1].id,
          )}
          isFullscreen={itemFullscreen.value}
          onToggleFullscreen={() => (itemFullscreen.value = !itemFullscreen.value)}
          onClose={() => (itemSelected.value = null)}
          contentRef={contentRef}
          language={language.value}
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
        {showModal.value === "newfeed" && (
          <NewFeed
            folders={folders.value}
            folderId={currentFeed.value?.folder_id ?? currentFolder.value?.id ?? null}
            onCreated={feed => {
              refreshFeeds();
              refreshStats();
              showModal.value = "";
              feedSelected.value = "feed:" + feed.id;
            }}
            onFolderCreated={refreshFeeds}
          />
        )}
        {showModal.value === "shortcuts" && <Shortcuts />}
      </Modal>

      <ToastContainer />
    </div>
  );
}
