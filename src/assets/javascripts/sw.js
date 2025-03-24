const cacheName = "yarr-assets";
const assetsToCache = [
  "./",
  "./sw.js",
  "./manifest.json",
  "./static/stylesheets/bootstrap.min.css",
  "./static/stylesheets/app.css",
  "./static/javascripts/vue.min.js",
  "./static/javascripts/api.js",
  "./static/javascripts/app.js",
  "./static/javascripts/key.js"
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches
      .open(cacheName)
      .then((cache) => cache.addAll(assetsToCache))
      .catch((error) =>
        console.error("Error caching assets during install:", error)
      )
  );
  // Activate the service worker as soon as installation is complete.
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  // Remove any old caches and immediately take control.
  event.waitUntil(
    caches.keys().then((cacheKeys) =>
      Promise.all(
        cacheKeys.map((key) => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      )
    )
  );
  self.clients.claim();
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.match(event.request).then((cachedResponse) => {
      return cachedResponse || fetch(event.request);
    })
  );
});
