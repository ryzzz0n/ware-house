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
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	adapters.AutoMigrate(db)

	// product
	productRepo := adapters.NewGormProductRepository(db)
	productService := core.NewProductService(productRepo)
	productHandler := adapters.NewHttpProductHandler(productService)

	// auth
	userRepo := authadapter.NewGormUserRepository(db)
	userService := authcore.NewUserService(userRepo)
	userHandler := authadapter.NewHttpUserHandler(userService)

	// warehouse
	warehouseRepo := adapters.NewGormWarehouseRepository(db)
	warehouseService := core.NewWarehouseService(warehouseRepo)
	warehouseHandler := adapters.NewHttpWarehouseHandler(warehouseService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	authadapter.RegisterAuthMiddleware(app)
	productHandler.RegisterRoutes(app)
	userHandler.RegisterAuthRoutes(app)
	warehouseHandler.RegisterRoutes(app)

	fmt.Println("Server starting on :8000")
	log.Fatal(app.Listen(":8000"))
}
