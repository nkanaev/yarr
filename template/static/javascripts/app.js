'use strict';

var FILTERS = [
  {'title': 'All', 'value': 'all', 'icon': 'circle-full'},
  {'title': 'Unread', 'value': 'unread', 'icon': 'circle'},
  {'title': 'Starred', 'value': 'starred', 'icon': 'star'},
]


new Vue({
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
        {'id': '2', 'title': '/r/programming', 'folder_id': 1},
        {'id': '3', 'title': 'BBC', 'folder_id': 2},
        {'id': '4', 'title': 'The Guardian', 'folder_id': 2},
        {'id': '5', 'title': 'Random Stuff', 'folder_id': null},
      ],
      'feedSelected': null,
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
  },
  methods: {
    toggleFolderExpanded: function(folder) {
      folder.is_expanded = !folder.is_expanded
    }
  }
})
