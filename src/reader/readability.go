// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package reader

import (
	"bytes"
	"fmt"
	"io"
	//"log"
	"math"
	"regexp"
	"strings"

	"github.com/nkanaev/yarr/src/htmlutil"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	defaultTagsToScore = "section,h2,h3,h4,h5,h6,p,td,pre,div"
)

var (
	divToPElementsRegexp = regexp.MustCompile(`(?i)<(a|blockquote|dl|div|img|ol|p|pre|table|ul)`)
	sentenceRegexp       = regexp.MustCompile(`\.( |$)`)

	blacklistCandidatesRegexp  = regexp.MustCompile(`(?i)popupbody|-ad|g-plus`)
	okMaybeItsACandidateRegexp = regexp.MustCompile(`(?i)and|article|body|column|main|shadow`)
	unlikelyCandidatesRegexp   = regexp.MustCompile(`(?i)banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|foot|header|legends|menu|modal|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote`)

	negativeRegexp = regexp.MustCompile(`(?i)hidden|^hid$|hid$|hid|^hid |banner|combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|modal|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|tool|widget|byline|author|dateline|writtenby|p-author`)
	positiveRegexp = regexp.MustCompile(`(?i)article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story`)
)

type candidate struct {
	selection *goquery.Selection
	score     float32
}

func (c *candidate) Node() *html.Node {
	return c.selection.Get(0)
}

type scorelist map[*html.Node]float32

func (c *candidate) String() string {
	id, _ := c.selection.Attr("id")
	class, _ := c.selection.Attr("class")

	if id != "" && class != "" {
		return fmt.Sprintf("%s#%s.%s => %f", c.Node().DataAtom, id, class, c.score)
	} else if id != "" {
		return fmt.Sprintf("%s#%s => %f", c.Node().DataAtom, id, c.score)
	} else if class != "" {
		return fmt.Sprintf("%s.%s => %f", c.Node().DataAtom, class, c.score)
	}

	return fmt.Sprintf("%s => %f", c.Node().DataAtom, c.score)
}

type candidateList map[*html.Node]*candidate

func (c candidateList) String() string {
	var output []string
	for _, candidate := range c {
		output = append(output, candidate.String())
	}

	return strings.Join(output, ", ")
}

// ExtractContent returns relevant content.
func ExtractContent(page io.Reader) (string, error) {
	document, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		return "", err
	}

	root := document.Get(0)
	for _, trash := range htmlutil.Query(root, "script,style") {
		if trash.Parent != nil {
			trash.Parent.RemoveChild(trash)
		}
	}

	transformMisusedDivsIntoParagraphs(root)
	removeUnlikelyCandidates(root)

	scores := getCandidates(root)
	//log.Printf("[Readability] Candidates: %v", candidates)

	best := getTopCandidate(scores)
	if best == nil {
		best = root
	}
	//log.Printf("[Readability] TopCandidate: %v", topCandidate)

	output := getArticle(root, best, scores)
	return output, nil
}

// Now that we have the top candidate, look through its siblings for content that might also be related.
// Things like preambles, content split by ads that we removed, etc.
func getArticle(root, best *html.Node, scores scorelist) string {
	selection := goquery.NewDocumentFromNode(root).FindNodes(best)
	topCandidate := &candidate{selection: selection, score: scores[best]}
	output := bytes.NewBufferString("<div>")
	siblingScoreThreshold := float32(math.Max(10, float64(topCandidate.score*.2)))

	topCandidate.selection.Siblings().Union(topCandidate.selection).Each(func(i int, s *goquery.Selection) {
		append := false
		node := s.Get(0)

		if node == topCandidate.Node() {
			append = true
		} else if score, ok := scores[node]; ok && score >= siblingScoreThreshold {
			append = true
		}

		if s.Is("p") {
			linkDensity := getLinkDensity(s.Get(0))
			content := s.Text()
			contentLength := len(content)

			if contentLength >= 80 && linkDensity < .25 {
				append = true
			} else if contentLength < 80 && linkDensity == 0 && sentenceRegexp.MatchString(content) {
				append = true
			}
		}

		if append {
			tag := "div"
			if s.Is("p") {
				tag = node.Data
			}

			html, _ := s.Html()
			fmt.Fprintf(output, "<%s>%s</%s>", tag, html, tag)
		}
	})

	output.Write([]byte("</div>"))
	return output.String()
}

func removeUnlikelyCandidates(root *html.Node) {
	body := htmlutil.Query(root, "body")
	if len(body) == 0 {
		return
	}
	for _, node := range htmlutil.Query(body[0], "*") {
		str := htmlutil.Attr(node, "class") + htmlutil.Attr(node, "id")

		blacklisted := (
			blacklistCandidatesRegexp.MatchString(str) ||
			(unlikelyCandidatesRegexp.MatchString(str) &&
			 !okMaybeItsACandidateRegexp.MatchString(str)))
		if blacklisted && node.Parent != nil {
			node.Parent.RemoveChild(node)
		}
	}
}

func getTopCandidate(scores scorelist) *html.Node {
	var best *html.Node
	var maxScore float32

	for node, score := range scores {
		if score > maxScore {
			best = node
			maxScore = score
		}
	}

	return best
}

// Loop through all paragraphs, and assign a score to them based on how content-y they look.
// Then add their score to their parent node.
// A score is determined by things like number of commas, class names, etc.
// Maybe eventually link density.
func getCandidates(root *html.Node) scorelist {
	scores := make(scorelist)
	for _, node := range htmlutil.Query(root, defaultTagsToScore) {
		text := htmlutil.Text(node)

		// If this paragraph is less than 25 characters, don't even count it.
		if len(text) < 25 {
			continue
		}

		parentNode := node.Parent
		grandParentNode := parentNode.Parent

		if _, found := scores[parentNode]; !found {
			scores[parentNode] = scoreNode(parentNode)
		}

		if grandParentNode != nil {
			if _, found := scores[grandParentNode]; !found {
				scores[grandParentNode] = scoreNode(grandParentNode)
			}
		}

		// Add a point for the paragraph itself as a base.
		contentScore := float32(1.0)

		// Add points for any commas within this paragraph.
		contentScore += float32(strings.Count(text, ",") + 1)

		// For every 100 characters in this paragraph, add another point. Up to 3 points.
		contentScore += float32(math.Min(float64(int(len(text)/100.0)), 3))

		scores[parentNode] += contentScore
		if grandParentNode != nil {
			scores[grandParentNode] += contentScore / 2.0
		}
	}

	// Scale the final candidates score based on link density. Good content
	// should have a relatively small link density (5% or less) and be mostly
	// unaffected by this operation
	for node, _ := range scores {
		scores[node] *= (1 - getLinkDensity(node))
	}

	return scores
}

func scoreNode(node *html.Node) float32 {
	var score float32

	switch node.Data {
	case "div":
		score += 5
	case "pre", "td", "blockquote", "img":
		score += 3
	case "address", "ol", "ul", "dl", "dd", "dt", "li", "form":
		score -= 3
	case "h1", "h2", "h3", "h4", "h5", "h6", "th":
		score -= 5
	}

	return score + getClassWeight(node)
}

// Get the density of links as a percentage of the content
// This is the amount of text that is inside a link divided by the total text in the node.
func getLinkDensity(n *html.Node) float32 {
	textLength := len(htmlutil.Text(n))
	if textLength == 0 {
		return 0
	}

	linkLength := 0.0
	for _, a := range htmlutil.Query(n, "a") {
		linkLength += float64(len(htmlutil.Text(a)))
	}

	return float32(linkLength) / float32(textLength)
}

// Get an elements class/id weight. Uses regular expressions to tell if this
// element looks good or bad.
func getClassWeight(node *html.Node) float32 {
	weight := 0
	class := htmlutil.Attr(node, "class")
	id := htmlutil.Attr(node, "id")

	if class != "" {
		if negativeRegexp.MatchString(class) {
			weight -= 25
		}

		if positiveRegexp.MatchString(class) {
			weight += 25
		}
	}

	if id != "" {
		if negativeRegexp.MatchString(id) {
			weight -= 25
		}

		if positiveRegexp.MatchString(id) {
			weight += 25
		}
	}

	return float32(weight)
}

func transformMisusedDivsIntoParagraphs(root *html.Node) {
	for _, node := range htmlutil.Query(root, "div") {
		if !divToPElementsRegexp.MatchString(htmlutil.InnerHTML(node)) {
			node.Data = "p"
		}
	}
}
