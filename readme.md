# Config

Config loads configuration settings into a struct from the following sources in this order.

1. default struct values
1. config.yml YAML file
1. environment variables

Each source successively overrides the previous ones.

## Requirements

### Go 1.11.5

[Download Go](https://golang.org/dl/)

### Dep 0.4.1

[Dep Releases](https://github.com/golang/dep/releases/)

## How to Use

See the the examples folder and the unit tests for how to define and consume your configuration.
