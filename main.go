package main

import (
	"fmt"
	"os"
	"database/sql"
	"github.com/pablomotte/gator_go/internal/config"
	"github.com/pablomotte/gator_go/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db *database.Queries
}

func main() {
	// Read the config file
	new_config, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config", err)
		os.Exit(1)
	}
	// Open the database
	db, err := sql.Open("postgres", new_config.DbURL)
	curr_state := state{
		config: &new_config,
		db: database.New(db),
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

