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
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo .env: %v", err)
	}

	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	// dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")

	cnn, err := db.DB_Connection(fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbPort, dbName))

	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	r := chi.NewRouter()

	// ConfiguraciOn de CORS

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Config Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Routers
	auth.RegisterAuthRoutes(r, cnn)
	services.ServicesRouter(r, cnn, redisClient)
	orders.OrderRoutes(r, cnn, redisClient)
	schedules.SchedulesRoutes(r, cnn)
	analytics.AnalyticsRoutes(r, cnn, redisClient)

	// Crear y configurar el servidor HTTP
	srv := &http.Server{
		Handler:      r,
		Addr:         ":9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Servidor iniciado en %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
