package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"interactive-scraper/database"
	"interactive-scraper/handlers"
	"interactive-scraper/services"
	"net/http"
)

func main() {
	database.InitDB()

	go services.StartRealMonitoring()
	r := gin.Default()

	store := cookie.NewStore([]byte("cti-secret-key"))
	r.Use(sessions.Sessions("ctisession", store))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/login", handlers.GetLogin)
	r.POST("/login", handlers.PostLogin)

	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/dashboard")
		})

		authorized.GET("/dashboard", handlers.GetDashboard)
		authorized.GET("/detail/:id", handlers.GetDataDetail)
		authorized.POST("/update-criticality", handlers.UpdateCriticality)
		authorized.GET("/logout", handlers.Logout)
	}

	r.Run(":8080")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}