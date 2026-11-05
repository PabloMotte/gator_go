package main

import (
	"fmt"
)

type command struct {
	name      string
	arguments []string
}
type gator_function func(*state, command) error
type commands struct {
	cmd_map map[string]gator_function
}
type command_handlers struct {
	cmd     command
	handler gator_function
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.cmd_map[cmd.name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}
	if err := fn(s, cmd); err != nil {
		return fmt.Errorf("Error in %q: %w", cmd.name, err)
	}
	return nil
}

func (c *commands) register(name string, f gator_function) {
	// handlerLogin(s *state, cmd command) error
	_, ok := c.cmd_map[name]
	if ok {
		fmt.Printf("Command %q already exists", name)
		return
	}
	c.cmd_map[name] = f
}

func (c *commands) listCommands() {
	fmt.Println("Registered commands:")
	for name, _ := range c.cmd_map {
		fmt.Printf(" - %q\n", name)
	}
}

func createCommands(curr_state state) commands {
	cmds := commands{
		cmd_map: make(map[string]gator_function),
	}
	for name, cmd := range defineCommands() {
		cmds.register(name, cmd.handler)
	}
	return cmds
}

func defineCommands() map[string]command_handlers {
	return map[string]command_handlers{
		"login": {
			cmd: command{
				name:      "login",
				arguments: []string{"username"},
			},
			handler: handlerLogin,
		},
		"register": {
			cmd: command{
				name:      "register",
				arguments: []string{"username"},
			},
			handler: handlerRegister,
		},
		"reset": {
			cmd: command{
				name:      "reset",
				arguments: []string{},
			},
			handler: handlerReset,
		},
		"users": {
			cmd: command{
				name:      "users",
				arguments: []string{},
			},
			handler: handlerUsers,
		},
		"agg": {
			cmd: command{
				name:      "agg",
				arguments: []string{"timeBetweenRequests"},
			},
			handler: handlerAgg,
		},
		"addfeed": {
			cmd: command{
				name:      "addfeed",
				arguments: []string{"name", "feedURL"},
			},
			handler: middlewareLoggedIn(handlerAddFeed),
		},
		"feeds": {
			cmd: command{
				name:      "feeds",
				arguments: []string{},
			},
			handler: handlerFeeds,
		},
		"follow": {
			cmd: command{
				name:      "follow",
				arguments: []string{"feedURL"},
			},
			handler: middlewareLoggedIn(handlerAddFollow),
		},
		"unfollow": {
			cmd: command{
				name:      "unfollow",
				arguments: []string{"feedURL"},
			},
			handler: middlewareLoggedIn(handlerDeleteFollow),
		},
		"following": {
			cmd: command{
				name:      "following",
				arguments: []string{"feedURL"},
			},
			handler: middlewareLoggedIn(handlerListFollowing),
		},
		"browse": {
			cmd: command{
				name:      "browse",
				arguments: []string{"limit", "offset"},
			},
			handler: middlewareLoggedIn(handlerPostBrowse),
		},
	}
}
