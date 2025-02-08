package model

import (
	"net/url"

	"github.com/golang-jwt/jwt"
)

type UserDetails struct {
	Name string `json:"name"`

	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Verified bool   `json:"verified"`
	UserType string `json:"user_type"`
}
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AdminClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type AdminOtp struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func (u *UserDetails) Valid() url.Values {
	err := url.Values{}

	if len(u.Password) < 6 {
		err.Add("password", "password must be greater than 6")
	}

	return err
}
