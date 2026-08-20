package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/content/htmlutil"
	"github.com/nkanaev/yarr/src/content/readability"
	"github.com/nkanaev/yarr/src/content/sanitizer"
	"github.com/nkanaev/yarr/src/content/silo"
	"github.com/nkanaev/yarr/src/server/auth"
	"github.com/nkanaev/yarr/src/server/gzip"
	"github.com/nkanaev/yarr/src/server/opml"
	"github.com/nkanaev/yarr/src/storage/model"
	"github.com/nkanaev/yarr/src/worker"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("failed to write JSON:", err)
	}
}

func writeHTML(w http.ResponseWriter, status int, tmpl *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("failed to write HTML:", err)
	}
}

func (s *Server) handler() http.Handler {
	staticFS := http.FileServer(http.FS(assets.StaticFS()))

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/{$}", s.handleIndex)
	publicMux.HandleFunc("/login", s.handleLogin)
	publicMux.HandleFunc("/static/{path...}", http.StripPrefix("/static/", staticFS).ServeHTTP)
	publicMux.HandleFunc("/fever/", s.handleFever)
	publicMux.HandleFunc("/manifest.json", s.handleManifest)

	secureMux := http.NewServeMux()
	secureMux.HandleFunc("/api/status", s.handleStatus)
	secureMux.HandleFunc("/api/folders", s.handleFolderList)
	secureMux.HandleFunc("/api/folders/{id}", s.handleFolder)
	secureMux.HandleFunc("/api/feeds", s.handleFeedList)
	secureMux.HandleFunc("/api/feeds/refresh", s.handleFeedRefresh)
	secureMux.HandleFunc("/api/feeds/errors", s.handleFeedErrors)
	secureMux.HandleFunc("/api/feeds/{id}", s.handleFeed)
	secureMux.HandleFunc("/api/items", s.handleItemList)
	secureMux.HandleFunc("/api/items/{id}", s.handleItem)
	secureMux.HandleFunc("/api/settings", s.handleSettings)
	secureMux.HandleFunc("/opml/import", s.handleOPMLImport)
	secureMux.HandleFunc("/opml/export", s.handleOPMLExport)
	secureMux.HandleFunc("/page", s.handlePageCrawl)
	secureMux.HandleFunc("/logout", s.handleLogout)

	var protected http.Handler = secureMux
	if s.Username != "" && s.Password != "" {
		protected = auth.Middleware(s.Username, s.Password, secureMux)
	}

	dispatch := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := secureMux.Handler(r)
		if pattern != "" {
			protected.ServeHTTP(w, r)
		} else {
			publicMux.ServeHTTP(w, r)
		}
	})

	if s.BasePath != "" {
		baseDispatch := dispatch
		dispatch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == s.BasePath {
				http.Redirect(w, r, s.BasePath+"/", http.StatusFound)
				return
			}
			if !strings.HasPrefix(r.URL.Path, s.BasePath) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			r2 := r.Clone(r.Context())
			r2.URL.Path = strings.TrimPrefix(r.URL.Path, s.BasePath)
			baseDispatch.ServeHTTP(w, r2)
		})
	}

	return gzip.Middleware(dispatch)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := false
	requiresAuth := false
	if s.Username == "" && s.Password == "" {
		isAuthenticated = true
	} else {
		requiresAuth = true
		isAuthenticated = auth.IsAuthenticated(r, s.Username, s.Password)
	}

	settings := s.db.GetSettings()
	if !isAuthenticated {
		settings = model.Settings{
			Language:  settings.Language,
			ThemeName: settings.ThemeName,
		}
	}

	writeHTML(w, http.StatusOK, assets.Templates().Lookup("index.html"), map[string]any{
		"settings":      settings.Map(),
		"authenticated": isAuthenticated,
		"requiresAuth":  requiresAuth,
	})
}

func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"$schema":     "https://json.schemastore.org/web-manifest-combined.json",
		"name":        "yarr!",
		"short_name":  "yarr",
		"description": "yet another rss reader",
		"display":     "standalone",
		"start_url":   "/" + strings.TrimPrefix(s.BasePath, "/"),
		"icons": []map[string]any{
			{
				"src":   s.BasePath + "/static/favicon.png",
				"sizes": "64x64",
				"type":  "image/png",
			},
		},
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"running": s.worker.FeedsPending(),
		"stats":   s.db.FeedStats(),
	})
}

func (s *Server) handleFolderList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.db.ListFolders())
	case http.MethodPost:
		var body FolderCreateForm
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(body.Title) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Folder title missing."})
			return
		}
		writeJSON(w, http.StatusCreated, s.db.CreateFolder(body.Title))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFolder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var body FolderUpdateForm
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.db.UpdateFolder(id, model.UpdateFolderParams{
			Title:      body.Title,
			IsExpanded: body.IsExpanded,
		})
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		s.db.DeleteFolder(id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFeedRefresh(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.worker.RefreshFeeds()
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFeedErrors(w http.ResponseWriter, r *http.Request) {
	errors := make(map[int64]string)
	states, err := s.db.ListFeedStates()
	if err == nil {
		for _, state := range states {
			if state.LastError != "" {
				errors[state.FeedID] = state.LastError
			}
		}
	}
	writeJSON(w, http.StatusOK, errors)
}

func (s *Server) handleFeedList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.db.ListFeeds())
	case http.MethodPost:
		var form FeedCreateForm
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := worker.DiscoverFeed(form.Url)
		switch {
		case err != nil:
			log.Printf("Faild to discover feed for %s: %s", form.Url, err)
			writeJSON(w, http.StatusOK, map[string]string{"status": "notfound"})
		case len(result.Sources) > 0:
			writeJSON(
				w,
				http.StatusOK,
				map[string]any{"status": "multiple", "choice": result.Sources},
			)
		case result.Feed != nil:
			title := result.Feed.Title
			if form.TitleOverride != "" {
				title = form.TitleOverride
			}
			feed := s.db.CreateFeed(model.CreateFeedParams{
				Title:    title,
				Link:     result.Feed.SiteURL,
				FeedLink: result.FeedLink,
				FolderID: form.FolderID,
			})
			items := worker.ConvertItems(result.Feed.Items, *feed)
			if len(items) > 0 {
				s.db.CreateItems(items)
			}
			s.worker.FindFeedFavicon(*feed)

			writeJSON(w, http.StatusOK, map[string]any{
				"status": "success",
				"feed":   feed,
			})
		default:
			writeJSON(w, http.StatusOK, map[string]string{"status": "notfound"})
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		feed := s.db.GetFeed(id)
		if feed == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		params := model.UpdateFeedParams{}
		if title, ok := body["title"]; ok {
			if reflect.TypeOf(title).Kind() == reflect.String {
				t := title.(string)
				params.Title = &t
			}
		}
		if f_id, ok := body["folder_id"]; ok {
			if f_id == nil {
				params.FolderID = model.SetNullable[int64](nil)
			} else if reflect.TypeOf(f_id).Kind() == reflect.Float64 {
				folderId := int64(f_id.(float64))
				params.FolderID = model.SetNullable(&folderId)
			}
		}
		if link, ok := body["feed_link"]; ok {
			if reflect.TypeOf(link).Kind() == reflect.String {
				l := link.(string)
				params.FeedLink = &l
			}
		}
		s.db.UpdateFeed(id, params)
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		s.db.DeleteFeed(id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		item := s.db.GetItem(id)
		if item == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// runtime fix for relative links
		if !htmlutil.IsAPossibleLink(item.Link) {
			if feed := s.db.GetFeed(item.FeedId); feed != nil {
				item.Link = htmlutil.AbsoluteUrl(item.Link, feed.Link)
			}
		}

		item.Content = sanitizer.Sanitize(item.Link, item.Content)
		for i, link := range item.MediaLinks {
			item.MediaLinks[i].Description = sanitizer.Sanitize(item.Link, link.Description)
		}

		writeJSON(w, http.StatusOK, item)
	case http.MethodPut:
		var body ItemUpdateForm
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Status != nil {
			s.db.UpdateItemStatus(id, *body.Status)
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleItemList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		perPage := 20
		query := r.URL.Query()

		filter := model.ItemFilter{}
		if folderID, err := strconv.ParseInt(query.Get("folder_id"), 10, 64); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := strconv.ParseInt(query.Get("feed_id"), 10, 64); err == nil {
			filter.FeedID = &feedID
		}
		if after, err := strconv.ParseInt(query.Get("after"), 10, 64); err == nil {
			filter.After = &after
		}
		if status := query.Get("status"); len(status) != 0 {
			statusValue := model.StatusValues[status]
			filter.Status = &statusValue
		}
		if search := query.Get("search"); len(search) != 0 {
			filter.Search = &search
		}
		newestFirst := query.Get("oldest_first") != "true"

		items := s.db.ListItems(filter, perPage+1, newestFirst, true)
		hasMore := false
		if len(items) == perPage+1 {
			hasMore = true
			items = items[:perPage]
		}

		for i, item := range items {
			if item.Title == "" {
				text := htmlutil.ExtractText(item.Content)
				items[i].Title = htmlutil.TruncateText(text, 140)
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"list":     items,
			"has_more": hasMore,
		})
	case http.MethodPut:
		filter := model.MarkFilter{}

		query := r.URL.Query()
		if folderID, err := strconv.ParseInt(query.Get("folder_id"), 10, 64); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := strconv.ParseInt(query.Get("feed_id"), 10, 64); err == nil {
			filter.FeedID = &feedID
		}
		s.db.MarkItemsRead(filter)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.db.GetSettings())
	case http.MethodPut:
		var params model.UpdateSettingsParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if s.db.UpdateSettings(params) {
			if params.RefreshRate != nil {
				s.worker.SetRefreshRate(s.db.GetSettings().RefreshRate)
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOPMLImport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		file, _, err := r.FormFile("opml")
		if err != nil {
			log.Print(err)
			return
		}
		doc, err := opml.Parse(file)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, f := range doc.Feeds {
			s.db.CreateFeed(model.CreateFeedParams{
				Title:    f.Title,
				Link:     f.SiteUrl,
				FeedLink: f.FeedUrl,
			})
		}
		for _, f := range doc.Folders {
			folder := s.db.CreateFolder(f.Title)
			for _, ff := range f.AllFeeds() {
				s.db.CreateFeed(model.CreateFeedParams{
					Title:    ff.Title,
					Link:     ff.SiteUrl,
					FeedLink: ff.FeedUrl,
					FolderID: &folder.Id,
				})
			}
		}

		s.worker.RefreshFeeds()

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOPMLExport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		filename := fmt.Sprintf("subscriptions_%s.opml", time.Now().Format("2006-01-02_15-04-05"))
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		doc := opml.Folder{}

		feedsByFolderID := make(map[int64][]*model.Feed)
		for _, feed := range s.db.ListFeeds() {
			if feed.FolderId == nil {
				doc.Feeds = append(doc.Feeds, opml.Feed{
					Title:   feed.Title,
					FeedUrl: feed.FeedLink,
					SiteUrl: feed.Link,
				})
			} else {
				id := *feed.FolderId
				feedsByFolderID[id] = append(feedsByFolderID[id], &feed)
			}
		}

		for _, folder := range s.db.ListFolders() {
			folderFeeds := feedsByFolderID[folder.Id]
			if len(folderFeeds) == 0 {
				continue
			}
			opmlfolder := opml.Folder{Title: folder.Title}
			for _, feed := range folderFeeds {
				opmlfolder.Feeds = append(opmlfolder.Feeds, opml.Feed{
					Title:   feed.Title,
					FeedUrl: feed.FeedLink,
					SiteUrl: feed.Link,
				})
			}
			doc.Folders = append(doc.Folders, opmlfolder)
		}

		w.Write([]byte(doc.OPML()))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handlePageCrawl(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	url = silo.RedirectURL(url)

	if content := silo.VideoIFrame(url); content != "" {
		writeJSON(w, http.StatusOK, map[string]string{
			"content": sanitizer.Sanitize(url, content),
		})
		return
	}
	if isInternalFromURL(url) {
		log.Printf("attempt to access internal IP %s from %s", url, r.RemoteAddr)
		return
	}

	body, err := worker.GetBody(url)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	content, err := readability.ExtractContent(strings.NewReader(body))
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{
			"content": "error: " + err.Error(),
		})
		return
	}
	content = sanitizer.Sanitize(url, content)
	writeJSON(w, http.StatusOK, map[string]string{
		"content": content,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		if auth.StringsEqual(username, s.Username) && auth.StringsEqual(password, s.Password) {
			auth.Authenticate(w, s.Username, s.Password, s.BasePath)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, s.BasePath)
	w.WriteHeader(http.StatusNoContent)
}
