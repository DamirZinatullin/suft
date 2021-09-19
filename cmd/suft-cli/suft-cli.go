package main

import (
	"bufio"
	"encoding/json"
	"github.com/urfave/cli"
	"golang.org/x/term"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"suft_sdk/internal/auth"
	"syscall"
	"time"
)

type userConfig struct {
	Token       auth.Token `json:"token"`
	DateRefresh time.Time  `json:"date_refresh"`
}

func writeConfig(token *auth.Token) error {
	var output *os.File
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configDir = path.Join(configDir, "suft")
	_ = os.Mkdir(configDir, fs.ModePerm)
	configPath := path.Join(configDir, "suft_config.json")
	_, err = os.Stat(configPath)
	if err != nil {
		output, err = os.Create(configPath)
		if err != nil {
			return err
		}
	} else {
		output, err = os.OpenFile(configPath, os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}
	}
	defer output.Close()
	jsonEncoder := json.NewEncoder(output)
	user := userConfig{
		Token:       *token,
		DateRefresh: time.Now(),
	}
	err = jsonEncoder.Encode(user)
	if err != nil {
		return err
	}
	return nil
}

func initSuft() error {
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
	err = writeConfig(token)
	if err != nil {
		return err
	}
	return nil
}

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
