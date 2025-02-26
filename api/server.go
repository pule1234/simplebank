package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/pule1234/simple_bank/db/sqlc"
	"github.com/pule1234/simple_bank/token"
	"github.com/pule1234/simple_bank/util"
)

type Server struct {
	config     util.Config //读取文件配置
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("currency", validCurrency)
	//}
	server.setupRouter()
	return server, nil
}

// 定义路由， 同时初始化router
func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker)) // 使用中间件进行认证
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// 错误处理
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
