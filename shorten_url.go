package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const shorten_page = "https://cutt.ly/scripts/shortenUrl.php"

func shorten(url string) string {
	body := `-----------------------------1585130260283869584310048763
Content-Disposition: form-data; name="url"

%s
-----------------------------1585130260283869584310048763`
	body = fmt.Sprintf(body, url)

	req, err := http.NewRequest("POST", shorten_page, bytes.NewBuffer([]byte(body)))

	req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------1585130260283869584310048763")
	req.Header.Set("Referer", "https://cutt.ly")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	newbody, _ := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(newbody))
}
