package dto

type GetTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type GetTokenOutput struct {
	Token string `json:"token"`
}
type CreateUserInput struct {
	GetTokenInput
}
