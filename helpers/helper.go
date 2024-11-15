package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kanhaiyagupta9045/car_management/databases"
	"github.com/kanhaiyagupta9045/car_management/models"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Id int
	jwt.StandardClaims
}

func GenerateAccessToken(id int) (string, error) {
	ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
	if ACCESS_TOKEN_SECRET == "" {
		return "", errors.New("access token secret is not set")
	}
	accessClaims := &SignedDetails{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", err
	}

	return token, nil
}
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

	check := true

	if err != nil {
		return false
	}
	return check
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("ACCESS_TOKEN_SECRET")
		if secretKey == "" {
			return nil, errors.New("SECRET_KEY not set in environment")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims or token is not valid")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
func GetUserFromCookie(c *gin.Context) (*models.User, error) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("token not found in the cookie")
	}
	claims, err := ValidateToken(token)
	if err != nil {
		return nil, err
	}
	query := "SELECT user_id, first_name, last_name, email, password FROM users WHERE user_id = ?"
	var user models.User
	err = databases.DB.QueryRow(query, claims.Id).Scan(
		&user.User_Id,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
