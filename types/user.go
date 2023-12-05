package types

import (
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username string
	Password string
	Token    string
}

func CreateUser(username, password string) (*User, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &User{
		Username: username,
		Password: password,
		Token:    tokenStr,
	}, nil
}
