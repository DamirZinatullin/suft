package suft_api

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	URI string = "https://dev.gnivc.ru/tools/suft/api/v1/"
)

type ClientInterface interface {
	GetAuthToken() error
	Authorize() error
	Schedules() ([]Schedule, error)
	AddSchedule(Schedule) error
	DetailSchedule(int) error
	EditSchedule(int, map[string]string) error
	UpdateSchedule(int, map[string]string) error
	Logout() error
}

type suftAPI struct {
	URN          string
	Method       string
	AccessToken  string
	RefreshToken string
}

func NewSuftAPI() *suftAPI {
	return &suftAPI{}
}

type Employee struct {
	//todo
}

type Period struct {

	//todo
}

type Schedule struct {
	author Employee
	id     int
	period Period
}

func (s *suftAPI) GetAuthToken() error {
	cli := &http.Client{
		Timeout: 10 * time.Second,
	}

	reqBody := bytes.NewBuffer([]byte(`{"username":"demo@example.com","password":"demo"}`))
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprint(URI, "security/authenticate"),
		reqBody,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Auth-method", "Password")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unable to get authentification tokens")
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "Access-token" {
			s.AccessToken = cookie.Value
		}
		if cookie.Name == "Refresh-token" {
			s.RefreshToken = cookie.Value
		}
	}

	return nil
}
