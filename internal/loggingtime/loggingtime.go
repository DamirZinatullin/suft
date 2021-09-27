package loggingtime

import (
	"suftsdk/internal/schedule"
)

type status string

const Approved status = "УТВ"
const Denied status = "ОТКЛ"
const Created status = "СЗ"
const ToApprove status = "НУ"
const ToApproveAgain status = "НП"
const Corrected status = "СК"

type TimeLoggerInterface interface {
	Validate()
}

type LoggingTime struct {
	AdminEmployee        schedule.Employee `json:"adminEmployee"`
	CommentAdminEmployee string            `json:"commentAdminEmployee"`
	CommentEmployee      string            `json:"commentEmployee"`
	Day1Time             float64           `json:"day1Time"`
	Day2Time             float64           `json:"day2Time"`
	Day3Time             float64           `json:"day3Time"`
	Day4Time             float64           `json:"day4Time"`
	Day5Time             float64           `json:"day5Time"`
	Day6Time             float64           `json:"day6Time"`
	Day7Time             float64           `json:"day7Time"`
	Id                   int               `json:"id"`
	ImportedFrom         string            `json:"importedFrom"`
	ProjectId            int               `json:"projectId"`
	StatusCode           status            `json:"statusCode"`
	Task                 string            `json:"task"`
	WorkKindId           int               `json:"workKindId"`
}

type AddLoggingTime struct {
	CommentEmployee string  `json:"commentEmployee"`
	Day1Time        float64 `json:"day1Time"`
	Day2Time        float64 `json:"day2Time"`
	Day3Time        float64 `json:"day3Time"`
	Day4Time        float64 `json:"day4Time"`
	Day5Time        float64 `json:"day5Time"`
	Day6Time        float64 `json:"day6Time"`
	Day7Time        float64 `json:"day7Time"`
	ProjectId       int     `json:"projectId"`
	Task            string  `json:"task"`
	WorkKindId      int     `json:"workKindId"`
}

type EditLoggingTime struct {
	CommentAdminEmployee string  `json:"commentAdminEmployee"`
	CommentEmployee      string  `json:"commentEmployee"`
	Day1Time             float64 `json:"day1Time"`
	Day2Time             float64 `json:"day2Time"`
	Day3Time             float64 `json:"day3Time"`
	Day4Time             float64 `json:"day4Time"`
	Day5Time             float64 `json:"day5Time"`
	Day6Time             float64 `json:"day6Time"`
	Day7Time             float64 `json:"day7Time"`
	ProjectId            int     `json:"projectId"`
	StatusCode           status  `json:"statusCode"`
	Task                 string  `json:"task"`
	WorkKindId           int     `json:"workKindId"`
}