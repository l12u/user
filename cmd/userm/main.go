package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/l12u/userm/internal/handler"
	"github.com/l12u/userm/internal/middleware"
	"k8s.io/klog"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	klog.Infoln("Hello World!")

	gin.DisableConsoleColor()
	r := gin.New()
	r.Use(middleware.Logger(3 * time.Second))
	r.Use(gin.Recovery())

	reqHandler := handler.NewRequestHandler()

	r.POST("/login", reqHandler.Login)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	_ = r.Run(":8090")
}
