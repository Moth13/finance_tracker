package server

import (
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/moth13/finance_tracker/internal/api"
	db "github.com/moth13/finance_tracker/internal/db/sqlc"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web"
)

// Server serves HTTP request for our banking service
type Server struct {
	router *gin.Engine

	apiServer *api.Server
	webServer *web.Server
}

// New server creates a new HTTP server dans setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		apiServer: &api.Server{
			Config: config,
			Store:  store,
		},
		webServer: &web.Server{
			Store: store,
		},
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.Default())

	server.apiServer.SetupRoutes(router)
	server.webServer.SetupRoutes(router)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
