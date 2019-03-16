package main

import (
	"fmt"
	"time"

	config "github.com/335is/config"
)

type cfg struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

func main() {
	c := cfg{}
	config.Load("MYAPP", "./cfg.yml", &c)
	fmt.Printf("Address: %s\n", c.Address)
	fmt.Printf("Timeout: %v\n", c.Timeout)
}
