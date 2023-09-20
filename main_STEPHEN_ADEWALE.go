package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

/**
* This function counts how many unique normalized valid URLs were passed to the function
*
* Accepts a list of URLs
*
* Example:
*
* input: ['https://example.com']
* output: 1
*
* Notes:
*  - assume none of the URLs have authentication information (username, password).
*
* Normalized URL:
*  - process in which a URL is modified and standardized: https://en.wikipedia.org/wiki/URL_normalization
*
#    For example.
#    These 2 urls are the same:
#    input: ["https://example.com", "https://example.com/"]
#    output: 1
#
#    These 2 are not the same:
#    input: ["https://example.com/", "http://example.com"]
#    output 2
#
#    These 2 are the same:
#    input: ["https://example.com?", "https://example.com"]
#    output: 1
#
#    These 2 are the same:
#    input: ["https://example.com?a=1&b=2", "https://example.com?b=2&a=1"]
#    output: 1
*/

var urlRegex = regexp.MustCompile(`^(?:([a-zA-Z]+)://)?((?:[a-zA-Z0-9-]+\.)+[a-zA-Z0-9]{2,})(:[0-9]+)?(/[a-zA-Z0-9%/^*()$@!_.-]*)?(?:\?([a-zA-Z0-9^*()$@!_%=&.-]*))?(#[a-zA-Z0-9^*()$@!_%&]+)?$`)
var reservedCharacters = map[string]string{"20": " ", "21": "!", "22": "\"", "23": "#", "24": "$", "25": "%", "26": "&", "27": "'", "28": "(", "29": ")", "2a": "*", "2b": "+", "2c": ",", "2d": "-", "2e": ".", "2f": "/", "3a": ":", "3b": ";", "3c": "<", "3d": "=", "3e": ">", "3f": "?", "40": "@", "5b": "[", "5c": "\\", "5d": "]", "5e": "^", "5f": "_", "60": "`", "7b": "{", "7c": "|", "7d": "}", "7e": "~"}

func CountUniqueUrls(urls []string) int {
	unique := make(map[string]struct{})
	for _, url := range urls {
		if !urlRegex.MatchString(url) {
			continue
		}
		url = normalize(url)
		unique[url] = struct{}{}
	}
	return len(unique)
}

func normalize(url string) string {
	url = strings.ToLower(url)                // convert to lowercase
	url = strings.TrimRight(url, "?/")        // remove trailing question mark and forward slash
	parts := urlRegex.FindStringSubmatch(url) // split the url into parts
	var (
		domain                       string
		port                         string
		directories                  string
		queries                      string
		protocol                     string
		result                       strings.Builder
		reservedCharsCheckStartIndex int
	)

	protocol = parts[1]
	domain = parts[2]
	port = parts[3]

	domain = strings.TrimLeft(domain, "www.") // remove www. from the domain

	if port != "" {
		port = port[1:]
		for scheme, portNumber := range map[string]string{"http": "80", "https": "443"} {
			if port == portNumber && scheme == protocol {
				port = ""
				break
			}
		}
	}

	directories = parts[4]

	queries = parts[5]

	// add the protocol
	if protocol != "" {
		result.WriteString(fmt.Sprintf("%s://", protocol))
	}

	// add the domain
	result.WriteString(domain)

	// add the port
	if port != "" {
		result.WriteString(fmt.Sprintf(":%s", port))
	}

	reservedCharsCheckStartIndex = result.Len()

	// remove double or triple forward slashes
	if directories != "" {
		files := strings.Split(directories, "/")[1:] // remove the first empty string

		for _, file := range files {
			if file == "." {
				continue
			}

			if file == ".." {
				res := result.String()
				result.Reset()
				result.WriteString(res[:strings.LastIndex(res, "/")])
				continue
			}

			result.WriteString(fmt.Sprintf("/%s", file))
		}
	}

	if queries != "" {
		q := strings.Split(queries, "&")

		sort.Strings(q)

		queries = strings.Join(q, "&")

		result.WriteString(fmt.Sprintf("?%s", queries))
	}

	partResult := result.String()[reservedCharsCheckStartIndex:]

	for i := 0; i < len(partResult); i++ {
		if string(partResult[i]) == "%" && i+2 < len(partResult) {
			if val, ok := reservedCharacters[partResult[i+1:i+3]]; ok {
				partResult = strings.Replace(partResult, partResult[i:i+3], val, 1)
				i--
			} else {
				partResult = strings.Replace(partResult, partResult[i+1:i+3], strings.ToUpper(partResult[i+1:i+3]), 1)
			}
		}
	}

	res := result.String()
	result.Reset()
	result.WriteString(fmt.Sprintf("%s%s", res[:reservedCharsCheckStartIndex], partResult))

	return result.String()
}

/**
 * This function counts how many unique normalized valid URLs were passed to the function per top level domain
 *
 * A top level domain is a domain in the form of example.com. Assume all top level domains end in .com
 * subdomain.example.com is not a top level domain.
 *
 * Accepts a list of URLs
 *
 * Example:
 *
 * input: ["https://example.com"]
 * output: Hash["example.com" => 1]
 *
 * input: ["https://example.com", "https://subdomain.example.com"]
 * output: Hash["example.com" => 2]
 *
 */

func CountUniqueUrlsPerTopLevelDomain(urls []string) map[string]int {
	tlds := make(map[string]int)
	tldList := make(map[string][]string)
	for _, url := range urls {
		if url != "" {
			normalized := normalize(urlRegex.FindStringSubmatch(strings.ToLower(url))[2])
			tld := strings.Join(strings.Split(normalized, ".")[len(strings.Split(normalized, "."))-2:], ".")
			tldList[tld] = append(tldList[tld], url)
		}
	}

	for tld, normalized := range tldList {
		tlds[tld] = CountUniqueUrls(normalized)
	}

	return tlds
}
