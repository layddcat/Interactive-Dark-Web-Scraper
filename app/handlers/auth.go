package handlers

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetLogin(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user") != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func PostLogin(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "cti123" {
		session.Set("user", username)
		session.Set("role", "Senior CTI Analyst") 
		session.Save()
		
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"error": "Geçersiz analist kimlik bilgileri. Erişim reddedildi.",
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}