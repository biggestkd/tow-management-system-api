package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tow-management-system-api/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserById(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, userId *string, update *model.User) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// PostUser POST /user
// Request: { "id": "...", "email": "..." }
// Response: 201 Empty body | 400 generic error text
func (h *UserHandler) PostUser(context *gin.Context) {
	log.Println("Running PostCompany")
	var body model.User
	if err := context.ShouldBindJSON(&body); err != nil {
		context.String(http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.CreateUser(context, &body); err != nil {
		log.Println(err)
		context.String(http.StatusBadRequest, "Something went wrong")
		return
	}

	context.Status(http.StatusCreated)
}

// GetUser GET /user/:userId
// Response: 200 User | 404 generic error text
func (h *UserHandler) GetUser(context *gin.Context) {
	userID := context.Param("userId")
	if userID == "" {
		context.String(http.StatusNotFound, "Something went wrong")
		return
	}

	var user = &model.User{
		ID: &userID,
	}

	u, err := h.userService.FindUserById(context, user)
	if err != nil || u == nil {
		context.String(http.StatusNotFound, "Something went wrong")
		return
	}

	context.JSON(http.StatusOK, u)
}

// PutUser handles: PUT /user/:userId
// Request BODY: partial User (only fields to change)
// Response: 204 | 400 invalid request | 404 not found (up to service to signal)
func (h *UserHandler) PutUser(ctx *gin.Context) {
	userID := ctx.Param("userId")

	var body model.User

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.String(http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.userService.UpdateUser(ctx, &userID, &body); err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusBadRequest, "Something went wrong")
		return
	}

	ctx.Status(http.StatusNoContent)
}
