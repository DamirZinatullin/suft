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
	err = client.DeleteLoggingTime(32907, 327702)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Временная затрата успешно удалена\n")
}
