package main

import (
	"flag"
	"log"

	"github.com/nikolaevv/my-investor/internal/server"
)

var (
	configPath = flag.String("conf", "./configs/app.json", "path to config file")
)

func main() {
	flag.Parse()

	s, err := server.New(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}
}
