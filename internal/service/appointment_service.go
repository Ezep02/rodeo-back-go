package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type AppointmentService struct {
	apptRepo domain.AppointmentRepository
	prodRepo domain.ProductRepository
}

func NewAppointmentService(apptRepo domain.AppointmentRepository, prodRepo domain.ProductRepository) *AppointmentService {
	return &AppointmentService{apptRepo, prodRepo}
}

func (s *AppointmentService) Schedule(ctx context.Context, appt *domain.Appointment) error {

	// 1. Validar que los productos existan
	for i, prod := range appt.Products {

		existingProd, err := s.prodRepo.GetByID(ctx, prod.ID)
		if err != nil {
			if err == domain.ErrNotFound {
				return errors.New("producto no encontrado")
			}
			return err
		}

		// Reemplazar el producto en la cita con el existente
		appt.Products[i] = *existingProd
	}

	// 4. Crear la cita
	return s.apptRepo.Create(ctx, appt)
}

func (s *AppointmentService) GetByID(ctx context.Context, id uint) (*domain.Appointment, error) {
	return s.apptRepo.GetByID(ctx, id)
}

func (s *AppointmentService) ListAll(ctx context.Context) ([]domain.Appointment, error) {
	return s.apptRepo.List(ctx)
}

func (s *AppointmentService) Update(ctx context.Context, id, slot_id uint, updatedAppt *domain.Appointment) error {

	// 1. Verificar que la cita exista
	existing, err := s.apptRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Manener el id original
	updatedAppt.ID = existing.ID

	// 4. Evitar citas duplicadas
	existingAppts, err := s.apptRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existing := range existingAppts {
		if existing.ID != id && s.appointmentsOverlap(existing, *updatedAppt) {
			return errors.New("ya existe una cita en esa fecha y hora")
		}
	}

	// 6. Actualizar la cita
	return s.apptRepo.Update(ctx, updatedAppt, slot_id)
}

func (s *AppointmentService) Cancel(ctx context.Context, id uint) error {
	// 1. Verificar que la cita exista
	_, err := s.apptRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			return errors.New("cita no encontrada")
		}
		return err
	}

	// 2. Eliminar la cita
	return s.apptRepo.Delete(ctx, id)
}

func (s *AppointmentService) GetTotalPrice(ctx context.Context, id uint) (float64, error) {

	// 1. Obtener la cita por ID
	appt, err := s.apptRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			return 0, errors.New("cita no encontrada")
		}
		return 0, err
	}

	// 2. Calcular el precio total de los productos
	totalPrice := 0.0

	for _, prod := range appt.Products {
		totalPrice += prod.Price
	}

	// 3. Devolver el precio total
	return totalPrice, nil
}

func (s *AppointmentService) GetByUserID(ctx context.Context, id uint) ([]domain.Appointment, error) {
	return s.apptRepo.GetByUserID(ctx, id)
}

// utilizad
func (s *AppointmentService) appointmentsOverlap(existing, newAppt domain.Appointment) bool {
	// Compara si las citas se superponen
	if existing.Slot.Date.Equal(newAppt.Slot.Date) && existing.Slot.Time == newAppt.Slot.Time {
		return true
	}
	return false
}
