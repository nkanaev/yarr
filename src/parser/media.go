package parser

import (
	"strings"
)

type media struct {
	MediaGroups       []mediaGroup       `xml:"http://search.yahoo.com/mrss/ group"`
	MediaContents     []mediaContent     `xml:"http://search.yahoo.com/mrss/ content"`
	MediaThumbnails   []mediaThumbnail   `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaDescriptions []mediaDescription `xml:"http://search.yahoo.com/mrss/ description"`
}

type mediaGroup struct {
	MediaContent      []mediaContent     `xml:"http://search.yahoo.com/mrss/ content"`
	MediaThumbnails   []mediaThumbnail   `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaDescriptions []mediaDescription `xml:"http://search.yahoo.com/mrss/ description"`
}

type mediaContent struct {
	MediaThumbnails  []mediaThumbnail `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaType        string           `xml:"type,attr"`
	MediaMedium      string           `xml:"medium,attr"`
	MediaURL         string           `xml:"url,attr"`
	MediaDescription mediaDescription `xml:"http://search.yahoo.com/mrss/ description"`
}

type mediaThumbnail struct {
	URL string `xml:"url,attr"`
}

type mediaDescription struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

func (m *media) firstMediaDescription() string {
	for _, d := range m.MediaDescriptions {
		return plain2html(d.Text)
	}
	for _, g := range m.MediaGroups {
		for _, d := range g.MediaDescriptions {
			return plain2html(d.Text)
		}
	}
	return ""
}

func (m *media) mediaLinks() []MediaLink {
	links := make([]MediaLink, 0)
	for _, thumbnail := range m.MediaThumbnails {
		links = append(links, MediaLink{URL: thumbnail.URL, Type: "image"})
	}
	for _, group := range m.MediaGroups {
		for _, thumbnail := range group.MediaThumbnails {
			links = append(links, MediaLink{
				URL:  thumbnail.URL,
				Type: "image",
			})
		}
		for _, content := range group.MediaContent {
			links = append(links, content.mediaLinks()...)
		}
	}
	for _, content := range m.MediaContents {
		links = append(links, content.mediaLinks()...)
	}
	if len(links) == 0 {
		return nil
	}
	return links
}

func (m mediaContent) mediaLinks() []MediaLink {
	links := make([]MediaLink, 0)
	if m.MediaURL != "" {
		url := m.MediaURL
		description := m.MediaDescription.Text
		if strings.HasPrefix(m.MediaType, "image/") {
			links = append(links, MediaLink{URL: url, Type: "image", Description: description})
		} else if strings.HasPrefix(m.MediaType, "audio/") {
			links = append(links, MediaLink{URL: url, Type: "audio", Description: description})
		} else if strings.HasPrefix(m.MediaType, "video/") {
			links = append(links, MediaLink{URL: url, Type: "video", Description: description})
		} else if m.MediaMedium == "image" || m.MediaMedium == "audio" || m.MediaMedium == "video" {
			links = append(
				links,
				MediaLink{URL: url, Type: m.MediaMedium, Description: description},
			)
		} else {
			if len(m.MediaThumbnails) > 0 {
				links = append(links, MediaLink{
					URL:  m.MediaThumbnails[0].URL,
					Type: "image",
				})
			}
		}
	}
	for _, thumbnail := range m.MediaThumbnails {
		links = append(links, MediaLink{
			URL:  thumbnail.URL,
			Type: "image",
		})
	}
	return links
}
