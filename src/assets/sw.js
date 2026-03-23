var CACHE_NAME = 'yarr-v1';
var STATIC_ASSETS = [
  './static/stylesheets/bootstrap.min.css',
  './static/stylesheets/app.css',
  './static/javascripts/vue.min.js',
  './static/javascripts/api.js',
  './static/javascripts/app.js',
  './static/javascripts/key.js',
  './static/graphicarts/favicon.svg',
  './static/graphicarts/favicon.png',
  './static/graphicarts/icon-192.png',
  './static/graphicarts/icon-512.png',
];

self.addEventListener('install', function(event) {
  event.waitUntil(
    caches.open(CACHE_NAME).then(function(cache) {
      return cache.addAll(STATIC_ASSETS);
    })
  );
  self.skipWaiting();
});

self.addEventListener('activate', function(event) {
  event.waitUntil(
    caches.keys().then(function(names) {
      return Promise.all(
        names
          .filter(function(name) { return name !== CACHE_NAME; })
          .map(function(name) { return caches.delete(name); })
      );
    })
  );
  self.clients.claim();
});

self.addEventListener('fetch', function(event) {
  var url = new URL(event.request.url);

  // Cache-first for static assets
  if (url.pathname.indexOf('/static/') !== -1) {
    event.respondWith(
      caches.match(event.request).then(function(cached) {
        return cached || fetch(event.request).then(function(response) {
          var clone = response.clone();
          caches.open(CACHE_NAME).then(function(cache) {
            cache.put(event.request, clone);
          });
          return response;
        });
      })
    );
    return;
  }

  // Network-first for everything else
  event.respondWith(
    fetch(event.request).catch(function() {
      return caches.match(event.request);
    })
  );
});
