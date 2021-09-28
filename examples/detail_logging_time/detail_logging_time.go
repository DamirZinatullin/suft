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

	loggingTime, err := client.DetailLoggingTime(32907, 327703)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Подробная информация по LoggingTime\n")
	fmt.Printf("%#v\n\n", *loggingTime)

}
