package services

import (
	"gelio/m/initializers"
	"gelio/m/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (UserService) GetUserWithId(userId interface{}) (models.User, error) {
	var User models.User
	err := initializers.DB.Get(&User, "select * from users where user_id = $1", userId)

	if err != nil {
		return models.User{}, err
	}

	return User, nil

}

func (UserService) GetUserWithName(userName string) (models.User, error) {
	var User models.User
	err := initializers.DB.Get(&User, "select * from users where username = $1", userName)

	if err != nil {
		return models.User{}, err
	}

	return User, nil

}
