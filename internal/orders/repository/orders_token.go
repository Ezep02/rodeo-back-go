package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
)

func (r *OrderRepository) SavingOrderToken(ctx context.Context, token string, order models.PendingOrderToken) error {
	orderCacheKey := fmt.Sprintf("order:%s", token)

	data, _ := json.Marshal(order)
	// Almacena el token en Redis sin necesidad de serializacion
	err := r.RedisConnection.Set(ctx, orderCacheKey, data, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) SearchingOrderToken(ctx context.Context, token string) (models.PendingOrderToken, error) {
	var (
		pendingOrder  models.PendingOrderToken
		orderCacheKey = fmt.Sprintf("order:%s", token)
	)

	cachedOrder, err := r.RedisConnection.Get(ctx, orderCacheKey).Result()
	if err != nil {
		return pendingOrder, err // Retorna un struct vacío y el error
	}

	// Intentamos deserializar el JSON
	if err := json.Unmarshal([]byte(cachedOrder), &pendingOrder); err != nil {
		return pendingOrder, err // Retorna error si falla el Unmarshal
	}

	return pendingOrder, nil
}
