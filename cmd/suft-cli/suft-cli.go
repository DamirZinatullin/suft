package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
	"suft_sdk/internal/auth"
	"syscall"
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
				reader := bufio.NewReader(os.Stdin)

				_, _ = os.Stdout.Write([]byte("Введите логин пользователя системы СУФТ:\n"))
				login, _ := reader.ReadString('\n')
				login = strings.Trim(login, "\n")

				_, _ = os.Stdout.Write([]byte("Введите пароль пользователя системы СУФТ:\n"))
				bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
				password := string(bytePassword)
				password = strings.Trim(password, "\n")

				token, err := auth.Authenticate(login, password)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(token.AccessToken, token.RefreshToken)
				//client, err := api.NewClient("demo@example.com", "demo")
				//fmt.Println(client,err)

				return nil
			}},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
