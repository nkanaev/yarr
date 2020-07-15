"use strict";

(function() {
  var api = function(method, endpoint, data) {
    return fetch(endpoint, {
      method: method,
      headers: {'content-type': 'application/json'},
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

  window.api = {
    feeds: {
      list: function() {
        return api('get', '/api/feeds').then(json)
      },
      create: function(data) {
        return api('post', '/api/feeds', data).then(json)
      },
      update: function(id, data) {
        return api('put', '/api/feeds/' + id, data)
      },
      delete: function(id) {
        return api('delete', '/api/feeds/' + id)
      },
      list_items: function(id) {
        return api('get', '/api/feeds/' + id + '/items').then(json)
      },
    },
    folders: {
      list: function() {
        return api('get', '/api/folders').then(json)
      },
      create: function(data) {
        return api('post', '/api/folders', data).then(json)
      },
      update: function(id, data) {
        return api('put', '/api/folders/' + id, data)
      },
      delete: function(id) {
        return api('delete', '/api/folders/' + id)
      },
      list_items: function(id) {
        return api('get', '/api/folders/' + id + '/items').then(json)
      }
    },
    items: {
      list: function(query) {
        return api('get', '/api/items' + param(query)).then(json)
      },
      update: function(id, data) {
        return api('put', '/api/items/' + id, data)
      },
      mark_read: function(query) {
        return api('put', '/api/items' + param(query))
      },
    },
    settings: {
      get: function() {
        return api('get', '/api/settings').then(json)
      },
      update: function(data) {
        return api('put', '/api/settings', data)
      },
    },
    status: function() {
      return api('get', '/api/status').then(json)
    },
    upload_opml: function(form) {
      return fetch('/opml/import', {
        method: 'post',
        body: new FormData(form),
      })
    },
    crawl: function(url) {
      return fetch('/page?url=' + url).then(function(res) {
        return res.text()
      })
    }
  }
})()
