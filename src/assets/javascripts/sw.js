'use strict';

self.addEventListener('install', e => {
/*
    e.waitUntil(
        caches.open('static').then(cache => {
            return cache.addAll([
                "/",

                "/static/javascripts/vue.min.js",
                "/static/javascripts/app.js	",
                "/static/javascripts/api.js",
                "/static/javascripts/key.js	",

                "/static/stylesheets/bootstrap.min.css",
                "/static/stylesheets/app.css",

                "/static/graphicarts/icon.png",
                "/static/graphicarts/icon.ico"
            ]);
        })
    );
*/
});

self.addEventListener('fetch', e => {
/*
    e.respondWith(
        caches.match(e.request).then(response => {
            return response || fetch(e.request);
        })
    );
*/
});
