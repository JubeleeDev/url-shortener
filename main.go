package main

import (
	"fmt"

	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

func printCode(length int) {
	code, err := shortener.GenerateCode(length)
	if err != nil {
		fmt.Println("failed to generate code:", err)
		return
	}

	fmt.Println("generated code:", code)
}

func main() {
	printCode(15)
	printCode(0)
	printCode(5)
	printCode(-3)
}
