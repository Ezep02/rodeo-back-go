package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/internal/orders/repository"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
}

func NewOrderService(ord_srv *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: ord_srv,
	}
}

func (s *OrderService) CreateNewOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return s.OrderRepo.CreatingNewOrder(ctx, order)
}

func (s *OrderService) GetOrderService(ctx context.Context, barberID int, limit int, offset int) ([]models.BarberPendingOrder, error) {
	return s.OrderRepo.GetBarberPendingOrders(ctx, barberID, limit, offset)
}

// obtener ordenes pendientes del cliente
func (s *OrderService) GetCustomerPendingOrder(ctx context.Context, userID int) ([]models.CustomerPendingOrder, error) {
	return s.OrderRepo.GettingCustomerPendingOrders(ctx, userID)
}

// almacenar el token transitoriamente para realizar el pago
func (s *OrderService) SetOrderToken(ctx context.Context, token string, order models.PendingOrderToken) error {
	return s.OrderRepo.SavingOrderToken(ctx, token, order)
}

func (s *OrderService) GetOrderByToken(ctx context.Context, token string) (models.PendingOrderToken, error) {
	return s.OrderRepo.SearchingOrderToken(ctx, token)
}

// reprogramar horario del turno
func (s *OrderService) UpdateScheduleOrder(ctx context.Context, schedule models.RescheduleRequest, user_id int) (*models.UpdatedCustomerPendingOrder, error) {
	return s.OrderRepo.ReschedulingDateTimeOrder(ctx, schedule, user_id)
}

// Refound
func (s *OrderService) NewRefound(ctx context.Context, refund models.RefundRequest) (*models.UpdatedCustomerPendingOrder, error) {
	return s.OrderRepo.CreatingRefund(ctx, refund)
}

func (s *OrderService) CheckOrderStatus(ctx context.Context, order_id int) (bool, error) {
	return s.OrderRepo.CheckingOrderStatus(ctx, order_id)
}

// Coupon
func (s *OrderService) GenerateCoupon(ctx context.Context, coupon models.Coupon) (models.Coupon, error) {
	return s.OrderRepo.CreatingCoupon(ctx, coupon)
}

func (s *OrderService) GetCustomerCoupons(ctx context.Context, user_id int) (*[]models.Coupon, error) {
	return s.OrderRepo.GettingCustomerCoupons(ctx, user_id)
}
