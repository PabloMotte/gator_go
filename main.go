package main

import (
	"fmt"
	"github.com/pablomotte/gator_go/internal/config"
)

func main() {
	fmt.Println("Hello, World!")
	config, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config", err)
	} else {
		print_config(config)
	}
	// config.SetUser("paul")
	// print_config(config)
	// config.SetUser("bob")
	// print_config(config)
	// config.SetUser("")
	// print_config(config)
}

func print_config(cfg config.Config) {
	for k,v := range []interface{}{cfg.DbURL, cfg.CurrentUserName} {
		fmt.Println(k, v)
	}
}
