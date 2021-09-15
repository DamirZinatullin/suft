package main

import (
	"log"
	"os"
	"suft_sdk/internal/http-client"
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
			Name:  "auth",
			Usage: "Получение авторизационных токенов",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				client := api.NewSuftAPI()
				client.
				suftAPI := http_client.NewHttpClient()
				err := suftAPI.AuthTokens()
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
