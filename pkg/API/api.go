package api

import (
	"fmt"
	"net/http"
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
	LoggingTimeList(scheduleId int, options *Options) error
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

func NewClient(email string, password string) (*client, error) {
	token, err := auth.Authenticate(email, password)
	if err != nil {
		return nil, fmt.Errorf("unable to get authentification tokens: %s", err)
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

func (c client) Schedules(options *Options) ([]schedule.Schedule, error) {
	panic("implement me")
}

func (c client) AddSchedule(periodId int) error {
	panic("implement me")
}

func (c client) DetailSchedule(scheduleId int) error {
	panic("implement me")
}

func (c client) LoggingTimeList(scheduleId int, options *Options) error {
	panic("implement me")
}

func (c client) AddLoggingTime(scheduleId int, loggingTime *logging_time.AddLoggingTime) error {
	panic("implement me")
}

func (c client) DetailLoggingTime(scheduleId int, loggingTimeId int) error {
	panic("implement me")
}

func (c client) EditLoggingTime(scheduleId int, loggingTimeId int, loggingTime *logging_time.EditLoggingTime) {
	panic("implement me")
}

func (c client) DeleteLoggingTime(scheduleId int, loggingTimeId int) error {
	panic("implement me")
}

func (c client) SubmitForApproveSchedule(scheduleId int, loggingTimeId int, status *schedule.EditStatusSchedule) error {
	panic("implement me")
}

func (c client) ApproveSchedule(scheduleId int, loggingTimeId int, status *schedule.EditStatusSchedule) error {
	panic("implement me")
}
