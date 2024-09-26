package main

import (
	"fmt"
	"github.com/Zelvalna/song-library/config"
	"github.com/Zelvalna/song-library/controllers"
	"github.com/Zelvalna/song-library/migrations"
	"github.com/Zelvalna/song-library/repositories"
	"github.com/Zelvalna/song-library/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Music Library API
// @version         1.0
// @description     API для управления онлайн-библиотекой песен.

// @host      localhost:8080
// @BasePath  /

// @schemes http
func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal("Ошибка загрузки конфигурации: ", err)
	}

	// Подключение к БД
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Ошибка подключения к БД: ", err)
	}

	// Миграции
	err = migrations.Migrate(db)
	if err != nil {
		logrus.Fatal("Ошибка миграции БД: ", err)
	}

	// Инициализация репозиториев и сервисов
	songRepo := repositories.NewSongRepository(db)
	songService := services.NewSongService(songRepo, cfg.APIURL)
	songController := controllers.NewSongController(songService)

	// Инициализация роутера
	router := gin.Default()

	// Маршруты
	router.GET("/songs", songController.GetSongs)
	router.POST("/songs", songController.AddSong)
	router.GET("/songs/:id/text", songController.GetSongText)
	router.PUT("/songs/:id", songController.UpdateSong)
	router.DELETE("/songs/:id", songController.DeleteSong)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	logrus.Info("Сервер запущен на порту ", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		logrus.Fatal("Ошибка запуска сервера: ", err)
	}
}
