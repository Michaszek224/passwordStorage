package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"passwordStorage/database"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func registerHandlerGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", nil)
}

func registerHandlerPost(ctx *gin.Context, db *sql.DB) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	err := database.InsertUser(username, password, db)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		users := database.GetUsers(db)
		ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": err.Error(),
			"users": users,
		})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/login")
}

func loginHandlerGet(ctx *gin.Context) {
	// users := database.GetUsers(db)
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func loginHandlerPost(ctx *gin.Context, db *sql.DB) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	user, err := database.AuthenicateUser(username, password, db)
	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	session := sessions.Default(ctx)
	session.Set("user", username)
	session.Set("userId", user.ID)
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/vault")
}

func homeHandlerGet(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")

	if user == nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/vault")
}

func authRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get("user")

		if user == nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
