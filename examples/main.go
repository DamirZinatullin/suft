package main

import (
	"fmt"
	"log"
	api "suft_sdk/pkg/API"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo")
	if err != nil {
		log.Fatalln("error in client creating", err)
	}
	fmt.Println(client.AccessToken, client.RefreshToken)
}
