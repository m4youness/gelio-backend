package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	c.Bind(&body)

	var User models.User

	err := initializers.DB.Get(&User, "select * from users where username = $1", body.Username)

	if err != nil {
		fmt.Println(err)
		return
	}

	Err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(body.Password))

	if Err != nil {
		fmt.Println(Err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": User.UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		log.Fatal(err)
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenstring, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"Successful": true})

}

func IsLoggedIn(c *gin.Context) {

	c.JSON(200, true)
}

func SignUp(c *gin.Context) {
	var body struct {
		UserName     string
		Password     string
		CreatedDate  string
		IsActive     bool
		ProfileImage string
		PersonID     int
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

	Res := initializers.DB.QueryRow("insert into users (username, password, created_date, is_active, profile_image, person_id) values ($1, $2, $3, $4, $5, $6) RETURNING user_id",
		body.UserName, hash, body.CreatedDate, body.IsActive, body.ProfileImage, body.PersonID)

	var userID int
	Err := Res.Scan(&userID)

	if Err != nil {
		fmt.Println(Err)
		return
	}

	c.JSON(200, userID)
}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"LoggedOut": true})
}

func GetUserId(c *gin.Context) {
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

		if user.UserID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.JSON(200, user.UserID)

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

	}

}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	err := initializers.DB.Get(&user, "select * from users where user_id = $1", id)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, user)
}

func DoesUserExist(c *gin.Context) {
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
