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
}

func main() {
	os.Setenv("MYAPP_ADDRESS", "http://example.com/funapp")
	os.Setenv("MYAPP_TIMEOUT", "1m30s")

	c := cfg{}
	config.Load("MYAPP", "", &c)
	fmt.Printf("Address: %s\n", c.Address)
	fmt.Printf("Timeout: %v\n", c.Timeout)
}
