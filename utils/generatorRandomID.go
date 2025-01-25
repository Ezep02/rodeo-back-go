package utils

import (
	"crypto/rand"
	"encoding/binary"
)

func GenerateRandomID() (uint, error) {
	var id uint32            // Usamos uint32 para asegurar compatibilidad y evitar desbordamientos
	bytes := make([]byte, 4) // 4 bytes = 32 bits (uint32)
	_, err := rand.Read(bytes)
	if err != nil {
		return 0, err
	}
	id = binary.BigEndian.Uint32(bytes) // Convertimos los bytes a un número
	return uint(id), nil
}
