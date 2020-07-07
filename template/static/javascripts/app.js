'use strict';

var vm = new Vue({
  el: '#app',
  created: function() {
    var vm = this
    api.settings.get().then(function(data) {
      vm.filterSelected = data.filter
      vm.refreshItems()
    })
    this.refreshFeeds()
  },
  data: function() {
    return {
      'filterSelected': null,
      'folders': [],
      'feeds': [],
      'feedSelected': null,
      'items': [],
      'itemSelected': null,
      'itemSelectedDetails': {},
      'settings': 'create',
      'loading': {newfeed: 0},
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
  },
  watch: {
    'filterSelected': function(newVal, oldVal) {
      if (oldVal === null) return  // do nothing, initial setup
      var vm = this
      api.settings.update({filter: newVal}).then(function() {
        vm.refreshItems()
      })
    },
    'feedSelected': function(newVal, oldVal) {
      this.refreshItems()
    },
    'itemSelected': function(newVal, oldVal) {
      this.itemSelectedDetails = this.itemsById[newVal]
      if (this.itemSelectedDetails.status == 'unread') {
        this.itemSelectedDetails.status = 'read'
        api.items.update(this.itemSelectedDetails.id, {status: this.itemSelectedDetails.status})
      }
    },
  },
  methods: {
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
      var promise = null
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
      api.items.list(query).then(function(items) {
        vm.items = items
      })
    },
    toggleFolderExpanded: function(folder) {
      folder.is_expanded = !folder.is_expanded
    },
    formatDate: function(datestr) {
      return new Date(datestr).toLocaleDateString(undefined, {year: "numeric", month: "long", day: "numeric"})
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
      } else if (item.status != 'starred') {
        item.status = 'starred'
      }
      api.items.update(item.id, {status: item.status})
    },
    toggleItemRead: function(item) {
      if (item.status == 'unread') {
        item.status = 'read'
      } else if (item.status == 'read') {
        item.status = 'unread'
      }
      api.items.update(item.id, {status: item.status})
    },
  }
})
