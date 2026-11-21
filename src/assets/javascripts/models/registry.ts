import { signal, computed, effect, createModel } from "@preact/signals";
import api from "../api";
import { to } from "../utils";
import type { Feed, Folder, FeedStat, ItemStatus } from "../api-types";
import type { TranslationKey } from "../i18n";

export interface FeedStatSummary {
  unread: number;
  starred: number;
}

export interface FeedTreeStats {
  feeds: Record<number, FeedStatSummary>;
  folders: Record<number, FeedStatSummary>;
  total: FeedStatSummary;
}

export const RegistryModel = createModel((onError?: (err: Error, key: TranslationKey) => void) => {
  const feeds = signal<Feed[]>([]);
  const folders = signal<Folder[]>([]);
  const feedStats = signal<Record<number, FeedStat>>({});
  const feedErrors = signal<Record<number, string>>({});
  const loadingFeeds = signal(0);
  const refreshAvailable = signal(true);

  let statsTimer: ReturnType<typeof setTimeout> | undefined;
  effect(() => () => clearTimeout(statsTimer));

  const feedsById = computed(() => {
    return feeds.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Feed>);
  });
  const foldersById = computed(() => {
    return folders.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Folder>);
  });
  const stats = computed((): FeedTreeStats => {
    const statsFeeds: Record<number, FeedStatSummary> = {};
    const statsFolders: Record<number, FeedStatSummary> = {};
    const statsTotal: FeedStatSummary = { unread: 0, starred: 0 };

    for (const feed of feeds.value) {
      const n = feedStats.value[feed.id];
      if (!n) continue;
      const fStat: FeedStatSummary = { unread: n.unread || 0, starred: n.starred || 0 };
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

  function adjustStats(feedId: number, from: ItemStatus, to: ItemStatus) {
    const cur = feedStats.value[feedId];
    if (!cur) return;

    let unread = cur.unread || 0;
    let starred = cur.starred || 0;
    if (from === "unread") unread -= 1;
    if (from === "starred") starred -= 1;
    if (to === "unread") unread += 1;
    if (to === "starred") starred += 1;

    feedStats.value = {
      ...feedStats.value,
      [feedId]: { ...cur, unread: Math.max(0, unread), starred: Math.max(0, starred) },
    };
  }

  async function refresh() {
    const [err, values] = await to(Promise.all([api.folders.list(), api.feeds.list()]));
    if (err) {
      onError?.(err, "fail_load");
      return;
    }
    folders.value = values[0];
    feeds.value = values[1];
  }

  async function refreshStats() {
    const [err, data] = await to(api.status());
    if (err) {
      onError?.(err, "fail_load");
      return;
    }

    loadingFeeds.value = data.running;
    refreshAvailable.value = data.refresh;
    if (data.running) statsTimer = setTimeout(refreshStats, 500);

    feedStats.value = data.stats.reduce(
      (acc, stat) => ({ ...acc, [stat.feed_id]: stat }),
      {} as Record<number, FeedStat>,
    );

    const [feedErr, errors] = await to(api.feeds.list_errors());
    if (feedErr) {
      onError?.(feedErr, "fail_load");
      return;
    }
    feedErrors.value = errors;
  }

  return {
    feeds,
    folders,
    feedErrors,
    loadingFeeds,
    refreshAvailable,
    feedsById,
    foldersById,
    stats,
    refresh,
    refreshStats,
    adjustStats,
  };
});

export type AppRegistry = InstanceType<typeof RegistryModel>;
