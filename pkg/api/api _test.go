package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"suft_sdk/internal/schedule"
	"testing"
)

var fakeSchedule1 = schedule.Schedule{
	Author: schedule.Employee{
		Email:      "test@gmail.com",
		FirstName:  "Ivan",
		Id:         5,
		LastName:   "Ivanovich",
		MiddleName: "Ivanov",
	},
	Id: 0,
	Period: schedule.Period{
		CloseDate:  "",
		EndDate:    "",
		Id:         5,
		StartDate:  "",
		WeekNumber: 2,
	},
	StatusCode: "22",
}

var fakeSchedule2 = schedule.Schedule{
	Author: schedule.Employee{
		Email:      "test@gmail.com",
		FirstName:  "Petrov",
		Id:         7,
		LastName:   "Petr",
		MiddleName: "Petrovich",
	},
	Id: 0,
	Period: schedule.Period{
		CloseDate:  "",
		EndDate:    "",
		Id:         5,
		StartDate:  "",
		WeekNumber: 2,
	},
	StatusCode: "25",
}

var GetRequireResp func() (*http.Response, error)

type mockedHttpClient struct{}

func (m *mockedHttpClient) Do(req *http.Request) (*http.Response, error) {
	return GetRequireResp()
}

func TestSchedulesSuccess(t *testing.T) {
	client, err := NewFakeClient()
	GetRequireResp = SuccessRespSchedules
	if err != nil {
		log.Fatalln(err)
	}
	schedules, err := client.Schedules(nil)
	require.NoError(t, err)
	assert.Equal(t, []schedule.Schedule{fakeSchedule1, fakeSchedule2}, schedules)
}

func TestSchedulesFail(t *testing.T) {
	client, err := NewFakeClient()
	GetRequireResp = FailRespSchedules
	if err != nil {
		log.Fatalln(err)
	}
	schedules, err := client.Schedules(nil)
	assert.Error(t, err)
	assert.Equal(t, []schedule.Schedule(nil), schedules)
}

func NewFakeClient() (*Client, error) {
	httpClient = new(mockedHttpClient)
	return &Client{
		BaseURL:      BaseURL,
		AccessToken:  "fake_access_token",
		RefreshToken: "fake_refresh_token",
		HttpClient:   httpClient,
	}, nil
}

func SuccessRespSchedules() (*http.Response, error) {
	schedules := []schedule.Schedule{fakeSchedule1, fakeSchedule2}
	respB, _ := json.Marshal(schedules)
	body := ioutil.NopCloser(bytes.NewReader([]byte(respB)))
	resp := http.Response{StatusCode: 200,
		Body: body}
	return &resp, nil
}




func FailRespSchedules() (*http.Response, error) {
	respB := []byte("Unauthorized")
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: http.StatusUnauthorized,
		Body: body}
	return &resp, nil
}