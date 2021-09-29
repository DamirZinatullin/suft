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
	LoggingTimeList(scheduleId ScheduleId, options *OptionsLT) ([]*LoggingTime, error)
	AddLoggingTime(scheduleId ScheduleId, loggingTime *AddLoggingTime) (*LoggingTime, error)
	DetailLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) (*LoggingTime, error)
	DeleteLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId) error
	SubmitForApproveSchedule(scheduleId ScheduleId) (*Schedule, error)
	ApproveLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, comment string) (*LoggingTime, error)
	DeclineLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, comment string) (*LoggingTime, error)
}

type Client struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
	HttpClient   HttpClient
}

func NewClient(email string, password string, options *OptionsNC) (API, error) {
	baseURL := BaseURL
	httpTimeout := 2 * time.Second
	if options != nil {
		if options.SuftAPIURL != "" {
			baseURL = options.SuftAPIURL
		}
		if options.HttpTimeout != 0 {
			httpTimeout = options.HttpTimeout
		}
	}

	authOptions := auth.Options{
		SuftAPIURL:  baseURL,
		HttpTimeout: httpTimeout,
	}

	token, err := auth.Authenticate(email, password, &authOptions)
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
	page := 0
	size := 5
	creatorApprover := Creator
	if options != nil {
		page = options.Page
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

func (c *Client) LoggingTimeList(scheduleId ScheduleId, options *OptionsLT) ([]*LoggingTime, error) {
	page := 0
	size := 5
	if options != nil {
		page = options.Page
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
	loggingTimes := make([]*LoggingTime, 1)
	err = json.Unmarshal(respB, &loggingTimes)
	if err != nil {
		log.Println("LoggingTimeList: unable to unmarshal response body:", err)
		return nil, err
	}
	for _, loggingTime := range loggingTimes {
		loggingTime.client = c
		loggingTime.scheduleId = scheduleId
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
	loggingTimeResp.client = c
	loggingTimeResp.scheduleId = scheduleId
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
	loggingTime.client = c
	loggingTime.scheduleId = scheduleId
	return &loggingTime, nil
}


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
		StatusCode StatusCode `json:"statusCode"`
	}{
		StatusCode: ToApprove,
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
	editLoggingTime := EditLoggingTime{
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
		log.Println("ApproveLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	loggingTimeResp.client = c
	loggingTimeResp.scheduleId = scheduleId
	return &loggingTimeResp, nil
}

func (c *Client) DeclineLoggingTime(scheduleId ScheduleId, loggingTimeId LoggingTimeId, comment string) (*LoggingTime, error) {
	editLoggingTime := EditLoggingTime{
		CommentAdminEmployee: comment,
		StatusCode:           Declined,
	}

	reqB, err := json.Marshal(&editLoggingTime)
	if err != nil {
		log.Println("DeclineLoggingTime: unable to marshal body:", err)
		return nil, err
	}

	URN := fmt.Sprintf("%s/%d/%s/%d", SchedulesURN, scheduleId, LoggingTimeURN, loggingTimeId)

	resp, err := c.doHTTP(http.MethodPatch, URN, reqB)
	if err != nil {
		log.Println("DeclineLoggingTime: doHTTP:", err)
		return nil, err
	}

	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("DeclineLoggingTime: unable to read response body:", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
	}

	loggingTimeResp := LoggingTime{}

	err = json.Unmarshal(respB, &loggingTimeResp)
	if err != nil {
		log.Println("DeclineLoggingTime: unable to unmarshal response body:", err)
		return nil, err
	}
	loggingTimeResp.client = c
	loggingTimeResp.scheduleId = scheduleId
	return &loggingTimeResp, nil
}

func (c *Client) doHTTP(httpMethod string, URN string, body []byte) (*http.Response, error) {
	reqBody := bytes.NewBuffer(body)
	req1, err := http.NewRequest(
		httpMethod,
		fmt.Sprint(BaseURL, URN),
		reqBody,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req1.Header.Add("Auth-method", "Password")
	req1.Header.Add("Content-Type", "application/json")
	req1.Header.Add("Accept-Charset", "UTF-8")

	cookieAccessToken := &http.Cookie{
		Name:  "Access-token",
		Value: c.AccessToken,
	}
	req1.AddCookie(cookieAccessToken)

	resp, err := c.HttpClient.Do(req1)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		tokens, err := auth.Refresh(c.RefreshToken, nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		c.AccessToken = tokens.AccessToken
		c.RefreshToken = tokens.RefreshToken

		cookieAccessToken.Value = tokens.AccessToken
		req2, err := http.NewRequest(
			httpMethod,
			fmt.Sprint(BaseURL, URN),
			reqBody,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		req2.Header.Add("Auth-method", "Password")
		req2.Header.Add("Content-Type", "application/json")
		req2.Header.Add("Accept-Charset", "UTF-8")
		req2.AddCookie(cookieAccessToken)
		resp, err = c.HttpClient.Do(req2)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return resp, nil
}
