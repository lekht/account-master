package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lekht/account-master/src/internal/model"
	"github.com/lekht/account-master/src/pkg/storage/mock"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// TODO: create by name search for auth
type Repository interface {
	Users() ([]model.Profile, error)
	UserByID(int) (model.Profile, error)
	CreateUser(model.Profile) (int, error)
	UpdateUser(int, model.Profile) error
	DeleteUser(int) error
}

type Router struct {
	repo Repository

	router *gin.Engine
}

func New(repo Repository) *Router {
	r := Router{
		repo:   repo,
		router: gin.New(),
	}

	r.router.Use(gin.Logger())
	r.router.Use(gin.Recovery())

	r.router.GET("/user", r.getUsers)
	r.router.GET("/user/:id", r.getUserById)
	r.router.POST("/user", basicAuthMiddleware(), r.createUser)
	r.router.PUT("/user/:id", basicAuthMiddleware(), r.updateUserById)
	r.router.DELETE("/user/:id", basicAuthMiddleware(), r.deleteUserById)

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &r
}

func (r *Router) GetRouter() *gin.Engine {
	return r.router
}

// createUser
// @Summary Create User
// @Description Create new user
// @Security BasicAuth
// @Accept  json
// @Produce  json
// @Param user body model.Profile true "Email, Username, Password, Admin"
// @Success 201
// @Failure 400
// @Failure 409
// @Header       all {string} string "header"
// @Router /user [post]
func (r *Router) createUser(c *gin.Context) {
	var usr model.Profile

	if err := c.ShouldBindJSON(&usr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong json"})
		return
	}

	id, err := r.repo.CreateUser(usr)
	if err != nil {
		if errors.Is(err, mock.ERROR_USER_EXISTS) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// getUsers
// @Summary Get Users
// @Description Get full users list
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 404
// @Header       all {string} string "header"
// @Router /user [get]
func (r *Router) getUsers(c *gin.Context) {
	users, err := r.repo.Users()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "users not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// getUserById
// @Summary Get User By ID
// @Description Get specific user by ID
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Header       all {string} string "header"
// @Router /user/{id} [get]
func (r *Router) getUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	u, err := r.repo.UserByID(id)
	if errors.Is(err, mock.ERROR_NO_USER_ID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server errror"})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       u.Id,
		"username": u.Username,
		"email":    u.Email,
		"admin":    u.Admin,
	})
}

// updateUserById()
// @Summary Update User
// @Description Update user by ID
// @Security BasicAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body model.Profile true "at least one field is reqired"
// @Success 200
// @Failure 400
// @Failure 404
// @Header       all {string} string "header"
// @Router /user/{id} [put]
func (r *Router) updateUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var u model.Profile
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err = r.repo.UpdateUser(id, u)
	if errors.Is(err, mock.ERROR_NO_USER_ID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

// deleteUserById()
// @Summary Delete User
// @Description Delete user by ID
// @Security BasicAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Header       all {string} string "header"
// @Router /user/{id} [delete]
func (r *Router) deleteUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = r.repo.DeleteUser(id)
	if errors.Is(err, mock.ERROR_NO_USER_ID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
