package api

import "suft_sdk/internal/http-client/schedule"

type API interface {
	Schedules() ([]schedule.Schedule, error)
	AddSchedule(schedule.Schedule) error
	UpdateSchedule(schedule.Schedule) error
	SubmitForApproveSchedule(schedule.Schedule) error
	ApproveSchedule(id int) error
}

type SuftAPI struct {
}

func NewSuftAPI() *SuftAPI {
	return &SuftAPI{}
}

func (s *SuftAPI) Schedules() ([]schedule.Schedule, error) {
	panic("implement me")
}

func (s *SuftAPI) AddSchedule(schedule *schedule.Schedule) error {
	panic("implement me")
}

func (s *SuftAPI) UpdateSchedule(schedule *schedule.Schedule) error {
	panic("implement me")
}

func (s *SuftAPI) SubmitForApproveSchedule(schedule *schedule.Schedule) error {
	panic("implement me")
}

func (s *SuftAPI) ApproveSchedule(id int) error {
	panic("implement me")
}
