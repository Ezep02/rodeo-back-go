package utils

import (
	"crypto/sha256"
	"encoding/binary"
)

func GenerateDeterministicUint(input string) uint {
	// Calcula el hash SHA-256 del input
	hash := sha256.Sum256([]byte(input))

	// Extrae los primeros 8 bytes del hash y los convierte en un uint64
	id := binary.BigEndian.Uint64(hash[:8])

	// Convierte uint64 a uint (esto depende de la arquitectura del sistema, pero es seguro)
	return uint(id)
}
