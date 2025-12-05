package authcore

import (
	"warehouse-app/database"
)

// Интерфейс сервиса пользователей
type UserService interface {
	CreateUser(user *database.User) error
	LoginUser(user *database.User) (string, error)
}

// Реализация сервиса пользователей
type userServiceImpl struct {
	// Поле repo — это репозиторий пользователей
	repo UserRepository
}

// Создаёт новый экземпляр UserService
func NewUserService(repo UserRepository) UserService {
	// Возвращаем указатель на новую структуру с репозиторием
	return &userServiceImpl{repo: repo}
}

// Вызывает репозиторий для создания пользователя
func (s *userServiceImpl) CreateUser(user *database.User) error {
	// Вызываем метод репозитория
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для аутентификации пользователя
func (s *userServiceImpl) LoginUser(user *database.User) (string, error) {
	// Вызываем метод репозитория — возвращает токен и ошибку
	t, err := s.repo.LoginUser(user)
	// Если ошибка — возвращаем пустую строку и ошибку
	if err != nil {
		return "", err
	}
	// Иначе — возвращаем токен и nil (успешно)
	return t, nil
}
