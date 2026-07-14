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
var vm;

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
  created: function () {
    vm = this;
    this.refreshStats()
      .then(this.refreshFeeds.bind(this))
      .then(this.refreshItems.bind(this, false));

    api.feeds.list_errors().then(function (errors) {
      vm.feed_errors = errors;
    });
    this.updateMetaTheme(app.settings.theme_name);
    this.$setLang(app.settings.language);

    // keep the theme-color meta tag in sync when the OS color scheme changes
    if (window.matchMedia) {
      this._colorSchemeMql = window.matchMedia("(prefers-color-scheme: dark)");
      this._colorSchemeHandler = function () {
        this.updateMetaTheme(this.theme.name);
      }.bind(this);
      this._colorSchemeMql.addEventListener("change", this._colorSchemeHandler);
    }
  },
  beforeUnmount: function () {
    if (this._colorSchemeMql) {
      this._colorSchemeMql.removeEventListener(
        "change",
        this._colorSchemeHandler,
      );
    }
  },
  mounted: function () {
    setupKeybindings(this);
  },
  data: function () {
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
    foldersWithFeeds: function () {
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
    feedsById: function () {
      return this.feeds.reduce(function (acc, f) {
        acc[f.id] = f;
        return acc;
      }, {});
    },
    foldersById: function () {
      return this.folders.reduce(function (acc, f) {
        acc[f.id] = f;
        return acc;
      }, {});
    },
    current: function () {
      var parts = (this.feedSelected || "").split(":", 2);
      var type = parts[0];
      var guid = parts[1];

      var folder = {},
        feed = {};

      if (type == "feed") feed = this.feedsById[guid] || {};
      if (type == "folder") folder = this.foldersById[guid] || {};

      return { type: type, feed: feed, folder: folder };
    },
    searchScope: function () {
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
    itemSelectedContent: function () {
      if (!this.itemSelected) return "";

      if (this.itemSelectedReadability) return this.itemSelectedReadability;

      return this.itemSelectedDetails.content || "";
    },
    contentImages: function () {
      if (!this.itemSelectedDetails) return [];
      return (this.itemSelectedDetails.media_links || []).filter(
        (l) => l.type === "image",
      );
    },
    contentAudios: function () {
      if (!this.itemSelectedDetails) return [];
      return (this.itemSelectedDetails.media_links || []).filter(
        (l) => l.type === "audio",
      );
    },
    contentVideos: function () {
      if (!this.itemSelectedDetails) return [];
      return (this.itemSelectedDetails.media_links || []).filter(
        (l) => l.type === "video",
      );
    },
    refreshRateTitle: function () {
      const entry = this.refreshRateOptions.find(
        (o) => o.value === this.refreshRate,
      );
      return entry ? entry.title : "0";
    },
  },
  watch: {
    theme: {
      deep: true,
      handler: function (theme) {
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
    filterSelected: function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      this.itemSelected = null;
      this.items = [];
      this.itemsHasMore = true;
      api.settings
        .update({ filter: newVal })
        .then(this.refreshItems.bind(this, false));
      this.computeStats();
    },
    feedSelected: function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      this.itemSelected = null;
      this.items = [];
      this.itemsHasMore = true;
      api.settings
        .update({ feed: newVal })
        .then(this.refreshItems.bind(this, false));
      if (this.$refs.itemlist) this.$refs.itemlist.scrollTop = 0;
    },
    itemSelected: function (newVal, oldVal) {
      this.itemSelectedReadability = "";
      if (newVal === null) {
        this.itemSelectedDetails = null;
        return;
      }
      if (this.$refs.content) this.$refs.content.scrollTop = 0;

      api.items.get(newVal).then(
        function (item) {
          this.itemSelectedDetails = item;
          if (this.itemSelectedDetails.status == "unread") {
            api.items
              .update(this.itemSelectedDetails.id, { status: "read" })
              .then(
                function () {
                  this.feedStats[this.itemSelectedDetails.feed_id].unread -= 1;
                  var itemInList = this.items.find(function (i) {
                    return i.id == item.id;
                  });
                  if (itemInList) itemInList.status = "read";
                  this.itemSelectedDetails.status = "read";
                }.bind(this),
              );
          }
        }.bind(this),
      );
    },
    itemSearch: debounce(function (newVal) {
      this.refreshItems();
    }, 500),
    itemSortNewestFirst: function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings
        .update({ sort_newest_first: newVal })
        .then(vm.refreshItems.bind(this, false));
    },
    feedListWidth: debounce(function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ feed_list_width: newVal });
    }, 1000),
    itemListWidth: debounce(function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ item_list_width: newVal });
    }, 1000),
    refreshRate: function (newVal, oldVal) {
      if (oldVal === undefined) return; // do nothing, initial setup
      api.settings.update({ refresh_rate: newVal });
    },
  },
  methods: {
    updateMetaTheme: function (theme) {
      if (theme == "system") {
        var dark =
          window.matchMedia &&
          window.matchMedia("(prefers-color-scheme: dark)").matches;
        theme = dark ? "night" : "light";
      }
      document.querySelector("meta[name='theme-color']").content =
        this.themeColors[theme];
    },
    refreshStats: function (loopMode) {
      return api.status().then(function (data) {
        if (loopMode && !vm.itemSelected) vm.refreshItems();

        vm.loading.feeds = data.running;
        if (data.running) {
          setTimeout(vm.refreshStats.bind(vm, true), 500);
        }
        vm.feedStats = data.stats.reduce(function (acc, stat) {
          acc[stat.feed_id] = stat;
          return acc;
        }, {});

        api.feeds.list_errors().then(function (errors) {
          vm.feed_errors = errors;
        });
      });
    },
    getItemsQuery: function () {
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
    refreshFeeds: function () {
      return Promise.all([api.folders.list(), api.feeds.list()]).then(
        function (values) {
          vm.folders = values[0];
          vm.feeds = values[1];
        },
      );
    },
    refreshItems: function (loadMore = false) {
      if (this.feedSelected === null) {
        vm.items = [];
        vm.itemsHasMore = false;
        return;
      }

      var query = this.getItemsQuery();
      if (loadMore) {
        query.after = vm.items[vm.items.length - 1].id;
      }

      this.loading.items = true;
      return api.items.list(query).then(function (data) {
        if (loadMore) {
          vm.items = vm.items.concat(data.list);
        } else {
          vm.items = data.list;
        }
        vm.itemsHasMore = data.has_more;
        vm.loading.items = false;

        // load more if there's some space left at the bottom of the item list.
        vm.$nextTick(function () {
          if (
            vm.itemsHasMore &&
            !vm.loading.items &&
            vm.itemListCloseToBottom()
          ) {
            vm.refreshItems(true);
          }
        });
      });
    },
    itemListCloseToBottom: function () {
      // approx. vertical space at the bottom of the list (loading el & paddings) when 1rem = 16px
      var bottomSpace = 70;
      var scale =
        (parseFloat(getComputedStyle(document.documentElement).fontSize) ||
          16) / 16;

      var el = this.$refs.itemlist;

      if (el.scrollHeight === 0) return false; // element is invisible (responsive design)

      var closeToBottom =
        el.scrollHeight - el.scrollTop - el.offsetHeight < bottomSpace * scale;
      return closeToBottom;
    },
    loadMoreItems: function (event, el) {
      if (!this.itemsHasMore) return;
      if (this.loading.items) return;
      if (this.itemListCloseToBottom()) return this.refreshItems(true);
      if (
        this.itemSelected &&
        this.itemSelected === this.items[this.items.length - 1].id
      )
        return this.refreshItems(true);
    },
    markItemsRead: function () {
      var query = this.getItemsQuery();
      api.items.mark_read(query).then(function () {
        vm.items = [];
        vm.itemsPage = { cur: 1, num: 1 };
        vm.itemSelected = null;
        vm.itemsHasMore = false;
        vm.refreshStats();
      });
    },
    toggleFolderExpanded: function (folder) {
      folder.is_expanded = !folder.is_expanded;
      api.folders.update(folder.id, { is_expanded: folder.is_expanded });
    },
    formatDate: function (datestr) {
      var options = {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      };
      return new Date(datestr).toLocaleDateString(undefined, options);
    },
    moveFeed: function (feed, folder) {
      var folder_id = folder ? folder.id : null;
      api.feeds.update(feed.id, { folder_id: folder_id }).then(function () {
        feed.folder_id = folder_id;
        vm.refreshStats();
      });
    },
    moveFeedToNewFolder: function (feed) {
      var title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;
      api.folders.create({ title: title }).then(function (folder) {
        api.feeds.update(feed.id, { folder_id: folder.id }).then(function () {
          vm.refreshFeeds().then(function () {
            vm.refreshStats();
          });
        });
      });
    },
    createNewFeedFolder: function () {
      var title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;
      api.folders.create({ title: title }).then(function (result) {
        vm.refreshFeeds().then(function () {
          vm.$nextTick(function () {
            if (vm.$refs.newFeedFolder) {
              vm.$refs.newFeedFolder.value = result.id;
            }
          });
        });
      });
    },
    renameFolder: function (folder) {
      var newTitle = prompt(this.$t("prompt_new_title"), folder.title);
      if (newTitle) {
        api.folders.update(folder.id, { title: newTitle }).then(
          function () {
            folder.title = newTitle;
            this.folders.sort(function (a, b) {
              return a.title.localeCompare(b.title);
            });
          }.bind(this),
        );
      }
    },
    deleteFolder: function (folder) {
      if (confirm(this.$t("confirm_delete", { name: folder.title }))) {
        api.folders.delete(folder.id).then(function () {
          vm.feedSelected = null;
          vm.refreshStats();
          vm.refreshFeeds();
        });
      }
    },
    updateFeedLink: function (feed) {
      var newLink = prompt(this.$t("prompt_feed_link"), feed.feed_link);
      if (newLink) {
        api.feeds.update(feed.id, { feed_link: newLink }).then(function () {
          feed.feed_link = newLink;
        });
      }
    },
    renameFeed: function (feed) {
      var newTitle = prompt(this.$t("prompt_new_title"), feed.title);
      if (newTitle) {
        api.feeds.update(feed.id, { title: newTitle }).then(function () {
          feed.title = newTitle;
        });
      }
    },
    deleteFeed: function (feed) {
      if (confirm(this.$t("confirm_delete", { name: feed.title }))) {
        api.feeds.delete(feed.id).then(function () {
          vm.feedSelected = null;
          vm.refreshStats();
          vm.refreshFeeds();
        });
      }
    },
    createFeed: function ($event) {
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
      api.feeds.create(data).then(function (result) {
        if (result.status === "success") {
          vm.refreshFeeds();
          vm.refreshStats();
          vm.settings = "";
          vm.feedSelected = "feed:" + result.feed.id;
        } else if (result.status === "multiple") {
          vm.feedNewChoice = result.choice;
          vm.feedNewChoiceSelected = result.choice[0].url;
        } else {
          alert("No feeds found at the given url.");
        }
        vm.loading.newfeed = false;
      });
    },
    toggleItemStatus: function (item, targetstatus, fallbackstatus) {
      var oldstatus = item.status;
      var newstatus =
        item.status !== targetstatus ? targetstatus : fallbackstatus;

      var updateStats = function (status, incr) {
        if (status == "unread" || status == "starred") {
          this.feedStats[item.feed_id][status] += incr;
        }
      }.bind(this);

      api.items.update(item.id, { status: newstatus }).then(
        function () {
          updateStats(oldstatus, -1);
          updateStats(newstatus, +1);

          var itemInList = this.items.find(function (i) {
            return i.id == item.id;
          });
          if (itemInList) itemInList.status = newstatus;
          item.status = newstatus;
        }.bind(this),
      );
    },
    toggleItemStarred: function (item) {
      this.toggleItemStatus(item, "starred", "read");
    },
    toggleItemRead: function (item) {
      this.toggleItemStatus(item, "unread", "read");
    },
    importOPML: function (event) {
      var input = event.target;
      var form = document.querySelector("#opml-import-form");
      this.$refs.menuDropdown.hide();
      api.upload_opml(form).then(function () {
        input.value = "";
        vm.refreshFeeds();
        vm.refreshStats();
      });
    },
    logout: function () {
      api.logout().then(function () {
        document.location.reload();
      });
    },
    toggleReadability: function () {
      if (this.itemSelectedReadability) {
        this.itemSelectedReadability = null;
        return;
      }
      var item = this.itemSelectedDetails;
      if (!item) return;
      if (item.link) {
        this.loading.readability = true;
        api.crawl(item.link).then(function (data) {
          vm.itemSelectedReadability = data && data.content;
          vm.loading.readability = false;
        });
      }
    },
    showSettings: function (settings) {
      this.settings = settings;

      if (settings === "create") {
        vm.feedNewChoice = [];
        vm.feedNewChoiceSelected = "";
      }
    },
    resizeFeedList: function (width) {
      this.feedListWidth = Math.min(Math.max(200, width), 700);
    },
    resizeItemList: function (width) {
      this.itemListWidth = Math.min(Math.max(200, width), 700);
    },
    resetFeedChoice: function () {
      this.feedNewChoice = [];
      this.feedNewChoiceSelected = "";
    },
    incrFont: function (x) {
      this.theme.size = +(this.theme.size + 0.1 * x).toFixed(1);
    },
    fetchAllFeeds: function () {
      if (this.loading.feeds) return;
      api.feeds.refresh().then(function () {
        vm.refreshStats();
      });
    },
    computeStats: function () {
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

        var n = vm.feedStats[feed.id][filter] || 0;

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
    navigateToItem: function (relativePosition) {
      let vm = this;
      if (vm.itemSelected == null) {
        // if no item is selected, select first
        if (vm.items.length !== 0) vm.itemSelected = vm.items[0].id;
        return;
      }

      var itemPosition = vm.items.findIndex(function (x) {
        return x.id === vm.itemSelected;
      });
      if (itemPosition === -1) {
        if (vm.items.length !== 0) vm.itemSelected = vm.items[0].id;
        return;
      }

      var newPosition = itemPosition + relativePosition;
      if (newPosition < 0 || newPosition >= vm.items.length) return;

      vm.itemSelected = vm.items[newPosition].id;

      vm.$nextTick(function () {
        var scroll = document.querySelector("#item-list-scroll");

        var handle = scroll.querySelector("input[type=radio]:checked");
        var target = handle && handle.parentElement;

        if (target && scroll) scrollto(target, scroll);

        vm.loadMoreItems();
      });
    },
    // navigation helper, navigate relative to selected feed
    navigateToFeed: function (relativePosition) {
      let vm = this;
      const navigationList = this.foldersWithFeeds
        .filter((folder) => !folder.id || !vm.mustHideFolder(folder))
        .map((folder) => {
          if (this.mustHideFolder(folder)) return [];
          const folds = folder.id ? [`folder:${folder.id}`] : [];
          const feeds =
            folder.is_expanded || !folder.id
              ? (folder.feeds || [])
                  .filter((f) => !vm.mustHideFeed(f))
                  .map((f) => `feed:${f.id}`)
              : [];
          return folds.concat(feeds);
        })
        .flat();
      navigationList.unshift("");

      var currentFeedPosition = navigationList.indexOf(vm.feedSelected);

      if (currentFeedPosition == -1) {
        vm.feedSelected = "";
        return;
      }

      var newPosition = currentFeedPosition + relativePosition;
      if (newPosition < 0 || newPosition >= navigationList.length) return;

      vm.feedSelected = navigationList[newPosition];

      vm.$nextTick(function () {
        var scroll = document.querySelector("#feed-list-scroll");

        var handle = scroll.querySelector("input[type=radio]:checked");
        var target = handle && handle.parentElement;

        if (target && scroll) scrollto(target, scroll);
      });
    },
    changeRefreshRate: function (offset) {
      const curIdx = this.refreshRateOptions.findIndex(
        (o) => o.value === this.refreshRate,
      );
      if (curIdx <= 0 && offset < 0) return;
      if (curIdx >= this.refreshRateOptions.length - 1 && offset > 0) return;
      this.refreshRate = this.refreshRateOptions[curIdx + offset].value;
    },
    mustHideFolder: function (folder) {
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
    mustHideFeed: function (feed) {
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
