import { signal, computed, createModel } from "@preact/signals";
import api from "../api";
import { to } from "../utils";
import type { Feed, Folder, FeedStat } from "../api-types";

export interface FeedStatSummary {
  unread: number;
  starred: number;
}

export interface FeedTreeStats {
  feeds: Record<number, FeedStatSummary>;
  folders: Record<number, FeedStatSummary>;
  total: FeedStatSummary;
}

export const RegistryModel = createModel((onError?: (err: unknown) => void) => {
  const feeds = signal<Feed[]>([]);
  const folders = signal<Folder[]>([]);
  const feedStats = signal<Record<number, FeedStat>>({});

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

  async function refresh() {
    const [err, values] = await to(Promise.all([api.folders.list(), api.feeds.list()]));
    if (err) {
      onError?.(err);
      return;
    }
    folders.value = values[0];
    feeds.value = values[1];
  }

  return { feeds, folders, feedStats, feedsById, foldersById, stats, refresh };
});

export type AppRegistry = InstanceType<typeof RegistryModel>;
