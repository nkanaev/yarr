// Atom 1.0 parser
package parser

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/nkanaev/yarr/src/content/htmlutil"
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
	Content   atomText  `xml:"http://www.w3.org/2005/Atom content"`
	OrigLink  string    `xml:"http://rssnamespace.org/feedburner/ext/1.0 origLink"`

	media
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

func (a *atomText) Text() string {
	if a.Type == "html" {
		return htmlutil.ExtractText(a.Data)
	} else if a.Type == "xhtml" {
		return htmlutil.ExtractText(a.XML)
	}
	return a.Data
}

func (a *atomText) String() string {
	data := a.Data
	if a.Type == "xhtml" {
		data = a.XML
	}
	return strings.TrimSpace(data)
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

	decoder := xmlDecoder(r)
	if err := decoder.Decode(&srcfeed); err != nil {
		return nil, err
	}

	dstfeed := &Feed{
		Title:   srcfeed.Title.String(),
		SiteURL: firstNonEmpty(srcfeed.Links.First("alternate"), srcfeed.Links.First("")),
	}
	for _, srcitem := range srcfeed.Entries {
		linkFromID := ""
		guidFromID := ""
		if htmlutil.IsAPossibleLink(srcitem.ID) {
			linkFromID = srcitem.ID
			guidFromID = srcitem.ID + "::" + srcitem.Updated
		}

		mediaLinks := srcitem.mediaLinks()

		link := firstNonEmpty(srcitem.OrigLink, srcitem.Links.First("alternate"), srcitem.Links.First(""), linkFromID)
		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:       firstNonEmpty(srcitem.ID, guidFromID, link),
			Date:       dateParse(firstNonEmpty(srcitem.Published, srcitem.Updated)),
			URL:        link,
			Title:      srcitem.Title.Text(),
			Content:    firstNonEmpty(srcitem.Content.String(), srcitem.Summary.String(), srcitem.firstMediaDescription()),
			MediaLinks: mediaLinks,
		})
	}
	return dstfeed, nil
}
