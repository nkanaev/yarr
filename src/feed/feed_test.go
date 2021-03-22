package feed

import "testing"

func TestSniff(t *testing.T) {
	testcases := [][2]string{
		{
			`<?xml version="1.0"?><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"></rdf:RDF>`,
			"rss",
		},
		{
			`<?xml version="1.0"?><rss version="2.0"><channel></channel></rss>`,
			"rss",
		},
		{
			`<?xml version="1.0" encoding="utf-8"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`,
			"atom",
		},
		{
			`{}`,
			"json",
		},
		{
			`<!DOCTYPE html><html><head><title></title></head><body></body></html>`,
			"",
		},
	}
	for _, testcase := range testcases {
		have, _ := sniff(testcase[0])
		want := testcase[1]
		if want != have {
			t.Log(testcase[0])
			t.Errorf("Invalid format: want=%#v have=%#v", want, have)
		}
	}
}
