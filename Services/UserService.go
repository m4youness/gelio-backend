package services

import (
	"gelio/m/initializers"
	"gelio/m/models"
)

type UserService struct{}

func (UserService) NewUser(username string, hash string, createdDate string, isActive bool, profileImageId int, personId int) models.User {
	return models.User{
		Username:       username,
		Password:       hash,
		CreatedDate:    createdDate,
		IsActive:       isActive,
		ProfileImageId: &profileImageId,
		PersonID:       personId,
	}
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

func (UserService) CreateUser(user models.User) (int, error) {
	res := initializers.DB.QueryRow("insert into users (username, password, created_date, is_active, profile_image_id, person_id) values ($1, $2, $3, $4, $5, $6) RETURNING user_id",
		user.Username, user.Password, user.CreatedDate, user.IsActive, user.ProfileImageId, user.PersonID)

	var userId int
	err := res.Scan(&userId)

	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (UserService) UpdateUser(Username string, ProfileImageId int, UserId int) error {
	_, err := initializers.DB.Exec("update users set username = $1, profile_image_id = $2 where user_id = $3", Username, ProfileImageId, UserId)

	if err != nil {
		return err
	}

	return nil
}

func (UserService) DeleteUser(userId interface{}) error {
	_, err := initializers.DB.Exec("update users set is_active = false where user_id = $1", userId)

	if err != nil {
		return err
	}

	return nil
}
