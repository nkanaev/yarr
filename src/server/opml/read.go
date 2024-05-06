package opml

import (
	"encoding/xml"
	"io"

	"golang.org/x/net/html/charset"
)

type opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Outlines []outline `xml:"body>outline"`
}

type outline struct {
	Type        string    `xml:"type,attr,omitempty"`
	Title       string    `xml:"text,attr"`
	Title2      string    `xml:"title,attr,omitempty"`
	FeedUrl     string    `xml:"xmlUrl,attr,omitempty"`
	SiteUrl     string    `xml:"htmlUrl,attr,omitempty"`
	CustomOrder string    `xml:"customOrder,attr,omitempty"`
	Outlines    []outline `xml:"outline,omitempty"`
}

func buildFolder(title string, outlines []outline) Folder {
	folder := Folder{Title: title}
	for _, outline := range outlines {
		if outline.Type == "rss" || outline.FeedUrl != "" {
			folder.Feeds = append(folder.Feeds, Feed{
				Title:       outline.Title,
				FeedUrl:     outline.FeedUrl,
				SiteUrl:     outline.SiteUrl,
				CustomOrder: outline.CustomOrder,
			})
		} else {
			title := outline.Title
			if title == "" {
				title = outline.Title2
			}
			subfolder := buildFolder(title, outline.Outlines)
			folder.Folders = append(folder.Folders, subfolder)
		}
	}
	return folder
}

func Parse(r io.Reader) (Folder, error) {
	val := new(opml)
	decoder := xml.NewDecoder(r)
	decoder.Entity = xml.HTMLEntity
	decoder.Strict = false
	decoder.CharsetReader = charset.NewReaderLabel

	err := decoder.Decode(&val)
	if err != nil {
		return Folder{}, err
	}
	return buildFolder("", val.Outlines), nil
}
