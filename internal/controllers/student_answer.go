package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tihai/internal/model"
	"tihai/internal/service"
	"time"
)

func CreateStudentAnswer(c *gin.Context) {
	var answer model.StudentAnswer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		return
	}
	answer.UserID = userID.(uint)
	answer.SubmitTime = time.Now()
	if score, rate, err := service.CreateStudentAnswer(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"result": score, "rate": rate})
	}
}

func GetStudentAnswerListForUser(c *gin.Context) {
	uid, _ := c.Get("uid")
	list, err := service.GetStudentAnswerList(uid.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetStudentAnswerListForQuestion(c *gin.Context) {
	questionId := c.Query("qid")
	qid, err := strconv.Atoi(questionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请先登录"})
		return
	}
	list, err := service.GetUserAnswerList(uid.(uint), uint(qid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
