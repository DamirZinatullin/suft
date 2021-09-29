package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeSchedule1 = Schedule{
	Author: Employee{
		Email:      "test@gmail.com",
		FirstName:  "Ivan",
		Id:         5,
		LastName:   "Ivanovich",
		MiddleName: "Ivanov",
	},
	Id: 0,
	Period: Period{
		CloseDate:  "",
		EndDate:    "",
		Id:         5,
		StartDate:  "",
		WeekNumber: 2,
	},
	StatusCode: "22",
}

var fakeSchedule2 = Schedule{
	Author: Employee{
		Email:      "test@gmail.com",
		FirstName:  "Petrov",
		Id:         7,
		LastName:   "Petr",
		MiddleName: "Petrovich",
	},
	Id: 0,
	Period: Period{
		CloseDate:  "",
		EndDate:    "",
		Id:         5,
		StartDate:  "",
		WeekNumber: 2,
	},
	StatusCode: "25",
}

var fakeLoggingTime1 = LoggingTime{
	scheduleId:           777,
	AdminEmployee:        Employee{},
	CommentAdminEmployee: "fake comment from Admin",
	CommentEmployee:      "fake comment from Employee",
	Day1Time:             1,
	Day2Time:             2,
	Day3Time:             1,
	Day4Time:             0,
	Day5Time:             1,
	Day6Time:             0,
	Day7Time:             0,
	Id:                   0,
	ImportedFrom:         "fake1",
	ProjectId:            0,
	StatusCode:           "fake1",
	Task:                 "fake1",
	WorkKindId:           0,
}

var fakeLoggingTime2 = LoggingTime{
	scheduleId:           777,
	AdminEmployee:        Employee{},
	CommentAdminEmployee: "fake comment from Admin2",
	CommentEmployee:      "fake comment from Employee2",
	Day1Time:             2,
	Day2Time:             1,
	Day3Time:             2,
	Day4Time:             0,
	Day5Time:             1,
	Day6Time:             0,
	Day7Time:             0,
	Id:                   0,
	ImportedFrom:         "fake2",
	ProjectId:            0,
	StatusCode:           "fake2",
	Task:                 "fake2",
	WorkKindId:           0,
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
	schedules, err := client.Schedules(&OptionsS{
		Page:            7,
		Size:            7,
		CreatorApprover: "fake",
	})
	require.NoError(t, err)
	assert.Equal(t, []*Schedule{&fakeSchedule1, &fakeSchedule2}, schedules)
}

func TestSchedulesUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	schedules, err := client.Schedules(nil)
	assert.Error(t, err)
	assert.Nil(t, schedules)
}

func TestSchedulesError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	schedules, err := client.Schedules(nil)
	require.Error(t, err)
	assert.Nil(t, schedules)
}

func TestAddScheduleSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeSchedule1
	expected.client = client
	GetRequireResp = SuccessRespAddSchedule
	schedule, err := client.AddSchedule(5)
	require.NoError(t, err)
	assert.Equal(t, expected, schedule)
}

func TestAddScheduleUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	schedule, err := client.AddSchedule(5)
	assert.Error(t, err)
	assert.Nil(t, schedule)
}

func TestAddScheduleError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	scheduleResp, err := client.AddSchedule(5)
	require.Error(t, err)
	assert.Nil(t, scheduleResp)
}

func TestDetailScheduleSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeSchedule1
	expected.client = client
	GetRequireResp = SuccessRespDetailSchedule
	schedule, err := client.DetailSchedule(777)
	require.NoError(t, err)
	assert.Equal(t, expected, schedule)
}

func TestDetailScheduleUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	schedule, err := client.DetailSchedule(777)
	assert.Error(t, err)
	assert.Nil(t, schedule)
}

func TestDetailScheduleError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	scheduleResp, err := client.DetailSchedule(777)
	require.Error(t, err)
	assert.Nil(t, scheduleResp)
}

func TestSubmitForApproveScheduleSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeSchedule1
	expected.client = client
	GetRequireResp = SuccessRespDetailSchedule
	scheduleResp, err := client.SubmitForApproveSchedule(777)
	require.NoError(t, err)
	assert.Equal(t, expected, scheduleResp)
}

func TestSubmitForApproveScheduleUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	scheduleResp, err := client.SubmitForApproveSchedule(777)
	assert.Error(t, err)
	assert.Nil(t, scheduleResp)
}

func TestSubmitForApproveScheduleError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	scheduleResp, err := client.SubmitForApproveSchedule(777)
	require.Error(t, err)
	assert.Nil(t, scheduleResp)
}

func TestLoggingTimeListSuccess(t *testing.T) {
	client, err := NewFakeClient()
	GetRequireResp = SuccessRespLoggingTimeList
	if err != nil {
		log.Fatalln(err)
	}
	loggingTimeList, err := client.LoggingTimeList(777, &OptionsLT{
		Page: 5,
		Size: 5,
	})
	require.NoError(t, err)
	assert.Equal(t, []*LoggingTime{&fakeLoggingTime1, &fakeLoggingTime2}, loggingTimeList)
}

func TestLoggingTimeListUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	loggingTimeList, err := client.LoggingTimeList(777, nil)
	assert.Error(t, err)
	assert.Nil(t, loggingTimeList)
}

func TestLoggingTimeListError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	loggingTimeList, err := client.LoggingTimeList(777, nil)
	require.Error(t, err)
	assert.Nil(t, loggingTimeList)
}

func TestAddLoggingTimeSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeLoggingTime1
	expected.client = client
	GetRequireResp = SuccessRespAddLoggingTime
	loggingTimeResp, err := client.AddLoggingTime(777, &AddLoggingTime{})
	require.NoError(t, err)
	assert.Equal(t, expected, loggingTimeResp)
}

func TestAddLoggingTimeUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	loggingTimeResp, err := client.AddLoggingTime(5, &AddLoggingTime{})
	assert.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}

func TestAddLoggingTimeError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	loggingTimeResp, err := client.AddLoggingTime(5, &AddLoggingTime{})
	require.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}

func TestDetailLoggingTimeSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeLoggingTime1
	expected.client = client
	GetRequireResp = SuccessRespDetailLoggingTime
	loggingTimeResp, err := client.DetailLoggingTime(777, 777)
	require.NoError(t, err)
	assert.Equal(t, &fakeLoggingTime1, loggingTimeResp)
}

func TestDetailLoggingTimeUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	loggingTimeResp, err := client.DetailLoggingTime(5, 777)
	assert.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}

func TestDetailLoggingTimeError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	loggingTimeResp, err := client.DetailLoggingTime(5, 777)
	require.Error(t, err)
	assert.Nil(t, loggingTimeResp)
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
	fakeSchedule1.client, _ = NewFakeClient()
	fakeSchedule2.client, _ = NewFakeClient()
	schedules := []*Schedule{&fakeSchedule1, &fakeSchedule2}
	respB, _ := json.Marshal(schedules)
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 200,
		Body: body}
	return &resp, nil
}

func UnauthorizedResp() (*http.Response, error) {
	respB := []byte("Unauthorized")
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: http.StatusUnauthorized,
		Body: body}
	return &resp, nil
}

func ErrorRespFromDoHttp() (*http.Response, error) {
	return nil, errors.New("error from doHTTP")
}

func SuccessRespAddSchedule() (*http.Response, error) {
	scheduleReq := fakeSchedule1
	// scheduleReq.client, _ = NewFakeClient()
	respB, _ := json.Marshal(scheduleReq)
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 201,
		Body: body}
	return &resp, nil
}

func SuccessRespDetailSchedule() (*http.Response, error) {
	scheduleReq := fakeSchedule1
	// scheduleReq.client, _ = NewFakeClient()
	respB, _ := json.Marshal(scheduleReq)
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 200,
		Body: body}
	return &resp, nil
}

func SuccessRespLoggingTimeList() (*http.Response, error) {
	fakeLoggingTime1.client, _ = NewFakeClient()
	fakeLoggingTime2.client, _ = NewFakeClient()
	loggingTimes := []*LoggingTime{&fakeLoggingTime1, &fakeLoggingTime2}
	respB, _ := json.Marshal(loggingTimes)
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 200,
		Body: body}
	return &resp, nil
}

func SuccessRespAddLoggingTime() (*http.Response, error) {
	loggingTime := fakeLoggingTime1
	// loggingTime.client, _ = NewFakeClient()
	respB, _ := json.Marshal(loggingTime)
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 201,
		Body: body}
	return &resp, nil
}

func SuccessRespDeleteLoggingTime() (*http.Response, error) {
	respB := []byte("OK")
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 200,
		Body: body}
	return &resp, nil
}

func SuccessRespDetailLoggingTime() (*http.Response, error) {
	loggingTime := fakeLoggingTime1
	// loggingTime.client, _ = NewFakeClient()
	respB, _ := json.Marshal(loggingTime)
	body := ioutil.NopCloser(bytes.NewReader(respB))
	resp := http.Response{StatusCode: 200,
		Body: body}
	return &resp, nil
}

// func TestEditLoggingTimeSuccess(t *testing.T) {
// 	client, err := NewFakeClient()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	GetRequireResp = SuccessRespDetailLoggingTime
// 	loggingTimeResp, err := client.EditLoggingTime(777, 777, &loggingtime.EditLoggingTime{})
// 	require.NoError(t, err)
// 	assert.Equal(t, &fakeLoggingTime1, loggingTimeResp)
// }

// func TestEditLoggingTimeUnauthorized(t *testing.T) {
// 	client, err := NewFakeClient()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	GetRequireResp = UnauthorizedResp
// 	loggingTimeResp, err := client.EditLoggingTime(777, 777, &loggingtime.EditLoggingTime{})
// 	assert.Error(t, err)
// 	assert.Nil(t, loggingTimeResp)
// }

// func TestEditLoggingTimeError(t *testing.T) {
// 	client, err := NewFakeClient()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	GetRequireResp = ErrorRespFromDoHttp
// 	loggingTimeResp, err := client.EditLoggingTime(777, 777, &loggingtime.EditLoggingTime{})
// 	require.Error(t, err)
// 	assert.Nil(t, loggingTimeResp)
// }

func TestDeleteLoggingTimeSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = SuccessRespDeleteLoggingTime
	err = client.DeleteLoggingTime(777, 777)
	require.NoError(t, err)

}

func TestDeleteLoggingTimeUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	err = client.DeleteLoggingTime(777, 777)
	assert.Error(t, err)
}

func TestDeleteLoggingTimeError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	err = client.DeleteLoggingTime(777, 777)
	require.Error(t, err)
}

func TestApproveLoggingTimeSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeLoggingTime1
	expected.client = client
	GetRequireResp = SuccessRespDetailLoggingTime
	loggingTimeResp, err := client.ApproveLoggingTime(777, 777, "ok")
	require.NoError(t, err)
	assert.Equal(t, expected, loggingTimeResp)
}

func TestApproveLoggingTimeUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	loggingTimeResp, err := client.ApproveLoggingTime(5, 777, "ok")
	assert.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}

func TestApproveLoggingTimeError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	loggingTimeResp, err := client.ApproveLoggingTime(5, 777, "ok")
	require.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}

func TestDeclineLoggingTimeSuccess(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	expected := &fakeLoggingTime1
	expected.client = client
	GetRequireResp = SuccessRespDetailLoggingTime
	loggingTimeResp, err := client.DeclineLoggingTime(777, 777, "fail")
	require.NoError(t, err)
	assert.Equal(t, expected, loggingTimeResp)
}

func TestDeclineLoggingTimeUnauthorized(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = UnauthorizedResp
	loggingTimeResp, err := client.DeclineLoggingTime(5, 777, "fail")
	assert.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}

func TestDeclineLoggingTimeError(t *testing.T) {
	client, err := NewFakeClient()
	if err != nil {
		log.Fatalln(err)
	}
	GetRequireResp = ErrorRespFromDoHttp
	loggingTimeResp, err := client.DeclineLoggingTime(5, 777, "fail")
	require.Error(t, err)
	assert.Nil(t, loggingTimeResp)
}
