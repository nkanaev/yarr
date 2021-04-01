package silo

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	youtubeFrame = `<iframe src="https://www.youtube.com/embed/%s" width="560" height="315" frameborder="0" allowfullscreen></iframe>`
	vimeoFrame = `<iframe src="https://player.vimeo.com/video/%s" width="640" height="360" frameborder="0" allowfullscreen></iframe>`
	vimeoRegex = regexp.MustCompile(`\/(\d+)$`)
)

func VideoIFrame(link string) string {
	l, err := url.Parse(link)
	if err != nil {
		return ""
	}

	youtubeID := ""
	if l.Host == "www.youtube.com" && l.Path == "/watch" {
		youtubeID = l.Query().Get("v")
	} else if l.Host == "youtu.be" {
		youtubeID = strings.TrimLeft(l.Path, "/")
	}
	if youtubeID != "" {
		return fmt.Sprintf(youtubeFrame, youtubeID)
	}

	if l.Host == "vimeo.com" {
		if matches := vimeoRegex.FindStringSubmatch(l.Path); len(matches) > 0 {
			return fmt.Sprintf(vimeoFrame, matches[1])	
		}
	}
	return ""
}
