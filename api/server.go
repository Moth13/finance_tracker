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
	router.GET("/", server.homePageHandler)

	router.GET("/new_line", server.newLinePageHandler)
	router.POST("/new_line", server.postNewLineHandler)

	router.GET("/about", server.aboutPageHandler)
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
	router.DELETE("/users/:id", server.deleteUser)
	// router.GET("/users/:id", server.getUser)
	// router.GET("/users", server.listUsers)
	// router.POST("/users/login", server.loginUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PATCH("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/months", server.createMonth)
	router.GET("/months/:id", server.getMonth)
	router.GET("/months", server.listMonths)
	router.PATCH("/months/:id", server.updateMonth)
	router.DELETE("/months/:id", server.deleteMonth)

	router.POST("/years", server.createYear)
	router.GET("/years/:id", server.getYear)
	router.GET("/years", server.listYears)
	router.PATCH("/years/:id", server.updateYear)
	router.DELETE("/years/:id", server.deleteYear)

	router.POST("/categories", server.createCategory)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.listCategories)
	// router.UPDATE("/categories/:id", server.updateCategory)
	router.DELETE("/categories/:id", server.deleteCategory)

	router.POST("/lines", server.createLine)
	router.GET("/lines/:id", server.getLine)
	router.GET("/lines", server.listLines)
	router.PATCH("/lines/:id", server.updateLine)
	router.DELETE("/lines/:id", server.deleteLine)

	router.POST("/reclines", server.createRecLine)
	router.GET("/reclines/:id", server.getRecLine)
	router.GET("/reclines", server.listRecLines)
	// router.UPDATE("/reclines/:id", server.UpdateRecLine)
	router.DELETE("/reclines/:id", server.deleteRecLine)

	// router.GET("/stats/", server.getStats)

	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

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
