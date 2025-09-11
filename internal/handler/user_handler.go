package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iragsraghu/user-service/internal/models"
	"github.com/iragsraghu/user-service/internal/repo"
)

type UserHandler struct {
	Repo *repo.UserRepo
}

func NewUserHandler(r *repo.UserRepo) *UserHandler {
	return &UserHandler{Repo: r}
}

func (h *UserHandler) Create(c *gin.Context) {
	var in struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &models.User{Name: in.Name, Email: in.Email}
	if err := h.Repo.Create(c.Request.Context(), u); err != nil {
		// ✅ check for duplicate email
		if errors.Is(err, repo.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}

		// fallback to generic 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	u, err := h.Repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.Repo.List(c.Request.Context(), 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, users)
}
