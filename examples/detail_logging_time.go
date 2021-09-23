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

	loggingTime, err := client.DetailLoggingTime(32907, 327702)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Подробная информация по LoggingTime\n")
	fmt.Println(loggingTime)
}
