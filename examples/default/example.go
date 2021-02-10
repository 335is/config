package main

import (
	"fmt"
	"time"

	config "github.com/335is/config"
)

type cfg struct {
	Address string        `default:"http://website.org"`
	Timeout time.Duration `default:"22m33s"`
}

func main() {
	c := cfg{}
	config.Load("", &c)
	s, _ := config.ToYaml(&c)
	fmt.Printf("%s\n", s)
}
