package authadapter

import (
	"os"
	"time"

	authcore "warehouse-app/auth/authCore"
	"warehouse-app/database"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GORM-реализация репозитория пользователей
type gormUserRepository struct {
	db *gorm.DB
}

// Создаёт новый экземпляр GORM-репозитория пользователей
func NewGormUserRepository(db *gorm.DB) authcore.UserRepository {
	// Возвращаем указатель на новую структуру с подключением к БД
	return &gormUserRepository{db: db}
}

// Создаёт пользователя: хеширует пароль и сохраняет в БД
func (r *gormUserRepository) CreateUser(user *database.User) error {
	// Хешируем пароль с помощью bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// Если ошибка — возвращаем её
	if err != nil {
		return err
	}
	// Присваиваем захешированный пароль пользователю
	user.Password = string(hashedPassword)
	// Вызываем метод Create GORM — создаём пользователя в БД
	if result := r.db.Create(user); result != nil {
		// Если была ошибка — возвращаем её
		return result.Error
	}
	return nil
}

// Аутентифицирует пользователя: проверяет email и пароль, возвращает JWT-токен
func (r *gormUserRepository) LoginUser(user *database.User) (string, error) {
	// Создаём новую переменную пользователя
	selectedUser := new(database.User)
	// Ищем пользователя в БД по email
	if result := r.db.Where("email = ?", user.Email).First(selectedUser); result.Error != nil {
		// Если не найден — возвращаем ошибку
		return "", result.Error
	}
	// Сравниваем введённый пароль с захешированным
	err := bcrypt.CompareHashAndPassword([]byte(selectedUser.Password), []byte(user.Password))
	// Если ошибка — возвращаем её
	if err != nil {
		return "", err
	}
	// Читаем секретный ключ из переменной окружения
	jwtSecretKey := os.Getenv("JWT_SECRETKEY")
	// Создаём новый JWT-токен с алгоритмом HS256
	token := jwt.New(jwt.SigningMethodHS256)

	// Получаем claims — данные внутри токена
	claims := token.Claims.(jwt.MapClaims)
	// Устанавливаем ID пользователя
	claims["user_id"] = selectedUser.ID
	// Устанавливаем время истечения токена (72 часа)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Подписываем токен с помощью секретного ключа
	t, err := token.SignedString([]byte(jwtSecretKey))
	// Если ошибка — возвращаем её
	if err != nil {
		return "", err
	}
	// Возвращаем токен и nil (успешно)
	return t, nil
}
