package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)


func main() {
	app := cli.NewApp()
	app.Name = "SUFT CLI"
	app.Usage = "CLI предоставляет возможность взаимодействия с API СУФТ (системы учета фактических трудозатрат)"
	Flags := []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Инициализация клиента",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				err := initSuft()
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
