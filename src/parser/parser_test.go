package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/expr-lang/expr"
)

// TestParseCases discovers every fixture under the "fixtures" directory and
// runs it as a subtest. Each fixture is a feed document (xml or json) whose
// leading comment block describes the case and lists assertions. The assertions
// are "@"-prefixed expressions evaluated with expr-lang against the parsed
// *Feed (exposed as the `feed` variable). The subtest name is the fixture path
// relative to "fixtures", e.g. "rss/basic.xml".
func TestParseCases(t *testing.T) {
	const root = "fixtures"
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		if strings.HasSuffix(path, ".xml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("cannot walk %s: %v", root, err)
	}
	if len(files) == 0 {
		t.Fatalf("no testcase files found under %s", root)
	}

	for _, f := range files {
		name := strings.TrimPrefix(filepath.ToSlash(f), root+"/")
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				runTestCaseFile(t, f)
			})
		})
	}
}

func runTestCaseFile(t *testing.T, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	comment, rest, ok := extractComment(string(data))
	if !ok {
		t.Fatalf("no leading comment block in %s", path)
	}
	assertions := extractAssertions(comment)

	var feed *Feed
	feed, err = ParseAndFix(strings.NewReader(rest), "", "")
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	env := map[string]any{
		"feed": feed,
		"time": map[string]any{
			"UTC":       time.UTC,
			"FixedZone": time.FixedZone,
			"Now":       time.Now,
			"Date": func(year, month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
				return time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
			},
		},
	}
	for _, assertion := range assertions {
		prg, cerr := expr.Compile(assertion, expr.Env(env))
		if cerr != nil {
			t.Fatalf("compile assertion %q: %v", assertion, cerr)
		}
		out, rerr := expr.Run(prg, env)
		if rerr != nil {
			t.Fatalf("run assertion %q: %v", assertion, rerr)
		}
		if !out.(bool) {
			t.Logf("feed=%#v\n", feed)
			t.Errorf("assertion failed: @ %s\n", assertion)
		}
	}
}

func extractComment(data string) (comment, rest string, ok bool) {
	s := strings.TrimLeft(data, " \t\r\n\ufeff")
	switch {
	case strings.HasPrefix(s, "<!--"):
		end := strings.Index(s, "-->")
		if end < 0 {
			return
		}
		return s[4:end], s[end+3:], true
	case strings.HasPrefix(s, "/*"):
		end := strings.Index(s, "*/")
		if end < 0 {
			return
		}
		return s[2:end], s[end+2:], true
	}
	return
}

func extractAssertions(comment string) (assertions []string) {
	indicator := "@"
	for line := range strings.SplitSeq(comment, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if expr, ok := strings.CutPrefix(trimmed, indicator); ok {
			assertions = append(assertions, strings.TrimSpace(expr))
			continue
		}
	}
	return
}
