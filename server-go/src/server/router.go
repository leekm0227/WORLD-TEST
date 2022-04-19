package server

import (
	"fmt"

	"AAA/src/server/api"
	"AAA/src/server/db"
	"AAA/src/server/world"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Run(port string) {
	db.Conn()
	go world.Run()

	router := gin.Default()
	pprof.Register(router)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE")
		c.Next()
	})

	// API
	apiRouter := router.Group("/api")
	{
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/signUp", api.SignUpHandler)
			authRouter.POST("/signIn", api.SignInHandler)
		}
	}

	// WS
	router.GET("/world", world.Handler)
	router.Run(fmt.Sprintf("localhost:%s", port))
}
