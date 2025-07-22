package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var superusers = gin.Accounts{
	"admin": "password",
}

func basicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "autentication required"})
			return
		}

		pwd, ok := superusers[user]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		if pwd != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// TODO: search through db by name. If ok -> add to hash

		c.Set("authUser", user)
		c.Next()
	}
}
