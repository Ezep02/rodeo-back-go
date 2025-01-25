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

func (sch_s *ScheduleService) CreateBarberSchedules(ctx context.Context, schedules *[]models.Schedule) (*[]models.ScheduleResponse, error) {
	return sch_s.Sch_repo.CreateNewSchedules(ctx, schedules)
}

func (sch_s *ScheduleService) DeleteSchedules(ctx context.Context, id []int) error {
	return sch_s.Sch_repo.DeleteSchedules(ctx, id)
}

func (sch_s *ScheduleService) GetAvailableSchedules(ctx context.Context, limit int, offset int) (*[]models.Schedule, error) {
	return sch_s.Sch_repo.GetAvailableSchedules(ctx, limit, offset)
}

func (sch_s *ScheduleService) GetBarberSchedules(ctx context.Context, id int, limit int, offset int) (*[]models.Schedule, error) {
	return sch_s.Sch_repo.GetSchedulesList(ctx, id, limit, offset)
}
func (sch_s *ScheduleService) GetScheduleByID(ctx context.Context, id int) (*models.Schedule, error) {
	return sch_s.Sch_repo.GetScheduleByID(ctx, id)
}

func (sch_s *ScheduleService) UpdateShiftAvailability(ctx context.Context, id int) error {
	return sch_s.Sch_repo.UpdateShiftAvailability(ctx, id)
}

func (sch_s *ScheduleService) GetTotaBarberCuts(ctx context.Context, barberID int) (*[]models.CutsQuantity, error) {
	return sch_s.Sch_repo.GetCutsQuantity(ctx, barberID)
}
