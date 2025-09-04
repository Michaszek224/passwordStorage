package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"passwordStorage/database"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	githubOauthConfig *oauth2.Config
	googleOauthConfig *oauth2.Config
)

func InitOAuthConfigs() {
	required := []string{"GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET", "GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			panic(fmt.Sprintf("Environment variable %s is not set", v))
		}
	}
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:2137/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:2137/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type GithubUser struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func githubLoginHandler(ctx *gin.Context) {
	url := githubOauthConfig.AuthCodeURL("state-random")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func githubCallbackHandler(ctx *gin.Context, db *sql.DB) {
	code := ctx.Query("code")
	token, err := githubOauthConfig.Exchange(ctx, code)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Failed to exchange token: %v", err)
		return
	}
	client := githubOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil || resp.StatusCode != 200 {
		ctx.String(http.StatusUnauthorized, "Failed to get user info")
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ghUser GithubUser
	err = json.Unmarshal(body, &ghUser)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to parse user info: %v", err)
		return
	}

	if ghUser.Email == "" {
		emailResp, emailErr := client.Get("https://api.github.com/user/emails")
		if emailErr == nil && emailResp.StatusCode == 200 {
			defer emailResp.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			emailBody, _ := io.ReadAll(emailResp.Body)
			if json.Unmarshal(emailBody, &emails) == nil {
				for _, e := range emails {
					if e.Primary && e.Verified {
						ghUser.Email = e.Email
						break
					}
				}
				if ghUser.Email == "" && len(emails) > 0 {
					ghUser.Email = emails[0].Email
				}
			}
		}
	}

	if ghUser.Email == "" {
		ctx.String(http.StatusInternalServerError, "Could not retrieve email from GitHub account.")
		return
	}

	id, err := database.FindOrCreateOAuthUser("github",
		fmt.Sprintf("%v", ghUser.ID),
		ghUser.Email,
		ghUser.Login,
		db,
	)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "DB error: %v", err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("user", ghUser.Login)
	session.Set("userId", int64(id))
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/vault")
}

func googleLoginHandler(ctx *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state-random")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func googleCallbackHandler(ctx *gin.Context, db *sql.DB) {
	code := ctx.Query("code")
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Failed to exchange token: %v", err)
		return
	}

	client := googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != 200 {
		ctx.String(http.StatusUnauthorized, "Failed to get user info")
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var gUser GoogleUser

	if err := json.Unmarshal(body, &gUser); err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to parse user info")
		return
	}

	id, err := database.FindOrCreateOAuthUser("google",
		gUser.ID,
		gUser.Email,
		gUser.Email, // using email also as username
		db,
	)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "DB error: %v", err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("user", gUser.Email)
	session.Set("userId", id)
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/vault")

}
