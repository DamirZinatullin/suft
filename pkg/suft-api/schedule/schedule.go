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
	Email string
	FirstName string
	Id int
	LastName string
	MiddleName string
}


type Period struct {
	CloseDate string
	EndDate string
	Id int
	StartDate string
	weekNumber int
}


type Schedule struct {
	author Employee
	id     int
	period Period
	StatusCode string
}
