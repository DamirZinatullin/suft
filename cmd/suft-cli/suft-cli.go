package main

import (
	"fmt"
	"log"
	"os"

	api "suft_sdk/pkg/API"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "SUFT CLI"
	app.Usage = "CLI предоставляет возможность взаимодействия с API СУФТ (системы учета фактических трудозатрат)"
	Flags := []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:  "schedules",
			Usage: "Получение списка расписаний",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				suftAPI := api.NewSuftAPI()
				schedules, err := suftAPI.Schedules()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(schedules)
				return nil
			},

		},
		{
			Name: "add",
			Usage: "Добавление расписания",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "update",
			Usage: "Изменение расписания",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "submit",
			Usage: "Отправка расписания на утвердение",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "approve",
			Usage: "Утверждение расписания",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
