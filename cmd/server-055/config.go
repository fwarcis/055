package main

import (
	"flag"
	"os"
)

var defaultConfig = Config{
	Address: ":5500",
}
var cfg = NewConfig()

func NewConfig() Config {
	cfg := defaultConfig
	flag.Parse()
	if addr := flag.Arg(0); addr != "" {
		cfg.Address = addr
	} else if addr := os.Getenv("O055_ADDRESS"); addr != "" {
		cfg.Address = addr
	}
	return cfg
}

type Config struct {
	Address string
}
