package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func getRSSFeed(data []byte) (*RSSFeed, error) {
	// decode the json data from the response
	// into a new locationAreas struct containing a list of areas
	var rss_feed RSSFeed
	err := xml.Unmarshal(data, &rss_feed)
	if err != nil {
		return nil, err
	}

	return &rss_feed, nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := NewClient(10 * time.Second)
	xmlBlob, err := client.GetResponseById(ctx, &feedURL, "")
	if err != nil {
		return nil, err
	}
	if len(xmlBlob) == 0 {
		return nil, fmt.Errorf("No XML data returned from feed %s", feedURL)
	}
	rss_feed, err := getRSSFeed(xmlBlob)
	if err != nil {
		return nil, err
	}
	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)
	for i, v := range rss_feed.Channel.Item {
		v.Title = html.UnescapeString(v.Title)
		v.Description = html.UnescapeString(v.Description)
		rss_feed.Channel.Item[i] = v
	}
	return rss_feed, nil
}

func handlerAgg(_ *state, cmd command) error {
	override_url := "https://www.wagslane.dev/index.xml"
	// if len(cmd.arguments) < 1 {
	// 	return fmt.Errorf("The %s command expects a URL", cmd.Name)
	// }
	// sel_url := cmd.arguments[0]
	sel_url := override_url
	rss_feed, err := fetchFeed(context.Background(), sel_url)
	if err != nil {
		return err
	}
	fmt.Println(rss_feed)
	return nil
}
