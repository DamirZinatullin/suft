package main

import (
	"fmt"
	"log"
	api "suft_sdk/pkg/API"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Access-token: %s \nRefresh-token: %s\n", client.AccessToken, client.RefreshToken)
}
