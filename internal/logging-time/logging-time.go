package logging_time

import (
	"suft_sdk/internal/schedule"
)

type TimeLoggerInterface interface {
	Validate()
}

type LoggingTime struct {
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

type AddLoggingTime struct {
	CommentEmployee string
	Day1Time        int
	Day2Time        int
	Day3Time        int
	Day4Time        int
	Day5Time        int
	Day6Time        int
	Day7Time        int
	ProjectId       int
	Task            string
	WorkKindId      int
}

type EditLoggingTime struct {
	CommentAdminEmployee string
	CommentEmployee      string
	Day1Time             int
	Day2Time             int
	Day3Time             int
	Day4Time             int
	Day5Time             int
	Day6Time             int
	Day7Time             int
	ProjectId            int
	StatusCode           string
	Task                 string
	WorkKindId           int
}
