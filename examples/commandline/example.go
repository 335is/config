package main

import (
	"fmt"
	"os"
	"time"

	config "github.com/335is/config"
)

type cfg struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Sub     sub
}

type sub struct {
	Enabled bool `yaml:"enabled"`
	Level   int  `yaml:"level"`
}

func main() {
	// override the command line for demonstration purposes
	os.Args = []string{
		"example",
		"Address=http://example.com",
		"Timeout=1h20m30s",
		"Sub.Enabled=true",
		"Sub.Level=42",
	}

	c := cfg{}

	err := config.FromArguments(os.Args[1:], &c)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", c)
}
