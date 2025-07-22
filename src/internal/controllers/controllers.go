package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lekht/account-master/src/internal/model"
	"github.com/lekht/account-master/src/pkg/storage/mock"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
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

	g := r.router.Group("/user")
	{
		g.POST("", r.createUser)
		g.GET("", r.getUsers)
	}

	// reqire basic access auth
	authGroup := g.Group("/:id", basicAuthMiddleware())
	{
		authGroup.GET("", r.getUserById)
		authGroup.PUT("", r.updateUserById)
		authGroup.DELETE("", r.deleteUserById)
	}

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &r
}

func (r *Router) GetRouter() *gin.Engine {
	return r.router
}

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
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

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
