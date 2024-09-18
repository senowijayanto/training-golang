package router

import (
	"training-golang/session-4-sample-separate-layer/step-4/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler handler.IUserHandler) {
	usersPublicEndpoint := r.Group("/users")

	usersPublicEndpoint.GET("/", userHandler.GetAllUsers)
}
