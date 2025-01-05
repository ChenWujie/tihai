package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tihai/internal/model"
	"tihai/internal/service"
)

func CreateQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("uid")
	question.TeacherID = uid.(uint)
	if err := service.CreateQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func UpdateQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UpdateQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": question})
}

func DeleteQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func GetQuestion(c *gin.Context) {
	t := c.Query("type")
	token := c.GetHeader("Authorization")
	var list []model.Question
	var err error
	if token == "" {
		list, err = service.FindListByGuest(t)
	} else {
		list, err = service.FindList(t, token)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func LikeQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not login"})
		return
	}
	result, data, err := service.LikeQuestion(uid.(uint), question.ID)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "result": result})
}

func SearchQuestion(c *gin.Context) {
	search := c.Query("query")
	res, err := service.SearchArticles(search)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
