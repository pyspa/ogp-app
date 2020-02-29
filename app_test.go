package main

import "testing"

func TestIsTLS(t *testing.T) {
	cases := []struct {
		url      string
		expected bool
	}{
		{url: "https://ogp.app", expected: true},
		{url: "http://ogp.app", expected: false},
		{url: "https://ogp.app/images", expected: true},
	}

	for _, c := range cases {
		if isTLS(c.url) != c.expected {
			t.Errorf("%s", c.url)
		}
	}
}
