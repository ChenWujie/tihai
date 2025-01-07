package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tihai/internal/model"
	"tihai/internal/service"
)

type RequestPaper struct {
	model.Paper
	QuestionIDS []uint `json:"question_ids"`
}

func CreatePaper(c *gin.Context) {
	var requestPaper RequestPaper
	if err := c.ShouldBindJSON(&requestPaper); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	paper := requestPaper.Paper
	uid, _ := c.Get("uid")
	paper.UserID = uid.(uint)
	err := service.CreatePaper(paper, requestPaper.QuestionIDS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paper.ID})
}
