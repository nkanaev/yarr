import i18n from "./i18n";
import api from "./api";
import template from "./templates/app.html" with { type: "text" };
import icons from "./icons";
import { setupKeybindings } from "./key";
import { scrollto, debounce, dateRepr } from "./utils";
import drag from "./components/drag";
import dropdown from "./components/dropdown";
import modal from "./components/modal";
import relativeTime from "./components/relative-time";
import icon from "./components/icon";
import scrollDir from "./directives/scroll";
import focusDir from "./directives/focus";
import { defineComponent } from "vue";

var app = window.app;

type Theme = "system" | "light" | "sepia" | "night";

var TITLE = document.title;

export default defineComponent({
  template: template,
  components: {
    "v-drag": drag,
    "v-dropdown": dropdown,
    "v-modal": modal,
    "v-relative-time": relativeTime,
    "v-icon": icon,
  },
  directives: {
    scroll: scrollDir,
    focus: focusDir,
  },
  created() {
    this.refreshStats()
      .then(() => this.refreshFeeds())
      .then(() => this.refreshItems(false));

    api.feeds.list_errors().then((errors) => {
      this.feed_errors = errors;
    });
    this.updateMetaTheme(app.settings.theme_name);
    this.$setLang(app.settings.language);

    // keep the theme-color meta tag in sync when the OS color scheme changes
    if (window.matchMedia) {
      this._colorSchemeMql = window.matchMedia("(prefers-color-scheme: dark)");
      this._colorSchemeHandler = () => {
        this.updateMetaTheme(this.theme.name);
      };
      this._colorSchemeMql.addEventListener("change", this._colorSchemeHandler);
    }
  },
  beforeUnmount() {
    if (this._colorSchemeMql) {
      this._colorSchemeMql.removeEventListener(
        "change",
        this._colorSchemeHandler,
      );
    }
  },
  mounted() {
    setupKeybindings(this);
  },
  data() {
    var s = app.settings;
    return {
      filterSelected: s.filter,
      folders: [],
      feeds: [],
      feedSelected: s.feed,
      feedListWidth: s.feed_list_width || 300,
      feedNewChoice: [],
      feedNewChoiceSelected: "",
      items: [],
      itemsHasMore: true,
      itemSelected: null,
      itemSelectedDetails: null,
      itemSelectedReadability: "",
      itemSearch: "",
      itemSortNewestFirst: s.sort_newest_first,
      itemListWidth: s.item_list_width || 300,

      filteredFeedStats: {},
      filteredFolderStats: {},
      filteredTotalStats: null,

      settings: "",
      loading: {
        feeds: 0,
        newfeed: false,
        items: false,
        readability: false,
      },
      fonts: ["", "serif", "monospace"],
      feedStats: {},
      theme: {
        name: s.theme_name,
        font: s.theme_font,
        size: s.theme_size,
      },
      themeColors: {
        night: "#0e0e0e",
        sepia: "#f4f0e5",
        light: "#fff",
      },
      refreshRate: s.refresh_rate,
      authenticated: app.authenticated,
      requiresAuth: app.requiresAuth,
      feed_errors: {},

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
      ],
    };
  },
  computed: {
    foldersWithFeeds() {
      var feedsByFolders = this.feeds.reduce(function (folders, feed) {
        if (!folders[feed.folder_id]) folders[feed.folder_id] = [feed];
        else folders[feed.folder_id].push(feed);
        return folders;
      }, {});
      var folders = this.folders.slice().map(function (folder) {
        folder.feeds = feedsByFolders[folder.id];
        return folder;
      });
      folders.push({ id: null, feeds: feedsByFolders[null] });
      return folders;
    },
    feedsById() {
      return this.feeds.reduce(function (acc, f) {
        acc[f.id] = f;
        return acc;
      }, {});
    },
    foldersById() {
      return this.folders.reduce(function (acc, f) {
        acc[f.id] = f;
        return acc;
      }, {});
    },
    current() {
      var parts = (this.feedSelected || "").split(":", 2);
      var type = parts[0];
      var guid = parts[1];

      var folder = {},
        feed = {};

      if (type == "feed") feed = this.feedsById[guid] || {};
      if (type == "folder") folder = this.foldersById[guid] || {};

      return { type: type, feed: feed, folder: folder };
    },
    searchScope() {
      void this.language;
      var type = (this.feedSelected || "").split(":", 2)[0];
      if (type == "feed")
        return (
          (this.feedsById[this.feedSelected.split(":", 2)[1]] || {}).title || ""
        );
      if (type == "folder")
        return (
          (this.foldersById[this.feedSelected.split(":", 2)[1]] || {}).title ||
          ""
        );
      if (this.filterSelected == "unread") return this.$t("all_unread");
      if (this.filterSelected == "starred") return this.$t("all_starred");
      return this.$t("all_feeds");
    },
    itemSelectedContent() {
      if (!this.itemSelected) return "";

      if (this.itemSelectedReadability) return this.itemSelectedReadability;

      return this.itemSelectedDetails.content || "";
    },
    contentImages() {
      if (!this.itemSelectedDetails) return [];
      return (this.itemSelectedDetails.media_links || []).filter(
        (l) => l.type === "image",
      );
    },
    contentAudios() {
      if (!this.itemSelectedDetails) return [];
      return (this.itemSelectedDetails.media_links || []).filter(
        (l) => l.type === "audio",
      );
    },
    contentVideos() {
      if (!this.itemSelectedDetails) return [];
      return (this.itemSelectedDetails.media_links || []).filter(
        (l) => l.type === "video",
      );
    },
    refreshRateTitle() {
      const entry = this.refreshRateOptions.find(
        (o) => o.value === this.refreshRate,
      );
      return entry ? entry.title : "0";
    },
  },
  watch: {
    theme: {
      deep: true,
      handler(theme) {
        this.updateMetaTheme(theme.name);
        document.body.classList.value = "theme-" + theme.name;
        api.settings.update({
          theme_name: theme.name,
          theme_font: theme.font,
          theme_size: theme.size,
        });
      },
    },
    feedStats: {
      deep: true,
      handler: debounce(function () {
        var title = TITLE;
        var unreadCount = Object.values(this.feedStats).reduce(function (
          acc,
          stat,
        ) {
          return acc + stat.unread;
        }, 0);
        if (unreadCount) {
          title += " (" + unreadCount + ")";
        }
        document.title = title;
        this.computeStats();
      }, 500),
    },
    filterSelected(newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      this.itemSelected = null;
      this.items = [];
      this.itemsHasMore = true;
      api.settings
        .update({ filter: newVal })
        .then(() => this.refreshItems(false));
      this.computeStats();
    },
    feedSelected(newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      this.itemSelected = null;
      this.items = [];
      this.itemsHasMore = true;
      api.settings
        .update({ feed: newVal })
        .then(() => this.refreshItems(false));
      if (this.$refs.itemlist) this.$refs.itemlist.scrollTop = 0;
    },
    itemSelected(newVal, oldVal) {
      this.itemSelectedReadability = "";
      if (newVal === null) {
        this.itemSelectedDetails = null;
        return;
      }
      if (this.$refs.content) this.$refs.content.scrollTop = 0;

      api.items.get(newVal).then(
        (item) => {
          this.itemSelectedDetails = item;
          if (this.itemSelectedDetails.status == "unread") {
            api.items
              .update(this.itemSelectedDetails.id, { status: "read" })
              .then(
                () => {
                  this.feedStats[this.itemSelectedDetails.feed_id].unread -= 1;
                  var itemInList = this.items.find(function (i) {
                    return i.id == item.id;
                  });
                  if (itemInList) itemInList.status = "read";
                  this.itemSelectedDetails.status = "read";
                },
              );
          }
        },
      );
    },
    itemSearch: debounce(function (newVal) {
      this.refreshItems();
    }, 500),
    itemSortNewestFirst(newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings
        .update({ sort_newest_first: newVal })
        .then(() => this.refreshItems(false));
    },
    feedListWidth: debounce(function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ feed_list_width: newVal });
    }, 1000),
    itemListWidth: debounce(function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ item_list_width: newVal });
    }, 1000),
    refreshRate(newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ refresh_rate: newVal });
    },
  },
  methods: {
    updateMetaTheme(theme: Theme) {
      if (theme == "system") {
        var dark =
          window.matchMedia &&
          window.matchMedia("(prefers-color-scheme: dark)").matches;
        theme = dark ? "night" : "light";
      }
      document.querySelector("meta[name='theme-color']").content =
        this.themeColors[theme];
    },
    refreshStats(loopMode?: boolean) {
      return api.status().then((data) => {
        if (loopMode && !this.itemSelected) this.refreshItems();

        this.loading.feeds = data.running;
        if (data.running) {
          setTimeout(() => this.refreshStats(true), 500);
        }
        this.feedStats = data.stats.reduce((acc, stat) => {
          acc[stat.feed_id] = stat;
          return acc;
        }, {});

        api.feeds.list_errors().then((errors) => {
          this.feed_errors = errors;
        });
      });
    },
    getItemsQuery() {
      var query = {};
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
    refreshFeeds() {
      return Promise.all([api.folders.list(), api.feeds.list()]).then(
        (values) => {
          this.folders = values[0];
          this.feeds = values[1];
        },
      );
    },
    refreshItems(loadMore = false) {
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
      return api.items.list(query).then((data) => {
        if (loadMore) {
          this.items = this.items.concat(data.list);
        } else {
          this.items = data.list;
        }
        this.itemsHasMore = data.has_more;
        this.loading.items = false;

        // load more if there's some space left at the bottom of the item list.
        this.$nextTick(() => {
          if (
            this.itemsHasMore &&
            !this.loading.items &&
            this.itemListCloseToBottom()
          ) {
            this.refreshItems(true);
          }
        });
      });
    },
    itemListCloseToBottom() {
      // approx. vertical space at the bottom of the list (loading el & paddings) when 1rem = 16px
      var bottomSpace = 70;
      var scale =
        (parseFloat(getComputedStyle(document.documentElement).fontSize) ||
          16) / 16;

      var el = this.$refs.itemlist as HTMLElement;

      if (!el || el.scrollHeight === 0) return false; // element is invisible (responsive design)

      var closeToBottom =
        el.scrollHeight - el.scrollTop - el.offsetHeight < bottomSpace * scale;
      return closeToBottom;
    },
    loadMoreItems() {
      if (!this.itemsHasMore) return;
      if (this.loading.items) return;
      if (this.itemListCloseToBottom()) return this.refreshItems(true);
      if (
        this.itemSelected &&
        this.itemSelected === this.items[this.items.length - 1].id
      )
        return this.refreshItems(true);
    },
    markItemsRead() {
      var query = this.getItemsQuery();
      api.items.mark_read(query).then(() => {
        this.items = [];
        this.itemsPage = { cur: 1, num: 1 };
        this.itemSelected = null;
        this.itemsHasMore = false;
        this.refreshStats();
      });
    },
    toggleFolderExpanded(folder) {
      folder.is_expanded = !folder.is_expanded;
      api.folders.update(folder.id, { is_expanded: folder.is_expanded });
    },
    formatDate(datestr: string) {
      var options = {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      };
      return new Date(datestr).toLocaleDateString(undefined, options);
    },
    moveFeed(feed, folder) {
      var folder_id = folder ? folder.id : null;
      api.feeds.update(feed.id, { folder_id: folder_id }).then(() => {
        feed.folder_id = folder_id;
        this.refreshStats();
      });
    },
    moveFeedToNewFolder(feed) {
      var title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;
      api.folders.create({ title: title }).then((folder) => {
        api.feeds.update(feed.id, { folder_id: folder.id }).then(() => {
          this.refreshFeeds().then(() => {
            this.refreshStats();
          });
        });
      });
    },
    createNewFeedFolder() {
      var title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;
      api.folders.create({ title: title }).then((result) => {
        this.refreshFeeds().then(() => {
          this.$nextTick(() => {
            if (this.$refs.newFeedFolder) {
              this.$refs.newFeedFolder.value = result.id;
            }
          });
        });
      });
    },
    renameFolder(folder) {
      var newTitle = prompt(this.$t("prompt_new_title"), folder.title);
      if (newTitle) {
        api.folders.update(folder.id, { title: newTitle }).then(
          () => {
            folder.title = newTitle;
            this.folders.sort(function (a, b) {
              return a.title.localeCompare(b.title);
            });
          },
        );
      }
    },
    deleteFolder(folder) {
      if (confirm(this.$t("confirm_delete", { name: folder.title }))) {
        api.folders.delete(folder.id).then(() => {
          this.feedSelected = null;
          this.refreshStats();
          this.refreshFeeds();
        });
      }
    },
    updateFeedLink(feed) {
      var newLink = prompt(this.$t("prompt_feed_link"), feed.feed_link);
      if (newLink) {
        api.feeds.update(feed.id, { feed_link: newLink }).then(function () {
          feed.feed_link = newLink;
        });
      }
    },
    renameFeed(feed) {
      var newTitle = prompt(this.$t("prompt_new_title"), feed.title);
      if (newTitle) {
        api.feeds.update(feed.id, { title: newTitle }).then(function () {
          feed.title = newTitle;
        });
      }
    },
    deleteFeed(feed) {
      if (confirm(this.$t("confirm_delete", { name: feed.title }))) {
        api.feeds.delete(feed.id).then(() => {
          this.feedSelected = null;
          this.refreshStats();
          this.refreshFeeds();
        });
      }
    },
    createFeed($event) {
      var form = $event.target;
      var data = {
        url: form.querySelector("input[name=url]").value,
        folder_id:
          parseInt(form.querySelector("select[name=folder_id]").value) || null,
      };
      if (this.feedNewChoiceSelected) {
        var choice = this.feedNewChoice.find(
          (c) => c.url === this.feedNewChoiceSelected,
        );
        data.url = this.feedNewChoiceSelected;
        if (choice && choice.title_override)
          data.title_override = choice.title_override;
      }
      this.loading.newfeed = true;
      api.feeds.create(data).then((result) => {
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
        this.loading.newfeed = false;
      });
    },
    toggleItemStatus(item, targetstatus, fallbackstatus) {
      var oldstatus = item.status;
      var newstatus =
        item.status !== targetstatus ? targetstatus : fallbackstatus;

      var updateStats = (status, incr) => {
        if (status == "unread" || status == "starred") {
          this.feedStats[item.feed_id][status] += incr;
        }
      };

      api.items.update(item.id, { status: newstatus }).then(
        () => {
          updateStats(oldstatus, -1);
          updateStats(newstatus, +1);

          var itemInList = this.items.find(function (i) {
            return i.id == item.id;
          });
          if (itemInList) itemInList.status = newstatus;
          item.status = newstatus;
        },
      );
    },
    toggleItemStarred(item) {
      this.toggleItemStatus(item, "starred", "read");
    },
    toggleItemRead(item) {
      this.toggleItemStatus(item, "unread", "read");
    },
    importOPML(event) {
      var input = event.target;
      var form = document.querySelector("#opml-import-form");
      this.$refs.menuDropdown.hide();
      api.upload_opml(form).then(() => {
        input.value = "";
        this.refreshFeeds();
        this.refreshStats();
      });
    },
    logout() {
      api.logout().then(() => {
        document.location.reload();
      });
    },
    toggleReadability() {
      if (this.itemSelectedReadability) {
        this.itemSelectedReadability = null;
        return;
      }
      var item = this.itemSelectedDetails;
      if (!item) return;
      if (item.link) {
        this.loading.readability = true;
        api.crawl(item.link).then((data) => {
          this.itemSelectedReadability = data && data.content;
          this.loading.readability = false;
        });
      }
    },
    showSettings(settings) {
      this.settings = settings;

      if (settings === "create") {
        this.feedNewChoice = [];
        this.feedNewChoiceSelected = "";
      }
    },
    resizeFeedList(width) {
      this.feedListWidth = Math.min(Math.max(200, width), 700);
    },
    resizeItemList(width) {
      this.itemListWidth = Math.min(Math.max(200, width), 700);
    },
    resetFeedChoice() {
      this.feedNewChoice = [];
      this.feedNewChoiceSelected = "";
    },
    incrFont(x) {
      this.theme.size = +(this.theme.size + 0.1 * x).toFixed(1);
    },
    fetchAllFeeds() {
      if (this.loading.feeds) return;
      api.feeds.refresh().then(() => {
        this.refreshStats();
      });
    },
    computeStats() {
      var filter = this.filterSelected;
      if (!filter) {
        this.filteredFeedStats = {};
        this.filteredFolderStats = {};
        this.filteredTotalStats = null;
        return;
      }

      var statsFeeds = {},
        statsFolders = {},
        statsTotal = 0;

      for (var i = 0; i < this.feeds.length; i++) {
        var feed = this.feeds[i];
        if (!this.feedStats[feed.id]) continue;

        var n = this.feedStats[feed.id][filter] || 0;

        if (!statsFolders[feed.folder_id]) statsFolders[feed.folder_id] = 0;

        statsFeeds[feed.id] = n;
        statsFolders[feed.folder_id] += n;
        statsTotal += n;
      }

      this.filteredFeedStats = statsFeeds;
      this.filteredFolderStats = statsFolders;
      this.filteredTotalStats = statsTotal;
    },
    // navigation helper, navigate relative to selected item
    navigateToItem(relativePosition: number) {
      let vm = this;
      if (this.itemSelected == null) {
        // if no item is selected, select first
        if (this.items.length !== 0) this.itemSelected = this.items[0].id;
        return;
      }

      var itemPosition = this.items.findIndex((x) => {
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

        var handle = scroll.querySelector("input[type=radio]:checked");
        var target = handle && handle.parentElement;

        if (target && scroll) scrollto(target, scroll);

        this.loadMoreItems();
      });
    },
    // navigation helper, navigate relative to selected feed
    navigateToFeed(relativePosition: number) {
      let vm = this;
      const navigationList = this.foldersWithFeeds
        .filter((folder) => !folder.id || !this.mustHideFolder(folder))
        .map((folder) => {
          if (this.mustHideFolder(folder)) return [];
          const folds = folder.id ? [`folder:${folder.id}`] : [];
          const feeds =
            folder.is_expanded || !folder.id
              ? (folder.feeds || [])
                  .filter((f) => !this.mustHideFeed(f))
                  .map((f) => `feed:${f.id}`)
              : [];
          return folds.concat(feeds);
        })
        .flat();
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

        var handle = scroll.querySelector("input[type=radio]:checked");
        var target = handle && handle.parentElement;

        if (target && scroll) scrollto(target, scroll);
      });
    },
    changeRefreshRate(offset: number) {
      const curIdx = this.refreshRateOptions.findIndex(
        (o) => o.value === this.refreshRate,
      );
      if (curIdx <= 0 && offset < 0) return;
      if (curIdx >= this.refreshRateOptions.length - 1 && offset > 0) return;
      this.refreshRate = this.refreshRateOptions[curIdx + offset].value;
    },
    mustHideFolder(folder) {
      return (
        this.filterSelected &&
        !(
          this.current.folder.id == folder.id ||
          this.current.feed.folder_id == folder.id
        ) &&
        !this.filteredFolderStats[folder.id] &&
        (!this.itemSelectedDetails ||
          (this.feedsById[this.itemSelectedDetails.feed_id] || {}).folder_id !=
            folder.id)
      );
    },
    mustHideFeed(feed) {
      return (
        this.filterSelected &&
        !(this.current.feed.id == feed.id) &&
        !this.filteredFeedStats[feed.id] &&
        (!this.itemSelectedDetails ||
          this.itemSelectedDetails.feed_id != feed.id)
      );
    },
    changeLanguage(lang) {
      this.$setLang(lang);
      this.language = lang;
      api.settings.update({ language: lang });
    },
  },
});
