package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"main/reddit"
	"main/storage"
	"net/url"
	"os"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

var auth *storage.StoredAuth

func doInit(c context.Context) error {
	a, err := storage.GetAuth(c)
	if err != nil {
		panic(err)
	}

	if a == nil {
		log.Debugf(c, "%v", "Logging in")

		// Log in
		resp, err := reddit.Authorize(
			c,
			os.Getenv("USERNAME"),
			os.Getenv("PASSWORD"),
		)

		if err != nil {
			return err
		}

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Load response
		log.Debugf(c, "%v", resp.Header)
		log.Debugf(c, "%v", resp)
		log.Debugf(c, "%v", string(b))

		type loginJSON struct {
			JSON struct {
				Data struct {
					Modhash string `json:"modhash"`
					Cookie  string `json:"cookie"`
				} `json:"data"`
			} `json:"json"`
		}

		var l loginJSON
		err = json.Unmarshal(b, &l)
		if err != nil {
			return err
		}

		// Save
		data := l.JSON.Data
		a = &storage.StoredAuth{
			data.Modhash,
			data.Cookie,
		}
		err = storage.SetAuth(c, a)
		if err != nil {
			return err
		}
	} else {
		log.Debugf(c, "%v", "Reusing old login")
	}

	auth = a
	return nil
}

func doPost(c context.Context, simulate bool) error {
	subreddit := os.Getenv("subreddit")

	// Get previous offset
	o, err := storage.GetOffset(c)
	if err != nil {
		return err
	}

	// Fetch new links
	resp, err := reddit.GetNew(c, o.Last, subreddit)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type respJSON struct {
		Data struct {
			Children []struct {
				Data struct {
					Name      string `json:"name"`
					Title     string `json:"title"`
					SelfText  string `json:"selftext"`
					Author    string `json:"author"`
					Subreddit string `json:"subreddit"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	var r respJSON
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	// Check if there are links
	links := r.Data.Children
	if len(links) == 0 {
		log.Debugf(c, "%v", "No links to parse!")
		return nil
	}

	log.Debugf(c, "Processing %v links", len(links))

	// Update the offset
	o.Last = links[0].Data.Name
	err = storage.SetOffset(c, o)
	if err != nil {
		return err
	}

	// Process new links
	var skipped []string

	for _, link := range links {
		l := link.Data
		// Ensure correct subreddit
		if l.Subreddit == subreddit {
			// Try to correct the user
			corrected := makeWording(l.Title, l.SelfText, l.Author)
			if corrected != "" {
				if simulate {
					log.Infof(c, "%v", corrected)
				} else {
					log.Warningf(c, "Deferred %v", corrected)
					taskqueue.Add(c, taskqueue.NewPOSTTask("/do/comment", url.Values{
						"parent": {l.Name},
						"text":   {corrected},
					}), "default")
				}
				continue
			}
		}
		skipped = append(skipped, l.Title)
	}

	if len(skipped) != 0 {
		log.Debugf(c, "[Skipped %v] %v", len(skipped), skipped)
	}
	return nil
}

func doComment(c context.Context, parent, text string) error {
	resp, err := reddit.PostComment(c, parent, text, auth.ModHash, auth.Session)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debugf(c, "%v", resp.StatusCode)
	log.Debugf(c, "%v", resp.Header)
	log.Debugf(c, "%v", string(b))

	return nil
}
