package parser

type media struct {
	MediaGroups []mediaGroup `xml:"http://search.yahoo.com/mrss/ group"`

	MediaThumbnails   []mediaThumbnail   `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaDescriptions []mediaDescription `xml:"http://search.yahoo.com/mrss/ description"`
}

type mediaGroup struct {
	MediaThumbnails   []mediaThumbnail   `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaDescriptions []mediaDescription `xml:"http://search.yahoo.com/mrss/ description"`
}

type mediaContent struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	FileSize string `xml:"fileSize,attr"`
	Medium   string `xml:"medium,attr"`
}

type mediaThumbnail struct {
	URL string `xml:"url,attr"`
}

type mediaDescription struct {
	Type        string `xml:"type,attr"`
	Description string `xml:",chardata"`
}

func (m *media) firstMediaThumbnail() string {
	for _, t := range m.MediaThumbnails {
		return t.URL
	}
	for _, g := range m.MediaGroups {
		for _, t := range g.MediaThumbnails {
			return t.URL
		}
	}
	return ""
}

func (m *media) firstMediaDescription() string {
	for _, d := range m.MediaDescriptions {
		return plain2html(d.Description)
	}
	for _, g := range m.MediaGroups {
		for _, d := range g.MediaDescriptions {
			return plain2html(d.Description)
		}
	}
	return ""
}
