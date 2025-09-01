package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func vaultHandler(ctx *gin.Context, db *sql.DB) {
	session := sessions.Default(ctx)
	user := session.Get("user")

	ctx.HTML(http.StatusOK, "vault.html", gin.H{
		"user": user,
	})
}
