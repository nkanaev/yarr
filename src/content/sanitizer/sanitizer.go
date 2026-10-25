package sanitizer

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/nkanaev/yarr/src/content/htmlutil"
	"golang.org/x/net/html"
)

var splitSrcsetRegex = regexp.MustCompile(`,\s+`)

// Sanitize returns safe HTML.
func Sanitize(baseURL, input string) string {
	doc, err := html.Parse(strings.NewReader("<body>" + input + "</body>"))
	if err != nil {
		return ""
	}

	var body *html.Node
	for root := doc.FirstChild; root != nil; root = root.NextSibling {
		if root.Type == html.ElementNode && root.Data == "html" {
			for node := root.FirstChild; node != nil; node = node.NextSibling {
				if node.Type == html.ElementNode && node.Data == "body" {
					body = node
				}
			}
		}
	}
	if body == nil {
		return ""
	}

	sanitizeChildren(body, baseURL)

	var buffer strings.Builder
	for node := body.FirstChild; node != nil; node = node.NextSibling {
		if err := html.Render(&buffer, node); err != nil {
			return ""
		}
	}
	return buffer.String()
}

func sanitizeChildren(parent *html.Node, baseURL string) {
	for node := parent.FirstChild; node != nil; {
		next := node.NextSibling
		switch node.Type {
		case html.ElementNode:
			sanitizeElement(parent, node, baseURL)
		case html.CommentNode, html.DoctypeNode:
			parent.RemoveChild(node)
		}
		node = next
	}
}

func sanitizeElement(parent, node *html.Node, baseURL string) {
	tag := node.Data

	if isBlockedTag(tag) {
		parent.RemoveChild(node)
		return
	}

	// Children are sanitized first so that descendants of elements that end
	// up being removed or promoted have already been cleaned up.
	sanitizeChildren(node, baseURL)

	if !isValidTag(tag) {
		promoteChildren(parent, node)
		return
	}

	attrNames := sanitizeAttributes(baseURL, node)
	if !hasRequiredAttributes(tag, attrNames) {
		if tag == "iframe" {
			// A blocked iframe should not have its inner content rendered.
			parent.RemoveChild(node)
			return
		}
		promoteChildren(parent, node)
		return
	}

	if isVideoIframe(node) {
		wrap := &html.Node{
			Type: html.ElementNode,
			Data: "div",
			Attr: []html.Attribute{{Key: "class", Val: "video-wrapper"}},
		}
		next := node.NextSibling
		parent.RemoveChild(node)
		parent.InsertBefore(wrap, next)
		wrap.AppendChild(node)
	}

	if node.Data == "iframe" {
		// An iframe element never has fallback content.
		removeChildren(node)
	}
}

func removeChildren(node *html.Node) {
	for node.FirstChild != nil {
		node.RemoveChild(node.FirstChild)
	}
}

// promoteChildren replaces node with its (already sanitized) children.
func promoteChildren(parent, node *html.Node) {
	for node.FirstChild != nil {
		child := node.FirstChild
		node.RemoveChild(child)
		parent.InsertBefore(child, node)
	}
	parent.RemoveChild(node)
}

func sanitizeAttributes(baseURL string, node *html.Node) []string {
	var attrNames []string
	var attrs []html.Attribute
	tagName := node.Data

	for _, attribute := range node.Attr {
		// attribute names are case-insensitive; the parser may preserve
		// the original casing, while the allowlist is lowercase.
		key := strings.ToLower(attribute.Key)
		value := attribute.Val

		if !isValidAttribute(tagName, key) {
			continue
		}

		if (tagName == "img" || tagName == "source") && key == "srcset" {
			value = sanitizeSrcsetAttr(baseURL, value)
		}

		if isExternalResourceAttribute(key) {
			if tagName == "iframe" {
				if isValidIframeSource(baseURL, attribute.Val) {
					value = attribute.Val
				} else {
					continue
				}
			} else if tagName == "img" && key == "src" && isValidDataAttribute(attribute.Val) {
				value = attribute.Val
			} else {
				value = htmlutil.AbsoluteUrl(value, baseURL)
				if value == "" {
					continue
				}
				if !hasValidURIScheme(value) || isBlockedResource(value) {
					continue
				}
			}
		}

		attrNames = append(attrNames, key)
		attrs = append(attrs, html.Attribute{Key: key, Val: value})
	}

	extraNames, extraAttrs := extraAttributes(tagName)
	attrNames = append(attrNames, extraNames...)
	attrs = append(attrs, extraAttrs...)

	node.Attr = attrs
	return attrNames
}

func extraAttributes(tagName string) ([]string, []html.Attribute) {
	switch tagName {
	case "a":
		return []string{"rel", "target", "referrerpolicy"}, []html.Attribute{
			{Key: "rel", Val: "noopener noreferrer"},
			{Key: "target", Val: "_blank"},
			{Key: "referrerpolicy", Val: "no-referrer"},
		}
	case "video", "audio":
		return []string{"controls"}, []html.Attribute{{Key: "controls"}}
	case "iframe":
		return []string{"sandbox", "loading"}, []html.Attribute{
			{Key: "sandbox", Val: "allow-scripts allow-same-origin allow-popups"},
			{Key: "loading", Val: "lazy"},
		}
	case "img":
		return []string{"loading", "referrerpolicy"}, []html.Attribute{
			{Key: "loading", Val: "lazy"},
			{Key: "referrerpolicy", Val: "no-referrer"},
		}
	default:
		return nil, nil
	}
}

func isVideoIframe(node *html.Node) bool {
	if node.Data != "iframe" {
		return false
	}
	videoWhitelist := map[string]bool{
		"player.bilibili.com":      true,
		"player.vimeo.com":         true,
		"www.dailymotion.com":      true,
		"www.youtube-nocookie.com": true,
		"www.youtube.com":          true,
	}
	for _, attr := range node.Attr {
		if attr.Key == "src" {
			domain := htmlutil.URLDomain(attr.Val)
			return videoWhitelist[domain]
		}
	}
	return false
}

func isValidTag(tagName string) bool {
	return allowedTags.has(tagName) || allowedSvgTags.has(tagName) || allowedSvgFilters.has(tagName)
}

func isValidAttribute(tagName, attributeName string) bool {
	if attrs, ok := allowedAttrs[tagName]; ok {
		return attrs.has(attributeName)
	}
	if allowedSvgTags.has(tagName) || allowedSvgFilters.has(tagName) {
		return allowedSvgAttrs.has(attributeName)
	}
	return false
}

func isExternalResourceAttribute(attribute string) bool {
	switch attribute {
	case "src", "href", "poster", "cite":
		return true
	default:
		return false
	}
}

func hasRequiredAttributes(tagName string, attributes []string) bool {
	elements := make(map[string][]string)
	elements["a"] = []string{"href"}
	elements["iframe"] = []string{"src"}
	elements["img"] = []string{"src"}
	elements["source"] = []string{"src", "srcset"}

	for element, attrs := range elements {
		if tagName == element {
			for _, attribute := range attributes {
				if slices.Contains(attrs, attribute) {
					return true
				}
			}
			return false
		}
	}
	return true
}

func hasValidURIScheme(src string) bool {
	scheme, _, _ := strings.Cut(src, ":")
	return allowedURISchemes.has(scheme)
}

func isBlockedResource(src string) bool {
	blacklist := []string{
		"feedsportal.com",
		"api.flattr.com",
		"stats.wordpress.com",
		"plus.google.com/share",
		"twitter.com/share",
		"feeds.feedburner.com",
	}
	return slices.ContainsFunc(blacklist, func(element string) bool {
		return strings.Contains(src, element)
	})
}

func isValidIframeSource(baseURL, src string) bool {
	whitelist := []string{
		"bandcamp.com",
		"cdn.embedly.com",
		"invidio.us",
		"player.bilibili.com",
		"player.vimeo.com",
		"soundcloud.com",
		"vk.com",
		"w.soundcloud.com",
		"www.dailymotion.com",
		"www.youtube-nocookie.com",
		"www.youtube.com",
	}

	domain := htmlutil.URLDomain(src)
	// allow iframe from same origin
	if htmlutil.URLDomain(baseURL) == domain {
		return true
	}
	return slices.Contains(whitelist, domain)
}

func isBlockedTag(tagName string) bool {
	switch tagName {
	case "noscript", "script", "style":
		return true
	}
	return false
}

func sanitizeSrcsetAttr(baseURL, value string) string {
	var sanitizedSources []string
	rawSources := splitSrcsetRegex.Split(value, -1)
	for _, rawSource := range rawSources {
		parts := strings.Split(strings.TrimSpace(rawSource), " ")
		nbParts := len(parts)

		if nbParts > 0 {
			sanitizedSource := parts[0]
			if !strings.HasPrefix(parts[0], "data:") {
				sanitizedSource = htmlutil.AbsoluteUrl(parts[0], baseURL)
				if sanitizedSource == "" {
					continue
				}
			}

			if nbParts == 2 && isValidWidthOrDensityDescriptor(parts[1]) {
				sanitizedSource += " " + parts[1]
			}

			sanitizedSources = append(sanitizedSources, sanitizedSource)
		}
	}
	return strings.Join(sanitizedSources, ", ")
}

func isValidWidthOrDensityDescriptor(value string) bool {
	if value == "" {
		return false
	}

	lastChar := value[len(value)-1:]
	if lastChar != "w" && lastChar != "x" {
		return false
	}

	_, err := strconv.ParseFloat(value[0:len(value)-1], 32)
	return err == nil
}

func isValidDataAttribute(value string) bool {
	var dataAttributeAllowList = []string{
		"data:image/avif",
		"data:image/apng",
		"data:image/png",
		"data:image/svg",
		"data:image/svg+xml",
		"data:image/jpg",
		"data:image/jpeg",
		"data:image/gif",
		"data:image/webp",
	}
	return slices.ContainsFunc(dataAttributeAllowList, func(prefix string) bool {
		return strings.HasPrefix(value, prefix)
	})
}
