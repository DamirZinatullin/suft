package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"suft_sdk/pkg/suft-api"
)

func main() {
	app := cli.NewApp()
	app.Name = "SUFT CLI"
	app.Usage = "CLI предоставляет возможность взаимодействия с API СУФТ (системы учета фактических трудозатрат)"
	Flags := []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:  "auth",
			Usage: "Получение авторизационных токенов",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				suftAPI, err := suft_api.NewSuftAPI("test")
				if err != nil {
					return err
				}
				err = suftAPI.GetAuthToken()
				if err != nil {
					return err
				}
				return nil
			}},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
