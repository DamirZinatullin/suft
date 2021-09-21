package main

import (
	"fmt"
	"log"
	logging_time "suft_sdk/internal/logging-time"
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
		fmt.Printf("%#v\n", schedule)
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
	fmt.Printf("%#v\n\n",*schedule)

	loggingTimes, err := client.LoggingTimeList(32907, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nСписок расписаний\n")
	for _, loggingTime := range loggingTimes {
		fmt.Println(loggingTime)
	}

	addLT := logging_time.AddLoggingTime{
		CommentEmployee: "hi",
		Day1Time:        4,
		Day2Time:        5,
		Day3Time:        3,
		Day4Time:        1,
		Day5Time:        5,
		Day6Time:        3,
		Day7Time:        2,
		ProjectId:       8028,
		Task:            "dfgf",
		WorkKindId:      21,
	}
	loggingTimeCreated, err := client.AddLoggingTime(32907, &addLT)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nДобавление LoggingTime\n")
	fmt.Printf("%#v\n\n",*loggingTimeCreated)

}
