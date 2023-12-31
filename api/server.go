package api

import (
	"fmt"
	db "github.com/MsN-12/simpleBank/db/sqlc"
	"github.com/MsN-12/simpleBank/token"
	"github.com/MsN-12/simpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}
	server.setupRouter()
	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("tokens/renew_access", server.renewAccessToken)

	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccounts)
	authRouter.POST("/transfers", server.createTransfer)

	server.router = router

}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
