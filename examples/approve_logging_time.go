package main

import (
	"fmt"
	"log"
	api "suft_sdk/pkg/api"
)

func main() {
	client, err := api.NewClient("nikonovov", "147753")
	if err != nil {
		log.Fatalln(err)
	}

	loggingTime, err := client.ApproveLoggingTime(32992, 327829)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Объект трудозатрат с изменённым полем statusCode")
	fmt.Printf("%#v\n\n", *loggingTime)

}
