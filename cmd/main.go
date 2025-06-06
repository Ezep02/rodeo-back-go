package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ezep02/rodeo/internal/analytics"
	"github.com/ezep02/rodeo/internal/auth"
	"github.com/ezep02/rodeo/internal/orders"
	"github.com/ezep02/rodeo/internal/schedules"
	"github.com/ezep02/rodeo/internal/services"
	"github.com/ezep02/rodeo/pkg/db"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	// Mux base
	mux := http.NewServeMux()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error cargando .env at main: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	// dbHost := viper.GetString("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	cnn, err := db.DB_Connection(fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbPort, dbName))
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Middleware CORS
	cors_handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(mux)

	// Config Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Routers
	auth.RegisterAuthRoutes(mux, cnn)

	services.ServicesRouter(mux, cnn, redisClient)
	orders.OrderRoutes(mux, cnn, redisClient)

	schedules.SchedulesRoutes(mux, cnn)
	analytics.AnalyticsRoutes(mux, cnn, redisClient)

	//

	// Crear y configurar el servidor HTTP
	srv := &http.Server{
		Handler:      cors_handler,
		Addr:         ":9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Servidor iniciado en %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
