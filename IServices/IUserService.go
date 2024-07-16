package IServices

import "gelio/m/models"

type IUserService interface {
	GetUserWithId(userId interface{}) (models.User, error)
	GetUserWithName(userName string) (models.User, error)
	CreateUser(user models.User) (int, error)
	UpdateUser(Username string, ProfileImageId int, UserId int) error
	DeleteUser(userId interface{}) error
	NewUser(username string, hash string, createdDate string, isActive bool, profileImageId int, personId int) models.User
}
