package main

import (
	"fmt"
	"log"
	"suftsdk/pkg/api"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo")
	if err != nil {
		log.Fatalln(err)
	}
	loggingTimes, err := client.LoggingTimeList(32907, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список временных затрат\n")
	for _, loggingTime := range loggingTimes {
		fmt.Printf("%#v\n\n", loggingTime)
	}
	loggingTimes, err = client.LoggingTimeList(32907, &api.OptionsLT{Size: 2, Page: 1})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nСписок временных затрат при передаче опций\n")
	for _, loggingTime := range loggingTimes {
		fmt.Printf("%#v\n\n", loggingTime)
	}

}
