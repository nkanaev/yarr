package opml

import (
	"encoding/xml"
	"io"
)

type opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Outlines []outline `xml:"body>outline"`
}

type outline struct {
	Type     string    `xml:"type,attr,omitempty"`
	Title    string    `xml:"text,attr"`
	FeedUrl  string    `xml:"xmlUrl,attr,omitempty"`
	SiteUrl  string    `xml:"htmlUrl,attr,omitempty"`
	Outlines []outline `xml:"outline,omitempty"`
}

func buildFolder(title string, outlines []outline) Folder {
	folder := Folder{Title: title}
	for _, outline := range outlines {
		if outline.Type == "rss" {
			folder.Feeds = append(folder.Feeds, Feed{
				Title:   outline.Title,
				FeedUrl: outline.FeedUrl,
				SiteUrl: outline.SiteUrl,
			})
		} else {
			subfolder := buildFolder(outline.Title, outline.Outlines)
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

	err := decoder.Decode(&val)
	if err != nil {
		return Folder{}, err
	}
	return buildFolder("", val.Outlines), nil
}
