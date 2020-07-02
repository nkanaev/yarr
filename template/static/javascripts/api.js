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

  window.api = {
    feeds: {
      list: function() {
        return api('get', '/api/feeds').then(json)
      },
      create: function(data) {
        return api('post', '/api/feeds', data).then(json)
      },
      delete: function(id) {
        return api('delete', '/api/feeds/' + id)
      }
    },
    folders: {
      list: function() {
        return api('get', '/api/folders').then(json)
      },
      create: function(data) {
        return api('post', '/api/folders', data).then(json)
      },
    }
  }
})()
