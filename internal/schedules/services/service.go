package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/schedules/models"
	"github.com/ezep02/rodeo/internal/schedules/repository"
)

type ScheduleService struct {
	Sch_repo *repository.SchedulesRepository
}

func NewOrderService(sch_repo *repository.SchedulesRepository) *ScheduleService {
	return &ScheduleService{
		Sch_repo: sch_repo,
	}
}

func (sch_s *ScheduleService) CreateNewSchedulService(ctx context.Context, schedul *[]models.ShiftRequest, userID int) (*[]models.ScheduleResponse, error) {
	return sch_s.Sch_repo.CreateNewSchedul(ctx, schedul, userID)
}

func (sch_s *ScheduleService) CreateNewShift(ctx context.Context, shift *[]models.Shift) (*[]models.Shift, error) {
	return sch_s.Sch_repo.CreateNewShift(ctx, shift)
}

func (sch_s *ScheduleService) DeleteShifts(ctx context.Context, ID_array []int) error {
	return sch_s.Sch_repo.DeleteShifts(ctx, ID_array)
}

func (sch_s *ScheduleService) UpdateShiftList(ctx context.Context, data *[]models.Shift) (*[]models.Shift, error) {
	return sch_s.Sch_repo.UpdateShift(ctx, data)
}

func (sch_s *ScheduleService) UpdateScheduleList(ctx context.Context, userID int, data *[]models.Schedule) (*[]models.Schedule, error) {
	return sch_s.Sch_repo.UpdateSchedules(ctx, userID, data)
}
