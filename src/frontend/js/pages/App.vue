<template>
  <div
    class="d-flex"
    :class="{
      'feed-selected': feedSelected !== null,
      'item-selected': itemSelected !== null,
    }">
    <!-- feed list -->
    <div
      id="col-feed-list"
      class="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
      :style="{ width: feedListWidth + 'px' }">
      <v-drag :width="feedListWidth" @resize="resizeFeedList"></v-drag>
      <div class="px-2 py-1 d-flex align-items-center">
        <v-icon class="mx-2" name="anchor" />
        <div class="flex-grow-1"></div>
        <button
          class="c-button-pill ms-1"
          :aria-pressed="filterSelected == 'unread'"
          :title="$t('unread')"
          @click="filterSelected = 'unread'">
          <v-icon name="circle-full" />
        </button>
        <button
          class="c-button-pill mx-1"
          :aria-pressed="filterSelected == 'starred'"
          :title="$t('starred')"
          @click="filterSelected = 'starred'">
          <v-icon name="star-full" />
        </button>
        <button
          class="c-button-pill me-1"
          :aria-pressed="filterSelected == ''"
          :title="$t('all')"
          @click="filterSelected = ''">
          <v-icon name="assorted" />
        </button>
        <div class="flex-grow-1"></div>
        <v-dropdown
          toggle-class="c-button-link c-button-pill px-2"
          ref="menuDropdown"
          drop="right"
          :title="$t('settings')">
          <template v-slot:button>
            <v-icon name="more-horizontal" />
          </template>

          <button class="c-dropdown-item w-100 text-start" @click="showSettings('create')">
            <v-icon class="me-1" name="plus" />
            {{ $t("new_feed") }}
          </button>
          <div class="c-dropdown-divider"></div>
          <button class="c-dropdown-item w-100 text-start" @click="fetchAllFeeds()">
            <v-icon class="me-1" name="rotate-cw" />
            {{ $t("refresh_feeds") }}
          </button>

          <div class="c-dropdown-divider"></div>

          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ $t("theme") }}
          </header>
          <div class="row text-center m-0">
            <button
              class="c-button-link theme-swatch col-3 px-0 rounded-0"
              :class="'theme-' + t"
              :title="t"
              :aria-label="t"
              :aria-pressed="theme.name == t"
              @click.stop="theme.name = t"
              v-for="t in ['light', 'sepia', 'night', 'system']"></button>
          </div>

          <div class="c-dropdown-divider"></div>

          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ $t("auto_refresh") }}
          </header>
          <div class="row text-center m-0">
            <button
              class="c-dropdown-item col-4 px-0"
              @click.stop="changeRefreshRate(-1)"
              :disabled="!refreshRate">
              <v-icon name="chevron-down" />
            </button>
            <div class="col-4 d-flex align-items-center justify-content-center user-select-none">
              {{ refreshRateTitle }}
            </div>
            <button
              class="c-dropdown-item col-4 px-0"
              @click.stop="changeRefreshRate(1)"
              :disabled="refreshRate === refreshRateOptions[refreshRateOptions.length - 1].value">
              <v-icon name="chevron-up" />
            </button>
          </div>

          <div class="c-dropdown-divider"></div>

          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ $t("show_first") }}
          </header>
          <div class="d-flex">
            <button
              class="flex-fill c-dropdown-item text-center"
              :aria-pressed="itemSortNewestFirst"
              :class="{ active: itemSortNewestFirst }"
              @click.stop="itemSortNewestFirst = true">
              {{ $t("new") }}
            </button>
            <button
              class="flex-fill c-dropdown-item text-center"
              :aria-pressed="!itemSortNewestFirst"
              :class="{ active: !itemSortNewestFirst }"
              @click.stop="itemSortNewestFirst = false">
              {{ $t("old") }}
            </button>
          </div>
          <div class="c-dropdown-divider"></div>
          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ $t("subscriptions") }}
          </header>
          <form enctype="multipart/form-data" tabindex="-1" ref="opmlInputForm">
            <input
              type="file"
              id="opml-import"
              @change="importOPML"
              name="opml"
              style="opacity: 0; width: 1px; height: 0; position: absolute; z-index: -1" />
            <label
              class="c-dropdown-item mb-0 cursor-pointer w-100"
              for="opml-import"
              @click.stop="">
              <v-icon class="me-1" name="download" />
              {{ $t("import") }}
            </label>
          </form>
          <a class="c-dropdown-item d-block text-start text-decoration-none" href="./opml/export">
            <v-icon class="me-1" name="upload" />
            {{ $t("export") }}
          </a>
          <div class="c-dropdown-divider"></div>
          <button class="c-dropdown-item w-100 text-start" @click="showSettings('shortcuts')">
            <v-icon class="me-1" name="help-circle" />
            {{ $t("shortcuts") }}
          </button>
          <div class="c-dropdown-divider"></div>
          <header class="c-dropdown-header" role="heading" aria-level="2">A / あ / 文</header>
          <div class="container">
            <div class="row">
              <button
                v-for="lang in languages"
                class="c-dropdown-item text-center col-3 px-0"
                :aria-label="lang.name"
                :aria-pressed="language === lang.code"
                :title="lang.name"
                @click.stop="changeLanguage(lang.code)">
                {{ lang.code }}
              </button>
            </div>
          </div>
          <div class="c-dropdown-divider" v-if="requiresAuth"></div>
          <button class="c-dropdown-item w-100 text-start" v-if="requiresAuth" @click="logout()">
            <v-icon class="me-1" name="log-out" />
            {{ $t("log_out") }}
          </button>
        </v-dropdown>
      </div>
      <div id="feed-list-scroll" class="p-2 overflow-auto border-top flex-grow-1">
        <v-feedtree
          :tree="feedTree"
          v-model="feedSelected"
          :filter-selected="filterSelected"
          :stats="stats"
          :feed-errors="feed_errors"
          @toggle-folder="toggleFolderExpanded" />
      </div>
      <div
        class="px-2 py-1 d-flex align-items-center border-top flex-shrink-0"
        v-if="loading.feeds">
        <span class="c-spinner mx-2"></span>
        <span class="text-truncate cursor-default user-select-none">{{
          $t("refreshing_progress", { count: loading.feeds })
        }}</span>
      </div>
    </div>
    <!-- item list -->
    <div
      id="col-item-list"
      class="vh-100 position-relative d-flex flex-column border-end flex-shrink-0"
      :style="{ width: itemListWidth + 'px' }">
      <v-drag :width="itemListWidth" @resize="resizeItemList"></v-drag>
      <div class="px-2 py-1 d-flex gap-1 align-items-center">
        <button
          class="c-button-pill d-md-none"
          @click="feedSelected = null"
          :title="$t('show_feeds')">
          <v-icon name="chevron-left" />
        </button>
        <div class="c-search flex-grow-1">
          <v-icon name="search" />
          <!-- id used by keybindings -->
          <input
            id="searchbar"
            class="d-block"
            v-model="itemSearch"
            :placeholder="$t('search_placeholder', { scope: searchScope })"
            @keydown.enter="($event.target as HTMLInputElement).blur()" />
        </div>
        <button
          class="c-button-pill"
          @click="markItemsRead()"
          v-if="filterSelected == 'unread'"
          :title="$t('mark_all_read')">
          <v-icon name="check" />
        </button>

        <button class="c-button-link c-button-pill px-2" v-if="!current.type" disabled>
          <v-icon name="more-horizontal" />
        </button>
        <v-dropdown
          toggle-class="c-button-link c-button-pill px-2"
          drop="right"
          :title="$t('feed_settings')"
          v-if="current?.feed?.id">
          <template v-slot:button>
            <v-icon name="more-horizontal" />
          </template>
          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ current?.feed?.title }}
          </header>
          <a
            class="c-dropdown-item d-block text-start text-decoration-none"
            :href="current?.feed?.link"
            rel="noopener noreferrer"
            target="_blank"
            referrerpolicy="no-referrer"
            v-if="current?.feed?.link">
            <v-icon class="me-1" name="globe" />
            {{ $t("website") }}
          </a>
          <a
            class="c-dropdown-item d-block text-start text-decoration-none"
            :href="current.feed.feed_link"
            rel="noopener noreferrer"
            target="_blank"
            referrerpolicy="no-referrer"
            v-if="current.feed.feed_link">
            <v-icon class="me-1" name="rss" />
            {{ $t("feed_link") }}
          </a>
          <div class="c-dropdown-divider" v-if="current.feed.link || current.feed.feed_link"></div>
          <button class="c-dropdown-item w-100 text-start" @click="renameFeed(current.feed)">
            <v-icon class="me-1" name="edit" />
            {{ $t("rename") }}
          </button>
          <button
            class="c-dropdown-item w-100 text-start"
            @click="updateFeedLink(current.feed)"
            v-if="current.feed.feed_link">
            <v-icon class="me-1" name="edit" />
            {{ $t("change_link") }}
          </button>
          <div class="c-dropdown-divider"></div>
          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ $t("move_to") }}
          </header>
          <template v-for="folder in folders">
            <button
              class="c-dropdown-item w-100 text-start"
              v-if="folder.id != current.feed.folder_id"
              @click="moveFeed(current.feed, folder.id)">
              <v-icon class="me-1" name="folder" />
              {{ folder.title }}
            </button>
          </template>
          <button
            class="c-dropdown-item w-100 text-start opacity-75"
            @click="moveFeed(current.feed, null)"
            v-if="current.feed.folder_id">
            <v-icon class="me-1" name="folder-minus" />
            ──
          </button>
          <button
            class="c-dropdown-item w-100 text-start opacity-75"
            @click="moveFeedToNewFolder(current.feed)">
            <v-icon class="me-1" name="folder-plus" />
            {{ $t("new_folder") }}
          </button>
          <div class="c-dropdown-divider"></div>
          <button
            class="c-dropdown-item w-100 text-start text-danger"
            @click.prevent="deleteFeed(current.feed)">
            <v-icon class="me-1" name="trash" />
            {{ $t("delete") }}
          </button>
        </v-dropdown>
        <v-dropdown
          toggle-class="c-button-link c-button-pill px-2"
          :title="$t('folder_settings')"
          drop="right"
          v-if="current?.folder?.id">
          <template v-slot:button>
            <v-icon name="more-horizontal" />
          </template>
          <header class="c-dropdown-header" role="heading" aria-level="2">
            {{ current?.folder?.title }}
          </header>
          <button class="c-dropdown-item w-100 text-start" @click="renameFolder(current.folder)">
            <v-icon class="me-1" name="edit" />
            {{ $t("rename") }}
          </button>
          <div class="c-dropdown-divider"></div>
          <button
            class="c-dropdown-item w-100 text-start text-danger"
            @click="deleteFolder(current.folder)">
            <v-icon class="me-1" name="trash" />
            {{ $t("delete") }}
          </button>
        </v-dropdown>
      </div>
      <div
        id="item-list-scroll"
        class="d-flex flex-column p-2 overflow-auto border-top flex-grow-1 gap-1"
        v-scroll="loadMoreItems"
        ref="itemlist">
        <div
          v-for="item in items"
          :key="item.id"
          class="c-listitem d-flex flex-column user-select-none"
          role="radio"
          :aria-checked="itemSelected === item.id"
          @click="itemSelected = item.id">
          <div
            style="line-height: 100%; opacity: 0.7; margin-bottom: 0.1rem"
            class="d-flex align-items-center">
            <transition name="indicator">
              <v-icon
                :small="true"
                class="me-1"
                name="circle-full"
                v-if="item.status == 'unread'" />
              <v-icon
                :small="true"
                class="me-1"
                name="star-full"
                v-else-if="item.status == 'starred'" />
            </transition>
            <small class="flex-fill text-truncate me-1">
              {{ (feedsById[item.feed_id] || {}).title }}
            </small>
            <small class="flex-shrink-0"
              ><v-relative-time v-bind:title="formatDate(item.date)" :val="item.date"
            /></small>
          </div>
          <div class="text-break">{{ item.title || $t("untitled") }}</div>
        </div>
        <div class="text-center my-3" v-if="itemsHasMore">
          <span class="c-spinner"></span>
        </div>
      </div>
      <div
        class="px-3 py-2 border-top text-danger text-break"
        v-if="current?.feed?.id && feed_errors[current.feed.id]">
        {{ feed_errors[current.feed.id] }}
      </div>
    </div>
    <!-- item show -->
    <div id="col-item" class="vh-100 d-flex flex-column w-100" style="min-width: 0">
      <div class="px-2 py-1 d-flex gap-1 align-items-center" v-if="itemSelectedDetails">
        <button
          class="c-button-pill"
          @click="toggleItemStarred(itemSelectedDetails)"
          :title="$t('mark_starred')">
          <v-icon name="star-full" v-if="itemSelectedDetails.status == 'starred'" />
          <v-icon name="star" v-else />
        </button>
        <button
          class="c-button-pill"
          :title="$t('mark_unread')"
          @click="toggleItemRead(itemSelectedDetails)">
          <v-icon name="circle-full" v-if="itemSelectedDetails.status == 'unread'" />
          <v-icon name="circle" v-else />
        </button>
        <v-dropdown toggle-class="c-button-pill px-2" drop="center" :title="$t('appearance')">
          <template v-slot:button>
            <v-icon name="sliders" />
          </template>

          <button
            class="c-dropdown-item w-100 text-start font-sans-serif"
            :aria-pressed="theme.font == ''"
            @click.stop="theme.font = ''">
            {{ $t("sans_serif") }}
          </button>
          <button
            class="c-dropdown-item w-100 text-start font-serif"
            :aria-pressed="theme.font == 'serif'"
            @click.stop="theme.font = 'serif'">
            {{ $t("serif") }}
          </button>
          <button
            class="c-dropdown-item w-100 text-start font-monospace"
            :aria-pressed="theme.font == 'monospace'"
            @click.stop="theme.font = 'monospace'">
            {{ $t("monospace") }}
          </button>

          <div class="d-flex text-center">
            <button
              class="c-dropdown-item flex-fill"
              style="font-size: 0.8rem"
              @click.stop="incrFont(-1)">
              A
            </button>
            <button
              class="c-dropdown-item flex-fill"
              style="font-size: 1.2rem"
              @click.stop="incrFont(1)">
              A
            </button>
          </div>
        </v-dropdown>
        <button
          class="c-button-pill"
          :aria-pressed="!!itemSelectedReadability"
          @click="toggleReadability()"
          :title="$t('read_here')">
          <v-icon :class="{ 'is-loading': loading.readability }" name="book-open" />
        </button>
        <a
          class="c-button-pill"
          :href="itemSelectedDetails.link"
          rel="noopener noreferrer"
          target="_blank"
          referrerpolicy="no-referrer"
          :title="$t('open_link')">
          <v-icon name="external-link" />
        </a>
        <div class="flex-grow-1"></div>
        <button
          class="c-button-pill"
          @click="navigateToItem(-1)"
          :title="$t('previous_article')"
          :disabled="!items.length || itemSelected == items[0].id">
          <v-icon name="chevron-left" />
        </button>
        <button
          class="c-button-pill"
          @click="navigateToItem(+1)"
          :title="$t('next_article')"
          :disabled="!items.length || itemSelected == items[items.length - 1].id">
          <v-icon name="chevron-right" />
        </button>
        <button class="c-button-pill" @click="itemSelected = null" :title="$t('close_article')">
          <v-icon name="x" />
        </button>
      </div>
      <div
        v-if="itemSelectedDetails"
        ref="content"
        class="content px-4 pt-3 pb-5 border-top overflow-auto"
        :class="{
          'font-sans-serif': theme.font == '',
          'font-serif': theme.font == 'serif',
          'font-monospace': theme.font == 'monospace',
        }"
        :style="{ 'font-size': theme.size + 'rem' }">
        <div class="content-wrapper">
          <h1>
            <b>{{ itemSelectedDetails.title || $t("untitled") }}</b>
          </h1>
          <div class="opacity-50">
            <div>
              <span
                class="cursor-pointer"
                @click="feedSelected = 'feed:' + (feedsById[itemSelectedDetails.feed_id] || {}).id">
                {{ (feedsById[itemSelectedDetails.feed_id] || {}).title }}
              </span>
            </div>
            <time>{{ formatDate(itemSelectedDetails.date) }}</time>
          </div>
          <hr />
          <div v-if="!itemSelectedReadability">
            <div v-if="contentImages.length">
              <figure v-for="media in contentImages">
                <img :src="media.url" loading="lazy" />
                <figcaption v-if="media.description">
                  {{ media.description }}
                </figcaption>
              </figure>
            </div>
            <audio class="w-100" controls v-for="media in contentAudios" :src="media.url"></audio>
            <video class="w-100" controls v-for="media in contentVideos" :src="media.url"></video>
          </div>
          <div v-html="itemSelectedContent"></div>
        </div>
      </div>
    </div>
    <v-modal :open="!!settings" @hide="settings = ''">
      <button
        class="c-button-link outline-none position-absolute top-0 end-0 p-2 m-2"
        style="line-height: 1"
        @click="settings = ''">
        <v-icon name="x" />
      </button>
      <div v-if="settings == 'create'" class="d-flex flex-column">
        <p class="cursor-default mb-3">
          <b>{{ $t("new_feed") }}</b>
        </p>
        <form @submit.prevent="createFeed($event)" class="d-flex flex-column">
          <label for="feed-url" class="mb-2">{{ $t("url") }}</label>
          <input
            id="feed-url"
            name="url"
            type="url"
            class="c-input"
            required
            autocomplete="off"
            :readonly="feedNewChoice.length > 0"
            placeholder="https://example.com/feed"
            v-focus />
          <label for="feed-folder" class="mb-2 mt-3">
            {{ $t("folder") }}
            <a
              href="#"
              class="float-end text-decoration-none"
              @click.prevent="createNewFeedFolder()"
              >{{ $t("new_folder") }}</a
            >
          </label>
          <select class="c-input" id="feed-folder" name="folder_id" ref="newFeedFolder">
            <option value="">---</option>
            <option
              :value="folder.id"
              v-for="folder in folders"
              :selected="
                folder.id === current?.feed?.folder_id || folder.id === current?.folder?.id
              ">
              {{ folder.title }}
            </option>
          </select>
          <div class="mt-3" v-if="feedNewChoice.length">
            <p class="mb-2">
              {{ $t("multiple_feeds_found") }}
              <a
                href="#"
                class="float-end text-decoration-none"
                @click.prevent="resetFeedChoice()"
                >{{ $t("cancel") }}</a
              >
            </p>
            <div class="d-flex flex-column gap-1">
              <div
                class="c-listitem d-flex flex-column user-select-none"
                role="radio"
                :aria-checked="feedNewChoiceSelected === choice.url"
                @click="feedNewChoiceSelected = choice.url"
                v-for="choice in feedNewChoice">
                <div class="text-truncate">{{ choice.title }}</div>
                <div class="text-truncate" :class="{ 'opacity-50': choice.title }">
                  {{ choice.url }}
                </div>
              </div>
            </div>
          </div>
          <button class="c-button mt-3" :disabled="loading.newfeed" type="submit">
            <span class="c-spinner" v-if="loading.newfeed"></span>
            <span v-else>{{ $t("add") }}</span>
          </button>
        </form>
      </div>
      <v-shortcuts v-else-if="settings == 'shortcuts'" />
    </v-modal>
    <v-toast ref="toast" />
  </div>
</template>

<script lang="ts">
import type { Lang } from "../i18n";
import api, { NetworkError, HTTPError } from "../api";
import { scrollto, debounce, debounceMixin, to } from "../utils";
import drag from "../components/drag.vue";
import dropdown from "../components/dropdown.vue";
import modal from "../components/modal.vue";
import shortcuts from "../components/shortcuts.vue";
import relativeTime from "../components/relative-time.vue";
import icon from "../components/icon.vue";
import feedTree from "../components/feedtree.vue";
import toast from "../components/toast.vue";
import type { FeedTreeNode, TreeFeedNode, TreeFolderNode } from "../components/feedtree.vue";
import scrollDir from "../directives/scroll";
import focusDir from "../directives/focus";
import { defineComponent } from "vue";
import type {
  Feed,
  Folder,
  Item,
  FeedStat,
  FeedLink,
  MediaLink,
  ItemStatus,
  FeedCreateData,
  ItemListQuery,
  ItemMarkQuery,
} from "../api-types";

var app = window.app;

declare module "vue" {
  interface ComponentCustomProperties {
    $refs: {
      itemlist: HTMLElement;
      content: HTMLElement;
      newFeedFolder: HTMLSelectElement;
      opmlInputForm: HTMLFormElement;
      menuDropdown: InstanceType<typeof dropdown>;
      toast: InstanceType<typeof toast>;
    };
  }
}

type Theme = "system" | "light" | "sepia" | "night";
type ThemeFont = "" | "serif" | "monospace";
type Filter = "" | "starred" | "unread";
type SettingsLanguage = { code: Lang; name: string };
type Stats = { unread: number; starred: number };

var TITLE = document.title;

export default defineComponent({
  mixins: [debounceMixin],
  components: {
    "v-drag": drag,
    "v-dropdown": dropdown,
    "v-modal": modal,
    "v-shortcuts": shortcuts,
    "v-relative-time": relativeTime,
    "v-icon": icon,
    "v-feedtree": feedTree,
    "v-toast": toast,
  },
  directives: {
    scroll: scrollDir,
    focus: focusDir,
  },
  async created() {
    this.updateMetaTheme();
    this.$t.set(document.documentElement.lang as Lang);

    // keep the theme-color meta tag in sync when the OS color scheme changes
    this._colorSchemeMql = window.matchMedia("(prefers-color-scheme: dark)");
    this._colorSchemeMql.addEventListener("change", this.updateMetaTheme);

    const [statsErr] = await to(this.refreshStats());
    if (statsErr) {
      this.$refs.toast.addToast(
        { title: this.$t("fail_load"), description: this.errDescription(statsErr) },
        { level: "fail", closeable: false },
      );
    }

    const [feedsErr] = await to(this.refreshFeeds());
    if (feedsErr) {
      this.$refs.toast.addToast(
        { title: this.$t("fail_load"), description: this.errDescription(feedsErr) },
        { level: "fail", closeable: false },
      );
    }

    if (!feedsErr) {
      this.refreshItems(false);
      this.computeStats();
    }
  },
  beforeUnmount() {
    this._colorSchemeMql?.removeEventListener("change", this.updateMetaTheme);
  },
  data() {
    var s = app.settings;
    return {
      filterSelected: s.filter as Filter,
      folders: [] as Folder[],
      feeds: [] as Feed[],
      feedSelected: s.feed,
      feedListWidth: s.feed_list_width || 300,
      feedNewChoice: [] as FeedLink[],
      feedNewChoiceSelected: "",
      items: [] as Item[],
      itemsHasMore: true,
      itemSelected: null as number | null,
      itemSelectedDetails: null as Item | null,
      itemSelectedReadability: "",
      itemSearch: "",
      itemSortNewestFirst: s.sort_newest_first as boolean,
      itemListWidth: s.item_list_width || 300,

      stats: { folders: {}, feeds: {}, total: { unread: 0, starred: 0 } } as {
        folders: Record<number, Stats>;
        feeds: Record<number, Stats>;
        total: Stats;
      },

      settings: "",
      loading: {
        feeds: 0,
        newfeed: false,
        items: false,
        readability: false,
      },
      fonts: ["", "serif", "monospace"] as ThemeFont[],
      feedStats: {} as Record<number, FeedStat>,
      theme: {
        name: s.theme_name as Theme,
        font: s.theme_font as ThemeFont,
        size: s.theme_size as number,
      },
      themeColors: {
        night: "#0e0e0e",
        sepia: "#f4f0e5",
        light: "#fff",
      },
      refreshRate: s.refresh_rate,
      authenticated: app.authenticated,
      requiresAuth: app.requiresAuth,
      feed_errors: {} as Record<number, string>,

      refreshRateOptions: [
        { title: "0", value: 0 },
        { title: "10m", value: 10 },
        { title: "30m", value: 30 },
        { title: "1h", value: 60 },
        { title: "2h", value: 120 },
        { title: "4h", value: 240 },
        { title: "12h", value: 720 },
        { title: "24h", value: 1440 },
      ],

      language: s.language,
      languages: [
        { code: "en", name: "English" },
        { code: "de", name: "Deutsch" },
        { code: "es", name: "Español" },
        { code: "fr", name: "Français" },
        { code: "ja", name: "日本語" },
        { code: "pt", name: "Português" },
        { code: "ru", name: "Русский" },
        { code: "zh", name: "简体中文" },
      ] as SettingsLanguage[],

      _colorSchemeMql: null as MediaQueryList | null,
    };
  },
  computed: {
    feedTree(): FeedTreeNode[] {
      const [rootFeeds, folderFeeds] = this.feeds.reduce(
        (acc, f) => {
          acc[f.folder_id === null ? 0 : 1].push(f);
          return acc;
        },
        [[] as Feed[], [] as Feed[]],
      );

      const byFolder: Record<number, Feed[]> = folderFeeds.reduce(
        (acc, f) => {
          (acc[f.folder_id as number] ||= []).push(f);
          return acc;
        },
        {} as Record<number, Feed[]>,
      );

      const feedNode = (feed: Feed): TreeFeedNode => ({ type: "feed", feed });

      return [
        ...this.folders
          .filter(folder => !this.mustHideFolder(folder))
          .map(folder => {
            const feeds = (byFolder[folder.id] || [])
              .filter(f => !this.mustHideFeed(f))
              .map(feedNode);
            return { type: "folder" as const, folder, feeds };
          }),
        ...rootFeeds.filter(f => !this.mustHideFeed(f)).map(feedNode),
      ];
    },
    feedsById(): Record<number, Feed> {
      return this.feeds.reduce((acc, f) => ({ ...acc, [f.id]: f }), {});
    },
    foldersById(): Record<number, Folder> {
      return this.folders.reduce((acc, f) => ({ ...acc, [f.id]: f }), {});
    },
    current(): { type: string; feed: Feed | null; folder: Folder | null } {
      var parts = (this.feedSelected || "").split(":", 2);
      var type = parts[0];
      var guid = parts[1];

      const feed = type == "feed" ? this.feedsById[guid] : null;
      const folder = type == "folder" ? this.foldersById[guid] : null;

      return { type: type, feed: feed, folder: folder };
    },
    searchScope(): string {
      var type = (this.feedSelected || "").split(":", 2)[0];
      if (type == "feed")
        return (this.feedsById[this.feedSelected.split(":", 2)[1]] || {}).title || "";
      if (type == "folder")
        return (this.foldersById[this.feedSelected.split(":", 2)[1]] || {}).title || "";
      if (this.filterSelected == "unread") return this.$t("all_unread");
      if (this.filterSelected == "starred") return this.$t("all_starred");
      return this.$t("all_feeds");
    },
    itemSelectedContent(): string {
      if (!this.itemSelected) return "";
      if (this.itemSelectedReadability) return this.itemSelectedReadability;
      return this.itemSelectedDetails?.content || "";
    },
    contentImages(): MediaLink[] {
      if (!this.itemSelectedDetails) return [] as MediaLink[];
      return (this.itemSelectedDetails.media_links || []).filter(l => l.type === "image");
    },
    contentAudios(): MediaLink[] {
      if (!this.itemSelectedDetails) return [] as MediaLink[];
      return (this.itemSelectedDetails.media_links || []).filter(l => l.type === "audio");
    },
    contentVideos(): MediaLink[] {
      if (!this.itemSelectedDetails) return [] as MediaLink[];
      return (this.itemSelectedDetails.media_links || []).filter(l => l.type === "video");
    },
    refreshRateTitle() {
      const entry = this.refreshRateOptions.find(o => o.value === this.refreshRate);
      return entry ? entry.title : "0";
    },
  },
  watch: {
    theme: {
      deep: true,
      handler(theme) {
        this.updateMetaTheme();
        (async () => {
          const [err, _] = await to(
            api.settings.update({
              theme_name: theme.name,
              theme_font: theme.font,
              theme_size: theme.size,
            }),
          );
          if (err) {
            this.$refs.toast.addToast(
              { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
              { level: "fail", closeable: false },
            );
          }
        })();
      },
    },
    feedStats: {
      deep: true,
      handler() {
        this.$debounce("watch:feedStats", this.computeStats, 500);
      },
    },
    async filterSelected(newVal, oldVal) {
      if (oldVal === undefined) return;
      this.itemSelected = null;
      this.items = [];
      this.itemsHasMore = true;

      const [err] = await to(api.settings.update({ filter: newVal }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.refreshItems(false);
      this.computeStats();
    },
    async feedSelected(newVal, oldVal) {
      if (oldVal === undefined) return;
      this.itemSelected = null;
      this.items = [];
      this.itemsHasMore = true;

      this.refreshItems(false);
      if (this.$refs.itemlist) this.$refs.itemlist.scrollTop = 0;

      const [err] = await to(api.settings.update({ feed: newVal }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
      }
    },
    async itemSelected(newVal, oldVal) {
      this.itemSelectedReadability = "";
      if (newVal === null) {
        this.itemSelectedDetails = null;
        return;
      }
      if (this.$refs.content) this.$refs.content.scrollTop = 0;

      const [itemErr, item] = await to(api.items.get(newVal));
      if (itemErr) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_load"), description: this.errDescription(itemErr) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.itemSelectedDetails = item;
      const details = this.itemSelectedDetails;
      if (details.status == "unread") {
        const [updateErr] = await to(api.items.update(details.id, { status: "read" }));
        if (updateErr) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_update_article"), description: this.errDescription(updateErr) },
            { level: "fail", closeable: false },
          );
          return;
        }
        this.feedStats[details.feed_id].unread -= 1;
        var itemInList = this.items.find(i => i.id == item.id);
        if (itemInList) itemInList.status = "read";
        details.status = "read";
      }
    },
    itemSearch() {
      this.$debounce("watch:itemSearch", this.refreshItems, 500);
    },
    async itemSortNewestFirst(newVal, oldVal) {
      if (oldVal === undefined) return;
      const [err] = await to(api.settings.update({ sort_newest_first: newVal }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.refreshItems(false);
    },
    feedListWidth: debounce(function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ feed_list_width: newVal });

      // const [err] = await to(api.settings.update({ feed_list_width: newVal }));
      // if (err) {
      //   this.$refs.toast.addToast(
      //     { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
      //     { level: "fail", closeable: false },
      //   );
      // }
    }, 1000),
    itemListWidth: debounce(function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ item_list_width: newVal });

      // const [err] = await to(api.settings.update({ item_list_width: newVal }));
      // if (err) {
      //   this.$refs.toast.addToast(
      //     { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
      //     { level: "fail", closeable: false },
      //   );
      // }
    }, 1000),
    refreshRate(newVal, oldVal) {
      if (oldVal === undefined) return;
      (async () => {
        const [err] = await to(api.settings.update({ refresh_rate: newVal }));
        if (err) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
            { level: "fail", closeable: false },
          );
        }
      })();
    },
  },
  methods: {
    updateMetaTheme() {
      let theme = this.theme.name;
      if (theme == "system") {
        var dark = window?.matchMedia("(prefers-color-scheme: dark)").matches;
        theme = dark ? "night" : "light";
      }
      const metaTag: HTMLMetaElement | null = document.querySelector("meta[name='theme-color']");
      metaTag && (metaTag.content = this.themeColors[theme]);

      document.documentElement.dataset.theme = this.theme.name;
    },
    async refreshStats(loopMode?: boolean) {
      const [err, data] = await to(api.status());
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_load"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }

      if (loopMode && !this.itemSelected) this.refreshItems();

      this.loading.feeds = data.running;
      if (data.running) {
        setTimeout(() => this.refreshStats(true), 500);
      }
      this.feedStats = data.stats.reduce((acc, stat) => ({ ...acc, [stat.feed_id]: stat }), {});

      const [feedErr, errors] = await to(api.feeds.list_errors());
      if (feedErr) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_load"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.feed_errors = errors;
    },
    getItemsQuery(): ItemListQuery {
      var query: ItemListQuery = {};
      if (this.feedSelected) {
        var parts = this.feedSelected.split(":", 2);
        var type = parts[0];
        var guid = parts[1];
        if (type == "feed") {
          query.feed_id = guid;
        } else if (type == "folder") {
          query.folder_id = guid;
        }
      }
      if (this.filterSelected) {
        query.status = this.filterSelected;
      }
      if (this.itemSearch) {
        query.search = this.itemSearch;
      }
      if (!this.itemSortNewestFirst) {
        query.oldest_first = true;
      }
      return query;
    },
    async refreshFeeds() {
      const [err, values] = await to(Promise.all([api.folders.list(), api.feeds.list()]));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_load"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.folders = values[0];
      this.feeds = values[1];
    },
    async refreshItems(loadMore = false) {
      if (this.feedSelected === null) {
        this.items = [];
        this.itemsHasMore = false;
        return;
      }

      var query = this.getItemsQuery();
      if (loadMore) {
        query.after = this.items[this.items.length - 1].id;
      }

      this.loading.items = true;
      const [err, data] = await to(api.items.list(query));
      this.loading.items = false;
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_load"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }

      if (loadMore) {
        this.items = this.items.concat(data.list);
      } else {
        this.items = data.list;
      }
      this.itemsHasMore = data.has_more;

      // load more if there's some space left at the bottom of the item list.
      this.$nextTick(() => {
        if (this.itemsHasMore && !this.loading.items && this.itemListCloseToBottom()) {
          this.refreshItems(true);
        }
      });
    },
    itemListCloseToBottom() {
      // approx. vertical space at the bottom of the list (loading el & paddings) when 1rem = 16px
      var bottomSpace = 70;
      var scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;

      var el = this.$refs.itemlist;

      if (!el || el.scrollHeight === 0) return false; // element is invisible (responsive design)

      var closeToBottom = el.scrollHeight - el.scrollTop - el.offsetHeight < bottomSpace * scale;
      return closeToBottom;
    },
    loadMoreItems() {
      if (!this.itemsHasMore) return;
      if (this.loading.items) return;
      if (this.itemListCloseToBottom()) return this.refreshItems(true);
      if (this.itemSelected && this.itemSelected === this.items[this.items.length - 1].id)
        return this.refreshItems(true);
    },
    async markItemsRead() {
      const markQuery = this.getItemsQuery();
      const query: ItemMarkQuery = {
        folder_id: markQuery.folder_id,
        feed_id: markQuery.feed_id,
      };
      const [err] = await to(api.items.mark_read(query));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_update_article"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.items = [];
      this.itemSelected = null;
      this.itemsHasMore = false;
      this.refreshStats();
    },
    async toggleFolderExpanded(folder: Folder) {
      folder.is_expanded = !folder.is_expanded;
      const [err] = await to(api.folders.update(folder.id, { is_expanded: folder.is_expanded }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_folder"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
      }
    },
    formatDate(datestr: string) {
      const options: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      };
      return new Date(datestr).toLocaleDateString(undefined, options);
    },
    async moveFeed(feed: Feed, folder_id: number | null) {
      const [err] = await to(api.feeds.update(feed.id, { folder_id }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_feed"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      feed.folder_id = folder_id;
      this.refreshStats();
    },
    async moveFeedToNewFolder(feed: Feed) {
      const title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;

      const [folderErr, folder] = await to(api.folders.create({ title: title }));
      if (folderErr) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_folder"), description: this.errDescription(folderErr) },
          { level: "fail", closeable: false },
        );
        return;
      }

      const [updateErr] = await to(api.feeds.update(feed.id, { folder_id: folder.id }));
      if (updateErr) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_feed"), description: this.errDescription(updateErr) },
          { level: "fail", closeable: false },
        );
        return;
      }

      const [feedsErr] = await to(this.refreshFeeds());
      if (!feedsErr) this.refreshStats();
    },
    async createNewFeedFolder() {
      const title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;

      const [folderErr, result] = await to(api.folders.create({ title: title }));
      if (folderErr) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_folder"), description: this.errDescription(folderErr) },
          { level: "fail", closeable: false },
        );
        return;
      }

      const [feedsErr] = await to(this.refreshFeeds());
      if (!feedsErr) {
        this.$nextTick(() => {
          if (this.$refs.newFeedFolder) {
            this.$refs.newFeedFolder.value = String(result.id);
          }
        });
      }
    },
    async renameFolder(folder: Folder) {
      const newTitle = prompt(this.$t("prompt_new_title"), folder.title);
      if (newTitle) {
        const [err] = await to(api.folders.update(folder.id, { title: newTitle }));
        if (err) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_save_folder"), description: this.errDescription(err) },
            { level: "fail", closeable: false },
          );
          return;
        }
        folder.title = newTitle;
        this.folders.sort((a, b) => a.title.localeCompare(b.title));
      }
    },
    async deleteFolder(folder: Folder) {
      if (confirm(this.$t("confirm_delete", { name: folder.title }))) {
        const [err] = await to(api.folders.delete(folder.id));
        if (err) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_save_folder"), description: this.errDescription(err) },
            { level: "fail", closeable: false },
          );
          return;
        }
        this.feedSelected = null;
        this.refreshStats();
        this.refreshFeeds();
      }
    },
    async updateFeedLink(feed: Feed) {
      const newLink = prompt(this.$t("prompt_feed_link"), feed.feed_link);
      if (newLink !== null) {
        const [err] = await to(api.feeds.update(feed.id, { feed_link: newLink }));
        if (err) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_save_feed"), description: this.errDescription(err) },
            { level: "fail", closeable: false },
          );
          return;
        }
        feed.feed_link = newLink;
      }
    },
    async renameFeed(feed: Feed) {
      const newTitle = prompt(this.$t("prompt_new_title"), feed.title);
      if (newTitle) {
        const [err] = await to(api.feeds.update(feed.id, { title: newTitle }));
        if (err) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_save_feed"), description: this.errDescription(err) },
            { level: "fail", closeable: false },
          );
          return;
        }
        feed.title = newTitle;
      }
    },
    async deleteFeed(feed: Feed) {
      if (confirm(this.$t("confirm_delete", { name: feed.title }))) {
        const [err] = await to(api.feeds.delete(feed.id));
        if (err) {
          this.$refs.toast.addToast(
            { title: this.$t("fail_save_feed"), description: this.errDescription(err) },
            { level: "fail", closeable: false },
          );
          return;
        }
        this.feedSelected = null;
        this.refreshStats();
        this.refreshFeeds();
      }
    },
    async createFeed($event: Event) {
      var form = $event.target as HTMLFormElement;
      var data: FeedCreateData = {
        url: (form.querySelector("input[name=url]") as HTMLInputElement).value,
        folder_id:
          parseInt((form.querySelector("select[name=folder_id]") as HTMLSelectElement).value) ||
          null,
      };
      if (this.feedNewChoiceSelected) {
        var choice = this.feedNewChoice.find(c => c.url === this.feedNewChoiceSelected);
        data.url = this.feedNewChoiceSelected;
        if (choice && choice.title_override) data.title_override = choice.title_override;
      }
      this.loading.newfeed = true;
      const [err, result] = await to(api.feeds.create(data));
      this.loading.newfeed = false;
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_feed"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      if (result.status === "success") {
        this.refreshFeeds();
        this.refreshStats();
        this.settings = "";
        this.feedSelected = "feed:" + result.feed.id;
      } else if (result.status === "multiple") {
        this.feedNewChoice = result.choice;
        this.feedNewChoiceSelected = result.choice[0].url;
      } else {
        alert("No feeds found at the given url.");
      }
    },
    async toggleItemStatus(item: Item, targetstatus: ItemStatus) {
      const fallbackstatus: ItemStatus = "read";
      const oldstatus = item.status;
      const newstatus = item.status !== targetstatus ? targetstatus : fallbackstatus;

      const updateStats = (status: ItemStatus, incr: number) => {
        if (status == "unread" || status == "starred") {
          this.feedStats[item.feed_id][status] += incr;
        }
      };

      const [err] = await to(api.items.update(item.id, { status: newstatus }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_update_article"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }

      updateStats(oldstatus, -1);
      updateStats(newstatus, +1);

      var itemInList = this.items.find(i => i.id == item.id);
      if (itemInList) itemInList.status = newstatus;
      item.status = newstatus;
    },
    toggleItemStarred(item: Item) {
      this.toggleItemStatus(item, "starred");
    },
    toggleItemRead(item: Item) {
      this.toggleItemStatus(item, "unread");
    },
    async importOPML(event: Event) {
      const input = event.target as HTMLInputElement;
      const form = this.$refs.opmlInputForm;
      this.$refs.menuDropdown.hide();
      const [err] = await to(api.upload_opml(form));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_import"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      input.value = "";
      this.refreshFeeds();
      this.refreshStats();
    },
    async logout() {
      const [err] = await to(api.logout());
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_logout"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      document.location.reload();
    },
    async toggleReadability() {
      if (this.itemSelectedReadability) {
        this.itemSelectedReadability = "";
        return;
      }
      var item = this.itemSelectedDetails;
      if (!item?.link) return;
      this.loading.readability = true;
      const [err, data] = await to(api.crawl(item!.link));
      this.loading.readability = false;

      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_readability"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
      } else {
        this.itemSelectedReadability = data?.content || "";
      }
    },
    showSettings(settings: string) {
      this.settings = settings;

      if (settings === "create") {
        this.feedNewChoice = [];
        this.feedNewChoiceSelected = "";
      }
    },
    resizeFeedList(width: number) {
      this.feedListWidth = Math.min(Math.max(200, width), 700);
    },
    resizeItemList(width: number) {
      this.itemListWidth = Math.min(Math.max(200, width), 700);
    },
    resetFeedChoice() {
      this.feedNewChoice = [];
      this.feedNewChoiceSelected = "";
    },
    incrFont(x: number) {
      this.theme.size = +(this.theme.size + 0.1 * x).toFixed(1);
    },
    async fetchAllFeeds() {
      if (this.loading.feeds) return;
      const [err] = await to(api.feeds.refresh());
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_refresh"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      this.refreshStats();
    },
    computeStats() {
      let statsFeeds: Record<number, Stats> = {},
        statsFolders: Record<number, Stats> = {},
        statsTotal: Stats = { unread: 0, starred: 0 };

      for (var i = 0; i < this.feeds.length; i++) {
        const feed = this.feeds[i];
        if (!this.feedStats[feed.id]) continue;

        const n = this.feedStats[feed.id];
        const feedStats = { unread: n.unread || 0, starred: n.starred || 0 };

        statsFeeds[feed.id] = feedStats;

        if (feed.folder_id !== null) {
          if (!statsFolders[feed.folder_id])
            statsFolders[feed.folder_id] = { unread: 0, starred: 0 };
          statsFolders[feed.folder_id].unread += feedStats.unread;
          statsFolders[feed.folder_id].starred += feedStats.starred;
        }

        statsTotal.unread += feedStats.unread;
        statsTotal.starred += feedStats.starred;
      }

      this.stats = { feeds: statsFeeds, folders: statsFolders, total: statsTotal };

      const unread = this.stats.total.unread;
      document.title = TITLE + (unread ? ` (${unread})` : "");
    },
    // navigation helper, navigate relative to selected item
    navigateToItem(relativePosition: number) {
      let vm = this;
      if (this.itemSelected == null) {
        // if no item is selected, select first
        if (this.items.length !== 0) this.itemSelected = this.items[0].id;
        return;
      }

      var itemPosition = this.items.findIndex(x => {
        return x.id === this.itemSelected;
      });
      if (itemPosition === -1) {
        if (this.items.length !== 0) this.itemSelected = this.items[0].id;
        return;
      }

      var newPosition = itemPosition + relativePosition;
      if (newPosition < 0 || newPosition >= this.items.length) return;

      this.itemSelected = this.items[newPosition].id;

      this.$nextTick(() => {
        var scroll = document.querySelector("#item-list-scroll");
        var handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);

        this.loadMoreItems();
      });
    },
    // navigation helper, navigate relative to selected feed
    navigateToFeed(relativePosition: number) {
      const navigationList: string[] = [];
      for (const node of this.feedTree) {
        if (node.type === "folder") {
          navigationList.push("folder:" + node.folder.id);
          if (node.folder.is_expanded) {
            for (const feedNode of node.feeds) {
              navigationList.push("feed:" + feedNode.feed.id);
            }
          }
        } else {
          navigationList.push("feed:" + node.feed.id);
        }
      }
      navigationList.unshift("");

      var currentFeedPosition = navigationList.indexOf(this.feedSelected);

      if (currentFeedPosition == -1) {
        this.feedSelected = "";
        return;
      }

      var newPosition = currentFeedPosition + relativePosition;
      if (newPosition < 0 || newPosition >= navigationList.length) return;

      this.feedSelected = navigationList[newPosition];

      this.$nextTick(() => {
        var scroll = document.querySelector("#feed-list-scroll");
        var handle = scroll?.querySelector('[aria-checked="true"]');
        if (handle && scroll) scrollto(handle, scroll);
      });
    },
    changeRefreshRate(offset: number) {
      const curIdx = this.refreshRateOptions.findIndex(o => o.value === this.refreshRate);
      if (curIdx <= 0 && offset < 0) return;
      if (curIdx >= this.refreshRateOptions.length - 1 && offset > 0) return;
      this.refreshRate = this.refreshRateOptions[curIdx + offset].value;
    },
    mustHideFolder(folder: Folder): boolean {
      return !!(
        this.filterSelected &&
        !(this.current?.folder?.id === folder.id || this.current?.feed?.folder_id == folder.id) &&
        !this.stats.folders[folder.id]?.[this.filterSelected] &&
        (!this.itemSelectedDetails ||
          (this.feedsById[this.itemSelectedDetails.feed_id] || {}).folder_id != folder.id)
      );
    },
    mustHideFeed(feed: Feed): boolean {
      return !!(
        this.filterSelected &&
        !(this.current?.feed?.id === feed.id) &&
        !this.stats.feeds[feed.id]?.[this.filterSelected] &&
        (!this.itemSelectedDetails || this.itemSelectedDetails.feed_id != feed.id)
      );
    },
    async changeLanguage(lang: Lang) {
      this.$t.set(lang);
      this.language = lang;
      const [err] = await to(api.settings.update({ language: lang }));
      if (err) {
        this.$refs.toast.addToast(
          { title: this.$t("fail_save_settings"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
      }
    },
    errDescription(err: unknown): string | undefined {
      if (err instanceof HTTPError)
        return this.$t("error_server", { code: err.status, text: err.statusText });
      if (err instanceof NetworkError) return this.$t("error_network");
      return undefined;
    },
  },
});
</script>
