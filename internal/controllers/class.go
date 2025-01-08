package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tihai/internal/model"
	"tihai/internal/service"
)

func CreateClass(c *gin.Context) {
	var class model.Class
	if err := c.ShouldBind(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("uid")
	class.AdminID = uid.(uint)
	cid, err := service.CreateClass(class)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	class.ID = cid
	c.JSON(http.StatusOK, gin.H{"data": class})
}

func JoinClass(c *gin.Context) {
	var class model.Class
	if err := c.ShouldBind(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("uid")
	if err := service.JoinClass(class, uid.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func DeleteClass(c *gin.Context) {
	var class model.Class
	if err := c.ShouldBind(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("uid")
	if err := service.DeleteClass(class, uid.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
