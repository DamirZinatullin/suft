package main

import "fmt"

type Employee struct {
	//todo
}

type Period struct {

	//todo
}

type Schedule struct {
	author Employee
	id int
	period Period
}



type ClientInterface interface {
	Authorize() error
	Schedules() ([]Schedule, error)
	AddSchedule(schedule Schedule) error
	DetailSchedule(id int) error
	EditSchedule(id int, fields map[string]string) error
	UpdateSchedule(id int, fields map[string]string) error
	Logout() error
}

func main() {
	fmt.Println("hello")
}
