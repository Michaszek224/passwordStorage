package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"passwordStorage/database"

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

	ctx.Redirect(http.StatusSeeOther, "/")
}

func loginHandlerGet(ctx *gin.Context, db *sql.DB) {
	users := database.GetUsers(db)
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"users": users,
	})
}

func loginHandlerPost(ctx *gin.Context, db *sql.DB) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	err := database.AuthenicateUser(username, password, db)
	users := database.GetUsers(db)
	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid username or password",
			"users": users,
		})
		return
	}

	ctx.HTML(http.StatusOK, "vault.html", nil)
}
