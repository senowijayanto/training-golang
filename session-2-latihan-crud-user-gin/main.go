package main

import (
	"training-golang/session-2-latihan-crud-user-gin/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Router
	router.SetupRouter(r)

	r.Run("localhost:8080")
}
