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

	addLT := api.AddLoggingTime{
		CommentEmployee: "hi",
		Day1Time:        1,
		Day2Time:        1,
		Day3Time:        0,
		Day4Time:        1,
		Day5Time:        0,
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
	fmt.Printf("Добавлен LoggingTime:\n")
	fmt.Printf("%#v\n\n", *loggingTimeCreated)
}
