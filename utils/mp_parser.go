package utils

import (
	"fmt"
)

type Metadata struct {
	SlotID            uint
	PaymentPercentage int
	UserID            uint
	CouponCode        string
}

type SurchargeMetadata struct {
	OldSlotId uint
	NewSlotId uint
	ApptId    uint
}

func MetadataParser(metadata map[string]any) (*Metadata, error) {
	// slot_id
	rawSlotID, ok := metadata["slot_id"]
	if !ok {
		return nil, fmt.Errorf("faltante: metadata.slot_id")
	}

	slotStr, ok := rawSlotID.(float64)
	if !ok {
		return nil, fmt.Errorf("metadata.slot_id debe ser string no vacío")
	}

	// payment_percentage
	rawPayment, ok := metadata["payment_percentage"]
	if !ok {
		return nil, fmt.Errorf("faltante: metadata.payment_percentage")
	}

	paymentStr, ok := rawPayment.(float64)
	if !ok {
		return nil, fmt.Errorf("metadata.payment_percentage debe ser distinto de vacio")
	}

	if paymentStr != 50 && paymentStr != 100 {
		return nil, fmt.Errorf("payment_percentage debe ser 50 o 100")
	}

	// user_id
	var userID uint = 0
	if raw, ok := metadata["user_id"]; ok {
		// JSON numbers are float64 by default
		if floatVal, ok := raw.(float64); ok {
			userID = uint(floatVal)
		}
	}

	// Validar Time
	return &Metadata{
		SlotID:            uint(slotStr),
		PaymentPercentage: int(paymentStr),
		UserID:            userID,
		CouponCode:        metadata["coupon_code"].(string), // Optional, can be empty
	}, nil
}

func SurchargeMetadataParcer(metadata map[string]any) (*SurchargeMetadata, error) {

	var (
		oldSlotId uint = 0
		newSlotId uint = 0
		apptId    uint = 0
	)

	if raw, ok := metadata["old_slot_id"]; ok {
		// JSON numbers are float64 by default
		if floatVal, ok := raw.(float64); ok {
			oldSlotId = uint(floatVal)
		}
	}

	if raw, ok := metadata["new_slot_id"]; ok {
		// JSON numbers are float64 by default
		if floatVal, ok := raw.(float64); ok {
			newSlotId = uint(floatVal)
		}
	}

	if raw, ok := metadata["appt_id"]; ok {
		// JSON numbers are float64 by default
		if floatVal, ok := raw.(float64); ok {
			apptId = uint(floatVal)
		}
	}

	// Validar Time
	return &SurchargeMetadata{
		OldSlotId: oldSlotId,
		NewSlotId: newSlotId,
		ApptId:    apptId,
	}, nil
}
