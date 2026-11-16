import { useEffect, useRef, useMemo, useCallback } from "preact/hooks";
import { useSignal, useComputed, useModel } from "@preact/signals";
import { useI18n } from "../i18n";
import api, { NetworkError, HTTPError } from "../api";
import { scrollto, debounce, to } from "../utils";
import Modal from "../components/Modal";
import Shortcuts from "../components/Shortcuts";
import Icon from "../components/Icon";
import ToastContainer, { showToast } from "../components/Toast";
import ColumnItem from "../components/ColumnItem";
import ColumnItemList from "../components/ColumnItemList";
import ColumnFeedList from "../components/ColumnFeedList";
import type { ColumnFeedListHandle } from "../components/ColumnFeedList";
import { setupKeybindings, KeyActions } from "../key";
import { SettingsModel } from "../models/settings";
import type { Theme } from "../models/settings";
import { RegistryModel } from "../models/registry";
import type { Item, FeedStat, ItemStatus, ItemListQuery } from "../api-types";

const THEME_COLORS: Record<string, string> = {
  night: "#0e0e0e",
  sepia: "#f4f0e5",
  light: "#fff",
};

const TITLE = document.title;

export default function App() {
  const { t, setLang } = useI18n();
  const appSettings = window.app?.settings || {};

  // Feeds & folders
  const feedErrors = useSignal<Record<number, string>>({});

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
  const refreshAvailable = useSignal(true);
  const loadingFeeds = useSignal(0);
  const loadingItems = useSignal(false);
  const loadingReadability = useSignal(false);

  // Element refs
  const itemListRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const feedListRef = useRef<ColumnFeedListHandle>(null);
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

  const registry = useModel(
    () =>
      new RegistryModel(err =>
        showToast({ title: t("fail_load"), description: errDescription(err) }, { level: "fail" }),
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

  useEffect(() => {
    setLang(settings.values.language.value);
  }, [settings.values.language.value, setLang]);

  // Computed signals
  const selectedFeedId = useComputed(() => itemSelectedDetails.value?.feed_id ?? null);

  // Keep document title in sync with unread count
  useEffect(() => {
    const unread = registry.stats.value.total.unread;
    document.title = TITLE + (unread ? ` (${unread})` : "");
  }, [registry.stats.value.total.unread]);

  // Refresh feeds & stats
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

      registry.feedStats.value = data.stats.reduce(
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
    registry.refresh();
  }, [refreshStats, registry]);

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
          const cur = registry.feedStats.value[item.feed_id];
          if (cur) {
            registry.feedStats.value = {
              ...registry.feedStats.value,
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

    const cur = registry.feedStats.value[item.feed_id];
    if (cur) {
      let unread = cur.unread || 0;
      let starred = cur.starred || 0;
      if (oldstatus === "unread") unread -= 1;
      if (oldstatus === "starred") starred -= 1;
      if (newstatus === "unread") unread += 1;
      if (newstatus === "starred") starred += 1;
      registry.feedStats.value = {
        ...registry.feedStats.value,
        [item.feed_id]: { ...cur, unread, starred },
      };
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
    nextFeed: () => feedListRef.current?.navigate(+1),
    prevFeed: () => feedListRef.current?.navigate(-1),
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
      <ColumnFeedList
        settings={settings}
        registry={registry}
        feedErrors={feedErrors.value}
        loadingFeeds={loadingFeeds.value}
        refreshAvailable={refreshAvailable.value}
        selectedFeedId={selectedFeedId}
        requiresAuth={Boolean(window.app?.requiresAuth)}
        handleRef={feedListRef}
        refreshStats={refreshStats}
        describeError={errDescription}
        onShowShortcuts={() => (showModal.value = "shortcuts")}
      />

      {/* Item list column */}
      <ColumnItemList
        settings={settings}
        registry={registry}
        feedErrors={feedErrors.value}
        items={items}
        itemsHasMore={itemsHasMore}
        loadingItems={loadingItems.value}
        itemSelected={itemSelected}
        itemSearch={itemSearch}
        itemListRef={itemListRef}
        onScroll={onItemListScroll}
        onMarkAllRead={markItemsRead}
        refreshStats={refreshStats}
        describeError={errDescription}
      />

      {/* Item reader column */}
      <ColumnItem
        item={itemSelectedDetails.value}
        feedTitle={
          (registry.feedsById.value[itemSelectedDetails.value?.feed_id ?? -1] || {}).title || ""
        }
        feedId={(registry.feedsById.value[itemSelectedDetails.value?.feed_id ?? -1] || {}).id}
        settings={settings}
        readabilityContent={itemSelectedReadability.value}
        loadingReadability={loadingReadability.value}
        onToggleReadability={toggleReadability}
        onToggleStarred={item => toggleItemStatus(item, "starred")}
        onToggleRead={item => toggleItemStatus(item, "unread")}
        onNavigate={navigateToItem}
        canPrev={Boolean(items.value.length && itemSelected.value !== items.value[0].id)}
        canNext={Boolean(
          items.value.length && itemSelected.value !== items.value[items.value.length - 1].id,
        )}
        isFullscreen={itemFullscreen.value}
        onToggleFullscreen={() => (itemFullscreen.value = !itemFullscreen.value)}
        onClose={() => (itemSelected.value = null)}
        contentRef={contentRef}
      />

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
