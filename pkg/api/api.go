package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"suftsdk/internal/auth"
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

type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

var httpClient HttpClient

// опции для функции NewClient
type OptionsNC struct {
	SuftAPIURL  string
	HttpTimeout time.Duration
}

// опции для метода Schedules
type OptionsS struct {
	Page            int  `json:"page"`
	Size            int  `json:"size"`
	CreatorApprover Role `json:"creator_approver"`
}

// опции для метода LoggingTimeList
type OptionsLT struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ScheduleId int
type LoggingTimeId int
type PeriodId int

type API interface {
	Schedules(options *OptionsS) ([]*Schedule, error)
	AddSchedule(periodId PeriodId) (*Schedule, error)
	DetailSchedule(scheduleId ScheduleId) (*Schedule, error)
	LoggingTimeList(scheduleId ScheduleId, options *OptionsLT) ([]LoggingTime, error)
	AddLoggingTime(scheduleId ScheduleId, loggingTime *AddLoggingTime) (*LoggingTime, error)
	DetailLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) (*LoggingTime, error)
	DeleteLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) error
	SubmitForApproveSchedule(scheduleId ScheduleId) (*Schedule, error)
	ApproveLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, comment string) (*LoggingTime, error)
}

type Client struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
	HttpClient   HttpClient
}

func NewClient(email string, password string, options *OptionsNC) (API, error) {
	baseURL := BaseURL
	httpTimeout := time.Minute
	if options != nil {
		if options.SuftAPIURL != "" {
			baseURL = options.SuftAPIURL
		}
		if options.HttpTimeout != 0 {
			httpTimeout = options.HttpTimeout
		}
	}

	token, err := auth.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	httpClient = &http.Client{
		Timeout: httpTimeout,
	}
	return &Client{
		BaseURL:      baseURL,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		HttpClient:   httpClient,
	}, nil
}

func (c *Client) Schedules(options *OptionsS) ([]*Schedule, error) {
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
		if options.Page != 0 {
			page = options.Page
		}
		if options.Size != 0 {
			size = options.Size
		}
		if options.CreatorApprover != "" {
			creatorApprover = options.CreatorApprover
		}
	}
	URN := fmt.Sprint(SchedulesURN, "?page=", page, "&size=", size, "&creatorApprover=", creatorApprover)
	resp, err := c.doHTTP(http.MethodGet, URN, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Schedules: unable to read response body:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}
	schedulesResp := make([]*Schedule, 1)
	err = json.Unmarshal(respB, &schedulesResp)
	if err != nil {
		log.Println("Schedules: unable to unmarshal response body:", err)
		return nil, err
	}

	for _, item := range schedulesResp {
		item.client = c
	}
	return schedulesResp, nil
}

func (c *Client) AddSchedule(periodId PeriodId) (*Schedule, error) {
	peiodIdStruct := struct {
		PeriodId PeriodId `json:"periodId"`
	}{PeriodId: periodId}
	reqB, err := json.Marshal(peiodIdStruct)
	if err != nil {
		return nil, err
	}
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
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(respB))
	}

	schedule := Schedule{}
	err = json.Unmarshal(respB, &schedule)
	if err != nil {
		log.Println("AddSchedule: unable to unmarshal response body:", err)
		return nil, err
	}
	schedule.client = c
	return &schedule, nil
}

func (c *Client) DetailSchedule(scheduleId ScheduleId) (*Schedule, error) {

	URN := fmt.Sprintf("%s/%d", SchedulesURN, scheduleId)
	resp, err := c.doHTTP(http.MethodGet, URN, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("DetailSchedule: unable to read response body:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}
	schedule := Schedule{}
	err = json.Unmarshal(respB, &schedule)
	if err != nil {
		log.Println("DetailSchedule: unable to unmarshal response body:", err)
		return nil, err
	}
	schedule.client = c

	return &schedule, nil
}

func (c *Client) LoggingTimeList(scheduleId ScheduleId, options *OptionsLT) ([]LoggingTime, error) {
	page := 1
	size := 5
	if options != nil {
		if options.Page != 0 {
			page = options.Page
		}
		if options.Size != 0 {
			size = options.Size
		}
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
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}
	loggingTimes := []LoggingTime{}
	err = json.Unmarshal(respB, &loggingTimes)
	if err != nil {
		log.Println("LoggingTimeList: unable to unmarshal response body:", err)
		return nil, err
	}

	return loggingTimes, nil

}

func (c *Client) AddLoggingTime(scheduleId ScheduleId, loggingTime *AddLoggingTime) (*LoggingTime, error) {

	reqB, err := json.Marshal(loggingTime)
	if err != nil {
		return nil, err
	}

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
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(respB))
	}

	loggingTimeResp := LoggingTime{}
	err = json.Unmarshal(respB, &loggingTimeResp)
	if err != nil {
		log.Println("AddLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	return &loggingTimeResp, nil

}

func (c *Client) DetailLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) (*LoggingTime, error) {
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
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}
	loggingTime := LoggingTime{}
	err = json.Unmarshal(respB, &loggingTime)
	if err != nil {
		log.Println("DetailLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	return &loggingTime, nil
}

// func (c *Client) EditLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, loggingTime *loggingtime.EditLoggingTime) (*loggingtime.LoggingTime, error) {
// 	reqB, err := json.Marshal(loggingTime)
// 	if err != nil {
// 		log.Println("EditLoggingTime: unable to marshal request body:", err)
// 		return nil, err
// 	}
// 	URN := fmt.Sprintf("%s/%d/%s/%d", SchedulesURN, scheduleId, LoggingTimeURN, loggingTimeId)
// 	resp, err := c.doHTTP(http.MethodPatch, URN, reqB)
// 	if err != nil {
// 		log.Println("EditLoggingTime: doHTTP:", err)
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	respB, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("EditLoggingTime: unable to read response body:", err)
// 		return nil, err
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errors.New(string(respB))
// 	}
// 	loggingTimeResp := loggingtime.LoggingTime{}
// 	err = json.Unmarshal(respB, &loggingTimeResp)
// 	if err != nil {
// 		log.Println("EditLoggingTime: unable to unmarshal response body:", err)
// 		return nil, err
// 	}
// 	return &loggingTimeResp, nil
// }

func (c *Client) DeleteLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) error {
	URN := fmt.Sprintf("%s/%d/%s/%d", SchedulesURN, scheduleId, LoggingTimeURN, loggingTimeId)
	resp, err := c.doHTTP(http.MethodDelete, URN, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("DeleteLoggingTime: unable to read response body:", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(respB))
	}
	return nil
}

func (c *Client) SubmitForApproveSchedule(scheduleId ScheduleId) (*Schedule, error) {
	URN := fmt.Sprintf("%s/%d", SchedulesURN, scheduleId)

	statusCodeStruct := struct {
		statusCode StatusCode `json:"statusCode"`
	}{
		statusCode: ToApprove,
	}
	reqB, err := json.Marshal(statusCodeStruct)
	if err != nil {
		log.Println("SubmitForApproveSchedule: unable to marshal body:", err)
		return nil, err
	}

	resp, err := c.doHTTP(http.MethodPatch, URN, reqB)
	if err != nil {
		log.Println("SubmitForApproveSchedule: doHTTP:", err)
		return nil, err
	}
	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("SubmitForApproveSchedule: unable to read response body:", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}

	schedule := Schedule{}

	err = json.Unmarshal(respB, &schedule)
	if err != nil {
		log.Println("SubmitForApproveSchedule: unable to unmarshal response body:", err)
		return nil, err
	}
	schedule.client = c

	return &schedule, nil
}

func (c *Client) ApproveLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, comment string) (*LoggingTime, error) {
	editLoggingTime := struct {
		CommentAdminEmployee string     `json:"commentAdminEmployee"`
		CommentEmployee      string     `json:"commentEmployee"`
		Day1Time             float64    `json:"day1Time"`
		Day2Time             float64    `json:"day2Time"`
		Day3Time             float64    `json:"day3Time"`
		Day4Time             float64    `json:"day4Time"`
		Day5Time             float64    `json:"day5Time"`
		Day6Time             float64    `json:"day6Time"`
		Day7Time             float64    `json:"day7Time"`
		ProjectId            int        `json:"projectId"`
		StatusCode           StatusCode `json:"statusCode"`
		Task                 string     `json:"task"`
		WorkKindId           int                    `json:"workKindId"`
	}{
		CommentAdminEmployee: comment,
		StatusCode:           Approved,
	}

	reqB, err := json.Marshal(&editLoggingTime)
	if err != nil {
		log.Println("ApproveLoggingTime: unable to marshal body:", err)
		return nil, err
	}

	URN := fmt.Sprintf("%s/%d/%s/%d", SchedulesURN, scheduleId, LoggingTimeURN, loggingTimeId)

	resp, err := c.doHTTP(http.MethodPatch, URN, reqB)
	if err != nil {
		log.Println("ApproveLoggingTime: doHTTP:", err)
		return nil, err
	}

	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("ApproveLoggingTime: unable to read response body:", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}

	loggingTimeResp := LoggingTime{}

	err = json.Unmarshal(respB, &loggingTimeResp)
	if err != nil {
		log.Println("EditLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	return &loggingTimeResp, nil
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
