package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const shorten_page = "https://api-ssl.bitly.com/v4/shorten"
const token_not_mine = "c28a95a72f7b149061cfac25417628733195dde2"

func shorten(url string) string {
	body := `{
  "long_url": "%s",
  "group_guid": "Bj4r63akxYF"
}`

	body = fmt.Sprintf(body, url)

	req, err := http.NewRequest("POST", shorten_page, bytes.NewBuffer([]byte(body)))

	req.Header.Set("Authorization", "Bearer "+token_not_mine)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	newbody, _ := ioutil.ReadAll(resp.Body)
	for _, line := range strings.Split(string(newbody), ",") {
		if strings.Contains(line, "link") {
			link := strings.SplitN(line, ":", 2)[1]
			return link[1 : len(link)-1]
		}
	}

	return string(newbody)
}
