package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/renamrgb/go-expert-apis/internal/dto"
	"github.com/renamrgb/go-expert-apis/internal/entity"
	"github.com/renamrgb/go-expert-apis/internal/infra/database"
)

type Error struct {
	Message string `json:"message"`
} // @name ErrorResponse

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(userDB database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{UserDB: userDB, Jwt: jwt, JwtExpiresIn: jwtExpiresIn}
}

// GetJWT godoc
// @Summary GET a user JWT
// @Description GET a user JWT
// @Tags users
// @Accept json
// @Produce json
// @param   request body dto.GetJWTInput true "user credentials"
// @Success 200 {object} dto.GetJWTOutput
// @failure 404 {object} Error
// @failure 500 {object} Error
// @Router /user/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var input dto.GetJWTInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(input.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	if !user.ValidatePassword(input.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.Jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
	w.WriteHeader(http.StatusOK)
}

// Create User
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @param   request body dto.CreateUserInput true "user request"
// @Success 201
// @Failure 500 {object} Error
// @Router /user [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userInput dto.CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
