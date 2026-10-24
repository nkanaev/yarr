package parser

import "testing"

func TestDateParseNamedZoneOffset(t *testing.T) {
	cases := []struct {
		in       string
		wantUnix int64
	}{
		{"Mon, 02 Jan 2006 15:04:05 GMT", 1136214245},
		{"Mon, 02 Jan 2006 15:04:05 +0000", 1136214245},
		{"Mon, 02 Jan 2006 15:04:05 WET", 1136214245},
		{"Mon, 02 Jan 2006 15:04:05 EST", 1136232245},
		{"Mon, 02 Jan 2006 15:04:05 EDT", 1136228645},
		{"Mon, 02 Jan 2006 15:04:05 CET", 1136210645},
		{"Mon, 02 Jan 2006 15:04:05 JST", 1136181845},
		{"Thu, 17 Sep 2026 17:19:50 EDT", 1789679990},
	}
	for _, c := range cases {
		if got := dateParse(c.in).Unix(); got != c.wantUnix {
			t.Errorf("dateParse(%q).Unix() = %d, want %d", c.in, got, c.wantUnix)
		}
	}

	if got := dateParse("Mon, 02 Jan 2006 15:04:05 -0700 MST").Unix(); got != 1136239445 {
		t.Errorf("adjacent-offset layout regressed: got unix=%d, want 1136239445", got)
	}
	if got := dateParse("Mon, 02 Jan 2006 15:04:05 MST").Unix(); got != 1136214245 {
		t.Errorf("ambiguous MST path changed: got unix=%d, want 1136214245", got)
	}
}
