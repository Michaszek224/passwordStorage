package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func vaultHandler(ctx *gin.Context, db *sql.DB) {
	ctx.HTML(http.StatusOK, "vault.html", nil)
}
