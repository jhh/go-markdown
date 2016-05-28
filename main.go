package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open("test.md")
	if err != nil {
		log.Fatal(err)
	}
	_, err = NewMeta(f)
	if err != nil {
		log.Fatal(err)
	}
}
