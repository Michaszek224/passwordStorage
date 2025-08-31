package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"passwordStorage/database"

	"github.com/gin-gonic/gin"
)

func RoutesHandler(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) { indexPageHandler(ctx, db) })
	router.POST("/submit", func(ctx *gin.Context) { submitHandler(ctx, db) })

	return router
}

func indexPageHandler(ctx *gin.Context, db *sql.DB) {
	users := database.GetUsers(db)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"users": users,
	})
}

func submitHandler(ctx *gin.Context, db *sql.DB) {
	nickname := ctx.PostForm("nickname")
	password := ctx.PostForm("password")

	if nickname == "" || password == "" {
		users := database.GetUsers(db)
		ctx.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Nickname nad password are required",
			"users": users,
		})
		return
	}

	_, err := db.Exec(`INSERT INTO user(nickname, password) VALUES (?,?)`, nickname, password)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		users := database.GetUsers(db)
		ctx.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Could not save data",
			"users": users,
		})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
