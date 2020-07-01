'use strict';

var FILTERS = [
  {'title': 'All', 'value': 'all', 'icon': 'circle-full'},
  {'title': 'Unread', 'value': 'unread', 'icon': 'circle'},
  {'title': 'Starred', 'value': 'starred', 'icon': 'star'},
]

var vm = new Vue({
  el: '#app',
  created: function() {
    var vm = this
    Promise
      .all([api.folders.list(), api.feeds.list()])
      .then(function(values) {
        vm.folders = values[0]
        vm.feeds = values[1]
      })
  },
  data: function() {
    return {
      'filters': FILTERS,
      'filterSelected': 'all',
      'folders': [],
      'feeds': [],
      'feedSelected': null,
      'items': [], 
      'itemSelected': null,
      'settings': 'manage',
      'newFolderTitle': null,
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
    'feedSelected': function(newVal, oldVal) {
      if (newVal === null) return
      var parts = newVal.split(':', 2)
      var type = parts[0]
      var guid = parts[1]
    },
    'itemSelected': function(newVal, oldVal) {
      this.itemSelectedDetails = this.itemsById[newVal]
    },
  },
  methods: {
    toggleFolderExpanded: function(folder) {
      folder.is_expanded = !folder.is_expanded
    },
    formatDate: function(timestamp_s) {
      var d = new Date(timestamp_s * 1000)
      return d.getDate() + '/' + d.getMonth() + '/' + d.getFullYear()
    },
    moveFeed: function(feed, folder) {
      feed.folder_id = folder ? folder.id : null
    },
    newFolderCreate: function() {
      this.folders.push({
        id: Math.random() * 10000,
        title: this.newFolderTitle,
        is_expanded: true,
      })
    },
    renameFeed: function(feed) {
      var newTitle = prompt('Enter new title', feed.title)
      feed.title = newTitle
    },
    deleteFeed: function(feed) {
      if (confirm('Are you sure you want to delete ' + feed.title + '?')) {
        this.feeds = this.feeds.filter(function(f) { f.id != feed.id })
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
          vm.$bvModal.hide('settings-modal')
        }
        vm.loading.newfeed = false
      })
    },
  }
})
