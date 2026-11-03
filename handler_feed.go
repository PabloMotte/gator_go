package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/pablomotte/gator_go/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("The addfeed command expects a name and url")
	}
	sel_user := s.config.CurrentUserName
	user, err := s.db.GetUser(context.Background(), sel_user)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %q does not exist", sel_user)
		}
		return err
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
			return fmt.Errorf("feed %q already exsts", crt_feed.Url)
		}
		return err
	}
	// fmt.Printf("New feed: %v\n", feed)
	fmt.Printf("Feed %q (%s) for %s created\n", feed.Name, feed.Url, user.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("The feeds command does not expect any parameters")
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
		fmt.Printf("%d: %s (%s) %q\n", k+1, v.Name, v.Url, v.Name_2)
	}
	return nil
}
