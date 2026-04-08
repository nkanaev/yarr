'use strict';

var TITLE = document.title

function scrollto(target, scroll) {
  var padding = 10
  var targetRect = target.getBoundingClientRect()
  var scrollRect = scroll.getBoundingClientRect()

  // target
  var relativeOffset = targetRect.y - scrollRect.y
  var absoluteOffset = relativeOffset + scroll.scrollTop

  if (padding <= relativeOffset && relativeOffset + targetRect.height <= scrollRect.height - padding) return

  var newPos = scroll.scrollTop
  if (relativeOffset < padding) {
    newPos = absoluteOffset - padding
  } else {
    newPos = absoluteOffset - scrollRect.height + targetRect.height + padding
  }
  scroll.scrollTop = Math.round(newPos)
}

var debounce = function(callback, wait) {
  var timeout
  return function() {
    var ctx = this, args = arguments
    clearTimeout(timeout)
    timeout = setTimeout(function() {
      callback.apply(ctx, args)
    }, wait)
  }
}

Vue.directive('scroll', {
  inserted: function(el, binding) {
    el.addEventListener('scroll', debounce(function(event) {
      binding.value(event, el)
    }, 200))
  },
})

Vue.directive('focus', {
  inserted: function(el) {
    el.focus()
  }
})

Vue.component('drag', {
  props: ['width'],
  template: '<div class="drag"></div>',
  mounted: function() {
    var self = this
    var startX = undefined
    var initW = undefined
    var onMouseMove = function(e) {
      var offset = e.clientX - startX
      var newWidth = initW + offset
      self.$emit('resize', newWidth)
    }
    var onMouseUp = function(e) {
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }
    this.$el.addEventListener('mousedown', function(e) {
      startX = e.clientX
      initW = self.width
      document.addEventListener('mousemove', onMouseMove)
      document.addEventListener('mouseup', onMouseUp)
    })
  },
})

Vue.component('dropdown', {
  props: ['class', 'toggle-class', 'ref', 'drop', 'title'],
  data: function() {
    return {open: false}
  },
  template: `
    <div class="dropdown" :class="$attrs.class">
      <button ref="btn" @click="toggle" :class="btnToggleClass" :title="$props.title"><slot name="button"></slot></button>
      <div ref="menu" class="dropdown-menu" :class="{show: open}"><slot v-if="open"></slot></div>
    </div>
  `,
  computed: {
    btnToggleClass: function() {
      var c = this.$props.toggleClass || ''
      c += ' dropdown-toggle dropdown-toggle-no-caret'
      c += this.open ? ' show' : ''
      return c.trim()
    }
  },
  methods: {
    toggle: function(e) {
      this.open ? this.hide() : this.show()
    },
    show: function(e) {
      this.open = true
      this.$refs.menu.style.top = this.$refs.btn.offsetHeight + 'px'
      var drop = this.$props.drop

      if (drop === 'right') {
        this.$refs.menu.style.left = 'auto'
        this.$refs.menu.style.right = '0'
      } else if (drop === 'center') {
        this.$nextTick(function() {
          var btnWidth = this.$refs.btn.getBoundingClientRect().width
          var menuWidth = this.$refs.menu.getBoundingClientRect().width
          this.$refs.menu.style.left = '-' + ((menuWidth - btnWidth) / 2) + 'px'
        }.bind(this))
      }

      document.addEventListener('click', this.clickHandler)
    },
    hide: function() {
      this.open = false
      document.removeEventListener('click', this.clickHandler)
    },
    clickHandler: function(e) {
      var dropdown = e.target.closest('.dropdown')
      if (dropdown == null || dropdown != this.$el) return this.hide()
      if (e.target.closest('.dropdown-item') != null) return this.hide()
    }
  },
})

Vue.component('modal', {
  props: ['open'],
  template: `
    <div class="modal custom-modal" tabindex="-1" v-if="$props.open">
      <div class="modal-dialog">
        <div class="modal-content" ref="content">
          <div class="modal-body">
            <slot v-if="$props.open"></slot>
          </div>
        </div>
      </div>
    </div>
  `,
  data: function() {
    return {opening: false}
  },
  watch: {
    'open': function(newVal) {
      if (newVal) {
        this.opening = true
        document.addEventListener('click', this.handleClick)
      } else {
        document.removeEventListener('click', this.handleClick)
      }
    },
  },
  methods: {
    handleClick: function(e) {
      if (this.opening) {
        this.opening = false
        return
      }
      if (e.target.closest('.modal-content') == null) this.$emit('hide')
    },
  },
})

function dateRepr(d) {
  var sec = (new Date().getTime() - d.getTime()) / 1000
  var neg = sec < 0
  var out = ''

  sec = Math.abs(sec)
  if (sec < 2700)  // less than 45 minutes
    out = Math.round(sec / 60) + 'm'
  else if (sec < 86400)  // less than 24 hours
    out = Math.round(sec / 3600) + 'h'
  else if (sec < 604800)  // less than a week
    out = Math.round(sec / 86400) + 'd'
  else
    out = d.toLocaleDateString(undefined, {year: "numeric", month: "long", day: "numeric"})

  if (neg) return '-' + out
  return out
}

Vue.component('relative-time', {
  props: ['val'],
  data: function() {
    var d = new Date(this.val)
    return {
      'date': d,
      'formatted': dateRepr(d),
      'interval': null,
    }
  },
  template: '<time :datetime="val">{{ formatted }}</time>',
  mounted: function() {
    this.interval = setInterval(function() {
      this.formatted = dateRepr(this.date)
    }.bind(this), 600000)  // every 10 minutes
  },
  destroyed: function() {
    clearInterval(this.interval)
  },
})

var vm = new Vue({
  created: function() {
    // AI-related instance state (not reactive, no need for Vue to track)
    this._chatHistory = []
    this._chatReader = null
    this._briefingReader = null
    this._aiPollTimer = null
    this._lastAiDetail = ''

    this.refreshStats()
      .then(this.refreshFeeds.bind(this))
      .then(this.refreshItems.bind(this, false))

    api.feeds.list_errors().then(function(errors) {
      vm.feed_errors = errors
    })
    this.updateMetaTheme(app.settings.theme_name)
  },
  mounted: function() {
    if (this.aiEnabled) {
      this.checkTaskStatus()
    }
  },
  data: function() {
    var s = app.settings
    return {
      'filterSelected': s.filter,
      'folders': [],
      'feeds': [],
      'feedSelected': s.feed,
      'feedListWidth': s.feed_list_width || 300,
      'feedNewChoice': [],
      'feedNewChoiceSelected': '',
      'items': [],
      'itemsHasMore': true,
      'itemSelected': null,
      'itemSelectedDetails': null,
      'itemSelectedReadability': '',
      'itemSearch': '',
      'itemSortNewestFirst': s.sort_newest_first,
      'itemListWidth': s.item_list_width || 300,

      'filteredFeedStats': {},
      'filteredFolderStats': {},
      'filteredTotalStats': null,

      'settings': '',
      'loading': {
        'feeds': 0,
        'newfeed': false,
        'items': false,
        'readability': false,
      },
      'fonts': ['', 'serif', 'monospace'],
      'feedStats': {},
      'theme': {
        'name': s.theme_name,
        'font': s.theme_font,
        'size': s.theme_size,
      },
      'themeColors': {
        'night': '#0e0e0e',
        'sepia': '#f4f0e5',
        'light': '#fff',
      },
      'refreshRate': s.refresh_rate,
      'authenticated': app.authenticated,
      'feed_errors': {},

      'refreshRateOptions': [
        { title: "0", value: 0 },
        { title: "10m", value: 10 },
        { title: "30m", value: 30 },
        { title: "1h", value: 60 },
        { title: "2h", value: 120 },
        { title: "4h", value: 240 },
        { title: "12h", value: 720 },
        { title: "24h", value: 1440 },
      ],

      // AI feature flags
      'aiEnabled': !!(app.aiEnabled),

      // Density toggle
      'densityCompact': false,

      // Reader/focus mode
      'readerMode': false,

      // Topics/clusters
      'topicsActive': false,
      'topicsLoading': false,
      'topicsLoaded': false,
      'topicClusters': [],
      'topicTags': [],

      'topicsHealth': {},
      'selectedTopic': null,

      // AI task status
      'aiTaskActive': false,
      'aiTaskText': '',
    }
  },
  computed: {
    foldersWithFeeds: function() {
      var feedsByFolders = this.feeds.reduce(function(folders, feed) {
        if (!folders[feed.folder_id])
          folders[feed.folder_id] = [feed]
        else
          folders[feed.folder_id].push(feed)
        return folders
      }, {})
      var folders = this.folders.slice().map(function(folder) {
        folder.feeds = feedsByFolders[folder.id]
        return folder
      })
      folders.push({id: null, feeds: feedsByFolders[null]})
      return folders
    },
    feedsById: function() {
      return this.feeds.reduce(function(acc, f) { acc[f.id] = f; return acc }, {})
    },
    foldersById: function() {
      return this.folders.reduce(function(acc, f) { acc[f.id] = f; return acc }, {})
    },
    current: function() {
      var parts = (this.feedSelected || '').split(':', 2)
      var type = parts[0]
      var guid = parts[1]

      var folder = {}, feed = {}

      if (type == 'feed')
        feed = this.feedsById[guid] || {}
      if (type == 'folder')
        folder = this.foldersById[guid] || {}

      return {type: type, feed: feed, folder: folder}
    },
    itemSelectedContent: function() {
      if (!this.itemSelected) return ''

      if (this.itemSelectedReadability)
        return this.itemSelectedReadability

      return this.itemSelectedDetails.content || ''
    },
    contentImages: function() {
      if (!this.itemSelectedDetails) return []
      return (this.itemSelectedDetails.media_links || []).filter(l => l.type === 'image')
    },
    contentAudios: function() {
      if (!this.itemSelectedDetails) return []
      return (this.itemSelectedDetails.media_links || []).filter(l => l.type === 'audio')
    },
    contentVideos: function() {
      if (!this.itemSelectedDetails) return []
      return (this.itemSelectedDetails.media_links || []).filter(l => l.type === 'video')
    },
    refreshRateTitle: function () {
      const entry = this.refreshRateOptions.find(o => o.value === this.refreshRate)
      return entry ? entry.title : '0'
    },
  },
  watch: {
    'theme': {
      deep: true,
      handler: function(theme) {
        this.updateMetaTheme(theme.name)
        document.body.classList.value = 'theme-' + theme.name
        api.settings.update({
          theme_name: theme.name,
          theme_font: theme.font,
          theme_size: theme.size,
        })
      },
    },
    'feedStats': {
      deep: true,
      handler: debounce(function() {
        var title = TITLE
        var unreadCount = Object.values(this.feedStats).reduce(function(acc, stat) {
          return acc + stat.unread
        }, 0)
        if (unreadCount) {
          title += ' ('+unreadCount+')'
        }
        document.title = title
        this.computeStats()
      }, 500),
    },
    'filterSelected': function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      this.itemSelected = null
      this.items = []
      this.itemsHasMore = true
      api.settings.update({filter: newVal}).then(this.refreshItems.bind(this, false))
      this.computeStats()
    },
    'feedSelected': function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      // topic selections are handled by selectTopic(), not the normal feed flow
      if (newVal && newVal.indexOf('topic:') === 0) return
      this.itemSelected = null
      this.items = []
      this.itemsHasMore = true
      api.settings.update({feed: newVal}).then(this.refreshItems.bind(this, false))
      if (this.$refs.itemlist) this.$refs.itemlist.scrollTop = 0
    },
    'itemSelected': function(newVal, oldVal) {
      this.itemSelectedReadability = ''
      // reset reading progress
      var prog = document.getElementById('reading-progress')
      if (prog) prog.style.width = '0'
      if (newVal === null) {
        this.itemSelectedDetails = null
        return
      }
      if (this.$refs.content) this.$refs.content.scrollTop = 0

      api.items.get(newVal).then(function(item) {
        this.itemSelectedDetails = item
        if (this.itemSelectedDetails.status == 'unread') {
          api.items.update(this.itemSelectedDetails.id, {status: 'read'}).then(function() {
            this.feedStats[this.itemSelectedDetails.feed_id].unread -= 1
            var itemInList = this.items.find(function(i) { return i.id == item.id })
            if (itemInList) itemInList.status = 'read'
            this.itemSelectedDetails.status = 'read'
          }.bind(this))
        }
      }.bind(this))
    },
    'itemSearch': debounce(function(newVal) {
      this.refreshItems()
    }, 500),
    'itemSortNewestFirst': function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      api.settings.update({sort_newest_first: newVal}).then(vm.refreshItems.bind(this, false))
    },
    'feedListWidth': debounce(function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      api.settings.update({feed_list_width: newVal})
    }, 1000),
    'itemListWidth': debounce(function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      api.settings.update({item_list_width: newVal})
    }, 1000),
    'refreshRate': function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      api.settings.update({refresh_rate: newVal})
    },
  },
  methods: {
    updateMetaTheme: function(theme) {
      document.querySelector("meta[name='theme-color']").content = this.themeColors[theme]
    },
    refreshStats: function(loopMode) {
      return api.status().then(function(data) {
        if (loopMode && !vm.itemSelected) vm.refreshItems()

        vm.loading.feeds = data.running
        if (data.running) {
          setTimeout(vm.refreshStats.bind(vm, true), 500)
        }
        vm.feedStats = data.stats.reduce(function(acc, stat) {
          acc[stat.feed_id] = stat
          return acc
        }, {})

        api.feeds.list_errors().then(function(errors) {
          vm.feed_errors = errors
        })
      })
    },
    getItemsQuery: function() {
      var query = {}
      if (this.feedSelected) {
        var parts = this.feedSelected.split(':', 2)
        var type = parts[0]
        var guid = parts[1]
        if (type == 'feed') {
          query.feed_id = guid
        } else if (type == 'folder') {
          query.folder_id = guid
        }
      }
      if (this.filterSelected) {
        query.status = this.filterSelected
      }
      if (this.itemSearch) {
        query.search = this.itemSearch
      }
      if (!this.itemSortNewestFirst) {
        query.oldest_first = true
      }
      return query
    },
    refreshFeeds: function() {
      return Promise
        .all([api.folders.list(), api.feeds.list()])
        .then(function(values) {
          vm.folders = values[0]
          vm.feeds = values[1]
        })
    },
    refreshItems: function(loadMore = false) {
      if (this.feedSelected === null) {
        vm.items = []
        vm.itemsHasMore = false
        return
      }

      var query = this.getItemsQuery()
      if (loadMore) {
        query.after = vm.items[vm.items.length-1].id
      }

      this.loading.items = true
      return api.items.list(query).then(function(data) {
        if (loadMore) {
          vm.items = vm.items.concat(data.list)
        } else {
          vm.items = data.list
        }
        vm.itemsHasMore = data.has_more
        vm.loading.items = false

        // load more if there's some space left at the bottom of the item list.
        vm.$nextTick(function() {
          if (vm.itemsHasMore && !vm.loading.items && vm.itemListCloseToBottom()) {
            vm.refreshItems(true)
          }
        })
      })
    },
    itemListCloseToBottom: function() {
      // approx. vertical space at the bottom of the list (loading el & paddings) when 1rem = 16px
      var bottomSpace = 70
      var scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16

      var el = this.$refs.itemlist

      if (el.scrollHeight === 0) return false  // element is invisible (responsive design)

      var closeToBottom = (el.scrollHeight - el.scrollTop - el.offsetHeight) < bottomSpace * scale
      return closeToBottom
    },
    loadMoreItems: function(event, el) {
      if (!this.itemsHasMore) return
      if (this.loading.items) return
      if (this.itemListCloseToBottom()) return this.refreshItems(true)
      if (this.itemSelected && this.itemSelected === this.items[this.items.length - 1].id) return this.refreshItems(true)
    },
    markItemsRead: function() {
      var query = this.getItemsQuery()
      api.items.mark_read(query).then(function() {
        vm.items = []
        vm.itemsPage = {'cur': 1, 'num': 1}
        vm.itemSelected = null
        vm.itemsHasMore = false
        vm.refreshStats()
      })
    },
    toggleFolderExpanded: function(folder) {
      folder.is_expanded = !folder.is_expanded
      api.folders.update(folder.id, {is_expanded: folder.is_expanded})
    },
    formatDate: function(datestr) {
      var options = {
        year: "numeric", month: "long", day: "numeric",
        hour: '2-digit', minute: '2-digit',
      }
      return new Date(datestr).toLocaleDateString(undefined, options)
    },
    moveFeed: function(feed, folder) {
      var folder_id = folder ? folder.id : null
      api.feeds.update(feed.id, {folder_id: folder_id}).then(function() {
        feed.folder_id = folder_id
        vm.refreshStats()
      })
    },
    moveFeedToNewFolder: function(feed) {
      var title = prompt('Enter folder name:')
      if (!title) return
      api.folders.create({'title': title}).then(function(folder) {
        api.feeds.update(feed.id, {folder_id: folder.id}).then(function() {
          vm.refreshFeeds().then(function() {
            vm.refreshStats()
          })
        })
      })
    },
    archiveFeed: function(feed) {
      api.feeds.update(feed.id, {archived: true}).then(function() {
        feed.archived = true
        vm.refreshStats()
      })
    },
    unarchiveFeed: function(feed) {
      api.feeds.update(feed.id, {archived: false}).then(function() {
        feed.archived = false
        vm.refreshStats()
      })
    },
    createNewFeedFolder: function() {
      var title = prompt('Enter folder name:')
      if (!title) return
      api.folders.create({'title': title}).then(function(result) {
        vm.refreshFeeds().then(function() {
          vm.$nextTick(function() {
            if (vm.$refs.newFeedFolder) {
              vm.$refs.newFeedFolder.value = result.id
            }
          })
        })
      })
    },
    renameFolder: function(folder) {
      var newTitle = prompt('Enter new title', folder.title)
      if (newTitle) {
        api.folders.update(folder.id, {title: newTitle}).then(function() {
          folder.title = newTitle
          this.folders.sort(function(a, b) {
            return a.title.localeCompare(b.title)
          })
        }.bind(this))
      }
    },
    deleteFolder: function(folder) {
      if (confirm('Are you sure you want to delete ' + folder.title + '?')) {
        api.folders.delete(folder.id).then(function() {
          vm.feedSelected = null
          vm.refreshStats()
          vm.refreshFeeds()
          vm.toast('Folder deleted')
        })
      }
    },
    updateFeedLink: function(feed) {
      var newLink = prompt('Enter feed link', feed.feed_link)
      if (newLink) {
        api.feeds.update(feed.id, {feed_link: newLink}).then(function() {
          feed.feed_link = newLink
        })
      }
    },
    renameFeed: function(feed) {
      var newTitle = prompt('Enter new title', feed.title)
      if (newTitle) {
        api.feeds.update(feed.id, {title: newTitle}).then(function() {
          feed.title = newTitle
        })
      }
    },
    deleteFeed: function(feed) {
      if (confirm('Are you sure you want to delete ' + feed.title + '?')) {
        api.feeds.delete(feed.id).then(function() {
          vm.feedSelected = null
          vm.refreshStats()
          vm.refreshFeeds()
          vm.toast('Feed deleted')
        })
      }
    },
    createFeed: function(event) {
      var form = event.target
      var data = {
        url: form.querySelector('input[name=url]').value,
        folder_id: parseInt(form.querySelector('select[name=folder_id]').value) || null,
      }
      if (this.feedNewChoiceSelected) {
        data.url = this.feedNewChoiceSelected
      }
      this.loading.newfeed = true
      api.feeds.create(data).then(function(result) {
        if (result.status === 'success') {
          vm.refreshFeeds()
          vm.refreshStats()
          vm.settings = ''
          vm.feedSelected = 'feed:' + result.feed.id
          vm.toast('Feed added')
        } else if (result.status === 'multiple') {
          vm.feedNewChoice = result.choice
          vm.feedNewChoiceSelected = result.choice[0].url
        } else {
          vm.toast('No feeds found at the given url.', 'error')
        }
        vm.loading.newfeed = false
      })
    },
    toggleItemStatus: function(item, targetstatus, fallbackstatus) {
      var oldstatus = item.status
      var newstatus = item.status !== targetstatus ? targetstatus : fallbackstatus

      var updateStats = function(status, incr) {
        if ((status == 'unread') || (status == 'starred')) {
          this.feedStats[item.feed_id][status] += incr
        }
      }.bind(this)

      api.items.update(item.id, {status: newstatus}).then(function() {
        updateStats(oldstatus, -1)
        updateStats(newstatus, +1)

        var itemInList = this.items.find(function(i) { return i.id == item.id })
        if (itemInList) itemInList.status = newstatus
        item.status = newstatus
      }.bind(this))
    },
    toggleItemStarred: function(item) {
      this.toggleItemStatus(item, 'starred', 'read')
    },
    toggleItemRead: function(item) {
      this.toggleItemStatus(item, 'unread', 'read')
    },
    importOPML: function(event) {
      var input = event.target
      var form = document.querySelector('#opml-import-form')
      this.$refs.menuDropdown.hide()
      api.upload_opml(form).then(function() {
        input.value = ''
        vm.refreshFeeds()
        vm.refreshStats()
        vm.toast('OPML imported')
      })
    },
    logout: function() {
      api.logout().then(function() {
        document.location.reload()
      })
    },
    toggleReadability: function() {
      if (this.itemSelectedReadability) {
        this.itemSelectedReadability = null
        return
      }
      var item = this.itemSelectedDetails
      if (!item) return
      if (item.link) {
        this.loading.readability = true
        api.crawl(item.link).then(function(data) {
          vm.itemSelectedReadability = data && data.content
          vm.loading.readability = false
        })
      }
    },
    showSettings: function(settings) {
      this.settings = settings

      if (settings === 'create') {
        vm.feedNewChoice = []
        vm.feedNewChoiceSelected = ''
      }
    },
    resizeFeedList: function(width) {
      this.feedListWidth = Math.min(Math.max(200, width), 700)
    },
    resizeItemList: function(width) {
      this.itemListWidth = Math.min(Math.max(200, width), 700)
    },
    resetFeedChoice: function() {
      this.feedNewChoice = []
      this.feedNewChoiceSelected = ''
    },
    incrFont: function(x) {
      this.theme.size = +(this.theme.size + (0.1 * x)).toFixed(1)
    },
    fetchAllFeeds: function() {
      if (this.loading.feeds) return
      api.feeds.refresh().then(function() {
        vm.refreshStats()
      })
    },
    computeStats: function() {
      var filter = this.filterSelected
      if (!filter) {
        this.filteredFeedStats = {}
        this.filteredFolderStats = {}
        this.filteredTotalStats = null
        return
      }

      var statsFeeds = {}, statsFolders = {}, statsTotal = 0

      for (var i = 0; i < this.feeds.length; i++) {
        var feed = this.feeds[i]
        
        var n = 0
        if (filter === 'archived') {
          // For archived filter, show count of 1 if feed is archived, 0 otherwise
          n = feed.archived ? 1 : 0
        } else {
          // For other filters (unread, starred), use existing stats
          if (!this.feedStats[feed.id]) continue
          n = vm.feedStats[feed.id][filter] || 0
        }

        if (!statsFolders[feed.folder_id]) statsFolders[feed.folder_id] = 0

        statsFeeds[feed.id] = n
        statsFolders[feed.folder_id] += n
        statsTotal += n
      }

      this.filteredFeedStats = statsFeeds
      this.filteredFolderStats = statsFolders
      this.filteredTotalStats = statsTotal
    },
    // navigation helper, navigate relative to selected item
    navigateToItem: function(relativePosition) {
      let vm = this
      if (vm.itemSelected == null) {
        // if no item is selected, select first
        if (vm.items.length !== 0) vm.itemSelected = vm.items[0].id
        return
      }

      var itemPosition = vm.items.findIndex(function(x) { return x.id === vm.itemSelected })
      if (itemPosition === -1) {
        if (vm.items.length !== 0) vm.itemSelected = vm.items[0].id
        return
      }

      var newPosition = itemPosition + relativePosition
      if (newPosition < 0 || newPosition >= vm.items.length) return

      vm.itemSelected = vm.items[newPosition].id

      vm.$nextTick(function() {
        var scroll = document.querySelector('#item-list-scroll')

        var handle = scroll.querySelector('input[type=radio]:checked')
        var target = handle && handle.parentElement

        if (target && scroll) scrollto(target, scroll)

        vm.loadMoreItems()
      })
    },
    // navigation helper, navigate relative to selected feed
    navigateToFeed: function(relativePosition) {
      let vm = this
      const navigationList = this.foldersWithFeeds
        .filter(folder => !folder.id || !vm.mustHideFolder(folder))
        .map((folder) => {
          if (this.mustHideFolder(folder)) return []
          const folds = folder.id ? [`folder:${folder.id}`] : []
          const feeds = (folder.is_expanded || !folder.id)
            ? (folder.feeds || []).filter(f => !vm.mustHideFeed(f)).map(f => `feed:${f.id}`)
            : []
          return folds.concat(feeds)
        })
        .flat()
      navigationList.unshift('')

      var currentFeedPosition = navigationList.indexOf(vm.feedSelected)

      if (currentFeedPosition == -1) {
        vm.feedSelected = ''
        return
      }

      var newPosition = currentFeedPosition+relativePosition
      if (newPosition < 0 || newPosition >= navigationList.length) return

      vm.feedSelected = navigationList[newPosition]

      vm.$nextTick(function() {
        var scroll = document.querySelector('#feed-list-scroll')

        var handle = scroll.querySelector('input[type=radio]:checked')
        var target = handle && handle.parentElement

        if (target && scroll) scrollto(target, scroll)
      })
    },
    changeRefreshRate: function(offset) {
      const curIdx = this.refreshRateOptions.findIndex(o => o.value === this.refreshRate)
      if (curIdx <= 0 && offset < 0) return
      if (curIdx >= (this.refreshRateOptions.length - 1) && offset > 0) return
      this.refreshRate = this.refreshRateOptions[curIdx + offset].value
    },
    mustHideFolder: function (folder) {
      return this.filterSelected
        && !(this.current.folder.id == folder.id || this.current.feed.folder_id == folder.id)
        && !this.filteredFolderStats[folder.id]
        && (!this.itemSelectedDetails || (this.feedsById[this.itemSelectedDetails.feed_id] || {}).folder_id != folder.id)
    },
    mustHideFeed: function (feed) {
      // Archive visibility logic
      if (this.filterSelected === 'archived') {
        // In archived filter, hide non-archived feeds
        if (!feed.archived) return true
      } else if (this.filterSelected && this.filterSelected !== 'archived') {
        // In other filters (unread, starred), hide archived feeds
        if (feed.archived) return true
      }

      // Standard feed visibility logic
      return this.filterSelected
        && !(this.current.feed.id == feed.id)
        && !this.filteredFeedStats[feed.id]
        && (!this.itemSelectedDetails || this.itemSelectedDetails.feed_id != feed.id)
    },

    // ── Toast Notifications ────────────────────────────────────────────────
    toast: function(message, type) {
      var container = document.getElementById('toast-container')
      if (!container) return
      var el = document.createElement('div')
      el.className = 'toast' + (type === 'error' ? ' toast-error' : '')
      el.textContent = message
      container.appendChild(el)
      setTimeout(function() {
        el.classList.add('toast-out')
        setTimeout(function() { if (el.parentNode) el.parentNode.removeChild(el) }, 200)
      }, 3000)
    },

    // ── Density Toggle ─────────────────────────────────────────────────────
    toggleDensity: function() {
      this.densityCompact = !this.densityCompact
    },

    // ── Reader/Focus Mode ──────────────────────────────────────────────────
    toggleReaderMode: function() {
      this.readerMode = !this.readerMode
      var app = document.getElementById('app')
      if (app) app.classList.toggle('reader-mode', this.readerMode)
    },

    // ── Reading Progress Bar ───────────────────────────────────────────────
    updateReadingProgress: function(event) {
      var el = event.target
      var progress = document.getElementById('reading-progress')
      if (!progress) return
      var pct = el.scrollHeight <= el.clientHeight
        ? 0
        : (el.scrollTop / (el.scrollHeight - el.clientHeight)) * 100
      progress.style.width = Math.min(100, Math.max(0, pct)) + '%'
    },

    // ── Topics / Clusters ──────────────────────────────────────────────────
    toggleTopicsView: function() {
      this.topicsActive = !this.topicsActive
      if (this.topicsActive && !this.topicsLoaded) {
        this.loadTopics()
      }
    },
    loadTopics: function() {
      this.topicsLoading = true
      Promise.all([
        api.ai.clusters().catch(function() { return {clusters: []} }),
        api.ai.tags().catch(function() { return [] }),
        api.ai.health().catch(function() { return {} }),
      ]).then(function(results) {
        var clustersResp = results[0]
        var tagsResp = results[1]
        var healthResp = results[2]
        vm.topicClusters = clustersResp.clusters || []
        var clusterLabels = new Set(vm.topicClusters.map(function(c) { return c.label }))
        // tags not already covered by cluster labels
        var tags = Array.isArray(tagsResp) ? tagsResp : []
        vm.topicTags = tags.filter(function(t) { return !clusterLabels.has(t.tag) }).slice(0, 30)
        vm.topicsHealth = healthResp
        vm.topicsLoaded = true
        vm.topicsLoading = false
      }).catch(function() {
        vm.topicsLoading = false
      })
    },
    selectTopic: function(tag) {
      if (this.selectedTopic === tag) {
        this.selectedTopic = null
        return
      }
      this.selectedTopic = tag
      this.itemSelected = null
      this.items = []
      this.itemsHasMore = false
      this.feedSelected = 'topic:' + tag
      api.ai.articles(tag).then(function(articles) {
        var list = Array.isArray(articles) ? articles : []
        vm.items = list.map(function(a) {
          return {
            id: a.id,
            title: a.title || a.url || 'untitled',
            date: a.published,
            feed_id: null,
            status: 'read',
            _feedName: a.feed_name || '',
          }
        })
      }).catch(function(err) {
        console.error('selectTopic error:', err)
        vm.items = []
      })
    },

    // ── AI Task Status Polling ─────────────────────────────────────────────
    showAiStatus: function(text) {
      var bar = document.getElementById('ai-status')
      if (bar) bar.classList.add('active')
      var textEl = document.getElementById('ai-status-text')
      if (textEl) textEl.textContent = text
    },
    hideAiStatus: function() {
      var bar = document.getElementById('ai-status')
      if (bar) bar.classList.remove('active')
      clearInterval(this._aiPollTimer)
      this._aiPollTimer = null
    },
    startTaskPoll: function() {
      if (this._aiPollTimer) return
      this._aiPollTimer = setInterval(function() {
        api.ai.taskStatus().then(function(task) {
          if (task && task.type) {
            var detail = task.detail || 'running...'
            vm._lastAiDetail = detail
            vm.showAiStatus((task.type === 'reindex' ? 'Indexing' : 'Clustering') + ': ' + detail)
          } else {
            var msg = vm._lastAiDetail || 'AI task complete'
            vm.hideAiStatus()
            vm.toast(msg)
            vm.topicsLoaded = false
            if (vm.topicsActive) vm.loadTopics()
          }
        }).catch(function() { /* ignore */ })
      }, 5000)
    },
    checkTaskStatus: function() {
      api.ai.taskStatus().then(function(task) {
        if (task && task.type) {
          vm.showAiStatus((task.type === 'reindex' ? 'Indexing' : 'Clustering') + ': ' + (task.detail || 'running...'))
          vm.startTaskPoll()
        }
      }).catch(function() { /* ignore */ })
    },
    reindexArticles: function(event) {
      var btn = event && event.currentTarget ? event.currentTarget : null
      if (btn) { btn.disabled = true; btn.classList.add('loading') }
      this.showAiStatus('Indexing: starting...')
      api.ai.reindex().then(function() {
        vm.startTaskPoll()
      }).catch(function() {
        vm.hideAiStatus()
        vm.toast('Failed to start reindex', 'error')
        if (btn) { btn.disabled = false; btn.classList.remove('loading') }
      })
    },
    rebuildTopics: function(event) {
      var btn = event && event.currentTarget ? event.currentTarget : null
      if (btn) { btn.disabled = true; btn.classList.add('loading') }
      this.showAiStatus('Clustering: starting...')
      api.ai.recluster().then(function() {
        vm.startTaskPoll()
      }).catch(function() {
        vm.hideAiStatus()
        vm.toast('Failed to start clustering', 'error')
        if (btn) { btn.disabled = false; btn.classList.remove('loading') }
      })
    },

    // ── AI Chat ────────────────────────────────────────────────────────────
    openChat: function() {
      var panel = document.getElementById('chat-panel')
      if (panel) panel.classList.add('open')
      var input = document.getElementById('chat-input')
      if (input) setTimeout(function() { input.focus() }, 50)
    },
    closeChat: function() {
      var panel = document.getElementById('chat-panel')
      if (panel) panel.classList.remove('open')
      if (this._chatReader) {
        try { this._chatReader.cancel() } catch(e) {}
        this._chatReader = null
      }
    },
    sendChat: function(event) {
      event.preventDefault()
      var input = document.getElementById('chat-input')
      var query = (input ? input.value : '').trim()
      if (!query) return
      if (input) input.value = ''

      var messages = document.getElementById('chat-messages')
      if (!messages) return

      // append user message
      var userEl = document.createElement('div')
      userEl.className = 'chat-msg chat-msg-user'
      userEl.textContent = query
      messages.appendChild(userEl)

      // append assistant placeholder with streaming cursor
      var assistantEl = document.createElement('div')
      assistantEl.className = 'chat-msg chat-msg-assistant streaming-cursor'
      messages.appendChild(assistantEl)
      messages.scrollTop = messages.scrollHeight

      var history = this._chatHistory.slice()
      var fullResponse = ''

      api.ai.chat(query, history).then(function(response) {
        if (!response.ok) throw new Error('Chat request failed')
        var reader = response.body.getReader()
        vm._chatReader = reader
        var decoder = new TextDecoder()

        function read() {
          reader.read().then(function(result) {
            if (result.done) {
              assistantEl.classList.remove('streaming-cursor')
              vm._chatHistory = vm._chatHistory.concat([
                {role: 'user', content: query},
                {role: 'assistant', content: fullResponse},
              ]).slice(-6)
              return
            }
            var chunk = decoder.decode(result.value, {stream: true})
            var lines = chunk.split('\n')
            for (var i = 0; i < lines.length; i++) {
              var line = lines[i]
              if (line.startsWith('event: sources')) continue
              if (line.startsWith('data: ')) {
                var data = line.slice(6)
                if (data === '[DONE]') continue
                try {
                  var parsed = JSON.parse(data)
                  if (parsed.error) {
                    assistantEl.classList.remove('streaming-cursor')
                    assistantEl.textContent = 'Error: ' + parsed.error
                    return
                  }
                  if (parsed.sources && Array.isArray(parsed.sources)) {
                    var sourcesEl = document.createElement('div')
                    sourcesEl.className = 'chat-msg-sources'
                    parsed.sources.forEach(function(src, idx) {
                      var a = document.createElement('a')
                      a.href = src.url || '#'
                      a.textContent = '[' + (idx+1) + '] ' + (src.title || src.url || '')
                      a.target = '_blank'
                      a.rel = 'noopener noreferrer'
                      a.style.display = 'block'
                      sourcesEl.appendChild(a)
                    })
                    assistantEl.after(sourcesEl)
                    continue
                  }
                  if (typeof parsed === 'string') {
                    fullResponse += parsed
                    assistantEl.textContent = fullResponse
                  }
                } catch(e) {
                  // plain text token
                  fullResponse += data
                  assistantEl.textContent = fullResponse
                }
              }
            }
            messages.scrollTop = messages.scrollHeight
            read()
          }).catch(function() {
            assistantEl.classList.remove('streaming-cursor')
          })
        }
        read()
      }).catch(function(err) {
        assistantEl.classList.remove('streaming-cursor')
        assistantEl.textContent = 'Error: ' + err.message
      })
    },

    // ── AI Briefing ────────────────────────────────────────────────────────
    openBriefing: function() {
      var panel = document.getElementById('briefing-panel')
      if (panel) panel.classList.add('open')
      this.loadBriefing()
    },
    closeBriefing: function() {
      var panel = document.getElementById('briefing-panel')
      if (panel) panel.classList.remove('open')
      if (this._briefingReader) {
        try { this._briefingReader.cancel() } catch(e) {}
        this._briefingReader = null
      }
    },
    loadBriefing: function() {
      var sinceEl = document.getElementById('briefing-since')
      var since = sinceEl ? sinceEl.value : '7d'
      var content = document.getElementById('briefing-content')
      if (!content) return

      // show streaming cursor
      content.innerHTML = '<div class="streaming-cursor"></div>'
      var cursorEl = content.firstChild

      var fullText = ''

      api.ai.briefing(since).then(function(response) {
        if (!response.ok) throw new Error('Briefing request failed')
        var reader = response.body.getReader()
        vm._briefingReader = reader
        var decoder = new TextDecoder()

        function read() {
          reader.read().then(function(result) {
            if (result.done) {
              // convert basic markdown to HTML
              var html = fullText
                .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
                .replace(/\[(\d+)\]/g, '<sup>[$1]</sup>')
                .replace(/^### (.*)/gm, '<h3>$1</h3>')
                .replace(/^## (.*)/gm, '<h2>$1</h2>')
                .replace(/\n\n/g, '</p><p>')
              content.innerHTML = '<p>' + html + '</p>'
              return
            }
            var chunk = decoder.decode(result.value, {stream: true})
            var lines = chunk.split('\n')
            for (var i = 0; i < lines.length; i++) {
              var line = lines[i]
              if (line.startsWith('data: ')) {
                var data = line.slice(6)
                if (data === '[DONE]') continue
                fullText += data
                cursorEl.textContent = fullText
              }
            }
            read()
          }).catch(function() {
            if (cursorEl) cursorEl.classList.remove('streaming-cursor')
          })
        }
        read()
      }).catch(function(err) {
        content.innerHTML = '<p class="text-danger">Error: ' + err.message + '</p>'
      })
    },
  }
})

vm.$mount('#app')
