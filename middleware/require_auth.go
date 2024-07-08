package middleware

import (
	"fmt"
	util "gelio/m/Util"
	"gelio/m/initializers"
	"gelio/m/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		refreshTokenString, err := c.Cookie("RefreshToken")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("REFRESH_SECRET")), nil
		})

		if err != nil || !refreshToken.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
			var user models.User
			err := initializers.DB.Get(&user, "SELECT * FROM users WHERE user_id = $1", refreshClaims["sub"])
			if err != nil || user.UserId == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			newAccessToken, err := util.CreateAccessToken(user.UserId)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.SetCookie("Authorization", newAccessToken, 15*60, "", "", false, true)
			c.Set("user", user)
			c.Next()
			return
		}

	} else {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			fmt.Println("Error parsing token:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || float64(time.Now().Unix()) > exp {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired token"})
				return
			}

			var user models.User
			err := initializers.DB.Get(&user, "SELECT * FROM users WHERE user_id = $1", claims["sub"])
			if err != nil || user.UserId == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}

}
