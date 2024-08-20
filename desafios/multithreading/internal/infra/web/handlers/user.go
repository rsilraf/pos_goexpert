package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/dto"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/entity"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/infra/db"
)

type UserHandler struct {
	DAO db.UserDAOInterface
}

func NewUserHandler(dao db.UserDAOInterface) *UserHandler {
	return &UserHandler{
		DAO: dao,
	}
}

// GetToken		godoc
// @Summary		Get a JWT
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		request		 body 		dto.GetTokenInput		true	"token request"
// @Success		200			{object}	dto.GetTokenOutput
// @Failure		404			{object}	Error
// @Failure		500			{object}	Error
// @Router		/token		[post]
func (h *UserHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	// context
	jwt := r.Context().Value("token").(*jwtauth.JWTAuth)
	expiresIn := r.Context().Value("TTL").(int)

	// input
	var input dto.GetTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// find user
	user, err := h.DAO.FindByEmail(input.Email)
	if err != nil || !user.ValidatePassword(input.Password) {
		sendError(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	// token string
	_, token, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(expiresIn)).Unix(),
	})

	// token output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&dto.GetTokenOutput{Token: token})

}

// Create user	godoc
// @Summary		Create user
// @Tags		user
// @Accept		json
// @Param		request	body	dto.CreateUserInput true	"user request"
// @Success		201
// @Failure		400		{object}	Error
// @Failure		404		{object}	Error
// @Failure		500		{object}	Error
// @Router		/users	[post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// find user
	user, _ := h.DAO.FindByEmail(input.Email)
	if user != nil {
		sendError(w, http.StatusBadRequest, errors.New("user already exists"))
		return
	}

	user, err = entity.NewUser(input.Email, input.Password)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	if h.DAO.Create(user) != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
