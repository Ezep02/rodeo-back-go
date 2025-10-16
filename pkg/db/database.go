package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB_Connection establece la conexión a la base de datos y realiza las migraciones.
func DB_Connection(dbConn string) (*gorm.DB, error) {

	// Establecer conexión a la base de datos
	connection, err := gorm.Open(mysql.Open(dbConn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("[DB] error al conectar: %w", err) // Retornar el error sin terminar el programa
	}

	log.Println("[DB]: Successful connection")

	return connection, nil
}
