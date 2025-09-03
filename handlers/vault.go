package handlers

import (
	"database/sql"
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

	siteData, err := database.GetSiteData(db, userId)
	if err != nil {
		log.Fatalf("Eror fetching data from db: %v", err)
	}

	ctx.HTML(http.StatusOK, "vault.html", gin.H{
		"user":     user,
		"siteData": siteData,
	})
}

func addSite(ctx *gin.Context, db *sql.DB) {
	site := ctx.PostForm("site")
	password := ctx.PostForm("password")
	notes := ctx.PostForm("notes")

	sessions := sessions.Default(ctx)
	user := sessions.Get("user")
	userId := sessions.Get("userId").(int)

	siteData, err := database.GetSiteData(db, userId)
	if err != nil {
		log.Fatalf("Error fetching data from db: %v", err)
	}
	err = database.InsertSiteData(db, userId, site, password, notes)
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

func editSite(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.HTML(http.StatusOK, "edit.html", gin.H{
		"id": id,
	})
}

func editSiteConfirm(ctx *gin.Context, db *sql.DB) {
	sessions := sessions.Default(ctx)
	userId := sessions.Get("userId").(int)
	id := ctx.Param("id")

	siteData, err := database.GetSiteData(db, userId)
	if err != nil {
		log.Fatalf("error getting site data: %v", err)
	}

	password := ctx.PostForm("password")
	siteName := ctx.PostForm("site")
	notes := ctx.PostForm("notes")

	err = database.EditData(db, userId, id, password, siteName, notes)
	if err != nil {
		user := sessions.Get("user")
		ctx.HTML(http.StatusBadRequest, "vault.html", gin.H{
			"user":     user,
			"siteData": siteData,
			"error":    err,
		})
	}

	ctx.Redirect(http.StatusSeeOther, "/vault")
}

func deleteSite(ctx *gin.Context, db *sql.DB) {
	sessions := sessions.Default(ctx)
	userId := sessions.Get("userId").(int)

	siteId := ctx.Param("id")

	err := database.DeleteData(db, userId, siteId)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func copyPassword(ctx *gin.Context, db *sql.DB) {
	sessions := sessions.Default(ctx)
	userId := sessions.Get("userId").(int)

	siteId := ctx.Param("id")
	password, err := database.GetPassword(db, userId, siteId)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.String(http.StatusOK, password)
}
