'use strict';

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

function dateRepr(d) {
  var sec = (new Date().getTime() - d.getTime()) / 1000
  if (sec < 2700)  // less than 45 minutes
    return Math.round(sec / 60) + 'm'
  else if (sec < 86400)  // less than 24 hours
    return Math.round(sec / 3600) + 'h'
  else if (sec < 604800)  // less than a week
    return Math.round(sec / 86400) + 'd'
  else
    return d.toLocaleDateString(undefined, {year: "numeric", month: "long", day: "numeric"})
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
  template: '<time :datetime="val">{{formatted}}</time>',
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
  el: '#app',
  created: function() {
    var vm = this
    api.settings.get().then(function(data) {
      vm.feedSelected = data.feed
      vm.filterSelected = data.filter
      vm.itemSortNewestFirst = data.sort_newest_first
      vm.refreshItems()
    })
    this.refreshFeeds()
    this.refreshStats()
  },
  data: function() {
    return {
      'filterSelected': null,
      'folders': [],
      'feeds': [],
      'feedSelected': null,
      'items': [],
      'itemsPage': {
        'cur': 1,
        'num': 1,
      },
      'itemSelected': null,
      'itemSelectedDetails': {},
      'itemSelectedReadability': '',
      'itemSearch': '',
      'itemSortNewestFirst': null,
      'settings': 'create',
      'loading': {
        'newfeed': false,
        'items': false,
      },
      'feedStats': {},
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
    filteredFeedStats: function() {
      var filter = this.filterSelected
      if (filter != 'unread' && filter != 'starred') return {}

      var feedStats = this.feedStats
      return this.feeds.reduce(function(acc, feed) {
        if (feedStats[feed.id]) acc[feed.id] = vm.feedStats[feed.id][filter]
        return acc
      }, {})
    },
    filteredFolderStats: function() {
      var filter = this.filterSelected
      if (filter != 'unread' && filter != 'starred') return {}

      var feedStats = this.filteredFeedStats
      return this.feeds.reduce(function(acc, feed) {
        if (!acc[feed.folder_id]) acc[feed.folder_id] = 0
        if (feedStats[feed.id]) acc[feed.folder_id] += feedStats[feed.id]
        return acc
      }, {})
    },
    totalStats: function() {
      return Object.values(this.feedStats).reduce(function(acc, stat) {
        acc.unread += stat.unread
        acc.starred += stat.starred
        return acc
      }, {unread: 0, starred: 0})
    },
  },
  watch: {
    'filterSelected': function(newVal, oldVal) {
      if (oldVal === null) return  // do nothing, initial setup
      api.settings.update({filter: newVal}).then(this.refreshItems.bind(this))
    },
    'feedSelected': function(newVal, oldVal) {
      if (oldVal === null) return  // do nothing, initial setup
      api.settings.update({feed: newVal}).then(this.refreshItems.bind(this))
    },
    'itemSelected': function(newVal, oldVal) {
      this.itemSelectedReadability = ''
      this.itemSelectedDetails = this.itemsById[newVal]
      if (this.itemSelectedDetails.status == 'unread') {
        this.itemSelectedDetails.status = 'read'
        this.feedStats[this.itemSelectedDetails.feed_id].unread -= 1
        api.items.update(this.itemSelectedDetails.id, {status: this.itemSelectedDetails.status})
      }
    },
    'itemSearch': debounce(function(newVal) {
      if (newVal) {
        this.refreshItems()
      }
    }, 500),
    'itemSortNewestFirst': function(newVal, oldVal) {
      if (oldVal === null) return
      api.settings.update({sort_newest_first: newVal}).then(this.refreshItems.bind(this))
    },
  },
  methods: {
    refreshStats: function() {
      var vm = this
      api.status().then(function(data) {
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
      var vm = this
      Promise
        .all([api.folders.list(), api.feeds.list()])
        .then(function(values) {
          vm.folders = values[0]
          vm.feeds = values[1]
        })
    },
    refreshItems: function() {
      var query = this.getItemsQuery()
      this.loading.items = true
      var vm = this
      api.items.list(query).then(function(data) {
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
      var vm = this
      var query = this.getItemsQuery()
      api.items.mark_read(query).then(function() {
        vm.items = []
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
    createFolder: function(event) {
      var form = event.target
      var titleInput = form.querySelector('input[name=title]')
      var data = {'title': titleInput.value}
      var vm = this
      api.folders.create(data).then(function(result) {
        vm.folders.push(result)
        titleInput.value = ''
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
      var vm = this
      if (confirm('Are you sure you want to delete ' + folder.title + '?')) {
        api.folders.delete(folder.id).then(function() {
          vm.refresh()
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
        var vm = this
        api.feeds.delete(feed.id).then(function() {
          api.feeds.list().then(function(feeds) {
            vm.feeds = feeds
          })
        })
      }
    },
    createFeed: function(event) {
      var form = event.target
      var data = {
        url: form.querySelector('input[name=url]').value,
        folder_id: parseInt(form.querySelector('select[name=folder_id]').value) || null,
      }
      this.loading.newfeed = true
      var vm = this
      api.feeds.create(data).then(function(result) {
        if (result.status === 'success') {
          api.feeds.list().then(function(feeds) {
            vm.feeds = feeds
          })
          vm.$bvModal.hide('settings-modal')
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
      var vm = this
      var input = event.target
      var form = document.querySelector('#opml-import-form')
      api.upload_opml(form).then(function() {
        input.value = ''
        vm.refreshFeeds()
      })
    },
    getReadable: function(item) {
      if (item.link) {
        var vm = this
        api.crawl(item.link).then(function(body) {
          if (!body.length) return
          var doc = new DOMParser().parseFromString(body, 'text/html')
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
    },
  }
})
