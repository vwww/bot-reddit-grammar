package reddit

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func userAgent() string { return os.Getenv("USERAGENT") }

func GetNew(last, subreddit string) (*http.Response, error) {
	url := fmt.Sprintf("https://www.reddit.com/r/%v/new.json?limit=100", subreddit)
	if last != "" {
		url += "&before=" + last
	}
	url += fmt.Sprintf("&%v", rand.Float64())

	log.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent())

	return http.DefaultClient.Do(req)
}

func Authorize(username, password string) (*http.Response, error) {
	v := url.Values{}
	v.Set("api_type", "json")
	v.Set("user", username)
	v.Set("passwd", password)
	payload := v.Encode()

	log.Println(payload)

	req, err := http.NewRequest("POST", "https://ssl.reddit.com/api/login", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent())

	return http.DefaultClient.Do(req)
}

func PostComment(parent, text string, modhash, session string) (*http.Response, error) {
	v := url.Values{}
	v.Set("api_type", "json")
	v.Set("text", text)
	v.Set("thing_id", parent)
	v.Set("uh", modhash)
	payload := v.Encode()

	log.Println(payload)

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/comment", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Cookie", fmt.Sprintf("reddit_session=%v; Domain=reddit.com; Path=/; HttpOnly", url.QueryEscape(session)))

	return http.DefaultClient.Do(req)
}
