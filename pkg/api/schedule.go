package api

type StatusCode string

const (
	Approved  StatusCode = "УТВ"
	Declined  StatusCode = "ОТКЛ"
	Created   StatusCode = "СЗ"
	ToApprove StatusCode = "НУ"
)

type Employee struct {
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	Id         int    `json:"id"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type Period struct {
	CloseDate  string `json:"closeDate"`
	EndDate    string `json:"endDate"`
	Id         int    `json:"id"`
	StartDate  string `json:"startDate"`
	WeekNumber int    `json:"weekNumber"`
}

type Schedule struct {
	client     *Client
	Author     Employee `json:"author"`
	Id         int      `json:"id"`
	Period     Period   `json:"period"`
	StatusCode string   `json:"statusCode"`
}

func (s *Schedule) SubmitForApproveSchedule() (*Schedule, error) {
	scheduleId := ScheduleId(s.Id)
	scheduleResp, err := s.client.SubmitForApproveSchedule(scheduleId)
	if err != nil {
		return nil, err
	}
	return scheduleResp, nil
}
