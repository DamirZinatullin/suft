package main

import (
	"fmt"
	"log"
	"suftsdk/pkg/api"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo", nil)
	if err != nil {
		log.Fatalln(err)
	}

	schedule, err := client.DetailSchedule(32884)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Детализация расписания\n")
	fmt.Printf("%#v\n", *schedule)
}
