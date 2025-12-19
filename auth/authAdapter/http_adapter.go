// auth/authAdapter/http_adapter.go
package authadapter

import (
	"os"
	"time"

	authcore "warehouse-app/auth/authCore"
	"warehouse-app/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type httpUserHandler struct {
	service authcore.UserService
}

func NewHttpUserHandler(service authcore.UserService) *httpUserHandler {
	return &httpUserHandler{service: service}
}

// --- НОВОЕ: Регистрация маршрутов ---

func (h *httpUserHandler) RegisterAuthRoutes(app *fiber.App) {
	app.Post("/register", h.CreateUserFiber)
	app.Post("/login", h.LoginUserFiber)
}

// --- НОВОЕ: Middleware ---

func RegisterAuthMiddleware(app *fiber.App) {
	app.Use("/product", AuthRequired)
	app.Use("/category", AuthRequired)
	app.Use("/supplier", AuthRequired)
}

// --- Обработчики (без изменений, но для полноты) ---

func (h *httpUserHandler) CreateUserFiber(c *fiber.Ctx) error {
	user := new(database.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Register Successful"})
}

func (h *httpUserHandler) LoginUserFiber(c *fiber.Ctx) error {
	user := new(database.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	token, err := h.service.LoginUser(user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login Successful"})
}

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := os.Getenv("JWT_SECRETKEY")
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}
