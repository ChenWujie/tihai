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
	err := service.CreatePaper(&paper, requestPaper.QuestionIDS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paper.ID})
}

func DeletePaper(c *gin.Context) {
	var requestPaper RequestPaper
	if err := c.ShouldBindJSON(&requestPaper); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.DeletePaper(requestPaper.Paper)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func UpdatePaper(c *gin.Context) {
	var requestPaper RequestPaper
	if err := c.ShouldBindJSON(&requestPaper); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("uid")
	paper := requestPaper.Paper
	paper.UserID = uid.(uint)
	err := service.UpdatePaper(paper, requestPaper.QuestionIDS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paper.ID})
}

func GetPapers(c *gin.Context) {
	uid, _ := c.Get("uid")
	paper, err := service.GetPaper(uid.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"paper": paper})
}

func AssignPapers(c *gin.Context) {
	type temp struct {
		ClassIds []uint `json:"class_ids"`
		PaperId  uint   `json:"paper_id"`
	}
	uid, _ := c.Get("uid")
	var t temp
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.AssignPapers(uid.(uint), t.PaperId, t.ClassIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// GetClassPapers 查询当前用户所加入的班级拥有的试卷
func GetClassPapers(c *gin.Context) {
	uid, _ := c.Get("uid")
	papers, err := service.QueryClassPapers(uid.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"papers": papers})
}

func AnswerPaper(c *gin.Context) {

}
