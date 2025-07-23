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
		// 1. Получаем Basic Auth данные
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			return
		}

		// 2. Ищем пользователя по имени
		user, err := r.repo.UserByName(username)
		if err != nil {
			if errors.Is(err, mock.ErrNoUsername) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid username or password",
				})
				return
			}

			log.Printf("Database error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// 3. Проверяем пароль (в реальном приложении используйте хеширование!)
		if user.Password != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}

		// 4. Сохраняем данные пользователя в контексте
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

		if isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		c.Next()
	}
}
