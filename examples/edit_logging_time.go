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
	addLT := logging_time.AddLoggingTime{
		CommentEmployee: "hi",
		Day1Time:        0,
		Day2Time:        0,
		Day3Time:        0,
		Day4Time:        1,
		Day5Time:        0,
		Day6Time:        0,
		Day7Time:        0,
		ProjectId:       8028,
		Task:            "dfgf",
		WorkKindId:      21,
	}

	loggingTimeCreated, err := client.AddLoggingTime(32907, &addLT)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Добавлен LoggingTime:\n")
	fmt.Printf("%#v\n\n", *loggingTimeCreated)

	loggingTimeEdit := logging_time.EditLoggingTime{
		CommentAdminEmployee: "Комментарий админа",
		CommentEmployee:      "Комментарий сотрудника",
		Day1Time:             2,
		Day2Time:             1,
		Day3Time:             0,
		Day4Time:             0,
		Day5Time:             0,
		Day6Time:             0,
		Day7Time:             0,
		ProjectId:            8028,
		StatusCode:           logging_time.Denied,
		Task:                 "",
		WorkKindId:           21,
	}

	loggingTime, err := client.EditLoggingTime(32907, 327801, &loggingTimeEdit)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Подробная информация по измененному LoggingTime\n")
	fmt.Println(loggingTime)
}