package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"passwordStorage/database"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func vaultHandler(ctx *gin.Context, db *sql.DB) {
	session := sessions.Default(ctx)
	user := session.Get("user")
	userId := session.Get("userId").(int)

	siteData := database.GetSiteData(db, userId)

	ctx.HTML(http.StatusOK, "vault.html", gin.H{
		"user":     user,
		"siteData": siteData,
	})
}

func addNewSite(ctx *gin.Context, db *sql.DB) {
	site := ctx.PostForm("site")
	password := ctx.PostForm("password")
	notes := ctx.PostForm("notes")

	sessions := sessions.Default(ctx)
	user := sessions.Get("user")
	userId := sessions.Get("userId").(int)

	siteData := database.GetSiteData(db, userId)
	fmt.Println(siteData)
	err := database.InsertSiteData(db, userId, site, password, notes)
	if err != nil {
		log.Printf("Error inserting new site: %v", err)
		ctx.HTML(http.StatusBadRequest, "vault.html", gin.H{
			"user":     user,
			"siteData": siteData,
			"error":    err.Error(),
		})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/vault")
}
