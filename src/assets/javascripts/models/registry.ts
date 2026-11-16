import { signal, computed, createModel } from "@preact/signals";
import api from "../api";
import { to } from "../utils";
import type { Feed, Folder } from "../api-types";

export const RegistryModel = createModel((onError?: (err: unknown) => void) => {
  const feeds = signal<Feed[]>([]);
  const folders = signal<Folder[]>([]);

  const feedsById = computed(() => {
    return feeds.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Feed>);
  });
  const foldersById = computed(() => {
    return folders.value.reduce((acc, f) => ((acc[f.id] = f), acc), {} as Record<number, Folder>);
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

  return { feeds, folders, feedsById, foldersById, refresh };
});

export type AppRegistry = InstanceType<typeof RegistryModel>;
