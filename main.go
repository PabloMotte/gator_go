package main

import (
	"fmt"
	"os"
	"github.com/pablomotte/gator_go/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	new_config, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config", err)
		os.Exit(1)
	}
	curr_state := state{
		config: &new_config,
	}
	curr_commands := createCommands(curr_state)
	num_args := len(os.Args)
	if num_args < 2 {
		fmt.Println("Required: Command name")
		os.Exit(1)
	}
	var cmd command
	cmd_name := os.Args[1]
	if num_args < 3 {
		cmd = command{
			name: cmd_name,
			arguments: []string{},
		}
	} else {
		cmd = command{
			name: cmd_name,
			arguments: os.Args[2:],
		}
	}
	err = curr_commands.run(&curr_state, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func print_config(cfg config.Config) {
	for k,v := range []interface{}{cfg.DbURL, cfg.CurrentUserName} {
		fmt.Println(k, v)
	}
}

