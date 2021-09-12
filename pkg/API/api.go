package api

import "suft_sdk/internal/http-client/schedule"

type API interface {
	Schedules() ([]schedule.Schedule, error)
	AddSchedule(schedule.Schedule) error
	UpdateSchedule(schedule.Schedule) error
	SubmitForApproveSchedule(schedule.Schedule) error
	ApproveSchedule(id int) error
}
