package helpers

import "github.com/ezep02/rodeo/internal/booking/domain/booking"

func CalculateConsequences(isWithin24h bool, isPaymentComplete string) *booking.CancelationResponse {

	isComplete := isPaymentComplete == "total"

	response := &booking.CancelationResponse{
		Canceled:     false,
		LosesDeposit: false,
	}

	// DENTRO DE 24 HORAS
	if isWithin24h {
		if isComplete {
			response.CouponPercent = 50
			response.RequiresCoupon = true
			response.Message = "La cancelación está dentro de las 24 horas. Recibirás un cupón del 50%."
		} else {
			response.CouponPercent = 0
			response.RequiresCoupon = false
			response.LosesDeposit = true
			response.Message = "La cancelación está dentro de las 24 horas. Perderás la seña abonada."
		}
		return response
	}

	// FUERA DE 24 HORAS
	if isComplete {
		response.CouponPercent = 75
		response.RequiresCoupon = true
		response.Message = "La cancelación está fuera de las 24 horas. Recibirás un cupón del 75%."
	} else {
		response.CouponPercent = 25
		response.RequiresCoupon = true
		response.Message = "La cancelación está fuera de las 24 horas. Recibirás un cupón del 25%."
	}

	return response
}
