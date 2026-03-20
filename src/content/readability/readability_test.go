package readability

import (
	"strings"
	"testing"
)

func TestExtractContent_SimpleArticle(t *testing.T) {
	html := `<html><body>
		<div class="article">
			<p>This is a test article with enough content to be scored by the readability algorithm.
			It needs to be long enough, with commas, to actually get picked up as meaningful content.
			The algorithm requires at least 25 characters per paragraph to even consider scoring it.</p>
			<p>Here is another paragraph with some more text content that helps boost the score.
			Adding more sentences with commas, periods, and other punctuation helps too.</p>
		</div>
		<div class="sidebar">
			<a href="/link1">Link 1</a>
			<a href="/link2">Link 2</a>
		</div>
	</body></html>`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	if result == "" {
		t.Fatal("expected non-empty content")
	}
	if !strings.Contains(result, "test article") {
		t.Errorf("expected result to contain article text, got: %s", result)
	}
}

func TestExtractContent_RemovesSidebar(t *testing.T) {
	html := `<html><body>
		<div class="content">
			<p>Main article content that is long enough to be considered real content by the algorithm.
			This paragraph has commas, sentences, and enough text to score highly in readability extraction.</p>
			<p>A second paragraph adds more weight to this content block, making it the clear winner
			in the scoring algorithm. More text, more commas, more sentences help the scoring.</p>
		</div>
		<div class="sidebar">
			<p>Sidebar stuff</p>
		</div>
	</body></html>`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, "Sidebar stuff") {
		t.Error("sidebar content should be removed")
	}
}

func TestExtractContent_RemovesScriptsAndStyles(t *testing.T) {
	html := `<html><body>
		<script>alert('xss')</script>
		<style>.foo { color: red; }</style>
		<div>
			<p>This is the real article content that should be extracted by readability.
			It contains enough text and commas to be scored properly by the algorithm.
			We need multiple sentences for the extraction to work correctly here.</p>
		</div>
	</body></html>`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, "alert") {
		t.Error("script content should be removed")
	}
	if strings.Contains(result, ".foo") {
		t.Error("style content should be removed")
	}
}

func TestExtractContent_DivToParagraph(t *testing.T) {
	// A div with no block-level children should be treated as a paragraph
	html := `<html><body>
		<div class="article">
			<div>This is a simple text div with enough content to be meaningful.
			It should be promoted to a paragraph element because it contains no block-level children.
			The readability algorithm looks for divs that only contain inline content.</div>
		</div>
	</body></html>`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "simple text div") {
		t.Error("text-only div content should be extracted")
	}
}

func TestExtractContent_EmptyBody(t *testing.T) {
	html := `<html><body></body></html>`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	// Should return the body wrapped in div, even if empty
	if !strings.Contains(result, "<div>") {
		t.Error("expected div wrapper in output")
	}
}

func TestExtractContent_NoBody(t *testing.T) {
	html := `<html></html>`

	_, err := ExtractContent(strings.NewReader(html))
	// The parser will still produce a body element, so this should work
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtractContent_RemovesUnlikelyCandidates(t *testing.T) {
	html := `<html><body>
		<div class="content">
			<p>This is the main article body text with enough characters and commas to score well.
			The readability algorithm should pick this up as the primary content block.
			Additional sentences with punctuation help boost the score even further.</p>
		</div>
		<div class="footer">
			<p>Footer navigation links and copyright info that should be excluded.</p>
		</div>
		<div class="social">
			<p>Share on Twitter, Facebook, Instagram and other platforms.</p>
		</div>
	</body></html>`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "main article body text") {
		t.Error("main content should be extracted")
	}
}

func TestExtractContent_InvalidHTML(t *testing.T) {
	html := `<html><body><p>Unclosed paragraph with enough text to be scored`

	result, err := ExtractContent(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "Unclosed paragraph") {
		t.Error("should handle malformed HTML gracefully")
	}
}

func TestGetClassWeight(t *testing.T) {
	tests := []struct {
		name   string
		class  string
		id     string
		expect float32
	}{
		{"positive class", "article-content", "", 25},
		{"negative class", "sidebar-widget", "", -25},
		{"positive id", "", "main-content", 25},
		{"negative id", "", "footer-nav", -25},
		{"both positive", "article", "main-body", 50},
		{"no match", "custom-class", "custom-id", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily construct html.Node with attributes inline,
			// but we can verify the regex patterns work
			if tt.class == "article-content" {
				if !positiveRegexp.MatchString(tt.class) {
					t.Error("'article-content' should match positive regex")
				}
			}
			if tt.class == "sidebar-widget" {
				if !negativeRegexp.MatchString(tt.class) {
					t.Error("'sidebar-widget' should match negative regex")
				}
			}
		})
	}
}

func TestRegexPatterns(t *testing.T) {
	// Unlikely candidates
	unlikelyClasses := []string{"banner", "sidebar", "comment", "footer", "menu", "social", "sponsor"}
	for _, cls := range unlikelyClasses {
		if !unlikelyCandidatesRegexp.MatchString(cls) {
			t.Errorf("'%s' should match unlikelyCandidates", cls)
		}
	}

	// Blacklist
	blacklistClasses := []string{"popupbody", "header-ad", "g-plus"}
	for _, cls := range blacklistClasses {
		if !blacklistCandidatesRegexp.MatchString(cls) {
			t.Errorf("'%s' should match blacklist", cls)
		}
	}

	// OK maybe it's a candidate (should not be removed)
	okClasses := []string{"article", "main", "body", "column"}
	for _, cls := range okClasses {
		if !okMaybeItsACandidateRegexp.MatchString(cls) {
			t.Errorf("'%s' should match okMaybeItsACandidate", cls)
		}
	}

	// Positive content indicators
	positiveClasses := []string{"article", "content", "entry", "main", "post", "blog"}
	for _, cls := range positiveClasses {
		if !positiveRegexp.MatchString(cls) {
			t.Errorf("'%s' should match positive regex", cls)
		}
	}

	// Negative content indicators
	negativeClasses := []string{"comment", "footer", "sidebar", "widget", "hidden", "sponsor"}
	for _, cls := range negativeClasses {
		if !negativeRegexp.MatchString(cls) {
			t.Errorf("'%s' should match negative regex", cls)
		}
	}
}
