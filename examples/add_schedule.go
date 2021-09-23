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

	periodId := api.PeriodId(353)
	schedule, err := client.AddSchedule(periodId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Добавлено расписание\n")
	fmt.Printf("%#v\n\n", *schedule)

}
