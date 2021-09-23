package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"suft_sdk/internal/auth"
	logging_time "suft_sdk/internal/logging-time"
	"suft_sdk/internal/schedule"
	"time"
)

const (
	BaseURL      string = "https://dev.gnivc.ru/tools/suft/api/v1/"
	SchedulesURN string = "api/v1/schedules"

	LoggingTimeURN string = "logging-times"

	Creator  Role = "creator"
	Approver Role = "approver"
)

type Role string

type Options struct {
	Page            int  `json:"page"`
	Size            int  `json:"size"`
	CreatorApprover Role `json:"creator_approver"`
}

type ScheduleId int
type LoggingTimeId int
type PeriodId int

type API interface {
	Schedules(options *Options) ([]schedule.Schedule, error)
	AddSchedule(periodId PeriodId) (*schedule.Schedule, error)
	DetailSchedule(scheduleId ScheduleId) (*schedule.Schedule, error)
	LoggingTimeList(scheduleId ScheduleId, options *Options) ([]logging_time.LoggingTime, error)
	AddLoggingTime(scheduleId ScheduleId, loggingTime *logging_time.AddLoggingTime) (*logging_time.LoggingTime, error)
	DetailLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) (*logging_time.LoggingTime, error)
	EditLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, loggingTime *logging_time.EditLoggingTime)
	DeleteLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) error
	SubmitForApproveSchedule(scheduleId ScheduleId, loggingTimeId LoggingTimeId, status *schedule.EditStatusSchedule) error
	ApproveSchedule(scheduleId ScheduleId, loggingTimeId LoggingTimeId, status *schedule.EditStatusSchedule) error
}

type Client struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
	HttpClient   *http.Client
}

func NewClient(email string, password string) (API, error) {
	token, err := auth.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:      BaseURL,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		HttpClient: &http.Client{
			Timeout: time.Minute,
		},
	}, nil
}

func (c *Client) Schedules(options *Options) ([]schedule.Schedule, error) {
	// если хочешь, можем досконально изучить пакет url и сделать элегантно,
	// однако на данном этапе предлагаю сделать наиболее доступным, простым и рабочим способом,
	// как я сделал ниже, а потом, если останется время, то переделаем с использованием пакета url.
	// reqURL, err := url.Parse(SchedulesURN)
	// if err != nil {
	// 	return nil, err
	// }
	// queryString := reqURL.Query()
	// queryString.Set("creatorApprover", "creator")
	// reqURL.RawQuery = queryString.Encode()
	// base, err := url.Parse(BaseURL)
	// if err != nil {
	// 	return nil, err
	// }
	// reqURL = base.ResolveReference(reqURL)
	page := 1
	size := 5
	creatorApprover := Creator
	if options != nil {
		page = options.Page
		size = options.Size
		creatorApprover = options.CreatorApprover
	}
	URN := fmt.Sprint(SchedulesURN, "?page=", page, "&size=", size, "&creatorApprover=", creatorApprover)
	resp, err := c.doHTTP(http.MethodGet, URN, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Schedules: unable to get schedules")
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Schedules: unable to read response body:", err)
		return nil, err
	}

	schedules := make([]schedule.Schedule, 1)
	err = json.Unmarshal(respB, &schedules)
	if err != nil {
		log.Println("Schedules: unable to unmarshal response body:", err)
		return nil, err
	}

	return schedules, nil
}

func (c *Client) AddSchedule(periodId PeriodId) (*schedule.Schedule, error) {

	reqS := fmt.Sprintf(
		`{
   "periodId": %d
}`, periodId)
	reqB := []byte(reqS)
	resp, err := c.doHTTP(http.MethodPost, SchedulesURN, reqB)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("AddSchedule: unable to read response body:", err)
		return nil, err
	}
	schedule := schedule.Schedule{}
	err = json.Unmarshal(respB, &schedule)
	if err != nil {
		log.Println("AddSchedule: unable to unmarshal response body:", err)
		return nil, err
	}
	return &schedule, nil
}

func (c *Client) DetailSchedule(scheduleId ScheduleId) (*schedule.Schedule, error) {

	URN := fmt.Sprintf("%s/%d", SchedulesURN, scheduleId)
	resp, err := c.doHTTP(http.MethodGet, URN, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("DetailSchedule: unable to get schedule")
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("DetailSchedule: unable to read response body:", err)
		return nil, err
	}
	schedule := schedule.Schedule{}
	err = json.Unmarshal(respB, &schedule)
	if err != nil {
		log.Println("DetailSchedule: unable to unmarshal response body:", err)
		return nil, err
	}

	return &schedule, nil
}

func (c *Client) LoggingTimeList(scheduleId ScheduleId, options *Options) ([]logging_time.LoggingTime, error) {
	page := 1
	size := 5
	if options != nil {
		page = options.Page
		size = options.Size
	}

	URN := fmt.Sprintf("%s/%d/%s?page=%d&size=%d", SchedulesURN, scheduleId, LoggingTimeURN, page, size)

	resp, err := c.doHTTP(http.MethodGet, URN, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("LoggingTimeList: unable to read response body:", err)
		return nil, err
	}

	loggingTimes := []logging_time.LoggingTime{}
	err = json.Unmarshal(respB, &loggingTimes)
	if err != nil {
		log.Println("LoggingTimeList: unable to unmarshal response body:", err)
		return nil, err
	}

	return loggingTimes, nil

}

func (c *Client) AddLoggingTime(scheduleId ScheduleId, loggingTime *logging_time.AddLoggingTime) (*logging_time.LoggingTime, error) {

	reqB, err := json.Marshal(loggingTime)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n%s\n\n", string(reqB))

	URN := fmt.Sprintf("%s/%d/%s", SchedulesURN, scheduleId, LoggingTimeURN)
	resp, err := c.doHTTP(http.MethodPost, URN, reqB)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("AddLoggingTime: unable to read response body:", err)
		return nil, err
	}
	fmt.Println(string(respB))
	loggingTimeResp := logging_time.LoggingTime{}
	err = json.Unmarshal(respB, &loggingTimeResp)
	if err != nil {
		log.Println("AddLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	return &loggingTimeResp, nil

}

func (c *Client) DetailLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) (*logging_time.LoggingTime, error) {
	URN := fmt.Sprintf("%s/%d/%s/%d", SchedulesURN, scheduleId, LoggingTimeURN, loggingTimeId)
	resp, err := c.doHTTP(http.MethodGet, URN, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("DetailLoggingTime: unable to read response body:", err)
		return nil, err
	}

	loggingTime := logging_time.LoggingTime{}
	err = json.Unmarshal(respB, &loggingTime)
	if err != nil {
		log.Println("DetailLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	return &loggingTime, nil
}

func (c *Client) EditLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, loggingTime *logging_time.EditLoggingTime) {
	panic("implement me")
}

func (c *Client) DeleteLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) error {
	panic("implement me")
}

func (c *Client) SubmitForApproveSchedule(scheduleId ScheduleId, loggingTimeId LoggingTimeId, status *schedule.EditStatusSchedule) error {
	panic("implement me")
}

func (c *Client) ApproveSchedule(scheduleId ScheduleId, loggingTimeId LoggingTimeId, status *schedule.EditStatusSchedule) error {
	panic("implement me")
}

func (c *Client) doHTTP(httpMethod string, URN string, body []byte) (*http.Response, error) {
	reqBody := bytes.NewBuffer(body)
	req, err := http.NewRequest(
		httpMethod,
		fmt.Sprint(BaseURL, URN),
		reqBody,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Auth-method", "Password")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Charset", "UTF-8")

	cookieAccessToken := &http.Cookie{
		Name:  "Access-token",
		Value: c.AccessToken,
	}
	req.AddCookie(cookieAccessToken)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}
