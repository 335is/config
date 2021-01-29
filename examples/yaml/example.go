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
	s, _ := config.ToYaml(&c)
	fmt.Printf("%s\n", s)
}
