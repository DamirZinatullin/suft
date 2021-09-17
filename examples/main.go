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
	// без геттеров в интерфейсе посмотреть токены не получится. можно посмотреть их в auth_test.go
	// fmt.Printf("Access-token: %s \nRefresh-token: %s\n", client.AccessToken, client.RefreshToken)

	schedules, err := client.Schedules(nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, schedule := range schedules {
		fmt.Println(schedule)
	}
}
