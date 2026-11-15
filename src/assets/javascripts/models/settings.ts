import { signal, createModel } from "@preact/signals";
import type { Signal } from "@preact/signals";
import type { Lang } from "../i18n";
import type { Settings, SettingsUpdateData } from "../api-types";
import { debounce } from "../utils";
import api from "../api";

export type Filter = "" | "starred" | "unread";
export type Theme = "system" | "light" | "sepia" | "night";
export type ThemeFont = "" | "serif" | "monospace";

export const SettingsModel = createModel(
  (initial: Partial<Settings> = {}, onSaveError?: (err: unknown) => void) => {
    const values = {
      filter: signal<Filter>((initial.filter as Filter) || ""),
      feed: signal<string | null>(initial.feed ?? null),
      feed_list_width: signal<number>(initial.feed_list_width || 300),
      item_list_width: signal<number>(initial.item_list_width || 300),
      sort_newest_first: signal<boolean>(initial.sort_newest_first ?? true),
      theme_name: signal<Theme>((initial.theme_name as Theme) || "system"),
      theme_font: signal<ThemeFont>((initial.theme_font as ThemeFont) || ""),
      theme_size: signal<number>(initial.theme_size || 1),
      refresh_rate: signal<number>(initial.refresh_rate || 0),
      language: signal<Lang>((initial.language as Lang) || "en"),
    } satisfies Record<keyof SettingsUpdateData, Signal<unknown>>;

    const save = (data: SettingsUpdateData) => {
      api.settings.update(data).catch(err => {
        onSaveError?.(err);
      });
    };

    const saveDebounced = debounce((data: SettingsUpdateData) => save(data), 1000);

    const assign = (patch: SettingsUpdateData) => {
      for (const key of Object.keys(patch) as (keyof SettingsUpdateData)[]) {
        const value = patch[key];
        if (value === undefined) continue;
        (values[key] as Signal<unknown>).value = value;
      }
    };

    return {
      values,

      update(partial: SettingsUpdateData) {
        assign(partial);

        const { feed_list_width, item_list_width, ...rest } = partial;
        if (feed_list_width !== undefined) saveDebounced({ feed_list_width });
        if (item_list_width !== undefined) saveDebounced({ item_list_width });

        if (Object.keys(rest).length > 0) save(rest);
      },
    };
  },
);

export type AppSettings = InstanceType<typeof SettingsModel>;
