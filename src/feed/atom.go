// Atom 1.0 parser
package feed

import (
	"encoding/xml"
	"html"
	"io"
	"strings"
)

type atomFeed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	ID      string      `xml:"id"`
	Title   atomText    `xml:"title"`
	Links   atomLinks   `xml:"link"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	ID        string    `xml:"id"`
	Title     atomText  `xml:"title"`
	Summary   atomText  `xml:"summary"`
	Published string    `xml:"published"`
	Updated   string    `xml:"updated"`
	Links     atomLinks `xml:"link"`
	Content   atomText  `xml:"content"`
}

type atomText struct {
	Type string `xml:"type,attr"`
	Data string `xml:",chardata"`
	XML  string `xml:",innerxml"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type atomLinks []atomLink

func (a *atomText) String() string {
	data := a.Data
	if a.Type == "xhtml" {
		data = a.XML
	}
	return html.UnescapeString(strings.TrimSpace(data))
}

func (links atomLinks) First(rel string) string {
	for _, l := range links {
		if l.Rel == rel {
			return l.Href
		}
	}
	return ""
}

func ParseAtom(r io.Reader) (*Feed, error) {
	f := atomFeed{}

	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&f); err != nil {
		return nil, err
	}

	feed := &Feed{
		Title:   f.Title.String(),
		SiteURL: first(f.Links.First("alternate"), f.Links.First("")),
		FeedURL: f.Links.First("self"),
	}
	for _, e := range f.Entries {
		date, _ := dateParse(first(e.Published, e.Updated))
		imageUrl := ""
		podcastUrl := ""

		feed.Items = append(feed.Items, Item{
			GUID:       first(e.ID),
			Date:       date,
			URL:        first(e.Links.First("alternate"), f.Links.First("")),
			Title:      e.Title.String(),
			Content:    e.Content.String(),
			ImageURL:   imageUrl,
			PodcastURL: podcastUrl,
		})
	}
	return feed, nil
}
