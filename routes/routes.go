package routes

import (
	"net/http"

	"github.com/minateegithub/go_microservice/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/enrollees", controllers.GetAllEnrollees)
	router.POST("/enrollee", controllers.CreateEnrollee)
	router.GET("/enrollee/:enrolleeId", controllers.GetSingleEnrollee)
	router.PUT("/enrollee/:enrolleeId", controllers.EditEnrollee)
	router.DELETE("/enrollee/:enrolleeId", controllers.DeleteEnrollee)
	router.POST("/enrollee/:enrolleeId/dependents", controllers.AddDependent)
	router.GET("/enrollee/:enrolleeId/dependents", controllers.GetDependents)
	router.PUT("/enrollee/:enrolleeId/dependents/:dependentId", controllers.EditDependent)
	router.DELETE("/enrollee/:enrolleeId/dependents/:dependentId", controllers.DeleteDependent)
	router.NoRoute(notFound)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
	return
}
