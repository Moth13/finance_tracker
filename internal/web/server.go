package web

import (
	"embed"
	"net/http"

	"github.com/a-h/templ"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/internal/db/sqlc"
)

type Server struct {
	Store db.Store
}

var staticFiles embed.FS

type User struct {
	Email    string
	Username string
}

func (server *Server) SetupRoutes(router *gin.Engine) {

	store := cookie.NewStore([]byte("temporary-secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 jours
		HttpOnly: true,
		Secure:   false, // true derrière HTTPS
	})
	router.Use(sessions.Sessions("app-session", store))

	router.StaticFS("/static/assets/css/", http.FS(staticFiles))

	//views
	views := router.Group("/views")
	views.GET("/", server.homePage)
	views.POST("/login", server.loginHandler)
	views.POST("/signup_user", server.signupUserPage)
	views.POST("/signup_account", server.signupAccountPage)
	views.GET("/logout", server.logoutHandler)

	auth := views.Group("/")
	auth.Use(AuthRequired())
	auth.GET("/lines", server.viewLinesPage)
	auth.GET("/lines/new", server.newViewLinePage)
	auth.GET("/lines/:id", server.getViewLinePage)
	auth.POST("/lines", server.postViewLine)
	auth.DELETE("/lines/:id", server.deleteViewLine)
	auth.PUT("/lines/:id", server.updateViewLine)

	auth.GET("/about", server.aboutPageHandler)
}

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}
