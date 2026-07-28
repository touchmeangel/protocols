package main

import (
	"log"

	"github.com/touchmeangel/protocols/defillama"
)

func main() {
	client, err := defillama.New()
	if err != nil {
		log.Fatal(err)
	}

	err = client.GetAllProtocols()
	if err != nil {
		log.Fatal(err)
	}
}
