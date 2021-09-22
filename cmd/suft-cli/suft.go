package main

import (

	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/term"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"suft_sdk/internal/auth"
	api "suft_sdk/pkg/api"
	"syscall"
	"time"
)

type userConfig struct {
	Token       auth.Token `json:"token"`
	DateRefresh time.Time  `json:"date_refresh"`
}

const configFileName string = "suft_config.json"
const configDirName string = "suft"


func main() {
	app := cli.NewApp()
	app.Name = "SUFT CLI"
	app.Usage = "CLI предоставляет возможность взаимодействия с api СУФТ (системы учета фактических трудозатрат)"
	Flags := []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Аутентификация клиента",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				err := loginSuft()
				if err != nil {
					return err
				}
				fmt.Println("Клиент успешно прошел аутентификацию")
				return nil
			}},
		{
			Name:  "logout",
			Usage: "Выход из клиента",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				err := logoutSuft()
				if err != nil {
					return err
				}
				fmt.Println("Успешный выход из клиента")
				return nil
			}},
		{
			Name:  "schedules",
			Usage: "Список расписаний",
			Flags: Flags,
			Action: func(c *cli.Context) error {
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				schedules, err := client.Schedules(nil)
				if err != nil {
					return err
				}

				for _, schedule := range schedules{
					scheduleJSON, err := json.Marshal(schedule)
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", scheduleJSON)
				}
				return nil
			}},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loginSuft() error {
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

func logoutSuft() error{
	_, err := configExists()
	if err != nil {
		return err
	}
	configPath, _ := configPath()
	err = os.Remove(configPath)
	if err != nil {
		return err
	}
	return nil
}

func writeConfig(token *auth.Token) error {
	var output *os.File
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configDir = path.Join(configDir, configDirName)
	_ = os.Mkdir(configDir, fs.ModePerm)
	configPath := path.Join(configDir, configFileName)
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

func configPath() (pathConfig string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := path.Join(configDir, configDirName, configFileName)
	return configPath, nil

}

func newClientFromConfig() (client api.API, err error) {
	_, err = configExists()
	if err != nil {
		return nil, errors.New("не инициализирован клиент, выполните команду login")
	}
	configPath, _ := configPath()
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	userConf := userConfig{}
	err = json.Unmarshal(data, &userConf)
	client = &api.Client{
		BaseURL:      api.BaseURL,
		AccessToken:  userConf.Token.AccessToken,
		RefreshToken: userConf.Token.RefreshToken,
		HttpClient: &http.Client{
			Timeout: time.Minute,
		},
	}

	return client, nil
}

func configExists() (bool, error) {
	configPath, err := configPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(configPath)
	if err != nil {
		return false, err
	}
	return true, nil

}

