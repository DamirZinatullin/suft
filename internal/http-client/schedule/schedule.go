package schedule

type EmployeeInterface interface {
	Validate()
}

type PeriodInterface interface {
	Validate()
}

type ScheduleInterface interface {
	Validate()
}

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
	Author     Employee `json:"author"`
	Id         int      `json:"id"`
	Period     Period   `json:"period"`
	StatusCode string   `json:"statusCode"`
}

type EditStatusSchedule struct {
	statusCode string
}