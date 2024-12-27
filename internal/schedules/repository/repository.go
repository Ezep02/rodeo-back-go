package repository

import (
	"context"
	"log"

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
func (sc *SchedulesRepository) CreateNewSchedul(ctx context.Context, schedules *[]models.ShiftRequest, barberID int, barberName string) (*[]models.ScheduleResponse, error) {
	var createdSchedules []models.ScheduleResponse

	for _, sch := range *schedules {
		// Crear un nuevo ScheduleRequest
		newSchedule := &models.Schedule{
			Barber_id:    barberID,
			Start_date:   sch.Start,
			End_date:     sch.End,
			Schedule_day: sch.Day,
		}

		// Guardar el nuevo ScheduleRequest en la base de datos
		result := sc.Connection.WithContext(ctx).Create(newSchedule)
		if result.Error != nil {
			return nil, result.Error
		}

		// Recuperar el ultimo ScheduleRequest creado
		var lastSchedule models.Schedule
		sc.Connection.WithContext(ctx).Last(&lastSchedule)

		// Crear un arreglo para almacenar los turnos asociados a este horario
		var shifts []models.Shift

		// Crear los Shifts asociados al ScheduleRequest
		for _, s := range sch.Shift_add {
			newShift := models.Shift{
				Schedule_id:     lastSchedule.ID, // Asociar el Shift con el ScheduleRequest
				Day:             sch.Day,
				Start_time:      s.Start,
				Created_by_name: barberName,
			}

			// Guardar el nuevo Shift en la base de datos
			result := sc.Connection.WithContext(ctx).Create(&newShift)
			if result.Error != nil {
				return nil, result.Error
			}

			// Agregar el Shift al arreglo de turnos (shifts)
			shifts = append(shifts, newShift)
		}

		// Crear un objeto ScheduleResponse con los turnos ya asociados
		createdSchedule := models.ScheduleResponse{
			Model: &gorm.Model{
				ID:        lastSchedule.ID,
				CreatedAt: lastSchedule.CreatedAt,
				UpdatedAt: lastSchedule.UpdatedAt,
				DeletedAt: lastSchedule.DeletedAt,
			},
			Barber_id: barberID,
			Start:     sch.Start,
			End:       sch.End,
			ShiftAdd:  shifts,
			Day:       sch.Day,
		}

		// Agregar el ScheduleResponse creado al arreglo de horarios
		createdSchedules = append(createdSchedules, createdSchedule)
	}

	// Devolver el arreglo de ScheduleResponse
	return &createdSchedules, nil
}

func (sc *SchedulesRepository) GetSchedules(ctx context.Context, barberID int) (*[]models.ScheduleResponse, error) {
	var schedules []models.Schedule

	// Obtener los schedules
	if err := sc.Connection.WithContext(ctx).
		Find(&schedules).Error; err != nil {
		log.Println("Error al obtener los schedules:", err)
		return nil, err
	}

	// Lista para almacenar las respuestas formateadas
	var response []models.ScheduleResponse

	// Recorrer cada schedule para obtener sus shifts
	for _, sch := range schedules {
		var scheduleShifts []models.Shift

		// Obtener los shifts asociados al schedule actual
		if err := sc.Connection.WithContext(ctx).
			Where("schedule_id = ?", sch.ID).
			Find(&scheduleShifts).Error; err != nil {
			log.Printf("Error al obtener los shifts para schedule ID %d: %v", sch.ID, err)
			return nil, err
		}

		// Crear la respuesta para cada schedule incluyendo los shifts filtrados
		response = append(response, models.ScheduleResponse{
			Model:     sch.Model,
			Barber_id: barberID,
			Start:     sch.Start_date,
			End:       sch.End_date,
			Day:       sch.Schedule_day,
			ShiftAdd:  scheduleShifts,
			ID:        int(sch.ID),
		})
	}

	return &response, nil
}

func (sc *SchedulesRepository) CreateNewShift(ctx context.Context, shifts *[]models.Shift) (*[]models.Shift, error) {
	var newShifts []models.Shift

	//Iniciar una transacción para asegurar atomicidad
	tx := sc.Connection.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Crear los Shifts asociados al ScheduleRequest
	for _, s := range *shifts {
		newShift := models.Shift{
			Schedule_id:     s.Schedule_id, // Asociar el Shift con el ScheduleRequest
			Day:             s.Day,
			Start_time:      s.Start_time,
			Available:       s.Available,
			Created_by_name: s.Created_by_name,
			ShiftStatus:     s.ShiftStatus,
		}

		// Guardar el nuevo Shift en la base de datos
		result := tx.Create(&newShift)
		if result.Error != nil {
			// Si ocurre un error, realizar rollback y devolver el error
			tx.Rollback()
			return nil, result.Error
		}

		// Agregar el Shift al arreglo de turnos (shifts)
		newShifts = append(newShifts, newShift)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newShifts, nil
}

func (sc *SchedulesRepository) DeleteShifts(ctx context.Context, ID_array []int) error {

	for _, ID := range ID_array {
		resultDelete := sc.Connection.WithContext(ctx).Where("id = ?", ID).Delete(&models.Shift{})
		//si algo sale mal, devuelve el error
		if resultDelete.Error != nil {
			return resultDelete.Error
		}
	}

	return nil
}

func (sc *SchedulesRepository) UpdateShift(ctx context.Context, data *[]models.Shift) (*[]models.Shift, error) {

	var updatedShiftList []models.Shift

	tx := sc.Connection.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, shift := range *data {
		updatedShift := models.Shift{
			Start_time:  shift.Start_time,
			Schedule_id: shift.Schedule_id,
			Day:         shift.Day,
		}

		result := tx.Where("id = ?", shift.ID).Updates(updatedShift)

		if result.Error != nil {
			log.Println("[UPDATE SHIFT] Error al actualizar registro, realizando rollback")
			tx.Rollback()
			return nil, result.Error
		}

		updatedShiftList = append(updatedShiftList, updatedShift)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &updatedShiftList, nil
}

func (sc *SchedulesRepository) UpdateSchedules(ctx context.Context, barberID int, data *[]models.Schedule) (*[]models.Schedule, error) {

	var updatedSchedulesList []models.Schedule

	tx := sc.Connection.WithContext(ctx).Begin()

	for _, sch := range *data {

		updatedSchedule := models.Schedule{
			Start_date:   sch.Start_date,
			Schedule_day: sch.Schedule_day,
			End_date:     sch.End_date,
			Barber_id:    barberID,
		}

		result := tx.Where("id = ?", sch.ID).Updates(updatedSchedule)

		if result.Error != nil {
			log.Println("[UPDATE SCHEDULES] Error al actualizar registro, realizando rollback")
			tx.Rollback()
			return nil, result.Error
		}

		updatedSchedulesList = append(updatedSchedulesList, updatedSchedule)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &updatedSchedulesList, nil
}

func (sc *SchedulesRepository) UpdateShiftByID(ctx context.Context, shiftID int) (*models.Shift, error) {
	// Actualizar el campo "available"
	result := sc.Connection.WithContext(ctx).
		Model(&models.Shift{}).
		Where("id = ?", shiftID).
		Update("available", false)

	// Verificar errores en la actualización
	if result.Error != nil {
		return nil, result.Error
	}

	// Buscar el registro actualizado
	ShiftResponse := &models.Shift{}
	findShift := sc.Connection.WithContext(ctx).Where("id = ?", shiftID).First(ShiftResponse)

	// Verificar errores en la búsqueda
	if findShift.Error != nil {
		return nil, findShift.Error
	}

	return ShiftResponse, nil
}
