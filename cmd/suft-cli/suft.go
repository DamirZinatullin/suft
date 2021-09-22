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
	"strconv"
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

var scheduleId string

func main() {
	app := cli.NewApp()
	app.Name = "SUFT CLI"
	app.Usage = "CLI предоставляет возможность взаимодействия с api СУФТ (системы учета фактических трудозатрат)"
	app.Commands = []cli.Command{
		{
			Name:     "login",
			Usage:    "Аутентификация клиента",
			Category: "Клиент",
			Action: func(c *cli.Context) error {
				err := loginSuft()
				if err != nil {
					return err
				}
				fmt.Println("Клиент успешно прошел аутентификацию")
				return nil
			}},
		{
			Name:     "logout",
			Usage:    "Выход из клиента",
			Category: "Клиент",
			Action: func(c *cli.Context) error {
				err := logoutSuft()
				if err != nil {
					return err
				}
				fmt.Println("Успешный выход из клиента")
				return nil
			}},
		{
			Name:     "schedules",
			Usage:    "Список расписаний",
			Category: "Расписания",
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				schedules, err := client.Schedules(nil)
				if err != nil {
					return err
				}

				for _, schedule := range schedules {
					scheduleJSON, err := json.Marshal(schedule)
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", scheduleJSON)
				}
				return nil
			}},
		{
			Name:  "schedule",
			Usage: "Детализация расписания",
			Category: "Расписания",
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				schedIdInt, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return errors.New("required to pass valid scheduleID")
				}
				schedId := api.ScheduleId(schedIdInt)

				schedule, err := client.DetailSchedule(schedId)
				if err != nil {
					return err
				}
				scheduleJSON, err := json.Marshal(schedule)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", scheduleJSON)

				return nil
			}},
		{
			Name:        "addSchedule",
			Usage:       "Добавление расписания",
			Description: "Для добавления расписания необходимо передать Id периуда",
			Category:    "Расписания",
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				periodIdInt, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return errors.New("required to pass valid periodID")
				}
				periodId := api.PeriodId(periodIdInt)

				schedule, err := client.AddSchedule(periodId)
				if err != nil {
					return err
				}
				scheduleJSON, err := json.Marshal(schedule)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", scheduleJSON)

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

func logoutSuft() error {
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

func refreshConfig() error {
	_, err := configExists()
	if err != nil {
		return errors.New("не инициализирован клиент, выполните команду login")
	}
	configPath, _ := configPath()
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	userConf := userConfig{}
	err = json.Unmarshal(data, &userConf)
	if userConf.DateRefresh.Add(time.Minute * 2).After(time.Now()) {
		return nil
	}
	token, err := auth.Refresh(userConf.Token.RefreshToken)
	if err != nil {
		return err
	}
	err = writeConfig(token)
	if err != nil {
		return err
	}
	return nil
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

func genScheduleFile(client api.Client) error {
	return nil
}
