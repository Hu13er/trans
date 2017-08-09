package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	extractRegexp = regexp.MustCompile(`<ul.*?class="result js-result".*?>(?P<result>.*?)<\/ul>`)
	translateURL  = "https://fastdic.com/word/%s"
)

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(0)
	}

	stmt := os.Args[1]

	html := getHTML(stmt)
	if html == "" {
		os.Exit(0)
	}

	trans := extract(html)
	fmt.Println(trans)
}

func extract(html string) string {
	var result string
	result = regexp.MustCompile("\n").ReplaceAllString(html, "")
	for k, v := range extractRegexp.SubexpNames() {
		if v == "result" {
			if sub := extractRegexp.FindStringSubmatch(result); sub != nil && len(sub) > k {
				result = sub[k]
			} else {
				result = "No result"
				break
			}
		}
	}

	result = regexp.MustCompile("<br>").ReplaceAllString(result, "\n")
	result = regexp.MustCompile("<.*?>").ReplaceAllString(result, "")

	return result
}

func getHTML(stmt string) string {
	url := fmt.Sprintf(translateURL, stmt)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(html)
}

func help() {
	help := "Usage: trans <statement>"
	fmt.Println(strings.TrimSpace(help))
}
