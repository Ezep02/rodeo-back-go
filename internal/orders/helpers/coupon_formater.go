package helpers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789232311223ACVVASD213ADX"

func GenerateCouponCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CouponFormater(refound_data models.RefundRequest, user_id int) (models.Coupon, error) {

	if refound_data.Order_id == 0 || refound_data.Shift_id == 0 {
		return models.Coupon{}, fmt.Errorf("invalid refund request: missing required fields")
	}

	code := GenerateCouponCode(10)

	return models.Coupon{
		Code:            strings.ToUpper(code),
		UserID:          user_id,
		DiscountPercent: refound_data.Refund_percentaje,
		Available:       true,
		Used:            false,
		CreatedAt:       time.Now(),
		AvailableToDate: time.Now().AddDate(0, 0, 7),
		Coupon_type:     refound_data.Refund_type,
	}, nil
}
