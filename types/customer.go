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

type CreateCustomer struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CreateCustomer) CryptPassword() error {
	cryptedPS, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		return err
	}
	c.Password = string(cryptedPS)
	return nil
}
