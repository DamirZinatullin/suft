package api

type LoggingTime struct {
	scheduleId           ScheduleId
	client               *Client
	AdminEmployee        Employee   `json:"adminEmployee"`
	CommentAdminEmployee string     `json:"commentAdminEmployee"`
	CommentEmployee      string     `json:"commentEmployee"`
	Day1Time             float64    `json:"day1Time"`
	Day2Time             float64    `json:"day2Time"`
	Day3Time             float64    `json:"day3Time"`
	Day4Time             float64    `json:"day4Time"`
	Day5Time             float64    `json:"day5Time"`
	Day6Time             float64    `json:"day6Time"`
	Day7Time             float64    `json:"day7Time"`
	Id                   int        `json:"id"`
	ImportedFrom         string     `json:"importedFrom"`
	ProjectId            int        `json:"projectId"`
	StatusCode           StatusCode `json:"statusCode"`
	Task                 string     `json:"task"`
	WorkKindId           int        `json:"workKindId"`
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

func (l *LoggingTime) ApproveLoggingTime(comment string) (*LoggingTime, error) {
	loggingTimeId := LoggingTimeId(l.Id)
	loggingTime, err := l.client.ApproveLoggingTime(l.scheduleId, loggingTimeId, comment)
	if err != nil {
		return nil, err
	}
	return loggingTime, nil
}

func (l *LoggingTime) DeleteLoggingTime() (err error) {
	loggingTimeId := LoggingTimeId(l.Id)
	err = l.client.DeleteLoggingTime(l.scheduleId, loggingTimeId)
	if err != nil {
		return err
	}
	return nil
}
