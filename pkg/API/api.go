package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"suft_sdk/internal/auth"
	"suft_sdk/internal/http-client/logging-time"
	"suft_sdk/internal/http-client/schedule"
	"time"
)

const (
	BaseURL string = "https://dev.gnivc.ru/tools/suft/api/v1/"
)

type Options struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type API interface {
	Schedules(options *Options) ([]schedule.Schedule, error)
	AddSchedule(periodId int) error
	DetailSchedule(scheduleId int) error
	LoggingTimeList(scheduleId int, options *Options) ([]logging_time.LoggingTime, error)
	AddLoggingTime(scheduleId int, loggingTime *logging_time.AddLoggingTime) error
	DetailLoggingTime(scheduleId int, loggingTimeId int) error
	EditLoggingTime(scheduleId int, loggingTimeId int, loggingTime *logging_time.EditLoggingTime)
	DeleteLoggingTime(scheduleId int, loggingTimeId int) error
	SubmitForApproveSchedule(scheduleId int, loggingTimeId int, status *schedule.EditStatusSchedule) error
	ApproveSchedule(scheduleId int, loggingTimeId int, status *schedule.EditStatusSchedule) error
}

type client struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
	request      *http.Request
	httpClient   *http.Client
}

func NewClient(email string, password string) (API, error) {
	token, err := auth.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	return &client{
		BaseURL:      BaseURL,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		request:      &http.Request{Header: map[string][]string{"Auth-method": []string{"password"},
			"Content-type": []string{"application/json; charset=UTF-8"}}},
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}, nil
}

func (c *client) Schedules(options *Options) ([]schedule.Schedule, error) {
	cli := c.httpClient
	req := c.request
	req.Method = http.MethodGet
	reqURL, err := url.Parse(fmt.Sprint(BaseURL, "api/v1/schedules?page=1&size=5&creatorApprover=creator"))
	if err != nil {
		return nil, err
	}
	req.URL = reqURL
	cookie1 := &http.Cookie{
		Name:  "Access-token",
		Value: c.AccessToken,
	}
	cookie2 := &http.Cookie{
		Name:  "Refresh-token",
		Value: c.RefreshToken,
	}
	req.AddCookie(cookie1)
	req.AddCookie(cookie2)

	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to get schedules")
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to read response body:", err)
		return nil, err
	}

	schedules := make([]schedule.Schedule, 1)
	err = json.Unmarshal(respB, &schedules)
	if err != nil {
		log.Println("unable to unmarshal response body:", err)
		return nil, err
	}

	return schedules, nil
}

func (c *client) AddSchedule(periodId int) error {
	panic("implement me")
}

func (c *client) DetailSchedule(scheduleId int) error {
	panic("implement me")
}

func (c *client) LoggingTimeList(scheduleId int, options *Options)([]logging_time.LoggingTime, error) {
	panic("implement me")
}

func (c *client) AddLoggingTime(scheduleId int, loggingTime *logging_time.AddLoggingTime) error {
	panic("implement me")
}

func (c *client) DetailLoggingTime(scheduleId int, loggingTimeId int) error {
	panic("implement me")
}

func (c *client) EditLoggingTime(scheduleId int, loggingTimeId int, loggingTime *logging_time.EditLoggingTime) {
	panic("implement me")
}

func (c *client) DeleteLoggingTime(scheduleId int, loggingTimeId int) error {
	panic("implement me")
}

func (c *client) SubmitForApproveSchedule(scheduleId int, loggingTimeId int, status *schedule.EditStatusSchedule) error {
	panic("implement me")
}

func (c *client) ApproveSchedule(scheduleId int, loggingTimeId int, status *schedule.EditStatusSchedule) error {
	panic("implement me")
}
