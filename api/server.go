package api

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/token"
	"github.com/moth13/finance_tracker/util"
)

// Server serves HTTP request for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

var staticFiles embed.FS

// New server creates a new HTTP server dans setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.Default())

	server.setupViewRoutes(router)
	server.setupApiRoutes(router)

	server.router = router
}

func (server *Server) setupViewRoutes(router *gin.Engine) {
	router.StaticFS("/static/css/", http.FS(staticFiles))

	//views
	router.GET("/", server.homePage)

	views := router.Group("/views")

	views.GET("/lines", server.getViewLinePage)
	views.GET("/lines/:id", server.getViewLinePage)
	views.POST("/lines", server.postViewLine)
	views.DELETE("/lines/:id", server.deleteViewLine)
	views.PUT("/lines/:id", server.updateViewLine)

	views.GET("/about", server.aboutPageHandler)
}

func (server *Server) render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func (server *Server) setupApiRoutes(router *gin.Engine) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/api").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/users", server.createUser)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	// authRoutes.GET("/users/:id", server.getUser)
	// authRoutes.GET("/users", server.listUsers)
	// authRoutes.POST("/users/login", server.loginUser)

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PATCH("/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	authRoutes.POST("/months", server.createMonth)
	authRoutes.GET("/months/:id", server.getMonth)
	authRoutes.GET("/months", server.listMonths)
	authRoutes.PATCH("/months/:id", server.updateMonth)
	authRoutes.DELETE("/months/:id", server.deleteMonth)

	authRoutes.POST("/years", server.createYear)
	authRoutes.GET("/years/:id", server.getYear)
	authRoutes.GET("/years", server.listYears)
	authRoutes.PATCH("/years/:id", server.updateYear)
	authRoutes.DELETE("/years/:id", server.deleteYear)

	authRoutes.POST("/categories", server.createCategory)
	authRoutes.GET("/categories/:id", server.getCategory)
	authRoutes.GET("/categories", server.listCategories)
	// authRoutes.UPDATE("/categories/:id", server.updateCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)

	authRoutes.POST("/lines", server.createLine)
	authRoutes.GET("/lines/:id", server.getLine)
	authRoutes.GET("/lines", server.listLines)
	authRoutes.PATCH("/lines/:id", server.updateLine)
	authRoutes.DELETE("/lines/:id", server.deleteLine)

	authRoutes.POST("/reclines", server.createRecLine)
	authRoutes.GET("/reclines/:id", server.getRecLine)
	authRoutes.GET("/reclines", server.listRecLines)
	// authRoutes.UPDATE("/reclines/:id", server.UpdateRecLine)
	authRoutes.DELETE("/reclines/:id", server.deleteRecLine)

	// api.GET("/stats/", server.getStats)

	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccounts)

	// authRoutes.POST("/transfers", server.createTransfer)
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
