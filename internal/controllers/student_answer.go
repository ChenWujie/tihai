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
	answer.UserId = userID.(uint)
	answer.SubmitTime = time.Now()
	if score, err := service.CreateStudentAnswer(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": score})
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
	list, err := service.GetUserAnswerList(uint(qid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
