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

type API interface {
	Schedules() ([]schedule.Schedule, error)
	AddSchedule(PeriodId int) error
	AddLoggingTime(loggingTime *logging_time.AddLoggingTime) error
	EditLoggingTime(loggingTime *logging_time.EditLoggingTime)
	DeleteLoggingTime(id int) error
	SubmitForApproveSchedule(status *schedule.EditStatusSchedule) error
	ApproveSchedule(status *schedule.EditStatusSchedule) error
}

type client struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
	HTTPClient   *http.Client
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
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}, nil
}

func (s *client) Schedules() ([]schedule.Schedule, error) {
	panic("implement me")
}

func (c *client) AddSchedule(PeriodId int) error {
	panic("implement me")
}

func (c *client) AddLoggingTime(loggingTime *logging_time.AddLoggingTime) error {
	panic("implement me")
}

func (c *client) EditLoggingTime(loggingTime *logging_time.EditLoggingTime) {
	panic("implement me")
}

func (c *client) DeleteLoggingTime(id int) error {
	panic("implement me")
}


func (c *client) SubmitForApproveSchedule(status *schedule.EditStatusSchedule) error {
	panic("implement me")
}

func (c *client) ApproveSchedule(status *schedule.EditStatusSchedule) error {
	panic("implement me")
}
