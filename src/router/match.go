package router

import "regexp"

func regexGroups(input string, regex *regexp.Regexp) map[string]string {
	groups := make(map[string]string)
	matches := regex.FindStringSubmatchIndex(input)
	for i, key := range regex.SubexpNames()[1:] {
		groups[key] = input[matches[i*2+2]:matches[i*2+3]]
	}
	return groups
}

func routeRegexp(route string) *regexp.Regexp {
	chunks := regexp.MustCompile(`[\*\:]\w+`)
	output := chunks.ReplaceAllStringFunc(route, func(m string) string {
		if m[0:1] == `*` {
			return "(?P<" + m[1:] + ">.+)"
		}
		return "(?P<" + m[1:] + ">[^/]+)"
	})
	output = "^" + output + "$"
	return regexp.MustCompile(output)
}
