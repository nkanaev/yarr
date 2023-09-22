const VERSION = "v2.4"
const APP_STATIC_RESOURCES = [
  "/",
  "/static/stylesheets/bootstrap.min.css",
  "/static/stylesheets/app.css",
  "/static/graphicarts/favicon.svg",
  "/static/graphicarts/favicon.png",
]
const CACHE_NAME = `yarr-${VERSION}`;

self.addEventListener("install", (e) => {
  e.waitUntil((async () => {
      const cache = await caches.open(CACHE_NAME)
      await cache.addAll(APP_STATIC_RESOURCES)
    })()
  )
})

// delete old caches on activate
self.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      const names = await caches.keys()
      await Promise.all(
        names.map((name) => {
          if (name !== CACHE_NAME) {
            return caches.delete(name)
          }
        }),
      )
      await clients.claim()
    })(),
  )
})
