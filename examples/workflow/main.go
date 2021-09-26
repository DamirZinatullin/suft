package main

import (
	"fmt"
	"log"
	logging_time "suft_sdk/internal/logging-time"
	api "suft_sdk/pkg/api"
)

func main() {
	client1, err := api.NewClient("pakua", "147753")
	if err != nil {
		log.Fatalln(err)
	}

	periodId := api.PeriodId(368)
	schedule, err := client1.AddSchedule(periodId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Добавлено расписание")
	fmt.Printf("%#v\n\n", *schedule)

	loggingTime := logging_time.AddLoggingTime{
		CommentEmployee: "test9",
		Day1Time:        1,
		Day2Time:        1,
		Day3Time:        0,
		Day4Time:        1,
		Day5Time:        0,
		Day6Time:        3,
		Day7Time:        2,
		ProjectId:       69753,
		Task:            "test9",
		WorkKindId:      21,
	}
	loggingTimeCreated, err := client1.AddLoggingTime(api.ScheduleId(schedule.Id), &loggingTime)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Добавлен LoggingTime:\n")
	fmt.Printf("%#v\n\n", *loggingTimeCreated)

	scheduleForApprove, err := client1.SubmitForApproveSchedule(api.ScheduleId(schedule.Id))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Расписание, отправленное на утверждение")
	fmt.Printf("%#v\n\n", *scheduleForApprove)

	client2, err := api.NewClient("nikonovov", "147753")
	if err != nil {
		log.Fatalln(err)
	}

	loggingTimeApproved, err := client2.ApproveLoggingTime(api.ScheduleId(scheduleForApprove.Id), api.LoggingTimeId(loggingTimeCreated.Id))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Утверждённый объект трудозатрат")
	fmt.Printf("%#v\n\n", *loggingTimeApproved)

	scheduleApproved, err := client2.DetailSchedule(api.ScheduleId(scheduleForApprove.Id))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Утверждённое расписание\n")
	fmt.Printf("%#v\n", *scheduleApproved)
}
