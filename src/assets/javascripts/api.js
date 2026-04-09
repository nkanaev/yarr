"use strict";

(function() {
  var xfetch = function(resource, init) {
    init = init || {}
    if (['post', 'put', 'delete'].indexOf(init.method) !== -1) {
      init['headers'] = init['headers'] || {}
      init['headers']['x-requested-by'] = 'yarr'
    }
    return fetch(resource, init)
  }
  var api = function(method, endpoint, data) {
    var headers = {'Content-Type': 'application/json'}
    return xfetch(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(data),
    })
  }

  var json = function(res) {
    return res.json()
  }

  var param = function(query) {
    if (!query) return ''
    return '?' + Object.keys(query).map(function(key) {
      return encodeURIComponent(key) + '=' + encodeURIComponent(query[key])
    }).join('&')
  }

  var topicFilters = function(filters) {
    var query = {}
    if (!filters) return query
    if (filters.status !== undefined && filters.status !== null && filters.status !== '') {
      query.status = filters.status
    }
    if (filters.since) {
      query.since = filters.since
    }
    return query
  }

  window.api = {
    feeds: {
      list: function() {
        return api('get', './api/feeds').then(json)
      },
      create: function(data) {
        return api('post', './api/feeds', data).then(json)
      },
      update: function(id, data) {
        return api('put', './api/feeds/' + id, data)
      },
      delete: function(id) {
        return api('delete', './api/feeds/' + id)
      },
      list_items: function(id) {
        return api('get', './api/feeds/' + id + '/items').then(json)
      },
      refresh: function() {
        return api('post', './api/feeds/refresh')
      },
      list_errors: function() {
        return api('get', './api/feeds/errors').then(json)
      },
    },
    folders: {
      list: function() {
        return api('get', './api/folders').then(json)
      },
      create: function(data) {
        return api('post', './api/folders', data).then(json)
      },
      update: function(id, data) {
        return api('put', './api/folders/' + id, data)
      },
      delete: function(id) {
        return api('delete', './api/folders/' + id)
      },
      list_items: function(id) {
        return api('get', './api/folders/' + id + '/items').then(json)
      }
    },
    items: {
      get: function(id) {
        return api('get', './api/items/' + id).then(json)
      },
      list: function(query) {
        return api('get', './api/items' + param(query)).then(json)
      },
      update: function(id, data) {
        return api('put', './api/items/' + id, data)
      },
      mark_read: function(query) {
        return api('put', './api/items' + param(query))
      },
    },
    settings: {
      get: function() {
        return api('get', './api/settings').then(json)
      },
      update: function(data) {
        return api('put', './api/settings', data)
      },
    },
    status: function() {
      return api('get', './api/status').then(json)
    },
    upload_opml: function(form) {
      return xfetch('./opml/import', {
        method: 'post',
        body: new FormData(form),
      })
    },
    logout: function() {
      return api('post', './logout')
    },
    crawl: function(url) {
      return api('get', './page?url=' + encodeURIComponent(url)).then(json)
    },
    ranking: {
      react: function(itemId, reaction) {
        return api('post', './api/reactions', {item_id: itemId, reaction: reaction})
      },
      getReaction: function(itemId) {
        return api('get', './api/reactions?item_id=' + itemId).then(json)
      },
      clickThrough: function(itemId) {
        return api('post', './api/click-throughs', {item_id: itemId})
      },
      readHere: function(itemId) {
        return api('post', './api/read-heres', {item_id: itemId})
      },
      list: function(query) {
        return api('get', './api/items/ranked' + param(query)).then(json)
      },
      preferences: function() {
        return api('get', './api/preferences').then(json)
      },
      resetPreferences: function() {
        return api('delete', './api/preferences')
      },
    },
    ai: {
      // Returns raw Response for SSE streaming (POST with JSON body)
      chat: function(query, history) {
        return xfetch('./api/ai/chat', {
          method: 'post',
          headers: {'Content-Type': 'application/json', 'x-requested-by': 'yarr'},
          body: JSON.stringify({query: query, history: history}),
        })
      },
      // Returns raw Response for SSE streaming (GET)
      briefing: function(since) {
        return xfetch('./api/ai/briefing?since=' + encodeURIComponent(since), {
          headers: {'Accept': 'text/event-stream'},
        })
      },
      clusters: function(filters) {
        return api('get', './api/ai/clusters' + param(topicFilters(filters))).then(json)
      },
      tags: function() {
        return api('get', './api/ai/tags').then(json)
      },
      health: function() {
        return api('get', './api/ai/health').then(json)
      },
      articles: function(tag, filters) {
        var query = topicFilters(filters)
        query.tag = tag
        return api('get', './api/ai/articles' + param(query)).then(json)
      },
      reindex: function() {
        return api('post', './api/ai/reindex')
      },
      recluster: function() {
        return api('post', './api/ai/recluster')
      },
      taskStatus: function() {
        return api('get', './api/ai/task-status').then(json)
      },
    }
  }
})()
