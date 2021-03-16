package router

import (
	"reflect"
	"testing"
)

func TestRouteRegexpPart(t *testing.T) {
	in := "/hello/:world"
	re := routeRegexp(in)

	pos := []string{
		"/hello/world",
		"/hello/1234",
		"/hello/bbc1",
	}
	for _, c := range pos {
		if !re.MatchString(c) {
			t.Errorf("%v must match %v", in, c)
		}
	}

	neg := []string{
		"/hello",
		"/hello/world/",
		"/sub/hello/123",
		"//hello/123",
		"/hello/123/hello/",
	}
	for _, c := range neg {
		if re.MatchString(c) {
			t.Errorf("%q must not match %q", in, c)
		}
	}
}

func TestRouteRegexpStar(t *testing.T) {
	in := "/hello/*world"
	re := routeRegexp(in)

	pos := []string{"/hello/world", "/hello/world/test"}
	for _, c := range pos {
		if !re.MatchString(c) {
			t.Errorf("%q must match %q", in, c)
		}
	}

	neg := []string{"/hello/", "/hello"}
	for _, c := range neg {
		if re.MatchString(c) {
			t.Errorf("%v must not match %v", in, c)
		}
	}
}

func TestRegexGroupsPart(t *testing.T) {
	re := routeRegexp("/foo/:bar/1/:baz")
	
	expect := map[string]string{"bar": "one", "baz": "two"}
	actual := regexGroups("/foo/one/1/two", re)

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expected: %q, actual: %q", expect, actual)
	}
}

func TestRegexGroupsStar(t *testing.T) {
	re := routeRegexp("/foo/*bar")
	
	expect := map[string]string{"bar": "bar/baz/"}
	actual := regexGroups("/foo/bar/baz/", re)

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expected: %q, actual: %q", expect, actual)
	}
}
