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

func (sch_s *ScheduleService) GetSchedules(ctx context.Context, id int) (*[]models.ScheduleResponse, error) {
	return sch_s.Sch_repo.GetSchedules(ctx, id)
}

func (sch_s *ScheduleService) CreateNewSchedulService(ctx context.Context, schedul *[]models.ShiftRequest, barberID int, barberName string) (*[]models.ScheduleResponse, error) {
	return sch_s.Sch_repo.CreateNewSchedul(ctx, schedul, barberID, barberName)
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

func (sch_s *ScheduleService) UpdateShiftByID(ctx context.Context, shift_id int) (*models.Shift, error) {
	return sch_s.Sch_repo.UpdateShiftByID(ctx, shift_id)
}

func (sch_s *ScheduleService) UpdateScheduleList(ctx context.Context, barberID int, data *[]models.Schedule) (*[]models.Schedule, error) {
	return sch_s.Sch_repo.UpdateSchedules(ctx, barberID, data)
}
