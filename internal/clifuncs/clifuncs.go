package clifuncs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"suftsdk/internal/auth"
	"suftsdk/pkg/api"
	"syscall"
	"time"
)

const configFileName string = "suft_config.json"
const configDirName string = "suft"
const loggingTimeFileName string = "logging_time.json"

type ClientBuilder interface {
	NewClient() (client api.API, err error)
}

type userConfig struct {
	Token       auth.Token
	DateRefresh time.Time
}

type ClientInit struct{}

func (c *ClientInit) NewClient() (client api.API, err error) {
	err = RefreshConfig()
	if err != nil {
		return nil, err
	}
	client, err = NewClientFromConfig()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewClientFromConfig() (client api.API, err error) {
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

func LoginSuft() error {
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

func LogoutSuft() error {
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

func RefreshConfig() error {
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

func GenLoggingTimeFile() (path string, err error) {
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

func LoggingTimeFromFile() (loggingTime *api.AddLoggingTime, err error) {
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
