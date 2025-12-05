package authadapter

import (
	"os"
	"time"
	authcore "warehouse-app/auth/authCore"
	"warehouse-app/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// HTTP-обработчик для пользователей
type httpUserHandler struct {
	// Поле service — это интерфейс UserService
	service authcore.UserService
}

// Создаёт новый HTTP-обработчик пользователей
func NewHttpUserHandler(service authcore.UserService) *httpUserHandler {
	// Возвращаем указатель на новый экземпляр обработчика
	return &httpUserHandler{service: service}
}

// Регистрирует нового пользователя
func (h *httpUserHandler) CreateUserFiber(c *fiber.Ctx) error {
	// Создаём новую переменную пользователя
	user := new(database.User)
	// Парсим тело запроса в структуру пользователя
	if err := c.BodyParser(user); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Вызываем метод сервиса для создания пользователя
	if err := h.service.CreateUser(user); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 201 и сообщение об успешной регистрации
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Register Successful"})
}

// Авторизует пользователя: возвращает JWT-токен
func (h *httpUserHandler) LoginUserFiber(c *fiber.Ctx) error {
	// Создаём новую переменную пользователя
	user := new(database.User)
	// Парсим тело запроса в структуру пользователя
	if err := c.BodyParser(user); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Вызываем метод сервиса для входа (возвращает токен и ошибку)
	token, err := h.service.LoginUser(user)
	// Если ошибка — возвращаем 401 и сообщение
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}

	// Устанавливаем токен в куки
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",                          // Название куки
		Value:    token,                          // Значение — токен
		Expires:  time.Now().Add(time.Hour * 72), // Время жизни — 72 часа
		HTTPOnly: true,                           // Кука недоступна из JS
	})
	// Возвращаем 200 и сообщение об успешном входе
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login Successful"})
}

// Middleware: проверяет JWT-токен в куках
func AuthRequired(c *fiber.Ctx) error {
	// Читаем JWT-токен из куки
	cookie := c.Cookies("jwt")
	// Читаем секретный ключ из переменной окружения
	jwtSecretKey := os.Getenv("JWT_SECRETKEY")
	// Парсим токен с помощью секретного ключа
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Возвращаем ключ
		return []byte(jwtSecretKey), nil
	})
	// Если ошибка или токен невалиден — возвращаем 401
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	// Иначе — продолжаем выполнение
	return c.Next()
}
