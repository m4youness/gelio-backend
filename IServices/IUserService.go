package IServices

import "gelio/m/models"

type IUserService interface {
	GetUserWithId(userId interface{}) (models.User, error)
	GetUserWithName(userName string) (models.User, error)
}
