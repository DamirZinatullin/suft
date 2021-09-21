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
	// без геттеров в интерфейсе посмотреть токены не получится. можно посмотреть их в auth_test.go
	// fmt.Printf("Access-token: %s \nRefresh-token: %s\n", client.AccessToken, client.RefreshToken)
	schedules, err := client.Schedules(nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nСписок расписаний\n")
	for _, schedule := range schedules {
		fmt.Printf("%#v\n",schedule)
	}
	schedule, err := client.DetailSchedule(32884)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nДетализация расписания\n")
	fmt.Printf("%#v\n",*schedule)

	periodId := api.PeriodId(353)
	schedule, err = client.AddSchedule(periodId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nДобавление расписания\n")
	fmt.Printf("%#v\n",*schedule)



}

