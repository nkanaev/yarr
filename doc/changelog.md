# upcoming

- (new) serve on unix socket (thanks to @rvighne)
- (fix) smooth scrolling on iOS (thanks to gatheraled)

# v2.5 (2025-03-26)

- (new) Fever API support (thanks to @icefed)
- (new) editable feed link (thanks to @adaszko)
- (new) switch to feed by clicking the title in the article page (thanks to @tarasglek for suggestion)
- (new) support multiple media links
- (new) next/prev article navigation buttons (thanks to @tillcash)
- (fix) duplicate articles caused by the same feed addition (thanks to @adaszko)
- (fix) relative article links (thanks to @adazsko for the report)
- (fix) atom article links stored in id element (thanks to @adazsko for the report)
- (fix) parsing atom feed titles (thanks to @wnh)
- (fix) sorting same-day batch articles (thanks to @lamescholar for the report)
- (fix) showing login page in the selected theme (thanks to @feddiriko for the report)
- (fix) parsing atom feeds with html elements (thanks to @tillcash & @toBeOfUse for the report, @krkk for the fix)
- (fix) parsing feeds with missing guids (thanks to @hoyii for the report)
- (fix) sending actual client version to servers (thanks to @aidanholm)
- (fix) error caused by missing config dir (thanks to @timster)
- (etc) load external images with no-referrer policy (thanks to @tillcash for the report)
- (etc) open external links with no-referrer policy (thanks to @donovanglover)
- (etc) show article content in the list if title is missing (thanks to @asimpson for suggestion)
- (etc) accessibility improvements (thanks to @tseykovets)

# v2.4 (2023-08-15)

- (new) ARM build support (thanks to @tillcash & @fenuks)
- (new) auth configuration via param or env variable (thanks to @pierreprinetti)
- (new) web app manifest for an app-like experience on mobile (thanks to @qbit)
- (fix) concurrency issue crashing the app (thanks to @quoing)
- (fix) favicon visibility in dark mode (thanks to @caycaycarly for the report)
- (fix) autoloading more articles not working in certain edge cases (thanks to @fenuks for the report)
- (fix) handle Google URL redirects in "Read Here" (thanks to @cubbei for discovery)
- (fix) handle failures to extract content in "Read Here" (thanks to @grigio for the report)
- (fix) article view width for high resolution screens (thanks to @whaler-ragweed for the report)
- (fix) make newly added feed searchable (thanks to @BMorearty for the report)
- (fix) feed/article selection accessibility via arrow keys (thanks to @grigio and @tillcash)
- (fix) keyboard shortcuts in Firefox (thanks to @kaloyan13)
- (fix) keyboard shortcuts in non-English layouts (thanks to @kaloyan13)
- (fix) sorting articles with timezone information (thanks to @x2cf)
- (fix) handling links set in guid only for certain feeds (thanks to @adaszko for the report)
- (fix) crashes caused by feed icon endpoint (thanks to @adaszko)

# v2.3 (2022-05-03)

- (fix) handling encodings (thanks to @f100024 & @fserb)
- (fix) parsing xml feeds with illegal characters (thanks to @stepelu for the report)
- (fix) old articles reappearing as unread (thanks to @adaszko for the report)
- (fix) item list scrolling issue on large screens (thanks to @bielej for the report)
- (fix) keyboard shortcuts color in dark mode (thanks to @John09f9 for the report)
- (etc) autofocus when adding a new feed (thanks to @lakuapik)

# v2.2 (2021-11-20)

- (fix) windows console support (thanks to @dufferzafar for the report)
- (fix) remove html tags from article titles (thanks to Alex Went for the report)
- (etc) autoselect current folder when adding a new feed (thanks to @krkk)
- (etc) folder/feed settings menu available across all filters

# v2.1 (2021-08-16)

- (new) configuration via env variables
- (fix) missing `content-type` headers (thanks to @verahawk for the report)
- (fix) handle opml files not following the spec (thanks to @huangnauh for the report)
- (fix) pagination in unread/starred feeds (thanks to @Farow for the report)
- (fix) handling feeds with non-utf8 encodings (thanks to @fserb for the report)
- (fix) errors caused by empty feeds (thanks to @decke)
- (fix) recognize all audio mime types as podcasts (thanks to @krkk)
- (fix) ui tweaks (thanks to @Farow)

# v2.0 (2021-04-18)

- (new) user interface tweaks
- (new) feed parser fully rewritten
- (new) show youtube/vimeo iframes in "read here"
- (new) keyboard shortcuts for article scrolling & toggling "read here"
- (new) more options for auto-refresh intervals
- (fix) `-base` not serving static files (thanks to @vfaronov)
- (etc) 3rd-party dependencies reduced to the bare minimum

special thanks to @tillcash for feedback & suggestions.

# v1.4 (2021-03-11)

- (new) keyboard shortcuts (thanks to @Duarte-Dias)
- (new) show podcast audio
- (fix) deleting feeds
- (etc) minor ui tweaks & changes

# v1.3 (2021-02-18)

- (fix) log out functionality if authentication is set
- (fix) import opml if authentication is set
- (fix) login page if authentication is set (thanks to @einschmidt)

# v1.2 (2021-02-11)

- (new) autorefresh rate
- (new) reduced bandwidth usage via stateful http headers `last-modified/etag`
- (new) show feed errors in feed management modal
- (new) `-open` flag for automatically opening the server url
- (new) `-base` flag for serving urls under non-root path (thanks to @hcl)
- (new) `-auth-file` flag for authentication
- (new) `-cert-file` & `-key-file` flags for TLS
- (fix) wrapping long words in the ui to prevent vertical scroll
- (fix) increased toolbar height in mobile/tablet layout (thanks to @einschmidt)

# v1.1 (2020-10-05)

- (new) responsive design
- (fix) server crash on favicon fetch timeout (reported by @minioin)
- (fix) handling byte order marks in feeds (reported by @ilaer)
- (fix) deleting a feed raises exception in the ui if the feed's items are shown.

# v1.0 (2020-09-24)

Initial Release
