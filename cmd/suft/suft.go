package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"suftsdk/internal/auth"
	"suftsdk/pkg/api"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/term"
)

type userConfig struct {
	Token       auth.Token `json:"token"`
	DateRefresh time.Time  `json:"date_refresh"`
}

const scheduleCategory string = "Расписания"
const loggingTimeCategory string = "Временные затраты"

const configFileName string = "suft_config.json"
const configDirName string = "suft"
const loggingTimeFileName string = "logging_time.json"

var scheduleId int
var loggingTimeId int
var periodId int
var page int
var size int
var role string
var editor string
var adminComment string

var scheduleIdFlag cli.Flag = cli.IntFlag{
	Name:        "schedule-id, scid",
	Usage:       "id расписания",
	Required:    true,
	Destination: &scheduleId,
}

var pageFlag cli.Flag = cli.IntFlag{
	Name:        "page, p",
	Usage:       "Страница отображения",
	Destination: &page,
}

var sizeFlag cli.Flag = cli.IntFlag{
	Name:        "size, s",
	Usage:       "Количество отображаемых элементов",
	Destination: &size,
}

var periodFlag cli.Flag = cli.IntFlag{
	Name:        "period-id, pid",
	Usage:       "id периуда",
	Required:    true,
	Destination: &periodId,
}

var loggingTimeIdFlag cli.Flag = cli.IntFlag{
	Name:        "logging-time-id, ltid",
	Usage:       "id временной затраты",
	Required:    true,
	Destination: &loggingTimeId,
}

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
			Category: scheduleCategory,
			Aliases:  []string{"scs"},
			Flags: []cli.Flag{
				pageFlag,
				sizeFlag,
				cli.StringFlag{
					Name:        "role, r",
					Usage:       "Роль клиента (approver или creator)",
					Destination: &role,
				},
			},
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				options := api.OptionsS{}
				if size != 0 {
					options.Size = size
				}
				if page != 0 {
					options.Page = page
				}
				if role != "" {
					clientRole := api.Role(role)
					options.CreatorApprover = clientRole
				}
				schedules, err := client.Schedules(&options)
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
			Name:     "schedule",
			Usage:    "Детализация расписания",
			Aliases:  []string{"sc"},
			Category: scheduleCategory,
			Flags: []cli.Flag{
				scheduleIdFlag,
			},
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				schedId := api.ScheduleId(scheduleId)

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
			Name:        "add-schedule",
			Usage:       "Добавление расписания",
			Description: "Для добавления расписания необходимо передать id периуда",
			Category:    scheduleCategory,
			Aliases:     []string{"as"},
			Flags: []cli.Flag{
				periodFlag,
			},
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}

				periodId := api.PeriodId(periodId)

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
		{
			Name:     "submit-for-approve",
			Usage:    "Отправить расписание на утверждение",
			Aliases:  []string{"s"},
			Category: scheduleCategory,
			Flags: []cli.Flag{
				scheduleIdFlag,
			},
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				schedId := api.ScheduleId(scheduleId)

				schedule, err := client.SubmitForApproveSchedule(schedId)
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
			Name:        "logging-times",
			Usage:       "Список временных затрат",
			Aliases:     []string{"lts"},
			Description: "Для вывода списка временных затрат необходимо передать id расписания",
			Flags: []cli.Flag{
				scheduleIdFlag,
				pageFlag,
				sizeFlag,
			},
			Category: loggingTimeCategory,
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				options := api.OptionsLT{}
				if size != 0 {
					options.Size = size
				}
				if page != 0 {
					options.Page = page
				}
				scheduleId := api.ScheduleId(scheduleId)
				loggingTimeList, err := client.LoggingTimeList(scheduleId, &options)
				if err != nil {
					return err
				}
				for _, schedule := range loggingTimeList {
					loggingTimeJSON, err := json.Marshal(schedule)
					if err != nil {
						return err
					}
					fmt.Printf("%s\n\n", loggingTimeJSON)
				}
				return nil
			}},
		{
			Name:        "logging-time",
			Usage:       "Детализация временной затраты",
			Description: "Для вывода временной затраты необходимо передать id расписания и id временой затраты",
			Aliases:     []string{"lt"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				loggingTimeIdFlag,
			},
			Category: loggingTimeCategory,
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				scheduleId := api.ScheduleId(scheduleId)
				loggingTimeId := api.LoggingTimeId(loggingTimeId)
				loggingTime, err := client.DetailLoggingTime(scheduleId, loggingTimeId)
				if err != nil {
					return err
				}
				loggingTimeJSON, err := json.Marshal(loggingTime)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n\n", loggingTimeJSON)
				return nil
			}},
		{
			Name:     "add-logging-time",
			Usage:    "Добавление временной затраты",
			Category: loggingTimeCategory,
			Aliases:  []string{"al"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				cli.StringFlag{
					Name:        "editor, e",
					Usage:       "Используемый текстовый редактор",
					Destination: &editor,
					Value:       "vim",
					EnvVar:      "EDITOR",
				},
			},
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				path, err := genLoggingTimeFile()
				if err != nil {
					return err
				}
				cmd := exec.Command(editor, path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				err = cmd.Run()
				if err != nil {
					return err
				}
				scheduleId := api.ScheduleId(scheduleId)
				loggingTime, err := loggingTimeFromFIle()
				if err != nil {
					return err
				}
				loggingTimeResp, err := client.AddLoggingTime(scheduleId, loggingTime)
				if err != nil {
					return err
				}
				LoggingTimeJSON, err := json.Marshal(loggingTimeResp)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", LoggingTimeJSON)

				return nil
			}},
		{
			Name:     "remove-logging-time",
			Usage:    "Удаление временной затраты",
			Category: loggingTimeCategory,
			Aliases:  []string{"rmlt"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				loggingTimeIdFlag,
			},
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				scheduleId := api.ScheduleId(scheduleId)
				loggingTimeId := api.LoggingTimeId(loggingTimeId)
				err = client.DeleteLoggingTime(scheduleId, loggingTimeId)
				if err != nil {
					return err
				}
				fmt.Println("Временная затрата успешно удалена")
				return nil
			}},
		{
			Name:        "approve-logging-time",
			Usage:       "Утверждение временной затраты",
			Description: "Для утверждения временной затраты необходимо передать id расписания и id временой затраты",
			Aliases:     []string{"aprv"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				loggingTimeIdFlag,
				cli.StringFlag{
					Name:        "comment, c",
					Usage:       "Комментарий руководителя/согласующего",
					Destination: &adminComment,
				},
			},
			Category: loggingTimeCategory,
			Action: func(c *cli.Context) error {
				err := refreshConfig()
				if err != nil {
					return err
				}
				client, err := newClientFromConfig()
				if err != nil {
					return err
				}
				scheduleId := api.ScheduleId(scheduleId)
				loggingTimeId := api.LoggingTimeId(loggingTimeId)
				loggingTime, err := client.ApproveLoggingTime(scheduleId, loggingTimeId, adminComment)
				if err != nil {
					return err
				}
				loggingTimeJSON, err := json.Marshal(loggingTime)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n\n", loggingTimeJSON)
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

	token, err := auth.Authenticate(login, password, nil)
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
		output, err = os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC, os.ModePerm)
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
	if err != nil {
		log.Fatalln(err)
	}
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
	if err != nil {
		return err
	}
	if userConf.DateRefresh.Add(time.Minute * 2).After(time.Now()) {
		return nil
	}
	token, err := auth.Refresh(userConf.Token.RefreshToken, nil)
	if err != nil {
		return errors.New("время сессии истекло, пройдите аутентификацию, выполнив команду login")
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

func genLoggingTimeFile() (path string, err error) {
	var output *os.File
	filePath, err := loggingTimeFilePath()
	fmt.Println(filePath)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(filePath)
	if err != nil {
		output, err = os.Create(filePath)
		if err != nil {
			return "", err
		}
	} else {
		output, err = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	defer output.Close()
	jsonEncoder := json.NewEncoder(output)
	loggingTime := api.AddLoggingTime{
		CommentEmployee: "",
		Day1Time:        0,
		Day2Time:        0,
		Day3Time:        0,
		Day4Time:        0,
		Day5Time:        0,
		Day6Time:        0,
		Day7Time:        0,
		ProjectId:       0,
		Task:            "",
		WorkKindId:      0,
	}
	err = jsonEncoder.Encode(loggingTime)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func loggingTimeFilePath() (filePath string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	filePath = path.Join(configDir, configDirName, loggingTimeFileName)
	return filePath, nil

}

func loggingTimeFromFIle() (loggingTime *api.AddLoggingTime, err error) {
	path, err := loggingTimeFilePath()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	logTime := api.AddLoggingTime{}
	err = json.Unmarshal(data, &logTime)
	if err != nil {
		log.Fatalln(err)
	}
	return &logTime, nil
}
