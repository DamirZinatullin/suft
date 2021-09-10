package main

import (
	"log"
	"os"
	suft_api "suft_sdk/pkg/suft-api"

	"github.com/urfave/cli"
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
				suftAPI := suft_api.NewSuftAPI()
				err := suftAPI.GetAuthTokens()
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
