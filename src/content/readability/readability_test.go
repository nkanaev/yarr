package readability

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractContent(t *testing.T) {
	files, err := filepath.Glob("testdata/*.html")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".html")
		if strings.HasSuffix(name, ".want") || strings.HasSuffix(name, ".have") {
			continue
		}

		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			have, err := ExtractContent(strings.NewReader(string(raw)))
			if err != nil {
				t.Fatal(err)
			}

			havePath := filepath.Join("testdata", name+".have.html")
			if err := os.WriteFile(havePath, []byte(have), 0o644); err != nil {
				t.Fatal(err)
			}

			wantPath := filepath.Join("testdata", name+".want.html")
			want, err := os.ReadFile(wantPath)
			if os.IsNotExist(err) {
				wantPath = filepath.Join("testdata", name+".want.html")
				want, err = os.ReadFile(wantPath)
				if os.IsNotExist(err) {
					if err := os.WriteFile(wantPath, []byte(have), 0o644); err != nil {
						t.Fatal(err)
					}
					return
				}
			}
			if err != nil {
				t.Fatal(err)
			}

			if strings.TrimSuffix(string(want), "\n") != strings.TrimSuffix(have, "\n") {
				t.Errorf("wrong output:\nwant: %s\nhave: %s", wantPath, havePath)
				return
			}

			if err := os.Remove(havePath); err != nil {
				t.Fatal(err)
			}
		})
	}
}
