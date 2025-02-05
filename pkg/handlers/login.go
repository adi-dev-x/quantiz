package handlers

type Register struct {
	Name string `json:"name" validate:"required"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}
type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type Otp struct {
	Email string `json:"email" validate:"required"`
	Otp   string `json:"otp" validate:"required"`
}
