package time_logger

import "suft_sdk/internal/http-client/schedule"

type TimeLoggerInterface interface {
	Validate()
}

var TimeLogger struct {
	ScheduleId           int `json:"-"`
	AdminEmployee        schedule.Employee
	CommentAdminEmployee string
	CommentEmployee      string
	Day1Time             int
	Day2Time             int
	Day3Time             int
	Day4Time             int
	Day5Time             int
	Day6Time             int
	Day7Time             int
	Id                   int
	ImportedFrom         string
	ProjectId            int
	StatusCode           string
	Task                 string
	WorkKindId           int
}
