package main

import (
	"log"

	"github.com/erwinvaneyk/go-version/cmd/goversion/cmd"
)

func main() {
	err := cmd.NewCmdRoot().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
