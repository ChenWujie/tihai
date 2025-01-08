package router

import (
	"github.com/gin-gonic/gin"
	"tihai/internal/controllers"
	"tihai/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/user")
	v1.POST("/login", controllers.Login)
	v1.POST("/register", controllers.Register)
	v1.Use(middleware.AuthMiddleWare())
	{
		v1.POST("/update", controllers.Update)
		v1.POST("/unlogin", controllers.UnLogin)
	}
	v2 := r.Group("/question")
	v2.GET("/get", controllers.GetQuestion)
	v2.GET("/search", controllers.SearchQuestion)
	v2.Use(middleware.AuthMiddleWare())
	{
		v2.POST("/like", controllers.LikeQuestion)
		v2.Use(middleware.TeacherMiddle())
		v2.POST("/create", controllers.CreateQuestion)
		v2.POST("/update", controllers.UpdateQuestion)
		v2.POST("/delete", controllers.DeleteQuestion)
	}
	v3 := r.Group("/answer")
	v3.Use(middleware.AuthMiddleWare())
	{
		v3.GET("/list", controllers.GetStudentAnswerListForUser)
		v3.POST("/create", controllers.CreateStudentAnswer)
		v3.GET("/get", controllers.GetStudentAnswerListForQuestion)
	}
	v4 := r.Group("/paper")
	v4.Use(middleware.AuthMiddleWare())
	{
		v4.GET("/get", controllers.GetPapers)
		v4.POST("/create", controllers.CreatePaper)
		v4.DELETE("/delete", controllers.DeletePaper)
		v4.POST("/update", controllers.UpdatePaper)
	}
	v5 := r.Group("/class")
	v5.Use(middleware.AuthMiddleWare())
	{
		v5.POST("/create", controllers.CreateClass)
		v5.DELETE("/delete", controllers.DeleteClass)
		v5.POST("/join", controllers.JoinClass)
	}
	return r
}
