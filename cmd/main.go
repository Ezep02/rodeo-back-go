package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	HTTP_ROUTER "github.com/ezep02/rodeo/internal/router"
	"github.com/redis/go-redis/v9"

	"github.com/ezep02/rodeo/pkg/db"
	"github.com/joho/godotenv"
)

func main() {

	// # Carga variables de entorno
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error cargando .env at main: %v", err)
	}

	// # Configuracion de la base de datos
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// # Configuracion de Cloudinary
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Conecta con la base de datos
	cnn, err := db.DB_Connection(fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbPort, dbName))
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err.Error())
	}

	// # Configuracion de Redis
	redisAddr := os.Getenv("REDIS_ADDR")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// # Inicializa Cloudinary
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Println("Error obteniendo variables de entorno del cache")
		return
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Println("Error iniciando cloudinary")
		return
	}

	// # Inicia el router
	r := HTTP_ROUTER.NewRouter(cnn, cld, redisClient)

	PORT := 9090
	log.Printf("Servidor iniciado en %d", PORT)
	r.Run(":" + fmt.Sprintf("%d", PORT))
}
