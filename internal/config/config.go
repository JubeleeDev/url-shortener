package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr   string
	CodeLength int
	DSN        string
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

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is empty")
	}

	return Config{CodeLength: codeLength, HTTPAddr: addr, DSN: dsn}, nil
}
