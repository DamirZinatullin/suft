package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"suftsdk/internal/clifuncs"
	"suftsdk/pkg/api"

	"github.com/urfave/cli"
)

const scheduleCategory string = "Расписания"
const loggingTimeCategory string = "Временные затраты"

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

var roleFlag cli.Flag = cli.StringFlag{
	Name:        "role, r",
	Usage:       "Роль клиента (approver или creator)",
	Destination: &role,
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

var editorFlag cli.Flag = cli.StringFlag{
	Name:        "editor, e",
	Usage:       "Используемый текстовый редактор",
	Destination: &editor,
	Value:       "vim",
	EnvVar:      "EDITOR",
}

var commentFlag cli.Flag = cli.StringFlag{
	Name:        "comment, c",
	Usage:       "Комментарий руководителя/согласующего",
	Destination: &adminComment,
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
			Action:   login,
		},
		{
			Name:     "logout",
			Usage:    "Выход из клиента",
			Category: "Клиент",
			Action:   logout,
		},
		{
			Name:     "schedules",
			Usage:    "Список расписаний",
			Category: scheduleCategory,
			Aliases:  []string{"scs"},
			Flags: []cli.Flag{
				pageFlag,
				sizeFlag,
				roleFlag,
			},
			Action: schedules,
		},
		{
			Name:     "schedule",
			Usage:    "Детализация расписания",
			Aliases:  []string{"sc"},
			Category: scheduleCategory,
			Flags: []cli.Flag{
				scheduleIdFlag,
			},
			Action: scheduleDetail,
		},
		{
			Name:        "add-schedule",
			Usage:       "Добавление расписания",
			Description: "Для добавления расписания необходимо передать id периуда",
			Category:    scheduleCategory,
			Aliases:     []string{"as"},
			Flags: []cli.Flag{
				periodFlag,
			},
			Action: addSchedule,
		},
		{
			Name:     "submit-for-approve",
			Usage:    "Отправить расписание на утверждение",
			Aliases:  []string{"s"},
			Category: scheduleCategory,
			Flags: []cli.Flag{
				scheduleIdFlag,
			},
			Action: submitForApprove,
		},

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
			Action:   loggingTimes,
		},
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
			Action:   loggingTimeDetail,
		},
		{
			Name:     "add-logging-time",
			Usage:    "Добавление временной затраты",
			Category: loggingTimeCategory,
			Aliases:  []string{"al"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				editorFlag,
			},
			Action: addLoggingTime,
		},
		{
			Name:     "remove-logging-time",
			Usage:    "Удаление временной затраты",
			Category: loggingTimeCategory,
			Aliases:  []string{"rmlt"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				loggingTimeIdFlag,
			},
			Action: removeLoggingTime,
		},
		{
			Name:        "approve-logging-time",
			Usage:       "Утверждение временной затраты",
			Description: "Для утверждения временной затраты необходимо передать id расписания и id временой затраты",
			Aliases:     []string{"aprv"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				loggingTimeIdFlag,
				commentFlag,
			},
			Category: loggingTimeCategory,
			Action:   approveLoggingTime,
		},
		{
			Name:        "decline-logging-time",
			Usage:       "Отклонение временной затраты",
			Description: "Для отклонения временной затраты необходимо передать id расписания и id временой затраты",
			Aliases:     []string{"dcl"},
			Flags: []cli.Flag{
				scheduleIdFlag,
				loggingTimeIdFlag,
				commentFlag,
			},
			Category: loggingTimeCategory,
			Action:   declineLoggingTime,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func login(c *cli.Context) error {
	err := clifuncs.LoginSuft()
	if err != nil {
		return err
	}
	fmt.Println("Клиент успешно прошел аутентификацию")
	return nil
}

func logout(c *cli.Context) error {
	err := clifuncs.LogoutSuft()
	if err != nil {
		return err
	}
	fmt.Println("Успешный выход из клиента")
	return nil
}

func schedules(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
	if err != nil {
		return err
	}
	options := api.OptionsS{}
	if size != 0 {
		options.Size = size
	}
	options.Page = page
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
}

func scheduleDetail(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
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
}

func addSchedule(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
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
}

func submitForApprove(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
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
}

func loggingTimes(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
	if err != nil {
		return err
	}
	options := api.OptionsLT{}
	if size != 0 {
		options.Size = size
	}
	options.Page = page
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
}

func loggingTimeDetail(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
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
}

func addLoggingTime(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
	if err != nil {
		return err
	}
	path, err := clifuncs.GenLoggingTimeFile()
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
	loggingTime, err := clifuncs.LoggingTimeFromFile()
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
}

func removeLoggingTime(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
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
}

func approveLoggingTime(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
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
}

func declineLoggingTime(c *cli.Context) error {
	err := clifuncs.RefreshConfig()
	if err != nil {
		return err
	}
	client, err := clifuncs.NewClientFromConfig()
	if err != nil {
		return err
	}
	scheduleId := api.ScheduleId(scheduleId)
	loggingTimeId := api.LoggingTimeId(loggingTimeId)
	loggingTime, err := client.DeclineLoggingTime(scheduleId, loggingTimeId, adminComment)
	if err != nil {
		return err
	}
	loggingTimeJSON, err := json.Marshal(loggingTime)

	if err != nil {
		return err
	}
	fmt.Printf("%s\n\n", loggingTimeJSON)
	return nil
}
