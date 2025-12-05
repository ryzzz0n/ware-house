package authcore

import (
	"warehouse-app/database"
)

type UserRepository interface {
	CreateUser(user *database.User) error
	LoginUser(user *database.User) (string, error)
}
