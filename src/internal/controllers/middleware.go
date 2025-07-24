package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lekht/account-master/src/pkg/storage/mock"
)

func (r *Router) basicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			return
		}

		user, err := r.repo.UserByName(username)
		if errors.Is(err, mock.ErrNoUsername) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// TODO: hash
		if user.Password != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}

		c.Set("username", user.Username)
		c.Set("isAdmin", user.Admin)

		c.Next()
	}
}

func isAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("isAdminMiddleware")
		isAdmin, ok := c.Get("isAdmin")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		log.Println("is admin: ", isAdmin.(bool))
		if !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		c.Next()
	}
}
