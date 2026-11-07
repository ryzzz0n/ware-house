package authadapter

import (
	"os"
	"strings"
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

func (h *httpUserHandler) RegisterAuthRoutes(app *fiber.App) {
	app.Post("/register", h.CreateUserFiber)
	app.Post("/login", h.LoginUserFiber)
}

// Все маршруты, кроме /register и /login, требуют JWT.
func RegisterAuthMiddleware(app *fiber.App) {
	protected := []string{
		"/product", "/category", "/supplier",
		"/zone", "/cell", "/stock",
		"/receipt", "/shipment", "/movement",
		"/inventory", "/processing",
	}
	for _, prefix := range protected {
		app.Use(prefix, AuthRequired)
	}
}

// credentials — тело запроса для /register и /login.
// Отдельная структура нужна, потому что у database.User поле Password
// помечено `json:"-"` (чтобы хеш не попадал в ответы API),
// и при парсинге напрямую в User пароль терялся.
type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func parseCredentials(c *fiber.Ctx) (*database.User, error) {
	var req credentials
	if err := c.BodyParser(&req); err != nil {
		return nil, err
	}
	return &database.User{Email: req.Email, Password: req.Password}, nil
}

func (h *httpUserHandler) CreateUserFiber(c *fiber.Ctx) error {
	user, err := parseCredentials(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}
	if err := h.service.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Register Successful"})
}

func (h *httpUserHandler) LoginUserFiber(c *fiber.Ctx) error {
	user, err := parseCredentials(c)
	if err != nil {
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login Successful",
		"token":   token,
	})
}

func AuthRequired(c *fiber.Ctx) error {
	tokenStr := c.Cookies("jwt")

	if tokenStr == "" {
		auth := c.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}
	}

	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	jwtSecretKey := os.Getenv("JWT_SECRETKEY")
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.Next()
}
