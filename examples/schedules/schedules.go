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

	schedules, err := client.Schedules(nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список расписаний\n")
	for _, schedule := range schedules {
		fmt.Printf("%#v\n", schedule)
	}

	schedules, err = client.Schedules(&api.OptionsS{Page: 1, Size: 2})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nСписок расписаний при передаче опций\n")
	for _, schedule := range schedules {
		fmt.Printf("%#v\n", schedule)
	}

	fmt.Printf("\nСписок расписаний при передаче опций после отправки на утверждение\n")
	for _, schedule := range schedules {
		scheduleResp, err := schedule.SubmitForApproveSchedule()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", scheduleResp)
	}

}
