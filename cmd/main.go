package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ezep02/rodeo/internal/analytics"
	"github.com/ezep02/rodeo/internal/auth"
	"github.com/ezep02/rodeo/internal/orders"
	"github.com/ezep02/rodeo/internal/schedules"
	"github.com/ezep02/rodeo/internal/services"
	"github.com/ezep02/rodeo/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo .env: %v", err)
	}

	DB_PASSWORD := viper.GetString("DB_PASSWORD")
	DB_NAME := viper.GetString("DB_NAME")

	cnn, err := db.DB_Connection(fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_PASSWORD, DB_NAME))
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	r := chi.NewRouter()

	// ConfiguraciOn de CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	auth.RegisterAuthRoutes(r, cnn)
	services.ServicesRouter(r, cnn)
	orders.OrderRoutes(r, cnn)
	schedules.SchedulesRoutes(r, cnn)
	analytics.AnalyticsRoutes(r, cnn)

	//Crear y configurar el servidor HTTP
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Servidor iniciado en %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
