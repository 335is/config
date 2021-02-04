package main

import (
	"fmt"
	"time"

	config "github.com/335is/config"
)

type cfg1 struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

type cfg2 struct {
	HTTP cfg1
}

func main() {
	c1 := cfg1{}
	config.Load("MYAPP1", "./cfg1.yml", &c1)
	s1, _ := config.ToYaml(&c1)
	fmt.Printf("%s\n", s1)

	c2 := cfg2{}
	config.Load("MYAPP2", "./cfg2.yml", &c2)
	s2, _ := config.ToYaml(&c2)
	fmt.Printf("%s\n", s2)
}
