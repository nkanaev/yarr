package silo

import "testing"

func TestYoutubeIframe(t *testing.T) {
	links := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ",
	}
	for _, link := range links {
		have := VideoIFrame(link)	
		want := `<iframe src="https://www.youtube.com/embed/dQw4w9WgXcQ" width="560" height="315" frameborder="0" allowfullscreen></iframe>`
		if have != want {
			t.Logf("want: %s", want)
			t.Logf("have: %s", have)
			t.Fail()
		}
	}
}

func TestVimeoIframe(t *testing.T) {
	links := []string{
		"https://vimeo.com/channels/staffpicks/526381128",
		"https://vimeo.com/526381128",
	}
	for _, link := range links {
		have := VideoIFrame(link)	
		want := `<iframe src="https://player.vimeo.com/video/526381128" width="640" height="360" frameborder="0" allowfullscreen></iframe>`
		if have != want {
			t.Logf("want: %s", want)
			t.Logf("have: %s", have)
			t.Fail()
		}
	}
}
