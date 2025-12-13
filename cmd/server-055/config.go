package main

import "os"

func NewConfig() Config {
	address := "[::]:8550"
	if len(os.Args) >= 2 {
		address = os.Args[1]
	} else if envAddr := os.Getenv("055_ADDRESS"); envAddr != "" {
		address = envAddr
	}

	return Config{address}
}

type Config struct {
	Address string
}
