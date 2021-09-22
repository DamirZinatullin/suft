package main

import (
	"fmt"
	"log"
	api "suft_sdk/pkg/api"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo")
	if err != nil {
		log.Fatalln(err)
	}

	schedules, err := client.Schedules(nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список расписаний\n")
	for _, schedule := range schedules {
		fmt.Printf("%#v\n", schedule)
	}
}
