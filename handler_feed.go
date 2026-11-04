package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/pablomotte/gator_go/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("The %s command expects a name and url", cmd.name)
	}
	crt_feed := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}
	feed, err := s.db.AddFeed(context.Background(), crt_feed)
	if err != nil {
		pqErr := pq.As(err, pqerror.UniqueViolation)
		if pqErr != nil {
			// return fmt.Errorf("feed %q already exsts", crt_feed.Url)
			cmd.arguments[0] = cmd.arguments[1]
			return handlerAddFollow(s, cmd, user)
		}
		return err
	}
	// fmt.Printf("New feed: %v\n", feed)
	fmt.Printf("Feed %q (%s) for %s created\n", feed.Name, feed.Url, user.Name)
	// Make sure to also add a "following" record
	cmd.arguments[0] = cmd.arguments[1]
	return handlerAddFollow(s, cmd, user)
}

func handlerListFollowing(s *state, cmd command, user database.User) error {
	// if len(cmd.arguments) < 0 {
	// 	return fmt.Errorf("The %s command expects a name", cmd.name)
	// }
	feeds, err := s.db.GetUserFollows(context.Background(), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %q does not follow any feeds", user.Name)
		}
		return err
	}
	// fmt.Printf("New feed: %v\n", feed)
	fmt.Printf("%s is following:\n", user.Name)
	for k, v := range feeds {
		fmt.Printf("%d: %s %q\n", k+1, v.FeedName, v.FeedUrl)
	}
	return nil
}

func handlerAddFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The %s command expects feed url", cmd.name)
	}
	sel_feed := cmd.arguments[0]
	feed, err := s.db.GetFeed(context.Background(), sel_feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed %q does not exist", sel_feed)
		}
		return err
	}
	crt_feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	follower, err := s.db.CreateFeedFollow(context.Background(), crt_feed_follow)
	if err != nil {
		pqErr := pq.As(err, pqerror.UniqueViolation)
		if pqErr != nil {
			return fmt.Errorf("user %s already following %q", user.Name, feed.Url)
		}
		return err
	}
	fmt.Printf("Feed %q now being followed by %s\n", feed.Url, follower.UserName)
	return nil
}

func handlerDeleteFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The %s command expects a feed url", cmd.name)
	}
	sel_feed := cmd.arguments[0]
	feed, err := s.db.GetFeed(context.Background(), sel_feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed %q does not exist", sel_feed)
		}
		return err
	}
	dlt_feed_follow := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if err := s.db.Unfollow(context.Background(), dlt_feed_follow); err != nil {
		pqErr := pq.As(err, pqerror.UniqueViolation)
		if pqErr != nil {
			return fmt.Errorf("user %s already following %q", user.Name, feed.Url)
		}
		return err
	}
	fmt.Printf("Feed %q not being followed by %s anymore\n", feed.Url, user.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("The %s command does not expect any parameters", cmd.name)
	}
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		pqErr := pq.As(err, pqerror.UniqueViolation)
		if pqErr != nil {
			return fmt.Errorf("No feeds exist")
		}
		return err
	}
	for k, v := range feeds {
		fmt.Printf("%d: %s (%s) %q\n", k+1, v.Name, v.Url, v.UserName)
	}
	return nil
}

func scrapeFeeds(s *state, cmd command) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	no_rows := false
	feed_data, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			no_rows = true
		} else {
			return err
		}
	}
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}
	fmt.Printf("Feed %q (%s) fetched", feed.Name, feed.Url)
	if !no_rows {
		no_rows = true // flip flag in case SQL didn't return "no rows" error
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		for k, v := range feed_data.Channel.Item {
			if no_rows {
				fmt.Printf(":\n")
				no_rows = false
			}
			fmt.Printf("- %d\t%s\t\n", k+1, v.Title)
		}
		w.Flush()
	}
	if no_rows {
		fmt.Printf(", but nothing to display\n")
	}
	return nil
}
