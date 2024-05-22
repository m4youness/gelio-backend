package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
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
		log.Fatal(err)
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

	c.JSON(200, gin.H{
		"Authorized": true,
	},
	)
}

func SignUp(c *gin.Context) {
	var body struct {
		UserName     string `json:"UserName"`
		Password     string `json:"Password"`
		CreatedDate  string `json:"CreatedDate"`
		IsActive     bool   `json:"IsActive"`
		ProfileImage string `json:"ProfileImage"`
		PersonID     int    `json:"PersonID"`
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

	c.JSON(200, gin.H{
		"UserID": userID,
	},
	)
}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"LoggedOut": true})
}
