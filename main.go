package main

import (
	"fmt"
	"log"
	"os"

	"warehouse-app/adapters"
	authadapter "warehouse-app/auth/authAdapter"
	authcore "warehouse-app/auth/authCore"
	"warehouse-app/core"
	"warehouse-app/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")      // default PostgreSQL port
	user := os.Getenv("DB_USER")         // as defined in docker-compose.yml
	password := os.Getenv("DB_PASSWORD") // as defined in docker-compose.yml
	dbname := os.Getenv("DB_NAME")       // as defined in docker-compose.yml

	// Формируем строку подключения к БД
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, portStr, user, password, dbname)

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Создаём репозитории, сервисы, обработчики
	productRepo := adapters.NewGormProductRepository(db)
	productService := core.NewProductService(productRepo)
	productHandler := adapters.NewHttpProductHandler(productService)

	userRepo := authadapter.NewGormUserRepository(db)
	userService := authcore.NewUserService(userRepo)
	userHandler := authadapter.NewHttpUserHandler(userService)

	// Автоматически создаём таблицы в БД
	db.AutoMigrate(&database.Product{}, &database.Category{}, &database.Supplier{}, &database.ProductCategory{}, &database.User{})
	fmt.Println("Automigrate Successful")

	// Создаём HTTP-приложение
	app := fiber.New()

	// Защищаем маршруты авторизацией
	app.Use("/product", authadapter.AuthRequired)
	app.Use("/category", authadapter.AuthRequired)
	app.Use("/supplier", authadapter.AuthRequired)

	// Регистрируем маршруты
	app.Post("/supplier", productHandler.CreateSupplierFiber)
	app.Post("/category", productHandler.CreateCategoryFiber)
	app.Get("/supplier", productHandler.GetAllSupplierFiber)
	app.Get("/category", productHandler.GetAllCategoryFiber)
	app.Post("/product", productHandler.CreateProductFiber)
	app.Get("/product/:id", productHandler.GetProductByIDFiber)
	app.Get("/product", productHandler.GetAllProductFiber)
	app.Get("/category/:id/product", productHandler.GetAllProductOfCategoryFiber)
	app.Get("/supplier/:id/product", productHandler.GetAllProductOfSupplierFiber)
	app.Put("/supplier/:id", productHandler.UpdateSupplierFiber)
	app.Put("/product/:id", productHandler.UpdateProductByIDFiber)
	app.Delete("/product/:id", productHandler.DeleteProductByIDFiber)
	app.Delete("/supplier/:id", productHandler.DeleteSupplierByIDFiber)
	app.Delete("/category/:id", productHandler.DeleteCategoryByIDFiber)
	app.Post("/register", userHandler.CreateUserFiber)
	app.Post("/login", userHandler.LoginUserFiber)

	app.Listen(":8000")
}
