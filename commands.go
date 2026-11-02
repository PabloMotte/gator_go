package main

import (
	"fmt"
)

type command struct {
	name string
	arguments []string
}
type gator_function func(*state, command) error
type commands struct {
	cmd_map map[string]gator_function
}
type command_handlers struct {
	cmd command
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
				name: "login",
				arguments: []string{"username",},
			},
			handler: handlerLogin,
		},
		"register": {
			cmd: command{
				name: "register",
				arguments: []string{"username",},
			},
			handler: handlerRegister,
		},
	}
}