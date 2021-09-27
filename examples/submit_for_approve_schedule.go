package main

import (
	"fmt"
	"log"
	"suftsdk/pkg/api"
)

func main() {
	client, err := api.NewClient("pakua", "147753", nil)
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
