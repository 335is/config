package main

import (
	"fmt"
	"os"
	"time"

	config "github.com/335is/config"
)

type cfg struct {
	Address string
	Timeout time.Duration
}

func main() {
	os.Setenv("ADDRESS", "http://example.com/funapp")
	os.Setenv("TIMEOUT", "1m30s")

	c := cfg{}
	config.Load("", &c)
	fmt.Printf("Address: %s\n", c.Address)
	fmt.Printf("Timeout: %v\n", c.Timeout)
}
