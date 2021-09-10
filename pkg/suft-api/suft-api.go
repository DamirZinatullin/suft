package suft_api

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"
	"suft_sdk/pkg/suft-api/schedule"
)

type ClientInterface interface {
	GetAuthToken() error
	Authorize() error
	Schedules() ([]schedule.Schedule, error)
	AddSchedule([]schedule.Schedule) error
	DetailSchedule(int) error
	EditSchedule(int, map[string]string) error
	UpdateSchedule(int, map[string]string) error
	Logout() error
}

type suftAPI struct {
	uri string
}

func NewSuftAPI(uri string) (*suftAPI, error) {
	return &suftAPI{
		uri: uri,
	}, nil
}



func (s *suftAPI) GetAuthToken() error {
	cli := &http.Client{
		Timeout: 10 * time.Second,
	}

	body := bytes.NewBuffer([]byte(`{"username":"demo@example.com","password":"demo"}`))
	req, err := http.NewRequest(
		http.MethodGet,
		s.uri,
		body,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Auth-method", "Password")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	if (resp.StatusCode / 100) != 2 {
		return errors.New("unable to get authentification tokens")
	}
	defer resp.Body.Close()
	return nil
}
