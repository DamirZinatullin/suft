package http_client

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"suft_sdk/internal/http-client/schedule"
	"time"
)

const (
	URI string = "https://dev.gnivc.ru/tools/suft/api/v1/"
)

type HttpClient interface {
	AuthTokens() error
	Authorize() error
	Schedules() ([]schedule.Schedule, error)
	AddSchedule([]schedule.Schedule) error
	DetailSchedule(int) error
	EditSchedule(int, map[string]string) error
	UpdateSchedule(int, map[string]string) error
	Logout() error
	GetURN() string
	GetMethod() string
	GetAccessToken() string
	GetRefreshToken() string
}

type httpClient struct {
	URN          string
	Method       string
	AccessToken  string
	RefreshToken string
}

func NewHttpClient() HttpClient {
	return &httpClient{}
}

func (s *httpClient) GetURN() string {
	return s.URN
}

func (s *httpClient) GetMethod() string {
	return s.Method
}

func (s *httpClient) GetAccessToken() string {
	return s.AccessToken
}

func (s *httpClient) GetRefreshToken() string {
	return s.RefreshToken
}

func (s *httpClient) AuthTokens() error {
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

func (s *httpClient) Authorize() error                            { return nil }
func (s *httpClient) Schedules() ([]schedule.Schedule, error)     { return nil, nil }
func (s *httpClient) AddSchedule([]schedule.Schedule) error       { return nil }
func (s *httpClient) DetailSchedule(int) error                    { return nil }
func (s *httpClient) EditSchedule(int, map[string]string) error   { return nil }
func (s *httpClient) UpdateSchedule(int, map[string]string) error { return nil }
func (s *httpClient) Logout() error                               { return nil }
