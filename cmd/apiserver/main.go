package main

import (
	"github.com/v1shn3vsk7/PlaylistAPI/internal/server"
	"log"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
