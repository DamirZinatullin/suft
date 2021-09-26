package main

import (
	"fmt"
	"log"
	"suftsdk/internal/loggingtime"
	"suftsdk/pkg/api"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo")
	if err != nil {
		log.Fatalln(err)
	}
	addLT := loggingtime.AddLoggingTime{
		CommentEmployee: "hi",
		Day1Time:        0,
		Day2Time:        0,
		Day3Time:        3,
		Day4Time:        0,
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

	logTimeId := api.LoggingTimeId(loggingTimeCreated.Id)
	err = client.DeleteLoggingTime(32907, logTimeId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Временная затрата успешно удалена\n")
}
