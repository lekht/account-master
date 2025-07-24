package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/lekht/account-master/src/internal/model"
	"github.com/lekht/account-master/src/pkg/storage/mock"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Repository interface {
	Users() ([]model.Profile, error)
	UserByID(uuid.UUID) (model.Profile, error)
	CreateUser(model.Profile) error
	UpdateUser(uuid.UUID, model.Profile) error
	DeleteUser(uuid.UUID) error
	UserByName(string) (model.Profile, error)
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

	authenticated := r.router.Group("/user", r.basicAuthMiddleware())
	{
		authenticated.GET("", r.getUsers)
		authenticated.GET("/:id", r.getUserById)
		authenticated.POST("", isAdminMiddleware(), r.createUser)
		authenticated.PUT("/:id", isAdminMiddleware(), r.updateUserById)
		authenticated.DELETE("/:id", isAdminMiddleware(), r.deleteUserById)
	}

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &r
}

func (r *Router) Router() *gin.Engine {
	return r.router
}

// createUser
//
//	@Summary		Create User
//	@Description	Create new user
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	AccountRequest	true	"Email, Username, Password, Admin"
//	@Success		201
//	@Failure		400
//	@Failure		409
//	@Header			all	{string}	string	"header"
//	@Router			/user [post]
func (r *Router) createUser(c *gin.Context) {
	var req AccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong json"})
		return
	}

	usr, err := requestToProfile(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong json"})
		return
	}

	err = r.repo.CreateUser(*usr)
	if err != nil {
		if errors.Is(err, mock.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	c.Status(http.StatusCreated)
}

// getUsers
//
//	@Summary		Get Users
//	@Description	Get full users list
//	@Header			all	{string}	string	"header"
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Return json
//	@Success		200
//	@Failure		404
//	@Router			/user [get]
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

	responses := make([]AccountResponse, 0, len(users))
	for _, u := range users {
		resp, err := profileToResponse(&u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		responses = append(responses, *resp)
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// getUserById
//
//	@Summary		Get User By ID
//	@Description	Get specific user by ID
//	@Header			all	{string}	string	"header"
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Return json
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Router			/user/{id} [get]
func (r *Router) getUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	u, err := r.repo.UserByID(id)
	if errors.Is(err, mock.ErrNoUserID) {
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
//
//	@Summary		Update User
//	@Description	Update user by ID
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"User ID"
//	@Param			user	body	AccountRequest	true	"request body"
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Header			all	{string}	string	"header"
//	@Router			/user/{id} [put]
func (r *Router) updateUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req AccountRequest
	if err = c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	u, err := requestToProfile(&req)

	err = r.repo.UpdateUser(id, *u)
	if errors.Is(err, mock.ErrNoUserID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

// deleteUserById()
//
//	@Summary		Delete User
//	@Description	Delete user by ID
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Header			all	{string}	string	"header"
//	@Router			/user/{id} [delete]
func (r *Router) deleteUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = r.repo.DeleteUser(id)
	if errors.Is(err, mock.ErrNoUserID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
