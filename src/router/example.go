package router
/*
func do() {
	server := NewServer(db, worker)

	router := NewRouter()

	router.Use(AuthMiddleware())
	router.Use(CorsMiddleware())

	router.For("/", server.index)
	router.For("/static/*path", server.static)
	router.For("/api/status", server.status)
	router.For("/api/folders", server.folderlist)
	router.For("/api/folders/:id", server.folder)
	router.For("/api/feeds", server.feedlist)
	router.For("/api/feeds/refresh", server.feedsRefresh)
	router.For("/api/feeds/errors", server.feedsErrors)
	router.For("/api/feeds/:id/icon", server.feedsIcons)
	router.For("/api/feeds/:id", server.feed)
	router.For("/api/items", server.itemlist)
	router.For("/api/items/:id", server.item)
	router.For("/api/settings", server.settings)
	router.For("/opml/import", server.opmlImport)
	router.For("/opml/export", server.opmlExport)
	router.For("/page", server.pagecrawl)
	router.For("/logout", server.logout)

	httpserver := &http.Server{Addr: server.Addr(), Handler: router}
	httpserver.ListenAndServe()
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
}
*/
