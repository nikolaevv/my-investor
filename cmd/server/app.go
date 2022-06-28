package main

import (
	"flag"
	"log"

	httpServer "github.com/nikolaevv/my-investor/internal/server/http"
)

var (
	configPath = flag.String("conf", "./configs/app.json", "path to config file")
)

func main() {
	flag.Parse()

	s, err := httpServer.New(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}
}
