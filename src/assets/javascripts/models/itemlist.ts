import { signal, effect, untracked, createModel } from "@preact/signals";
import api from "../api";
import { to } from "../utils";
import type { TranslationKey } from "../i18n";
import type { AppSettings } from "./settings";
import type { AppRegistry } from "./registry";
import type { Item, ItemListQuery, ItemStatus } from "../api-types";

export const ItemListModel = createModel(
  (
    settings: AppSettings,
    registry: AppRegistry,
    onError?: (err: Error, key: TranslationKey) => void,
  ) => {
    const items = signal<Item[]>([]);
    const hasMore = signal(true);
    const loading = signal(false);
    const search = signal("");
    const resetCount = signal(0);

    let activeReq = 0;

    function buildQuery(): ItemListQuery {
      const query: ItemListQuery = {};
      const selected = settings.values.feed.value;
      if (selected) {
        const [type, guid] = selected.split(":", 2);
        if (type === "feed") query.feed_id = guid;
        else if (type === "folder") query.folder_id = guid;
      }
      if (settings.values.filter.value) query.status = settings.values.filter.value;
      if (search.value) query.search = search.value;
      if (!settings.values.sort_newest_first.value) query.oldest_first = true;
      return query;
    }

    async function load(reset = true) {
      if (settings.values.feed.value === null) {
        items.value = [];
        hasMore.value = false;
        loading.value = false;
        return;
      }

      const reqId = ++activeReq;

      if (reset) {
        items.value = [];
        hasMore.value = true;
      }

      const query = buildQuery();
      if (!reset && items.value.length > 0) {
        query.after = items.value[items.value.length - 1].id;
      }

      loading.value = true;
      const [err, data] = await to(api.items.list(query));

      if (reqId !== activeReq) return;
      loading.value = false;

      if (err) {
        onError?.(err, "fail_load");
        return;
      }

      items.value = reset ? data.list : items.value.concat(data.list);
      hasMore.value = data.has_more;
      if (reset) resetCount.value++;
    }

    function loadMore() {
      if (!hasMore.value || loading.value) return;
      void load(false);
    }

    async function markAllRead(): Promise<boolean> {
      const query = buildQuery();
      const [err] = await to(
        api.items.mark_read({ folder_id: query.folder_id, feed_id: query.feed_id }),
      );
      if (err) {
        onError?.(err, "fail_update_article");
        return false;
      }
      items.value = [];
      hasMore.value = false;
      registry.refreshStats();
      return true;
    }

    function setStatus(id: number, status: ItemStatus) {
      items.value = items.value.map(i => (i.id === id ? { ...i, status } : i));
    }

    let lastFeed: string | null | undefined;
    let lastFilter: string | undefined;
    let lastNewestFirst: boolean | undefined;
    let lastSearch = "";
    let searchTimer: ReturnType<typeof setTimeout> | undefined;

    effect(() => {
      const feed = settings.values.feed.value;
      const filter = settings.values.filter.value;
      const newestFirst = settings.values.sort_newest_first.value;
      const searchTerm = search.value;

      const queryChanged =
        feed !== lastFeed || filter !== lastFilter || newestFirst !== lastNewestFirst;
      const searchChanged = searchTerm !== lastSearch;

      lastFeed = feed;
      lastFilter = filter;
      lastNewestFirst = newestFirst;
      lastSearch = searchTerm;

      clearTimeout(searchTimer);
      if (queryChanged) {
        untracked(() => void load(true));
      } else if (searchChanged) {
        searchTimer = setTimeout(() => untracked(() => void load(true)), 500);
      }
    });

    effect(() => () => clearTimeout(searchTimer));

    return { items, hasMore, loading, search, resetCount, loadMore, markAllRead, setStatus };
  },
);

export type AppItemList = InstanceType<typeof ItemListModel>;
