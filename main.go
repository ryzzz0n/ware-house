// main.go
package main

import (
	"fmt"
	"log"
	"os"

	"warehouse-app/adapters"
	authadapter "warehouse-app/auth/authAdapter"
	authcore "warehouse-app/auth/authCore"
	"warehouse-app/core"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading . .env file")
	}

	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, portStr, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Автомиграция — вынесена в адаптер
	adapters.AutoMigrate(db)

	// Сборка зависимостей
	productRepo := adapters.NewGormProductRepository(db)
	productService := core.NewProductService(productRepo)
	productHandler := adapters.NewHttpProductHandler(productService)

	userRepo := authadapter.NewGormUserRepository(db)
	userService := authcore.NewUserService(userRepo)
	userHandler := authadapter.NewHttpUserHandler(userService)

	// Запуск приложения
	app := fiber.New()

	// Регистрация middleware авторизации
	authadapter.RegisterAuthMiddleware(app)

	// Регистрация маршрутов
	productHandler.RegisterRoutes(app)
	userHandler.RegisterAuthRoutes(app)

	fmt.Println("Server starting on :8000")
	app.Listen(":8000")
}
