package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/pablomotte/gator_go/internal/database"
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

func handlerAgg(s *state, cmd command) error {
	// override_url := "https://www.wagslane.dev/index.xml"
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The %s command expects a duration between refreshes", cmd.name)
	}
	timeBetweenRequests := string(cmd.arguments[0])
	time_duration, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s (CTRL-C to end)\n", timeBetweenRequests)
	ticker := time.NewTicker(time_duration)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s, cmd); err != nil {
			return err
		}
	}
	// return nil
}

func handlerPostBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	var offset int32
	limit = 2
	offset = 0
	if len(cmd.arguments) < 1 {
		if len(cmd.arguments) < -1 { // optional parameter(s), so inhibit error message
			return fmt.Errorf("The %s command expects a limit of posts to view", cmd.name)
		}
	} else if len(cmd.arguments) >= 1 {
		limit_c, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			limit = 2
		} else {
			limit = int32(limit_c)
		}
		if len(cmd.arguments) > 1 {
			offset_c, err := strconv.Atoi(cmd.arguments[1])
			if err != nil {
				offset = 2
			} else {
				offset = int32(offset_c)
			}
		}
	}
	get_post_params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: offset,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), get_post_params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for k, v := range posts {
		fmt.Printf("- %3d\t%s\t\n", k+1, v.Title)
		fmt.Printf("     \t%s\t(%v)\t\n", v.Url, v.PublishedAt.Time.Format(time.UnixDate))
	}
	w.Flush()
	return nil
}
