package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) { // instead of having pointers `*RSSFeed` so that we can return nil, we can just have empty structs
	httpClient := http.Client{
		Timeout: 180 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err // can't use nil as RSSFeed value in return statement ... nil is a predeclared identifier representing the zero value for a pointer, channel, func, interface, map, or slice type
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body) // these slice of bytes, we want to read into these RSSFeeds
	if err != nil {
		return RSSFeed{}, err
	}
	rssFeed := RSSFeed{}

	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil // returning the newly populated rssFeed
}
