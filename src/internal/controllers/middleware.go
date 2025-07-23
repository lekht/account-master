package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var superusers = gin.Accounts{
	"su": "pwd",
}

func basicAuthMiddleware() gin.HandlerFunc {
	fmt.Println("BasicAuth middleware triggered!")
	return func(c *gin.Context) {
		fmt.Println("BasicAuth HandlerFunc!")
		user, password, ok := c.Request.BasicAuth()
		log.Printf("basic auth: name=%s pdw=%s ok=%t", user, password, ok)
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
			log.Printf("basic auth: name=%s password=%s pwd=%s", user, password, pwd)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// TODO: search through db by name. If ok -> add to hash

		// c.Set("authUser", user)
		c.Next()
	}
}
