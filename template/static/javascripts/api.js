"use strict";

(function() {
  var api = function(method, endpoint, data) {
    var promise = fetch(endpoint, {
      method: method,
      headers: {'content-type': 'application/json'},
      body: JSON.stringify(data),
    })
    return promise.then(function(res) {
      if (res.ok) return res.json()
    })
  }

  window.api = {
    feeds: {
      list: function() {
        return api('get', '/api/feeds')
      },
      create: function(data) {
        return api('post', '/api/feeds', data)
      },
    }
  }
})()
