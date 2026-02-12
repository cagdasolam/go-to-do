package main

import (
	"os"

	"example.com/mod/docs"
	"example.com/mod/internal/api"
	"example.com/mod/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GO TO DO API
// @version 1.0
// @description Basic todo API built with Go, Gin, and GORM
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// .env dosyasını yükle
	_ = godotenv.Load()

	// Logger kurulumu
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Veritabanı bağlantısı
	database, err := db.ConnectDB()
	if err != nil {
		sugar.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}
	sugar.Info("Veritabanına başarıyla bağlanıldı.")

	// Gin router kurulumu
	router := gin.Default()

	// Swagger ayarları - docs paketi üzerinden
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API rotaları
	api.SetupRoutes(router, sugar, database)

	// Sunucuyu başlatma
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	sugar.Infof("Sunucu %s portunda başlatılıyor...", port)
	if err := router.Run(":" + port); err != nil {
		sugar.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
