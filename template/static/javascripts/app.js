'use strict';

var FILTERS = [
  {'title': 'All', 'value': 'all', 'icon': 'circle-full'},
  {'title': 'Unread', 'value': 'unread', 'icon': 'circle'},
  {'title': 'Starred', 'value': 'starred', 'icon': 'star'},
]


var vm = new Vue({
  el: '#app',
  data: function() {
    return {
      'filters': FILTERS,
      'filterSelected': 'all',
      'folders': [
        {'id': 1, 'title': 'Tech', 'is_expanded': false},
        {'id': 2, 'title': 'News', 'is_expanded': true},
      ],
      'feeds': [
        {'id': '1', 'title': 'news.ycombinator.com', 'folder_id': 1},
        {'id': '2', 'title': 'daringfireball', 'folder_id': 1},
        {'id': '3', 'title': 'BBC', 'folder_id': 2},
        {'id': '4', 'title': 'The Guardian', 'folder_id': 2},
        {'id': '5', 'title': 'Random Stuff', 'folder_id': null},
      ],
      'feedSelected': null,
      'items': [
        {'id': '123', 'title': 'Apple Pulls Pocket Casts and Castro From Chinese App Store', 'status': 'unread', 'feed_id': 2, 'date': 1592250298},
        {'id': '456', 'title': 'On Apple Announcing the ARM Mac Transition at WWDC This Month', 'status': 'starred', 'feed_id': 2, 'date': 1592250298},
        {'id': '789', 'title': 'Marques Brownlee: ‘Reflecting on the Color of My Skin’', 'status': 'read', 'feed_id': 2, 'date': 1592250298},
      ],
      'itemSelected': null,
      'settingsShow': false,
      'settings': 'manage',
      'settingsManageDropdown': null,
      'newFolderTitle': null,
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
    'settingsShow': function(newVal) {
      if (newVal === true) {
        var vm = this
        var backdrop = document.createElement('div')
        backdrop.classList.add('modal-backdrop', 'fade', 'show')
        document.body.classList.add('modal-open')
        document.body.appendChild(backdrop)
      } else {
        document.body.classList.remove('modal-open')
        document.body.querySelector('.modal-backdrop').remove()
      }
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
      this.settingsManageDropdown = null
    },
    newFolderCreate: function() {
      this.folders.push({
        id: Math.random() * 10000,
        title: this.newFolderTitle,
        is_expanded: true,
      })
    },
  }
})
vm.settingsShow = true
