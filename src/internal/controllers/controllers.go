package controllers

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	g := r.Group("/user")
	{
		g.POST("", CreateUser)
		g.GET("", GetUsers)
		g.GET("/:id", GetUserById)
		g.PUT("/:id", UpdateUserById)
		g.DELETE("/:id", DeleteUserById)
	}

	return r
}

func CreateUser(c *gin.Context) {
}

func GetUsers(c *gin.Context) {
}

func GetUserById(c *gin.Context) {
}

func UpdateUserById(c *gin.Context) {
}

func DeleteUserById(c *gin.Context) {
}
