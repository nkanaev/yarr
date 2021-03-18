module github.com/nkanaev/yarr

go 1.16

require (
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/mmcdole/gofeed v1.0.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13
)

replace github.com/mmcdole/gofeed => ./src/gofeed
