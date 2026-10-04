import { useState, useEffect, useRef, useMemo, useCallback } from "preact/hooks";
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
  const appSettings = window.app.settings;

  // Settings & state
  const [filterSelected, setFilterSelected] = useState<Filter>(
    (appSettings.filter as Filter) || "",
  );
  const [theme, setTheme] = useState<{ name: Theme; font: ThemeFont; size: number }>({
    name: (appSettings.theme_name as Theme) || "system",
    font: (appSettings.theme_font as ThemeFont) || "",
    size: (appSettings.theme_size as number) || 1,
  });
  const [refreshRate, setRefreshRate] = useState<number>(appSettings.refresh_rate || 0);
  const [language, setLanguage] = useState<Lang>(appSettings.language || "en");

  // Feeds & folders
  const [folders, setFolders] = useState<Folder[]>([]);
  const [feeds, setFeeds] = useState<Feed[]>([]);
  const [feedSelected, setFeedSelected] = useState<string | null>(appSettings.feed ?? null);
  const [feedListWidth, setFeedListWidth] = useState<number>(appSettings.feed_list_width || 300);
  const [feedErrors, setFeedErrors] = useState<Record<number, string>>({});
  const [feedStats, setFeedStats] = useState<Record<number, FeedStat>>({});
  const [stats, setStats] = useState<{
    folders: Record<number, Stats>;
    feeds: Record<number, Stats>;
    total: Stats;
  }>({ folders: {}, feeds: {}, total: { unread: 0, starred: 0 } });

  // Items
  const [itemListWidth, setItemListWidth] = useState<number>(appSettings.item_list_width || 300);
  const [items, setItems] = useState<Item[]>([]);
  const [itemsHasMore, setItemsHasMore] = useState(true);
  const [itemSelected, setItemSelected] = useState<number | null>(null);
  const [itemSelectedDetails, setItemSelectedDetails] = useState<Item | null>(null);
  const [itemSelectedReadability, setItemSelectedReadability] = useState("");
  const [itemFullscreen, setItemFullscreen] = useState(false);
  const [itemSearch, setItemSearch] = useState("");
  const [itemSortNewestFirst, setItemSortNewestFirst] = useState<boolean>(
    appSettings.sort_newest_first ?? true,
  );

  // Status & modal
  const [showModal, setShowModal] = useState<"" | "shortcuts" | "newfeed">("");
  const [refreshAvailable, setRefreshAvailable] = useState(true);
  const [loadingFeeds, setLoadingFeeds] = useState(0);
  const [loadingItems, setLoadingItems] = useState(false);
  const [loadingReadability, setLoadingReadability] = useState(false);

  // Element refs
  const itemListRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const menuDropdownRef = useRef<DropdownHandle>(null);

  // Error description helper
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
    updateMetaTheme(theme.name);
    const mql = window.matchMedia?.("(prefers-color-scheme: dark)");
    const handler = () => updateMetaTheme(theme.name);
    mql?.addEventListener("change", handler);
    return () => mql?.removeEventListener("change", handler);
  }, [theme.name, updateMetaTheme]);

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
    setFeedListWidth(clamped);
    saveFeedListWidth(clamped);
  };

  const resizeItemList = (w: number) => {
    const clamped = Math.min(Math.max(200, w), 700);
    setItemListWidth(clamped);
    saveItemListWidth(clamped);
  };

  // Plain calculations
  const feedsById = useMemo(() => {
    return feeds.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Feed>);
  }, [feeds]);

  const foldersById = useMemo(() => {
    return folders.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Folder>);
  }, [folders]);

  const [currentType, currentGuidStr] = (feedSelected || "").split(":", 2);
  const currentFeed = currentType === "feed" ? feedsById[Number(currentGuidStr)] || null : null;
  const currentFolder =
    currentType === "folder" ? foldersById[Number(currentGuidStr)] || null : null;

  const searchScope =
    currentFeed?.title ||
    currentFolder?.title ||
    (filterSelected === "unread"
      ? t("all_unread")
      : filterSelected === "starred"
        ? t("all_starred")
        : t("all_feeds"));

  const refreshRateTitle = REFRESH_RATE_OPTIONS.find(o => o.value === refreshRate)?.title || "0";

  // Compute stats
  useEffect(() => {
    const statsFeeds: Record<number, Stats> = {};
    const statsFolders: Record<number, Stats> = {};
    const statsTotal: Stats = { unread: 0, starred: 0 };

    for (const feed of feeds) {
      const n = feedStats[feed.id];
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

    setStats({ feeds: statsFeeds, folders: statsFolders, total: statsTotal });
    const unread = statsTotal.unread;
    document.title = TITLE + (unread ? ` (${unread})` : "");
  }, [feeds, feedStats]);

  // Refresh feeds & stats
  const refreshFeeds = useCallback(async () => {
    const [err, values] = await to(Promise.all([api.folders.list(), api.feeds.list()]));
    if (err) {
      showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" });
      return;
    }
    setFolders(values[0]);
    setFeeds(values[1]);
  }, [t, errDescription]);

  const refreshStats = useCallback(
    async (loopMode?: boolean) => {
      const [err, data] = await to(api.status());
      if (err) {
        showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" });
        return;
      }

      setLoadingFeeds(data.running);
      setRefreshAvailable(data.refresh);
      if (data.running) setTimeout(() => refreshStats(true), 500);

      setFeedStats(
        data.stats.reduce(
          (acc, stat) => ({ ...acc, [stat.feed_id]: stat }),
          {} as Record<number, FeedStat>,
        ),
      );

      const [feedErr, errors] = await to(api.feeds.list_errors());
      if (feedErr) {
        showToast(
          { title: t("fail_load"), description: errDescription(feedErr) },
          { level: "fail" },
        );
        return;
      }
      setFeedErrors(errors);
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
      filterSelected &&
      currentFolder?.id !== folder.id &&
      currentFeed?.folder_id !== folder.id &&
      !stats.folders[folder.id]?.[filterSelected] &&
      feedsById[itemSelectedDetails?.feed_id ?? -1]?.folder_id !== folder.id,
    );

  const mustHideFeed = (feed: Feed) =>
    Boolean(
      filterSelected &&
      currentFeed?.id !== feed.id &&
      !stats.feeds[feed.id]?.[filterSelected] &&
      itemSelectedDetails?.feed_id !== feed.id,
    );

  const feedTree = useMemo((): FeedTreeNode[] => {
    const [rootFeeds, folderFeeds] = feeds.reduce(
      (acc, f) => (acc[f.folder_id === null ? 0 : 1].push(f), acc),
      [[] as Feed[], [] as Feed[]],
    );
    const byFolder: Record<number, Feed[]> = folderFeeds.reduce(
      (acc, f) => ((acc[f.folder_id as number] ||= []).push(f), acc),
      {} as Record<number, Feed[]>,
    );
    const feedNode = (feed: Feed): TreeFeedNode => ({ type: "feed", feed });

    return [
      ...folders
        .filter(folder => !mustHideFolder(folder))
        .map(folder => ({
          type: "folder" as const,
          folder,
          feeds: (byFolder[folder.id] || []).filter(f => !mustHideFeed(f)).map(feedNode),
        })),
      ...rootFeeds.filter(f => !mustHideFeed(f)).map(feedNode),
    ];
  }, [
    feeds,
    folders,
    filterSelected,
    currentFeed,
    currentFolder,
    stats,
    itemSelectedDetails,
    feedsById,
  ]);

  // Items loading & query
  const getItemsQuery = useCallback((): ItemListQuery => {
    const query: ItemListQuery = {};
    if (feedSelected) {
      const [type, guid] = feedSelected.split(":", 2);
      if (type === "feed") query.feed_id = guid;
      else if (type === "folder") query.folder_id = guid;
    }
    if (filterSelected) query.status = filterSelected;
    if (itemSearch) query.search = itemSearch;
    if (!itemSortNewestFirst) query.oldest_first = true;
    return query;
  }, [feedSelected, filterSelected, itemSearch, itemSortNewestFirst]);

  const itemsRef = useRef<Item[]>([]);
  itemsRef.current = items;
  const activeReqIdRef = useRef(0);

  const refreshItems = useCallback(
    async (loadMore = false) => {
      if (feedSelected === null) {
        setItems([]);
        setItemsHasMore(false);
        setLoadingItems(false);
        return;
      }

      const reqId = ++activeReqIdRef.current;

      if (!loadMore) {
        setItems([]);
        setItemsHasMore(true);
        if (itemListRef.current) itemListRef.current.scrollTop = 0;
      }

      const query = getItemsQuery();
      if (loadMore && itemsRef.current.length > 0) {
        query.after = itemsRef.current[itemsRef.current.length - 1].id;
      }

      setLoadingItems(true);
      const [err, data] = await to(api.items.list(query));

      if (reqId !== activeReqIdRef.current) return;
      setLoadingItems(false);

      if (err) {
        showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" });
        return;
      }

      setItems(prev => (loadMore ? prev.concat(data.list) : data.list));
      setItemsHasMore(data.has_more);

      if (data.has_more && !loadMore) {
        setTimeout(() => {
          const el = itemListRef.current;
          if (!el || el.scrollHeight === 0) return;
          const scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
          if (el.scrollHeight - el.scrollTop - el.offsetHeight < 70 * scale) {
            refreshItems(true);
          }
        }, 0);
      }
    },
    [feedSelected, getItemsQuery, t, errDescription],
  );

  useEffect(() => {
    setItemSelected(null);
    refreshItems(false);
  }, [feedSelected, filterSelected, itemSortNewestFirst]);

  useEffect(() => {
    const timer = setTimeout(() => refreshItems(false), 500);
    return () => clearTimeout(timer);
  }, [itemSearch]);

  const loadMoreItems = useCallback(() => {
    if (!itemsHasMore || loadingItems) return;
    const el = itemListRef.current;
    if (!el || el.scrollHeight === 0) return;
    const scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
    const closeToBottom = el.scrollHeight - el.scrollTop - el.offsetHeight < 70 * scale;
    if (
      closeToBottom ||
      (itemSelected && itemsRef.current.length && itemSelected === itemsRef.current[itemsRef.current.length - 1].id)
    ) {
      refreshItems(true);
    }
  }, [itemsHasMore, loadingItems, refreshItems, itemSelected]);

  const onItemListScroll = useMemo(() => debounce(loadMoreItems, 200), [loadMoreItems]);

  // Item details & auto-read
  useEffect(() => {
    setItemSelectedReadability("");
    if (itemSelected === null) {
      setItemSelectedDetails(null);
      setItemFullscreen(false);
      return;
    }

    if (contentRef.current) contentRef.current.scrollTop = 0;

    to(api.items.get(itemSelected)).then(([itemErr, item]) => {
      if (itemErr) {
        showToast(
          { title: t("fail_load"), description: errDescription(itemErr) },
          { level: "fail" },
        );
        return;
      }
      setItemSelectedDetails(item);

      if (item.status === "unread") {
        to(api.items.update(item.id, { status: "read" })).then(([updateErr]) => {
          if (updateErr) {
            showToast(
              { title: t("fail_update_article"), description: errDescription(updateErr) },
              { level: "fail" },
            );
            return;
          }
          setFeedStats(prev => {
            const cur = prev[item.feed_id];
            if (!cur) return prev;
            return {
              ...prev,
              [item.feed_id]: { ...cur, unread: Math.max(0, (cur.unread || 1) - 1) },
            };
          });
          setItems(prev => prev.map(i => (i.id === item.id ? { ...i, status: "read" } : i)));
          item.status = "read";
        });
      }
    });
  }, [itemSelected, t, errDescription]);

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

    setFeedStats(prev => {
      const cur = prev[item.feed_id];
      if (!cur) return prev;
      let unread = cur.unread || 0;
      let starred = cur.starred || 0;
      if (oldstatus === "unread") unread -= 1;
      if (oldstatus === "starred") starred -= 1;
      if (newstatus === "unread") unread += 1;
      if (newstatus === "starred") starred += 1;
      return { ...prev, [item.feed_id]: { ...cur, unread, starred } };
    });

    setItems(prev => prev.map(i => (i.id === item.id ? { ...i, status: newstatus } : i)));
    item.status = newstatus;
    setItemSelectedDetails(prev => (prev?.id === item.id ? { ...prev, status: newstatus } : prev));
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
    setItems([]);
    setItemSelected(null);
    setItemsHasMore(false);
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
    else setFeeds(prev => prev.map(f => (f.id === feed.id ? { ...f, title } : f)));
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
    else setFeeds(prev => prev.map(f => (f.id === feed.id ? { ...f, feed_link } : f)));
  };

  const moveFeed = async (feed: Feed, folder_id: number | null) => {
    const [err] = await to(api.feeds.update(feed.id, { folder_id }));
    if (err)
      showToast(
        { title: t("fail_save_feed"), description: errDescription(err) },
        { level: "fail" },
      );
    else {
      setFeeds(prev => prev.map(f => (f.id === feed.id ? { ...f, folder_id } : f)));
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
      setFeedSelected(null);
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
      setFolders(prev =>
        [...prev.map(f => (f.id === folder.id ? { ...f, title } : f))].sort((a, b) =>
          a.title.localeCompare(b.title),
        ),
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
      setFeedSelected(null);
      refreshStats();
      refreshFeeds();
    }
  };

  const toggleFolderExpanded = async (folder: Folder) => {
    const is_expanded = !folder.is_expanded;
    setFolders(prev => prev.map(f => (f.id === folder.id ? { ...f, is_expanded } : f)));
    const [err] = await to(api.folders.update(folder.id, { is_expanded }));
    if (err)
      showToast(
        { title: t("fail_save_folder"), description: errDescription(err) },
        { level: "fail" },
      );
  };

  const toggleReadability = async () => {
    if (itemSelectedReadability) {
      setItemSelectedReadability("");
      return;
    }
    if (!itemSelectedDetails?.link) return;
    setLoadingReadability(true);
    const [err, data] = await to(api.crawl(itemSelectedDetails.link));
    setLoadingReadability(false);
    if (err)
      showToast(
        { title: t("fail_readability"), description: errDescription(err) },
        { level: "fail" },
      );
    else setItemSelectedReadability(data?.content || "");
  };

  const fetchAllFeeds = async () => {
    if (loadingFeeds) return;
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
    const curIdx = REFRESH_RATE_OPTIONS.findIndex(o => o.value === refreshRate);
    if ((curIdx <= 0 && offset < 0) || (curIdx >= REFRESH_RATE_OPTIONS.length - 1 && offset > 0))
      return;
    const newRate = REFRESH_RATE_OPTIONS[curIdx + offset].value;
    setRefreshRate(newRate);
    updateSettings({ refresh_rate: newRate });
  };

  // Keyboard navigation
  const navigateToItem = useCallback(
    (dir: number) => {
      if (itemSelected === null) {
        if (items.length) setItemSelected(items[0].id);
        return;
      }
      const curIdx = items.findIndex(x => x.id === itemSelected);
      if (curIdx === -1) {
        if (items.length) setItemSelected(items[0].id);
        return;
      }
      const newIdx = curIdx + dir;
      if (newIdx < 0 || newIdx >= items.length) return;
      setItemSelected(items[newIdx].id);

      setTimeout(() => {
        const scroll = document.querySelector("#item-list-scroll");
        const handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);
        loadMoreItems();
      }, 0);
    },
    [itemSelected, items, loadMoreItems],
  );

  const navigateToFeed = useCallback(
    (dir: number) => {
      const list: string[] = [""];
      for (const node of feedTree) {
        if (node.type === "folder") {
          list.push("folder:" + node.folder.id);
          if (node.folder.is_expanded) {
            for (const f of node.feeds) list.push("feed:" + f.feed.id);
          }
        } else {
          list.push("feed:" + node.feed.id);
        }
      }

      const curIdx = list.indexOf(feedSelected || "");
      if (curIdx === -1) {
        setFeedSelected("");
        return;
      }
      const newIdx = curIdx + dir;
      if (newIdx < 0 || newIdx >= list.length) return;

      const nextVal = list[newIdx];
      setFeedSelected(nextVal);
      updateSettings({ feed: nextVal });

      setTimeout(() => {
        const scroll = document.querySelector("#feed-list-scroll");
        const handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);
      }, 0);
    },
    [feedTree, feedSelected, updateSettings],
  );

  // Keybindings
  const keyActionsRef = useRef<KeyActions>({} as KeyActions);
  keyActionsRef.current = {
    openItemLink: () => {
      if (itemSelectedDetails?.link)
        window.open(itemSelectedDetails.link, "_blank", "noopener,noreferrer");
    },
    toggleReadability,
    toggleItemRead: () => itemSelectedDetails && toggleItemStatus(itemSelectedDetails, "unread"),
    markAllRead: () => filterSelected === "unread" && markItemsRead(),
    toggleItemStarred: () =>
      itemSelectedDetails && toggleItemStatus(itemSelectedDetails, "starred"),
    nextItem: () => navigateToItem(1),
    previousItem: () => navigateToItem(-1),
    nextFeed: () => navigateToFeed(1),
    previousFeed: () => navigateToFeed(-1),
    closeItem: () => setItemSelected(null),
    showShortcuts: () => setShowModal("shortcuts"),
    showAll: () => {
      setFilterSelected("");
      updateSettings({ filter: "" });
    },
    showUnread: () => {
      setFilterSelected("unread");
      updateSettings({ filter: "unread" });
    },
    showStarred: () => {
      setFilterSelected("starred");
      updateSettings({ filter: "starred" });
    },
  };

  useEffect(() => setupKeybindings(keyActionsRef), []);

  return (
    <div
      className={`d-flex ${feedSelected !== null ? "feed-selected" : ""} ${
        itemSelected !== null ? "item-selected" : ""
      } ${itemFullscreen ? "item-fullscreen" : ""}`}>
      {/* Feed list column */}
      <div
        id="col-feed-list"
        className="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
        style={{ width: `${feedListWidth}px` }}>
        <Drag width={feedListWidth} onResize={resizeFeedList} />
        <div className="px-2 py-1 d-flex align-items-center">
          <Icon className="mx-2" name="anchor" />
          <div className="flex-grow-1" />
          <button
            type="button"
            className="c-button-pill ms-1"
            aria-pressed={filterSelected === "unread"}
            title={t("unread")}
            onClick={() => {
              setFilterSelected("unread");
              updateSettings({ filter: "unread" });
            }}>
            <Icon name="circle-full" />
          </button>
          <button
            type="button"
            className="c-button-pill mx-1"
            aria-pressed={filterSelected === "starred"}
            title={t("starred")}
            onClick={() => {
              setFilterSelected("starred");
              updateSettings({ filter: "starred" });
            }}>
            <Icon name="star-full" />
          </button>
          <button
            type="button"
            className="c-button-pill me-1"
            aria-pressed={filterSelected === ""}
            title={t("all")}
            onClick={() => {
              setFilterSelected("");
              updateSettings({ filter: "" });
            }}>
            <Icon name="assorted" />
          </button>
          <div className="flex-grow-1" />

          <DropdownSettings
            dropdownRef={menuDropdownRef}
            theme={theme.name}
            onThemeChange={th => {
              setTheme(prev => ({ ...prev, name: th }));
              updateSettings({ theme_name: th });
            }}
            refreshAvailable={refreshAvailable}
            refreshRate={refreshRate}
            refreshRateTitle={refreshRateTitle}
            maxRefreshRate={
              refreshRate === REFRESH_RATE_OPTIONS[REFRESH_RATE_OPTIONS.length - 1].value
            }
            onRefreshRateChange={changeRefreshRate}
            onFetchAllFeeds={fetchAllFeeds}
            onImportOPML={importOPML}
            onShowShortcuts={() => setShowModal("shortcuts")}
            onShowNewFeed={() => setShowModal("newfeed")}
            language={language}
            onLanguageChange={code => {
              setLang(code);
              setLanguage(code);
              updateSettings({ language: code });
            }}
            requiresAuth={window.app.requiresAuth}
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
            tree={feedTree}
            value={feedSelected}
            filterSelected={filterSelected}
            stats={stats}
            feedErrors={feedErrors}
            onSelect={val => {
              setFeedSelected(val);
              updateSettings({ feed: val });
            }}
            onToggleFolder={toggleFolderExpanded}
          />
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

      {/* Item list column */}
      <div
        id="col-item-list"
        className="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
        style={{ width: `${itemListWidth}px` }}>
        <Drag width={itemListWidth} onResize={resizeItemList} />
        <div className="px-2 py-1 d-flex gap-1 align-items-center">
          <button
            type="button"
            className="c-button-pill d-md-none"
            onClick={() => setFeedSelected(null)}
            title={t("show_feeds")}>
            <Icon name="chevron-left" />
          </button>
          <button
            type="button"
            className="c-button-pill"
            onClick={() => {
              const next = !itemSortNewestFirst;
              setItemSortNewestFirst(next);
              updateSettings({ sort_newest_first: next });
            }}
            title={`${t("show_first")}: ${itemSortNewestFirst ? t("new") : t("old")}`}>
            <Icon name={itemSortNewestFirst ? "sort-new-first" : "sort-old-first"} />
          </button>

          <div className="c-search flex-grow-1">
            <Icon name="search" />
            <input
              id="searchbar"
              type="search"
              className="d-block"
              value={itemSearch}
              placeholder={t("search_placeholder", { scope: searchScope })}
              onInput={e => setItemSearch((e.target as HTMLInputElement).value)}
              onKeyDown={e => e.key === "Enter" && (e.target as HTMLInputElement).blur()}
            />
          </div>

          {filterSelected === "unread" && (
            <button
              type="button"
              className="c-button-pill"
              onClick={markItemsRead}
              title={t("mark_all_read")}>
              <Icon name="check" />
            </button>
          )}

          {currentFeed ? (
            <DropdownFeed
              feed={currentFeed}
              folders={folders}
              onRename={renameFeed}
              onUpdateLink={updateFeedLink}
              onMove={moveFeed}
              onMoveToNewFolder={moveFeedToNewFolder}
              onDelete={deleteFeed}
            />
          ) : currentFolder ? (
            <DropdownFolder
              folder={currentFolder}
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
          {items.map(item => (
            <div
              key={item.id}
              className="c-listitem d-flex flex-column user-select-none"
              role="radio"
              aria-checked={itemSelected === item.id}
              onClick={() => setItemSelected(item.id)}>
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
                  {(feedsById[item.feed_id] || {}).title}
                </small>
                <small className="flex-shrink-0">
                  <RelativeTime val={item.date} locale={language} />
                </small>
              </div>
              <div className="text-break line-clamp-3">{item.title || t("untitled")}</div>
            </div>
          ))}
          {(itemsHasMore || loadingItems) && (
            <div className="text-center my-3">
              <span className="c-spinner" />
            </div>
          )}
        </div>

        {currentFeed && feedErrors[currentFeed.id] && (
          <div className="px-3 py-2 border-top text-danger text-break">
            {feedErrors[currentFeed.id]}
          </div>
        )}
      </div>

      {/* Item reader column */}
      <div id="col-item" className="vh-100 d-flex flex-column w-100" style={{ minWidth: 0 }}>
        <ArticleView
          item={itemSelectedDetails}
          feedTitle={(feedsById[itemSelectedDetails?.feed_id ?? -1] || {}).title || ""}
          feedId={(feedsById[itemSelectedDetails?.feed_id ?? -1] || {}).id}
          theme={theme}
          onFontChange={font => {
            setTheme(prev => ({ ...prev, font }));
            updateSettings({ theme_font: font });
          }}
          onFontSizeChange={delta => {
            const size = +(theme.size + 0.1 * delta).toFixed(1);
            setTheme(prev => ({ ...prev, size }));
            updateSettings({ theme_size: size });
          }}
          readabilityContent={itemSelectedReadability}
          loadingReadability={loadingReadability}
          onToggleReadability={toggleReadability}
          onToggleStarred={item => toggleItemStatus(item, "starred")}
          onToggleRead={item => toggleItemStatus(item, "unread")}
          onSelectFeed={id => setFeedSelected(`feed:${id}`)}
          onNavigate={navigateToItem}
          canPrev={Boolean(items.length && itemSelected !== items[0].id)}
          canNext={Boolean(items.length && itemSelected !== items[items.length - 1].id)}
          isFullscreen={itemFullscreen}
          onToggleFullscreen={() => setItemFullscreen(prev => !prev)}
          onClose={() => setItemSelected(null)}
          contentRef={contentRef}
          language={language}
        />
      </div>

      {/* Dialog modal */}
      <Modal open={showModal !== ""} onHide={() => setShowModal("")}>
        <button
          type="button"
          className="c-button-link outline-none position-absolute top-0 end-0 p-2 m-2"
          style={{ lineHeight: 1 }}
          onClick={() => setShowModal("")}>
          <Icon name="x" />
        </button>
        {showModal === "newfeed" && (
          <NewFeed
            folders={folders}
            folderId={currentFeed?.folder_id ?? currentFolder?.id ?? null}
            onCreated={feed => {
              refreshFeeds();
              refreshStats();
              setShowModal("");
              setFeedSelected("feed:" + feed.id);
            }}
            onFolderCreated={refreshFeeds}
          />
        )}
        {showModal === "shortcuts" && <Shortcuts />}
      </Modal>

      <ToastContainer />
    </div>
  );
}
