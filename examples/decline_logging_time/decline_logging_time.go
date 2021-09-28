package main

import (
	"fmt"
	"log"
	"suftsdk/pkg/api"
)

func main() {
	client, err := api.NewClient("nikonovov", "147753", nil)
	if err != nil {
		log.Fatalln(err)
	}

	loggingTime, err := client.DeclineLoggingTime(32992, 327829, "не принято")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Отклонённый объект трудозатрат")
	fmt.Printf("%#v\n\n", *loggingTime)

}
