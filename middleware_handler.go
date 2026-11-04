package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pablomotte/gator_go/internal/database"
)

// Ensure the user is logged in before calling the handler
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		sel_user := s.config.CurrentUserName
		user, err := s.db.GetUser(context.Background(), sel_user)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user %q does not exist", sel_user)
			}
			return err
		}
		return handler(s, cmd, user)
	}
}
