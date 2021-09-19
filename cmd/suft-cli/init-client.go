package main

import (
	"bufio"
	"encoding/json"
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

const configFileName string = "suft_config.json"
const configDirName string = "suft"

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


func configExists() (bool, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return false, err
	}
	configPath := path.Join(configDir, configDirName, configFileName)
	_, err = os.Stat(configPath)
	if err != nil {
		return false, err
	}
	return true, nil

}

