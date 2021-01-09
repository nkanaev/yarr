'use strict';

var TITLE = document.title

function authenticated() {
  return /auth=.+/g.test(document.cookie)

}

var FONTS = [
  "Arial",
  "Courier New",
  "Georgia",
  "Times New Roman",
  "Verdana",
]

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

var sanitize = function(content, base) {
  // WILD: `item.link` may be a relative link (or some nonsense)
  try { new URL(base) } catch(err) { base = null }

  var sanitizer = new DOMPurify
  sanitizer.addHook('afterSanitizeAttributes', function(node) {
    // set all elements owning target to target=_blank
    if ('target' in node)
      node.setAttribute('target', '_blank')
    // set non-HTML/MathML links to xlink:show=new
    if (!node.hasAttribute('target') && (node.hasAttribute('xlink:href') || node.hasAttribute('href')))
      node.setAttribute('xlink:show', 'new')

    // set absolute urls
    if (base && node.attributes.href && node.attributes.href.value)
      node.href = new URL(node.attributes.href.value, base).toString()
    if (base && node.attributes.src && node.attributes.src.value)
      node.src = new URL(node.attributes.src.value, base).toString()
  })
  return sanitizer.sanitize(content, {FORBID_TAGS: ['style'], FORBID_ATTR: ['style', 'class']})
}

function extensions(details, content) {
  const ytId = details.link.match(/youtube\.com\/watch\?v=(.*)$/)[1];
  if(ytId) {
    content+=`<iframe width="560" height="315" src="https://www.youtube.com/embed/${ytId}" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope" allowfullscreen></iframe>`;
  }
  return content;
}

Vue.use(VueLazyload)

Vue.directive('scroll', {
  inserted: function(el, binding) {
    el.addEventListener('scroll', debounce(function(event) {
      binding.value(event, el)
    }, 200))
  },
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
    this.refreshFeeds()
    this.refreshStats()
  },
  mounted: function() {
    this.$root.$on('bv::modal::hidden', function(bvEvent, modalId) {
      if (vm.settings == 'create') {
        vm.feedNewChoice = []
        vm.feedNewChoiceSelected = ''
      }
    })
  },
  data: function() {
    return {
      'filterSelected': undefined,
      'folders': [],
      'feeds': [],
      'feedSelected': undefined,
      'feedListWidth': undefined,
      'feedNewChoice': [],
      'feedNewChoiceSelected': '',
      'items': [],
      'itemsPage': {
        'cur': 1,
        'num': 1,
      },
      'itemSelected': null,
      'itemSelectedDetails': {},
      'itemSelectedReadability': '',
      'itemSearch': '',
      'itemSortNewestFirst': undefined,
      'itemListWidth': undefined,

      'filteredFeedStats': {},
      'filteredFolderStats': {},
      'filteredTotalStats': null,

      'settings': 'create',
      'loading': {
        'feeds': 0,
        'newfeed': false,
        'items': false,
        'readability': false,
      },
      'fonts': FONTS,
      'feedStats': {},
      'theme': {
        'name': 'light',
        'font': '',
        'size': 1,
      },
      'refreshRate': undefined,
      'authenticated': authenticated(),
      'feed_errors': {},
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
      return this.feeds.reduce(function(acc, feed) { acc[feed.id] = feed; return acc }, {})
    },
    itemsById: function() {
      return this.items.reduce(function(acc, item) { acc[item.id] = item; return acc }, {})
    },
    itemSelectedContent: function() {
      if (!this.itemSelected) return ''

      if (this.itemSelectedReadability)
        return this.itemSelectedReadability

      var content = ''
      if (this.itemSelectedDetails.content)
        content = this.itemSelectedDetails.content
      else if (this.itemSelectedDetails.description)
        content = this.itemSelectedDetails.description

      return extensions(this.itemSelectedDetails,sanitize(content, this.itemSelectedDetails.link))
    },
  },
  watch: {
    'theme': {
      deep: true,
      handler: function(theme) {
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
      api.settings.update({filter: newVal}).then(this.refreshItems.bind(this))
      this.itemSelected = null
      this.computeStats()
    },
    'feedSelected': function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      api.settings.update({feed: newVal}).then(this.refreshItems.bind(this))
      this.itemSelected = null
      if (this.$refs.itemlist) this.$refs.itemlist.scrollTop = 0
    },
    'itemSelected': function(newVal, oldVal) {
      this.itemSelectedReadability = ''
      if (newVal === null) {
        this.itemSelectedDetails = null
        return
      }
      if (this.$refs.content) this.$refs.content.scrollTop = 0

      this.itemSelectedDetails = this.itemsById[newVal]
      if (this.itemSelectedDetails.status == 'unread') {
        this.itemSelectedDetails.status = 'read'
        this.feedStats[this.itemSelectedDetails.feed_id].unread -= 1
        api.items.update(this.itemSelectedDetails.id, {status: this.itemSelectedDetails.status})
      }
    },
    'itemSearch': debounce(function(newVal) {
      this.refreshItems()
    }, 500),
    'itemSortNewestFirst': function(newVal, oldVal) {
      if (oldVal === undefined) return  // do nothing, initial setup
      api.settings.update({sort_newest_first: newVal}).then(this.refreshItems.bind(this))
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
    refreshStats: function(loopMode) {
      api.status().then(function(data) {
        if (loopMode && !vm.itemSelected) vm.refreshItems()

        vm.loading.feeds = data.running
        if (data.running) {
          setTimeout(vm.refreshStats.bind(vm, true), 500)
        }
        vm.feedStats = data.stats.reduce(function(acc, stat) {
          acc[stat.feed_id] = stat
          return acc
        }, {})
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
    refreshItems: function() {
      if (this.feedSelected === null) {
        vm.items = []
        vm.itemsPage = {'cur': 1, 'num': 1}
        return
      }
      var query = this.getItemsQuery()
      this.loading.items = true
      return api.items.list(query).then(function(data) {
        vm.items = data.list
        vm.itemsPage = data.page
        vm.loading.items = false
      })
    },
    loadMoreItems: function(event, el) {
      if (this.itemsPage.cur >= this.itemsPage.num) return
      if (this.loading.items) return
      var closeToBottom = (el.scrollHeight - el.scrollTop - el.offsetHeight) < 50
      if (closeToBottom) {
        this.loading.moreitems = true
        var query = this.getItemsQuery()
        query.page = this.itemsPage.cur + 1
        api.items.list(query).then(function(data) {
          vm.items = vm.items.concat(data.list)
          vm.itemsPage = data.page
          vm.loading.items = false
        })
      }
    },
    markItemsRead: function() {
      var query = this.getItemsQuery()
      api.items.mark_read(query).then(function() {
        vm.items = []
        vm.itemsPage = {'cur': 1, 'num': 1}
        vm.itemSelected = null
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
      })
    },
    moveFeedToNewFolder: function(feed) {
      var title = prompt('Enter folder name:')
      if (!title) return
      api.folders.create({'title': title}).then(function(folder) {
        api.feeds.update(feed.id, {folder_id: folder.id}).then(function() {
          vm.refreshFeeds()
        })
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
        })
      }
    },
    deleteFolder: function(folder) {
      if (confirm('Are you sure you want to delete ' + folder.title + '?')) {
        api.folders.delete(folder.id).then(function() {
          if (vm.feedSelected === 'folder:'+folder.id) {
            vm.items = []
            vm.feedSelected = ''
          }
          vm.refreshStats()
          vm.refreshFeeds()
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
          // unselect feed to prevent reading properties of null in template
          var isSelected = (vm.feedSelected === 'feed:'+feed.id
            || (feed.folder_id && vm.feedSelected === 'folder:'+feed.folder_id));
          if (isSelected) vm.feedSelected = null

          vm.refreshStats()
          vm.refreshFeeds()
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
          vm.$bvModal.hide('settings-modal')
        } else if (result.status === 'multiple') {
          vm.feedNewChoice = result.choice
          vm.feedNewChoiceSelected = result.choice[0].url
        } else {
          alert('No feeds found at the given url.')
        }
        vm.loading.newfeed = false
      })
    },
    toggleItemStarred: function(item) {
      if (item.status == 'starred') {
        item.status = 'read'
        this.feedStats[item.feed_id].starred -= 1
      } else if (item.status != 'starred') {
        item.status = 'starred'
        this.feedStats[item.feed_id].starred += 1
      }
      api.items.update(item.id, {status: item.status})
    },
    toggleItemRead: function(item) {
      if (item.status == 'unread') {
        item.status = 'read'
        this.feedStats[item.feed_id].unread -= 1
      } else if (item.status == 'read') {
        item.status = 'unread'
        this.feedStats[item.feed_id].unread += 1
      }
      api.items.update(item.id, {status: item.status})
    },
    importOPML: function(event) {
      var input = event.target
      var form = document.querySelector('#opml-import-form')
      this.$refs.menuDropdown.hide()
      api.upload_opml(form).then(function() {
        input.value = ''
        vm.refreshFeeds()
        vm.refreshStats()
      })
    },
    getReadable: function(item) {
      if (this.itemSelectedReadability) {
        this.itemSelectedReadability = null
        return
      }
      if (item.link) {
        this.loading.readability = true
        api.crawl(item.link).then(function(body) {
          vm.loading.readability = false
          if (!body.length) return
          var bodyClean = sanitize(body, item.link)
          var doc = new DOMParser().parseFromString(bodyClean, 'text/html')
          var parsed = new Readability(doc).parse()
          if (parsed && parsed.content) {
            vm.itemSelectedReadability = parsed.content
          }
        })
      }
    },
    showSettings: function(settings) {
      this.settings = settings
      this.$bvModal.show('settings-modal')

      if (settings === 'manage') {
        api.feeds.list_errors().then(function(errors) {
          vm.feed_errors = errors
        })
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
      api.feeds.refresh().then(this.refreshStats.bind(this))
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
        if (!this.feedStats[feed.id]) continue

        var n = vm.feedStats[feed.id][filter] || 0

        if (!statsFolders[feed.folder_id]) statsFolders[feed.folder_id] = 0

        statsFeeds[feed.id] = n
        statsFolders[feed.folder_id] += n
        statsTotal += n
      }

      this.filteredFeedStats = statsFeeds
      this.filteredFolderStats = statsFolders
      this.filteredTotalStats = statsTotal
    },
  }
})

api.settings.get().then(function(data) {
  vm.feedSelected = data.feed
  vm.filterSelected = data.filter
  vm.itemSortNewestFirst = data.sort_newest_first
  vm.feedListWidth = data.feed_list_width || 300
  vm.itemListWidth = data.item_list_width || 300
  vm.theme.name = data.theme_name
  vm.theme.font = data.theme_font
  vm.theme.size = data.theme_size
  vm.refreshRate = data.refresh_rate
  vm.refreshItems()
  vm.$mount('#app')
})
