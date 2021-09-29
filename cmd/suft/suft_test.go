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
type addScheduleFunc func() (*api.Schedule, error)
type loggingTimeListFunc func() ([]*api.LoggingTime, error)
type detailLoggingTimeFunc func() (*api.LoggingTime, error)
type deleteLoggingTimeFunc func() error
type submitForApproveScheduleFunc func() (*api.Schedule, error)
type approveLoggingTimeFunc func() (*api.LoggingTime, error)
type declineLoggingTimeFunc func() (*api.LoggingTime, error)
type addLoggingTimeFunc func() (*api.LoggingTime, error)

var respSchedules schedulesFunc
var respScheduleDetail scheduleDetailFunc
var respAddSchedule addScheduleFunc
var respLoggingTimeList loggingTimeListFunc
var respDetailLoggingTime detailLoggingTimeFunc
var respDeleteLoggingTime deleteLoggingTimeFunc
var respSubmitForApproveSchedule submitForApproveScheduleFunc
var respApproveLoggingTime approveLoggingTimeFunc
var respDeclineLoggingTime declineLoggingTimeFunc
var respAddLoggingTime addLoggingTimeFunc

var exitIndicator string

func TestCliFunc(t *testing.T) {
	clientConstructor = fakeClientInit{}
	app, err := cliFunc()
	require.NoError(t, err)
	app.ExitErrHandler = func(context *cli.Context, err error) {
		exitIndicator = "1"
	}
	t.Run("Успешный вызов Schedules", func(t *testing.T) {
		args := []string{"", "scs"}
		respSchedules = SuccessRespSchedules
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов Schedules с дополнительными аргументами", func(t *testing.T) {
		args := []string{"", "scs", "-p", "1", "-s", "2", "-r", "approver"}
		respSchedules = SuccessRespSchedules
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Вызов Schedules с лишним аргументом", func(t *testing.T) {
		args := []string{"", "scs", "fake"}
		respSchedules = SuccessRespSchedules
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове метода Schedules", func(t *testing.T) {
		args := []string{"", "scs"}
		respSchedules = ErrorRespSchedules
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов ScheduleDetail", func(t *testing.T) {
		args := []string{"", "sc", "-scid", "777"}
		respScheduleDetail = SuccessRespDetailSchedule
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове метода ScheduleDetail", func(t *testing.T) {
		args := []string{"", "sc", "-scid", "777"}
		respScheduleDetail = ErrorRespDetailSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача невалидного флага в ScheduleDetail", func(t *testing.T) {
		args := []string{"", "sc", "-scid", "ар"}
		respScheduleDetail = SuccessRespDetailSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача неправильного флага в ScheduleDetail", func(t *testing.T) {
		args := []string{"", "sc", "-fake", "777"}
		respScheduleDetail = SuccessRespDetailSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов addSchedule", func(t *testing.T) {
		args := []string{"", "as", "-pid", "777"}
		respAddSchedule = SuccessRespAddSchedule
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове метода addSchedule", func(t *testing.T) {
		args := []string{"", "as", "-pid", "777"}
		respAddSchedule = ErrorRespAddSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача невалидного флага в addSchedule", func(t *testing.T) {
		args := []string{"", "as", "-pid", "fake"}
		respAddSchedule = SuccessRespAddSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача неправильного флага в addSchedule", func(t *testing.T) {
		args := []string{"", "as", "-fake", "777"}
		respAddSchedule = SuccessRespAddSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов LoggingTimeList", func(t *testing.T) {
		args := []string{"", "lts", "-scid", "777"}
		respLoggingTimeList = SuccessRespLoggingTimeList
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов LoggingTimeList c доп аргументами", func(t *testing.T) {
		args := []string{"", "lts", "-scid", "777", "-s", "5", "-p", "1"}
		respLoggingTimeList = SuccessRespLoggingTimeList
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове LoggingTimeList", func(t *testing.T) {
		args := []string{"", "lts", "-scid", "777"}
		respLoggingTimeList = ErrorRespLoggingTimeList
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов DetailLoggingTime", func(t *testing.T) {
		args := []string{"", "lt", "-scid", "777", "-ltid", "777"}
		respDetailLoggingTime = SuccessRespDetailLoggingTime
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Вызов DetailLoggingTime без аргумента", func(t *testing.T) {
		args := []string{"", "lt", "-scid", "777"}
		respDetailLoggingTime = SuccessRespDetailLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове DetailLoggingTime", func(t *testing.T) {
		args := []string{"", "lt", "-scid", "777", "-ltid", "777"}
		respDetailLoggingTime = ErrorRespDetailLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов DeleteLoggingTime", func(t *testing.T) {
		args := []string{"", "rmlt", "-scid", "777", "-ltid", "777"}
		respDeleteLoggingTime = SuccessRespDeleteLoggingTime
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Вызов DeleteLoggingTime без аргумента", func(t *testing.T) {
		args := []string{"", "rmlt", "-scid", "777"}
		respDeleteLoggingTime = SuccessRespDeleteLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове DeleteLoggingTime", func(t *testing.T) {
		args := []string{"", "rmlt", "-scid", "777", "-ltid", "777"}
		respDeleteLoggingTime = ErrorRespDeleteLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов SubmitForApproveSchedule", func(t *testing.T) {
		args := []string{"", "s", "-scid", "777"}
		respSubmitForApproveSchedule = SuccessRespSubmitForApproveSchedule
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове метода SubmitForApproveSchedule", func(t *testing.T) {
		args := []string{"", "s", "-scid", "777"}
		respSubmitForApproveSchedule = ErrorRespSubmitForApproveSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача невалидного флага в SubmitForApproveSchedule", func(t *testing.T) {
		args := []string{"", "s", "-scid", "ар"}
		respSubmitForApproveSchedule = SuccessRespDetailSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача неправильного флага в SubmitForApproveSchedule", func(t *testing.T) {
		args := []string{"", "s", "-fake", "777"}
		respSubmitForApproveSchedule = SuccessRespDetailSchedule
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов ApproveLoggingTime", func(t *testing.T) {
		args := []string{"", "aprv", "-scid", "777", "-ltid", "777"}
		respApproveLoggingTime = SuccessRespApproveLoggingTime
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Вызов ApproveLoggingTime без аргумента", func(t *testing.T) {
		args := []string{"", "aprv", "-scid", "777"}
		respApproveLoggingTime = SuccessRespApproveLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове ApproveLoggingTime", func(t *testing.T) {
		args := []string{"", "aprv", "-scid", "777", "-ltid", "777"}
		respApproveLoggingTime = ErrorRespApproveLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов DeclineLoggingTime", func(t *testing.T) {
		args := []string{"", "dcl", "-scid", "777", "-ltid", "777"}
		respDeclineLoggingTime = SuccessRespDeclineLoggingTime
		err = app.Run(args)
		require.NoError(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Вызов DeclineLoggingTime без аргумента", func(t *testing.T) {
		args := []string{"", "dcl", "-scid", "777"}
		respDeclineLoggingTime = SuccessRespDeclineLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове DeclineLoggingTime", func(t *testing.T) {
		args := []string{"", "dcl", "-scid", "777", "-ltid", "777"}
		respDeclineLoggingTime = ErrorRespDeclineLoggingTime
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Успешный вызов AddLoggingTime", func(t *testing.T) {
		args := []string{"", "al", "-scid", "777"}
		respAddLoggingTime = SuccessRespAddLoggingTime
		err = app.Run(args)
		//Тест не может записать в открытый файл и выпадает ошибка
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Ошибка при вызове метода AddLoggingTime", func(t *testing.T) {
		args := []string{"", "al", "-scid", "777"}
		respAddLoggingTime = ErrorRespAddLoggingTime
		//Тест не может записать в открытый файл и выпадает ошибка
		err = app.Run(args)
		require.Error(t, err)
		assert.Equal(t, "1", exitIndicator)
		exitIndicator = ""
	})
	t.Run("Передача невалидного флага в AddLoggingTime", func(t *testing.T) {
		args := []string{"", "al", "-scid", "ар"}
		respAddLoggingTime = SuccessRespAddLoggingTime
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
	return respAddSchedule()
}

func (f *fakeClient) DetailSchedule(scheduleId api.ScheduleId) (*api.Schedule, error) {
	return respScheduleDetail()
}

func (f *fakeClient) LoggingTimeList(scheduleId api.ScheduleId, options *api.OptionsLT) ([]*api.LoggingTime, error) {
	return respLoggingTimeList()
}

func (f *fakeClient) AddLoggingTime(scheduleId api.ScheduleId, loggingTime *api.AddLoggingTime) (*api.LoggingTime, error) {
	return respAddLoggingTime()
}

func (f *fakeClient) DetailLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId) (*api.LoggingTime, error) {
	return respDetailLoggingTime()
}

func (f *fakeClient) DeleteLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId) error {
	return respDeleteLoggingTime()
}

func (f *fakeClient) SubmitForApproveSchedule(scheduleId api.ScheduleId) (*api.Schedule, error) {
	return respSubmitForApproveSchedule()
}

func (f *fakeClient) ApproveLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId, comment string) (*api.LoggingTime, error) {
	return respApproveLoggingTime()
}

func (f *fakeClient) DeclineLoggingTime(scheduleId api.ScheduleId, loggingTimeId api.LoggingTimeId, comment string) (*api.LoggingTime, error) {
	return respDeclineLoggingTime()
}


func SuccessRespSchedules() ([]*api.Schedule, error) {
	return []*api.Schedule{&fakeSchedule1, &fakeSchedule2}, nil
}

func ErrorRespSchedules() ([]*api.Schedule, error) {
	return nil, errors.New("error from schedules method")
}

func SuccessRespDetailSchedule() (*api.Schedule, error) {
	return &fakeSchedule1, nil
}

func ErrorRespDetailSchedule() (*api.Schedule, error) {
	return nil, errors.New("error from scheduleDetail method")
}

func SuccessRespAddSchedule() (*api.Schedule, error) {
	return &fakeSchedule1, nil
}

func ErrorRespAddSchedule() (*api.Schedule, error) {
	return nil, errors.New("error from addSchedule method")
}

func SuccessRespLoggingTimeList() ([]*api.LoggingTime, error) {
	return []*api.LoggingTime{&fakeLoggingTime1, &fakeLoggingTime2}, nil
}

func ErrorRespLoggingTimeList() ([]*api.LoggingTime, error) {
	return nil, errors.New("error from LoggingTimeList method")
}

func SuccessRespDetailLoggingTime() (*api.LoggingTime, error) {
	return &fakeLoggingTime1, nil
}

func ErrorRespDetailLoggingTime() (*api.LoggingTime, error) {
	return nil, errors.New("error from DeleteLoggingTime method")
}

func SuccessRespDeleteLoggingTime() error {
	return nil
}

func ErrorRespDeleteLoggingTime() error {
	return errors.New("error from DeleteLoggingTime method")
}

func SuccessRespSubmitForApproveSchedule() (*api.Schedule, error) {
	return &fakeSchedule1, nil
}

func ErrorRespSubmitForApproveSchedule() (*api.Schedule, error) {
	return nil, errors.New("error from SubmitForApproveSchedule method")
}

func SuccessRespApproveLoggingTime() (*api.LoggingTime, error) {
	return &fakeLoggingTime1, nil
}

func ErrorRespApproveLoggingTime() (*api.LoggingTime, error) {
	return nil, errors.New("error from ApproveLoggingTime method")
}

func SuccessRespDeclineLoggingTime() (*api.LoggingTime, error) {
	return &fakeLoggingTime1, nil
}

func ErrorRespDeclineLoggingTime() (*api.LoggingTime, error) {
	return nil, errors.New("error from DeclineLoggingTime method")
}

func SuccessRespAddLoggingTime() (*api.LoggingTime, error) {
	return &fakeLoggingTime1, nil
}

func ErrorRespAddLoggingTime() (*api.LoggingTime, error) {
	return nil, errors.New("error from AddLoggingTime method")
}
