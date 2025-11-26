package helpers

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCouponCode genera un código aleatorio de cupones de longitud n
func GenerateCouponCode(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, bigInt(len(charset)))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

// bigInt convierte un int a *big.Int
func bigInt(n int) *big.Int {
	return big.NewInt(int64(n))
}
