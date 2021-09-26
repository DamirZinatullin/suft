package main

import (
	"fmt"
	"log"
	api "suft_sdk/pkg/api"
)

func main() {
	client, err := api.NewClient("pakua", "147753")
	if err != nil {
		log.Fatalln(err)
	}

	schedule, err := client.SubmitForApproveSchedule(32992)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Расписание с изменённым полем statusCode")
	fmt.Printf("%#v\n", *schedule)
}
