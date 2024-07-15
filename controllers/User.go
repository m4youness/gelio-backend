package controllers

import (
	"encoding/json"
	"fmt"
	util "gelio/m/Util"
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct{}

func UserController() *User {
	return &User{}
}

func (User) SignIn(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	c.Bind(&body)

	var User models.User

	err := initializers.DB.Get(&User, "select * from users where username = $1", body.Username)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, false)
		return
	}

	Err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(body.Password))

	if Err != nil {
		fmt.Println(Err)
		c.JSON(http.StatusOK, false)
		return
	}

	accessToken, err := util.CreateAccessToken(User.UserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create access token"})
		return
	}

	refreshToken, err := util.CreateRefreshToken(User.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create refresh token"})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", accessToken, 15*60, "/", "glistening-respect-production.up.railway.app", true, true)
	c.SetCookie("RefreshToken", refreshToken, 7*24*60*60, "/", "glistening-respect-production.up.railway.app", true, true)

	c.JSON(http.StatusOK, true)

}

func (User) Register(c *gin.Context) {
	var body struct {
		UserName       string
		Password       string
		CreatedDate    string
		IsActive       bool
		ProfileImageId int
		PersonID       int
	}
	Error := c.Bind(&body)

	if Error != nil {
		fmt.Println(Error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		fmt.Println(err)
		return
	}

	Res := initializers.DB.QueryRow("insert into users (username, password, created_date, is_active, profile_image_id, person_id) values ($1, $2, $3, $4, $5, $6) RETURNING user_id",
		body.UserName, hash, body.CreatedDate, body.IsActive, body.ProfileImageId, body.PersonID)

	var userID int
	Err := Res.Scan(&userID)

	if Err != nil {
		fmt.Println(Err)
		return
	}

	c.JSON(200, userID)
}

func (User) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie("Authorization", "", -1, "/", "glistening-respect-production.up.railway.app", true, true)
	c.SetCookie("RefreshToken", "", -1, "/", "glistening-respect-production.up.railway.app", true, true)

	c.JSON(http.StatusOK, gin.H{"LoggedOut": true})
}

func (User) GetUserId(c *gin.Context) {
	tokenstring, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User

		initializers.DB.Get(&user, "select * from users where user_id = $1", claims["sub"])

		if user.UserId == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.JSON(200, user.UserId)

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

	}

}

func (User) GetUser(c *gin.Context) {
	id := c.Param("id")

	cachedUser, err := initializers.RedisClient.Get(initializers.Ctx, fmt.Sprintf("user:%s", id)).Result()

	if err == nil {
		var user models.User
		json.Unmarshal([]byte(cachedUser), &user)
		c.JSON(200, user)
		return
	}

	var user models.User

	err = initializers.DB.Get(&user, "select * from users where user_id = $1", id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	userData, _ := json.Marshal(user)
	err = initializers.RedisClient.Set(initializers.Ctx, fmt.Sprintf("user:%s", id), userData, time.Hour).Err()
	if err != nil {
		fmt.Println("Failed to cache user data:", err)
	}

	c.JSON(200, user)
}

func (User) DoesUserExist(c *gin.Context) {
	var body struct {
		UserName string
	}

	err := c.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var user models.User

	Err := initializers.DB.Get(&user, "select user_id from users where username = $1", body.UserName)

	if Err != nil {
		c.JSON(200, false)
		return
	}

	c.JSON(200, true)

}

func (User) MakeUserInActive(c *gin.Context) {
	id := c.Param("id")

	fmt.Println(id)

	_, err := initializers.DB.Exec("update users set is_active = false where user_id = $1", id)

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, true)

}

func (User) UserActivity(c *gin.Context) {
	username := c.Param("username")

	var User models.User

	err := initializers.DB.Get(&User, "select * from users where username = $1", username)

	// error means he is logged in
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// user is not logged in
	c.JSON(200, User.IsActive)

}

func (User) UpdateUser(c *gin.Context) {
	var body struct {
		Firstname      string
		Lastname       string
		Username       string
		Email          string
		Phonenumber    string
		CountryId      int
		GenderId       int
		ProfileImageId int
		UserId         int
		PersonId       int
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := initializers.DB.Exec("update people set first_name = $1, last_name = $2, email = $3, phone_number = $4, country_id = $5, gender_id = $6 where person_id = $7", body.Firstname, body.Lastname, body.Email, body.Phonenumber, body.CountryId, body.GenderId, body.PersonId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = initializers.DB.Exec("update users set username = $1, profile_image_id = $2 where user_id = $3")

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})

}

func (u *User) InitializeRoutes(r *gin.Engine) {
	r.POST("/Register", u.Register)
	r.POST("/SignIn", u.SignIn)
	r.GET("/Logout", middleware.RequireAuth, u.Logout)
	r.GET("/User/:id", middleware.RequireAuth, u.GetUser)
	r.POST("/User/Exists", u.DoesUserExist)
	r.GET("/User/Id", middleware.RequireAuth, u.GetUserId)
	r.GET("/User/InActive/:id", middleware.RequireAuth, u.MakeUserInActive)
	r.GET("/User/IsNotActive/:username", u.UserActivity)
	r.PUT("/User/Update", middleware.RequireAuth, u.UpdateUser)
}
