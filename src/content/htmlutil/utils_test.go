package htmlutil

import "testing"

func TestExtractText(t *testing.T) {
	testcases := [][2]string {
		{"hello", "<div>hello</div>"},
		{"hello world", "<div>hello</div> world"},
		{"helloworld", "<div>hello</div>world"},
		{"hello world", "hello <div>world</div>"},
		{"helloworld", "hello<div>world</div>"},
		{"hello world!", "hello <div>world</div>!"},
		{"hello world !", "hello <div>   world\r\n </div>!"},
	}
	for _, testcase := range testcases {
		want := testcase[0]
		base := testcase[1]
		have := ExtractText(base)
		if want != have {
			t.Logf("base: %#v\n", base)
			t.Logf("want: %#v\n", want)
			t.Logf("have: %#v\n", have)
			t.Fail()
		}
	}
}
