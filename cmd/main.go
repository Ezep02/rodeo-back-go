package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ezep02/rodeo/internal/repository"
	"github.com/ezep02/rodeo/internal/service"
	TransportHTTP "github.com/ezep02/rodeo/internal/transport/http"
	"github.com/redis/go-redis/v9"

	"github.com/ezep02/rodeo/pkg/db"
	"github.com/joho/godotenv"
)

func main() {

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

	// FIX V1

	// configuracion de redis
	redisAddr := os.Getenv("REDIS_ADDR")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,
	})

	// REPOS
	gormApptRepo := repository.NewGormAppointmentRepo(cnn, redisClient)
	gormProdRepo := repository.NewGormProductRepo(cnn, redisClient)
	gormAuthRepo := repository.NewGormAuthRepo(cnn)
	gormSlotRepo := repository.NewGormSlotRepo(cnn)
	gormReviewRepo := repository.NewGormReviewRepo(cnn, redisClient)
	gormAnalyticRepo := repository.NewGormAnalyticRepo(cnn)
	gormCouponRepo := repository.NewGormCouponRepo(cnn)
	gormInfoRepo := repository.NewGormInfoRepo(cnn)

	// SERVICES
	apptSvc := service.NewAppointmentService(gormApptRepo, gormProdRepo)
	prodSvc := service.NewProductService(gormProdRepo)
	authSvc := service.NewAuthRepository(gormAuthRepo)
	slotSvc := service.NewSlotService(gormSlotRepo)
	revSvc := service.NewReviewService(gormReviewRepo, gormApptRepo)
	analyticSvc := service.NewAnalyticService(gormAnalyticRepo)
	couponSvc := service.NewCouponService(gormCouponRepo)
	infoSvc := service.NewInfoRepository(gormInfoRepo)

	r := TransportHTTP.NewRouter(apptSvc, prodSvc, authSvc, slotSvc, revSvc, analyticSvc, couponSvc, infoSvc)

	PORT := 9090
	log.Printf("Servidor iniciado en %d", PORT)
	r.Run(":" + fmt.Sprintf("%d", PORT))
}
