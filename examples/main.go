package main

import (
	"fmt"
	"log"
	"suft_sdk/internal/auth"
	api "suft_sdk/pkg/API"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Access-token: %s \nRefresh-token: %s\n", client.AccessToken, client.RefreshToken)
	token, err := auth.Refresh(client.RefreshToken)
	fmt.Printf("Access-token: %s \nRefresh-token: %s\n", token.AccessToken, token.RefreshToken)
	schedules, err := client.Schedules(nil)
	for _, schedule := range schedules{
		fmt.Println(schedule)
	}
	fmt.Println(client.AccessToken == token.AccessToken)
	fmt.Println(client.RefreshToken == token.RefreshToken)
}