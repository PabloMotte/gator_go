package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The login handler expects a username")
	}
	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("Current user set to %q\n", s.config.CurrentUserName)
	return nil
}