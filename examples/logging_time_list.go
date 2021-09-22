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
	loggingTimes, err := client.LoggingTimeList(32907, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список расписаний\n")
	for _, loggingTime := range loggingTimes {
		fmt.Println(loggingTime)
	}
}
