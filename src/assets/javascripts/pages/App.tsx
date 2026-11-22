import { useEffect, useRef, useMemo, useCallback } from "preact/hooks";
import { useSignal, useComputed, useModel } from "@preact/signals";
import { useI18n } from "../i18n";
import type { TranslationKey } from "../i18n";
import api from "../api";
import { scrollto, debounce, to } from "../utils";
import Modal from "../components/Modal";
import Shortcuts from "../components/Shortcuts";
import Icon from "../components/Icon";
import ToastContainer, { showToast } from "../components/Toast";
import ColumnItemView from "../components/ColumnItemView";
import type { ColumnItemViewHandle } from "../components/ColumnItemView";
import ColumnItemList from "../components/ColumnItemList";
import ColumnFeedList from "../components/ColumnFeedList";
import type { ColumnFeedListHandle } from "../components/ColumnFeedList";
import { setupKeybindings, KeyActions } from "../key";
import { SettingsModel } from "../models/settings";
import type { Theme } from "../models/settings";
import { RegistryModel } from "../models/registry";
import { ItemListModel } from "../models/itemlist";
import type { Item, ItemStatus } from "../api-types";

const THEME_COLORS: Record<string, string> = {
  night: "#0e0e0e",
  sepia: "#f4f0e5",
  light: "#fff",
};

const TITLE = document.title;

export default function App() {
  const { setLang } = useI18n();
  const appSettings = window.app?.settings || {};

  // Items
  const itemSelected = useSignal<number | null>(null);
  const itemSelectedDetails = useSignal<Item | null>(null);

  // Status & modal
  const showModal = useSignal<"" | "shortcuts">("");

  // Element refs
  const itemListRef = useRef<HTMLDivElement>(null);
  const feedListRef = useRef<ColumnFeedListHandle>(null);
  const itemViewRef = useRef<ColumnItemViewHandle>(null);

  const showErrorToast = (err: Error, key: TranslationKey) =>
    showToast({ title: key, description: err }, { level: "fail" });

  const settings = useModel(() => new SettingsModel(appSettings, showErrorToast));
  const registry = useModel(() => new RegistryModel(showErrorToast));
  const itemList = useModel(() => new ItemListModel(settings, registry, showErrorToast));

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
  useEffect(() => {
    registry.refreshStats();
    registry.refresh();
  }, [registry]);

  // Reset selection when the item query changes
  useEffect(() => {
    itemSelected.value = null;
    if (itemListRef.current) itemListRef.current.scrollTop = 0;
  }, [
    settings.values.feed.value,
    settings.values.filter.value,
    settings.values.sort_newest_first.value,
  ]);

  const loadMoreItems = useCallback(() => {
    if (!itemList.hasMore.value || itemList.loading.value) return;
    const el = itemListRef.current;
    if (!el || el.scrollHeight === 0) return;
    const scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
    const closeToBottom = el.scrollHeight - el.scrollTop - el.offsetHeight < 70 * scale;
    const lastItem = itemList.items.value[itemList.items.value.length - 1];
    if (closeToBottom || (itemSelected.value !== null && lastItem?.id === itemSelected.value)) {
      itemList.loadMore();
    }
  }, [itemList]);

  // Fill the viewport after the first page loads
  useEffect(() => {
    if (itemList.resetCount.value === 0 || !itemList.hasMore.value) return;
    const timer = setTimeout(() => {
      const el = itemListRef.current;
      if (!el || el.scrollHeight === 0) return;
      const scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
      if (el.scrollHeight - el.scrollTop - el.offsetHeight < 70 * scale) {
        itemList.loadMore();
      }
    }, 0);
    return () => clearTimeout(timer);
  }, [itemList.resetCount.value]);

  const onItemListScroll = useMemo(() => debounce(loadMoreItems, 200), [loadMoreItems]);

  // Item details & auto-read
  useEffect(() => {
    if (itemSelected.value === null) {
      itemSelectedDetails.value = null;
      return;
    }

    to(api.items.get(itemSelected.value)).then(([itemErr, item]) => {
      if (itemErr) {
        showToast({ title: "fail_load", description: itemErr }, { level: "fail" });
        return;
      }
      itemSelectedDetails.value = item;

      if (item.status === "unread") {
        to(api.items.update(item.id, { status: "read" })).then(([updateErr]) => {
          if (updateErr) {
            showToast({ title: "fail_update_article", description: updateErr }, { level: "fail" });
            return;
          }
          registry.adjustStats(item.feed_id, "unread", "read");
          itemList.setStatus(item.id, "read");
          item.status = "read";
        });
      }
    });
  }, [itemSelected.value]);

  // Status toggle
  const toggleItemStatus = async (item: Item, targetstatus: ItemStatus) => {
    const oldstatus = item.status;
    const newstatus = item.status !== targetstatus ? targetstatus : "read";

    const [err] = await to(api.items.update(item.id, { status: newstatus }));
    if (err) {
      showToast({ title: "fail_update_article", description: err }, { level: "fail" });
      return;
    }

    registry.adjustStats(item.feed_id, oldstatus, newstatus);

    itemList.setStatus(item.id, newstatus);
    item.status = newstatus;
    if (itemSelectedDetails.value?.id === item.id) {
      itemSelectedDetails.value = { ...itemSelectedDetails.value, status: newstatus };
    }
  };

  const markItemsRead = async () => {
    const ok = await itemList.markAllRead();
    if (ok) itemSelected.value = null;
  };

  // Keyboard navigation
  const navigateToItem = useCallback(
    (dir: number) => {
      const list = itemList.items.value;
      if (itemSelected.value === null) {
        if (list.length) itemSelected.value = list[0].id;
        return;
      }
      const curIdx = list.findIndex(x => x.id === itemSelected.value);
      if (curIdx === -1) {
        if (list.length) itemSelected.value = list[0].id;
        return;
      }
      const newIdx = curIdx + dir;
      if (newIdx < 0 || newIdx >= list.length) return;
      itemSelected.value = list[newIdx].id;

      setTimeout(() => {
        const scroll = document.querySelector("#item-list-scroll");
        const handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);
        loadMoreItems();
      }, 0);
    },
    [itemList, loadMoreItems],
  );

  // Keybindings
  const keyActionsRef = useRef<KeyActions>({} as KeyActions);
  keyActionsRef.current = {
    enabled: () => showModal.value === "",
    openItemLink: () => {
      if (itemSelectedDetails.value?.link)
        window.open(itemSelectedDetails.value.link, "_blank", "noopener,noreferrer");
    },
    toggleReadability: () => itemViewRef.current?.toggleReadability(),
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
      }`}>
      {/* Feed list column */}
      <ColumnFeedList
        settings={settings}
        registry={registry}
        selectedFeedId={selectedFeedId}
        requiresAuth={Boolean(window.app?.requiresAuth)}
        handleRef={feedListRef}
        onShowShortcuts={() => (showModal.value = "shortcuts")}
      />

      {/* Item list column */}
      <ColumnItemList
        settings={settings}
        registry={registry}
        itemList={itemList}
        itemSelected={itemSelected}
        itemListRef={itemListRef}
        onScroll={onItemListScroll}
        onMarkAllRead={markItemsRead}
      />

      {/* Item reader column */}
      <ColumnItemView
        item={itemSelectedDetails.value}
        feedTitle={
          (registry.feedsById.value[itemSelectedDetails.value?.feed_id ?? -1] || {}).title || ""
        }
        feedId={(registry.feedsById.value[itemSelectedDetails.value?.feed_id ?? -1] || {}).id}
        settings={settings}
        handleRef={itemViewRef}
        onToggleStarred={item => toggleItemStatus(item, "starred")}
        onToggleRead={item => toggleItemStatus(item, "unread")}
        onNavigate={navigateToItem}
        canPrev={Boolean(
          itemList.items.value.length && itemSelected.value !== itemList.items.value[0].id,
        )}
        canNext={Boolean(
          itemList.items.value.length &&
          itemSelected.value !== itemList.items.value[itemList.items.value.length - 1].id,
        )}
        onClose={() => (itemSelected.value = null)}
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
