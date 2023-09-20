package main

import (
	"testing"
)

func testUnique(expected int, urls []string, t *testing.T) {

	unique := CountUniqueUrls(urls)

	if unique != expected {
		t.Errorf("Expected: %d\nGot: %d\n", expected, unique)
	}
}

func testUniqueTLD(expected map[string]int, url []string, t *testing.T) {
	unique := CountUniqueUrlsPerTopLevelDomain(url)

	if expected == nil && unique != nil {
		t.Error("Expected: nil\nGot: Not Nil\n")
	}

	if len(expected) != len(unique) {
		t.Error("Expected: Equal Length Result\nGot: Unequal Length Result\n")
	}

	for tld, count := range unique {
		if c, exists := expected[tld]; !exists || count != c {
			t.Errorf("Expected: %d\nGot: %d\n", c, count)
		}
	}
}

func Test_CountUnique(t *testing.T) {
	// test sorting query params
	urls := []string{"https://example.com?a=1&b=2", "https://example.com?b=2&a=1", "https://example.com?A=1&B=2"}

	testUnique(1, urls, t)

	// test empty list
	testUnique(0, []string{}, t)

	// test same domain (http vs. https)
	urls = []string{"https://www.google.com", "https://google.com", "https://www.google.com?", "hTTps://google.com", "http://google.com", "http://google.com/"}
	testUnique(2, urls, t)

	// test multiple unique urls
	urls = []string{"https://google.com", "https://google.org", "https://google.net"}
	testUnique(3, urls, t)

	// test same (duplicate) domain
	urls = []string{"https://google.com", "https://google.com", "https://google.com"}
	testUnique(1, urls, t)

	// test invalid urls
	urls = []string{"https://example.com", "invalid-url", "ftp://ftp.example.com"}
	testUnique(2, urls, t)

	// test urls with fragments
	urls = []string{"https://example.com#section1", "https://example.com#section2"}
	testUnique(1, urls, t)

	// test urls with non default port numbers
	urls = []string{"https://example.com:8080", "https://example.com"}
	testUnique(2, urls, t)

	// test urls with default port numbers
	urls = []string{"http://example.com:80", "http://example.com", "https://example.com:443", "https://example.com"}
	testUnique(2, urls, t)

	// test with dashes
	urls = []string{"https://xn--bcher-kva.ch", "https://example.com"}
	testUnique(2, urls, t)

	// test with query and fragment
	urls = []string{"https://example.com?a=1#section1", "https://example.com#section1?a=1"}
	testUnique(1, urls, t)

	// test url with incorrect default port number
	urls = []string{"http://example.com:443", "http://example.com", "https://example.com:80", "https://example.com"}
	testUnique(4, urls, t)

	// test same domain with path
	urls = []string{"https://www.google.com", "https://google.com", "https://google.com#fragment", "https://www.google.com?", "https://www.google.com/1", "https://google.com/1"}
	testUnique(2, urls, t)

	// test same domain with percent encoding
	urls = []string{"https://foo.bar/baz*", "https://foo.bar/baz%2A", "https://foo.bar/baz%2a", "https://example.com/encoded%20path"}
	testUnique(2, urls, t)

	// test removing dot-segments
	urls = []string{"https://example.com/foo/./bar/baz/../qux", "https://example.com/foo/bar/qux", "https://example.com/foo/./bar/baz/../../qux", "https://example.com/foo/qux"}
	testUnique(2, urls, t)

	// test special reserved characters
	urls = []string{"https://example.com?a=%20&b=%21", "https://example.com?a=%2520&b=%21"}
	testUnique(1, urls, t)

	// test different encodings
	urls = []string{"https://example.com/encoded%20path", "https://example.com/encoded%2520path"}
	testUnique(1, urls, t)

	// test trailing path
	urls = []string{"https://example.com/foo/bar", "https://example.com/foo/bar/"}
	testUnique(1, urls, t)

	urls = []string{"https://example.com?a=1&b=%24", "https://example.com?a=1&b=$"}
	// Expecting 1, as the query parameter should be case-sensitive
	testUnique(1, urls, t)

}

func Test_CountUniqueTLD(t *testing.T) {

	urls := []string{"https://example.com?a=1&b=2", "https://example.com?b=2&a=1", "https://www.example.com", "https://sub.example.com", "https://foo.com?a=1&b=2&c=3", "https://foo.com?b=2&a=1&c=3", "https://foo.com"}
	expected := map[string]int{"example.com": 3, "foo.com": 2}
	testUniqueTLD(expected, urls, t)

	// test different tlds
	urls = []string{"https://example.com", "https://example.org"}
	expected = map[string]int{"example.com": 1, "example.org": 1}
	testUniqueTLD(expected, urls, t)

	// test mixed case
	urls = []string{"https://Example.com", "HTTPS://example.COM"}
	expected = map[string]int{"example.com": 1}
	testUniqueTLD(expected, urls, t)

	// test empty list
	urls = []string{}
	expected = map[string]int{}
	testUniqueTLD(expected, urls, t)

	// test urls with different ports
	urls = []string{"https://example.com:8080", "https://example.com:8888"}
	expected = map[string]int{"example.com": 2}
	testUniqueTLD(expected, urls, t)

	// test with special characters
	urls = []string{"https://example.com?a=%20&b=%21", "https://example.com?a=%2520&b=%21"}
	expected = map[string]int{"example.com": 1}
	testUniqueTLD(expected, urls, t)

	// test with encoded path
	urls = []string{"https://example.com/encoded%20path", "https://example.com/encoded%2520path"}
	expected = map[string]int{"example.com": 1}

	// test different tld
	urls = []string{"https://sub1.example.com", "https://sub2.example.org"}
	expected = map[string]int{"example.com": 1, "example.org": 1}
}
