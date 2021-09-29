package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"suftsdk/pkg/api"
	"testing"
)

var fakeSchedule1 = api.Schedule{
	Author: api.Employee{
		Email:      "test@gmail.com",
		FirstName:  "Ivan",
		Id:         5,
		LastName:   "Ivanovich",
		MiddleName: "Ivanov",
	},
	Id: 0,
	Period: api.Period{
		CloseDate:  "",
		EndDate:    "",
		Id:         5,
		StartDate:  "",
		WeekNumber: 2,
	},
	StatusCode: "22",
}

var fakeSchedule2 = api.Schedule{
	Author: api.Employee{
		Email:      "test@gmail.com",
		FirstName:  "Petrov",
		Id:         7,
		LastName:   "Petr",
		MiddleName: "Petrovich",
	},
	Id: 0,
	Period: api.Period{
		CloseDate:  "",
		EndDate:    "",
		Id:         5,
		StartDate:  "",
		WeekNumber: 2,
	},
	StatusCode: "25",
}

var fakeLoggingTime1 = api.LoggingTime{
	AdminEmployee:        api.Employee{},
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

var fakeLoggingTime2 = api.LoggingTime{
	AdminEmployee:        api.Employee{},
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

type schedulesFunc func() ([]*api.Schedule, error)
type scheduleDetailFunc func() (*api.Schedule, error)

var respSchedules schedulesFunc
var respScheduleDetail scheduleDetailFunc

var exitIndicator string

func TestCliFunc(t *testing.T) {
	clientConstructor = fakeClientInit{}
	app, err := cliFunc()
	require.NoError(t, err)
	app.ExitErrHandler = func(context *cli.Context, err error) {
		exitIndicator = "1"
	}
	t.Run("Успешный вызов Schedules", func(t *testing.T){
		args := []string{"", "scs"}
		respSchedules = SuccessRespSchedules
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
	})
	t.Run("Вызов Schedules с лишним аргументом", func(t *testing.T){
		args := []string{"", "scs", "fake"}
		respSchedules = SuccessRespSchedules
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
	})
	t.Run("Ошибка при вызове метода Schedules", func(t *testing.T){
		args := []string{"", "scs"}
		respSchedules = ErrorRespSchedules
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов ScheduleDetail", func(t *testing.T){
		args := []string{"", "sc", "-scid" ,"777"}
		respScheduleDetail = SuccessRespScheduleDetail
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове метода ScheduleDetail", func(t *testing.T){
		args := []string{"", "sc", "-scid", "777"}
		respScheduleDetail = ErrorRespDetailSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача невалидного флага в ScheduleDetail", func(t *testing.T){
		args := []string{"", "sc", "-scid" ,"ар"}
		respScheduleDetail = SuccessRespScheduleDetail
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача неправильного флага в ScheduleDetail", func(t *testing.T){
		args := []string{"", "sc", "-fake" ,"777"}
		respScheduleDetail = SuccessRespScheduleDetail
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
}

type fakeClientInit struct{}

func (f fakeClientInit) NewClient() (client api.API, err error) {
	return &fakeClient{}, nil
}

type fakeClient struct {
}

func (f *fakeClient) Schedules(options *api.OptionsS) ([]*api.Schedule, error) {
	return respSchedules()
}

func (f *fakeClient) AddSchedule(periodId api.PeriodId) (*api.Schedule, error) {
	panic("implement me")
}

func (f *fakeClient) DetailSchedule(scheduleId api.ScheduleId) (*api.Schedule, error) {
	return respScheduleDetail()
}

func (f *fakeClient) LoggingTimeList(scheduleId api.ScheduleId, options *api.OptionsLT) ([]*api.LoggingTime, error) {
	panic("implement me")
}

func (f *fakeClient) AddLoggingTime(scheduleId api.ScheduleId, loggingTime *api.AddLoggingTime) (*api.LoggingTime, error) {
	panic("implement me")
}

func (f *fakeClient) DetailLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId) (*api.LoggingTime, error) {
	panic("implement me")
}

func (f *fakeClient) DeleteLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId) error {
	panic("implement me")
}

func (f *fakeClient) SubmitForApproveSchedule(scheduleId api.ScheduleId) (*api.Schedule, error) {
	panic("implement me")
}

func (f *fakeClient) ApproveLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId, comment string) (*api.LoggingTime, error) {
	panic("implement me")
}

func (f *fakeClient) DeclineLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId, comment string) (*api.LoggingTime, error) {
	panic("implement me")
}

func SuccessRespSchedules() ([]*api.Schedule, error) {
	return []*api.Schedule{&fakeSchedule1, &fakeSchedule2}, nil
}

func ErrorRespSchedules() ([]*api.Schedule, error) {
	return nil, errors.New("error from schedules method")
}

func SuccessRespScheduleDetail() (*api.Schedule, error) {
	return &fakeSchedule1, nil
}

func ErrorRespDetailSchedule() (*api.Schedule, error) {
	return nil, errors.New("error from scheduleDetail method")
}
