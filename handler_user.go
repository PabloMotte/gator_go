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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The %s command expects a username", cmd.name)
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

func handlerUsers(s *state, cmd command) error {
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("The %s command expects no parameters", cmd.name)
	}
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no users exist")
		}
		return err
	}
	if len(users) == 0 {
		fmt.Printf("No users registered\n")
	} else {
		for _, v := range users {
			tag := ""
			if s.config.CurrentUserName == v.Name {
				tag = " (current)"
			}
			fmt.Printf("* %s%s\n", v.Name, tag)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("The %s command expects no parameters", cmd.name)
	}
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no records found")
		}
		return err
	}
	if err := s.config.SetUser(""); err != nil {
		return err
	}
	fmt.Printf("Database reset\n")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The %s command expects a username", cmd.name)
	}
	crt_user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
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
