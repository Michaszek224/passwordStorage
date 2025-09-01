package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RoutesHandler(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) { loginHandlerGet(ctx, db) })
	router.GET("/register", registerHandlerGet)
	router.GET("/vault", func(ctx *gin.Context) { vaultHandler(ctx, db) })

	router.POST("/login", func(ctx *gin.Context) { loginHandlerPost(ctx, db) })
	router.POST("/register", func(ctx *gin.Context) { registerHandlerPost(ctx, db) })

	return router
}
