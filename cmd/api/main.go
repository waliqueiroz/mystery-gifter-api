package main

import (
	"log"

	"github.com/waliqueiroz/mystery-gifter-api/internal/infra"
)

func main() {
	if err := infra.Run(); err != nil {
		log.Fatal(err)
	}
}
