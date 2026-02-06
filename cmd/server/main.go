package main

import (
	"os"

	"example.com/mod/internal/api"
	"example.com/mod/internal/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {
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

	// Swagger ayarları
	docs.SwaggerInfo.Title = "Worker API"
	docs.SwaggerInfo.Description = "Bu API, arka plan görevlerini yönetmek için kullanılır."
	docs.SwaggerInfo.Version = "1.0"
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
