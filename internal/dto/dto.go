package dto

type CreateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
} // @name CreateProductRequest


type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}// @name CreateUserRequest

type GetJWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} // @name GetJWTRequest

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
} // @name GetJWTResponse
