package main

import (
	"log"

	"github.com/touchmeangel/protocols/config"
	"github.com/touchmeangel/protocols/defillama"
)

func main() {
	cfg := config.Load()

	client, err := defillama.New(cfg.Proxies)
	if err != nil {
		log.Fatal(err)
	}

	err = client.GetAllProtocols()
	if err != nil {
		log.Fatal(err)
	}
}
