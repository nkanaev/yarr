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
	srcfeed := atomFeed{}

	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&srcfeed); err != nil {
		return nil, err
	}

	dstfeed := &Feed{
		Title:   srcfeed.Title.String(),
		SiteURL: firstNonEmpty(srcfeed.Links.First("alternate"), srcfeed.Links.First("")),
	}
	for _, srcitem := range srcfeed.Entries {
		imageUrl := ""
		podcastUrl := ""

		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:       firstNonEmpty(srcitem.ID),
			Date:       dateParse(firstNonEmpty(srcitem.Published, srcitem.Updated)),
			URL:        firstNonEmpty(srcitem.Links.First("alternate"), srcfeed.Links.First("")),
			Title:      srcitem.Title.String(),
			Content:    srcitem.Content.String(),
			ImageURL:   imageUrl,
			PodcastURL: podcastUrl,
		})
	}
	return dstfeed, nil
}
