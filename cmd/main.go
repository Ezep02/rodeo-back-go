package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ezep02/rodeo/internal/auth"
	"github.com/ezep02/rodeo/internal/orders"
	"github.com/ezep02/rodeo/internal/services"
	"github.com/ezep02/rodeo/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	// enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cnn, _ := db.DB_Connection("root:7nc4381c4t@tcp(127.0.0.1:3306)/goMeli?charset=utf8mb4&parseTime=True&loc=Local")

	r := chi.NewRouter()

	// Configuración de CORS
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
