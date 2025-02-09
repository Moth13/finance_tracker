package api

import (
	"embed"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	// "github.com/techsool/simplebank/token"
)

// Server serves HTTP request for our banking service
type Server struct {
	config util.Config
	store  db.Store
	// tokenMaker token.Maker
	router *gin.Engine
}

var staticFiles embed.FS

// New server creates a new HTTP server dans setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannont create token maker %w", err)
	// }
	server := &Server{
		config: config,
		store:  store,
		// tokenMaker: tokenMaker,
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
	views.POST("/lines", server.postViewLine)
	views.DELETE("/lines/:id", server.deleteViewLine)

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

	api := router.Group("/api")

	api.POST("/users", server.createUser)
	api.DELETE("/users/:id", server.deleteUser)
	// api.GET("/users/:id", server.getUser)
	// api.GET("/users", server.listUsers)
	// api.POST("/users/login", server.loginUser)

	api.POST("/accounts", server.createAccount)
	api.GET("/accounts/:id", server.getAccount)
	api.GET("/accounts", server.listAccounts)
	api.PATCH("/accounts/:id", server.updateAccount)
	api.DELETE("/accounts/:id", server.deleteAccount)

	api.POST("/months", server.createMonth)
	api.GET("/months/:id", server.getMonth)
	api.GET("/months", server.listMonths)
	api.PATCH("/months/:id", server.updateMonth)
	api.DELETE("/months/:id", server.deleteMonth)

	api.POST("/years", server.createYear)
	api.GET("/years/:id", server.getYear)
	api.GET("/years", server.listYears)
	api.PATCH("/years/:id", server.updateYear)
	api.DELETE("/years/:id", server.deleteYear)

	api.POST("/categories", server.createCategory)
	api.GET("/categories/:id", server.getCategory)
	api.GET("/categories", server.listCategories)
	// api.UPDATE("/categories/:id", server.updateCategory)
	api.DELETE("/categories/:id", server.deleteCategory)

	api.POST("/lines", server.createLine)
	api.GET("/lines/:id", server.getLine)
	api.GET("/lines", server.listLines)
	api.PATCH("/lines/:id", server.updateLine)
	api.DELETE("/lines/:id", server.deleteLine)

	api.POST("/reclines", server.createRecLine)
	api.GET("/reclines/:id", server.getRecLine)
	api.GET("/reclines", server.listRecLines)
	// api.UPDATE("/reclines/:id", server.UpdateRecLine)
	api.DELETE("/reclines/:id", server.deleteRecLine)

	// api.GET("/stats/", server.getStats)

	// authRoutes := api.Group("/").Use(authMiddleware(server.tokenMaker))

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
