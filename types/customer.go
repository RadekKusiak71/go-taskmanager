package types

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *RegisterRequest) CryptPassword() error {
	cryptedPS, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		return err
	}
	c.Password = string(cryptedPS)
	return nil
}
