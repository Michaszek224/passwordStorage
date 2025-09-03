package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func storeOptionsConfig(store cookie.Store) {
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   300,
		Secure:   false, // change in produciton
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func RoutesHandler(db *sql.DB) *gin.Engine {
	_ = godotenv.Load()

	cookieSecret := os.Getenv("SECRET_COOKIE")
	if cookieSecret == "" {
		log.Fatal("secret cookie not set")
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte(cookieSecret))
	storeOptionsConfig(store)
	router.Use(sessions.Sessions("mysession", store))

	//unprotected get
	router.GET("/", homeHandlerGet)
	router.GET("/register", registerHandlerGet)
	router.GET("/login", func(ctx *gin.Context) { loginHandlerGet(ctx, db) })

	//protected rotues
	authorized := router.Group("/vault")
	authorized.Use(authRequired())
	{
		authorized.GET("/", func(ctx *gin.Context) { vaultHandler(ctx, db) })
		authorized.POST("/addSite", func(ctx *gin.Context) { addSite(ctx, db) })
		authorized.POST("/editSite", func(ctx *gin.Context) { editSite(ctx, db) })
		authorized.POST("/deleteSite", func(ctx *gin.Context) { deleteSite(ctx, db) })
	}

	//unprotected post
	router.POST("/login", func(ctx *gin.Context) { loginHandlerPost(ctx, db) })
	router.POST("/register", func(ctx *gin.Context) { registerHandlerPost(ctx, db) })

	return router
}
