package repository

import (
	"context"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/schedules/models"
	"gorm.io/gorm"
)

type SchedulesRepository struct {
	Connection *gorm.DB
}

func NewSchedulesRepository(db *gorm.DB) *SchedulesRepository {
	return &SchedulesRepository{
		Connection: db,
	}
}
func (sc *SchedulesRepository) CreateNewSchedules(ctx context.Context, schedules *[]models.Schedule) (*[]models.ScheduleResponse, error) {
	var schedulesResponse []models.ScheduleResponse

	tx := sc.Connection.WithContext(ctx)

	// Crear todos los registros en una sola operación
	if result := tx.Create(&schedules); result.Error != nil {
		log.Println("Error creando schedules:", result.Error)
		tx.Rollback()
		return nil, result.Error
	}

	// Llenar la respuesta con los datos insertados
	for _, schedule := range *schedules {
		schedulesResponse = append(schedulesResponse, models.ScheduleResponse{
			Model:             schedule.Model,
			Created_by_name:   schedule.Created_by_name,
			Barber_id:         schedule.Barber_id,
			Available:         schedule.Available,
			Schedule_day_date: schedule.Schedule_day_date,
			Start_time:        schedule.Start_time,
		})
	}

	return &schedulesResponse, nil
}

func (sc *SchedulesRepository) DeleteSchedules(ctx context.Context, ids []int) error {
	cnn := sc.Connection.WithContext(ctx)

	if result := cnn.Where("id IN ?", ids).Delete(&models.Schedule{}); result.Error != nil {
		log.Println("Error eliminando schedules:", result.Error)
		return result.Error
	}

	return nil
}

// Especifico para el admin panel
func (sc *SchedulesRepository) GetAvailableSchedules(ctx context.Context, limit int, offset int) (*[]models.Schedule, error) {

	var schedules []models.Schedule

	// Obtener los schedules creados por el barbero
	today := time.Now().Truncate(24 * time.Hour)
	if err := sc.Connection.WithContext(ctx).Where("schedule_day_date >= ?", today).Limit(limit).Find(&schedules).Error; err != nil {
		log.Println("Error al obtener los schedules:", err)
		return nil, err
	}

	return &schedules, nil
}

// Especifico para usuarios
func (sc *SchedulesRepository) GetSchedulesList(ctx context.Context, barberID int, limit int, offset int) (*[]models.Schedule, error) {
	var schedules []models.Schedule

	// Obtener los schedules creados por el barbero
	today := time.Now().Truncate(24 * time.Hour)
	if err := sc.Connection.WithContext(ctx).Where("barber_id = ? AND schedule_day_date >= ?", barberID, today).Limit(limit).Offset(offset).Find(&schedules).Error; err != nil {
		log.Println("Error al obtener los schedules:", err)
		return nil, err
	}

	return &schedules, nil
}

func (sc *SchedulesRepository) GetScheduleByID(ctx context.Context, id int) (*models.Schedule, error) {
	var updatedShift *models.Schedule

	result := sc.Connection.WithContext(ctx).Where("id = ?", id).Find(&updatedShift)

	if result.Error != nil {
		log.Printf("error %+v", result.Error)
		return nil, result.Error
	}

	return updatedShift, nil
}

func (sc *SchedulesRepository) UpdateShiftAvailability(ctx context.Context, id int) error {

	if err := sc.Connection.WithContext(ctx).Model(&models.Schedule{}).Where("id = ?", id).Update("available", false); err != nil {
		log.Println("Error updating schedule availability")
		log.Printf("error %+v", err.Error)
		return err.Error
	}

	return nil
}
