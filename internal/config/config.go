package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr   string
	CodeLength int
}

func Load() (Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	length := os.Getenv("CODE_LENGTH")
	if length == "" {
		length = "8"
	}

	codeLength, err := strconv.Atoi(length)
	if err != nil {
		return Config{}, err
	}

	if codeLength <= 0 {
		return Config{}, fmt.Errorf("invalid code length")
	}

	return Config{CodeLength: codeLength, HTTPAddr: addr}, nil
}
