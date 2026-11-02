package main

import (
	"context"
	"fmt"
	"time"
	"github.com/pablomotte/gator_go/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"database/sql"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The login command expects a username")
	}
	sel_user := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), sel_user)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %q does not exist", sel_user)
		}
		return err
	}
	if err := s.config.SetUser(sel_user); err != nil {
		return err
	}
	fmt.Printf("Current user set to %q\n", s.config.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The register command expects a username")
	}
	crt_user := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.arguments[0],
	}
	user, err := s.db.CreateUser(context.Background(), crt_user)
	if err != nil {
		pqErr := pq.As(err, pqerror.UniqueViolation)
		if pqErr != nil {
			return fmt.Errorf("user %q already exsts", crt_user.Name)
		}
		return err
	}
	// fmt.Printf("New user: %v\n", user)
	fmt.Printf("User %q created\n", user.Name)
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}