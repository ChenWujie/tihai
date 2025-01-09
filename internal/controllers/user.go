package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tihai/internal/model"
	"tihai/internal/service"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.Register(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := service.Login(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UnLogin(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.UnLogin(c.GetHeader("Authorization"))
	c.JSON(http.StatusOK, gin.H{"data": "退出登录！"})
}

func Update(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user.ID = t.(uint)
	token := c.GetHeader("Authorization")
	if err := service.Update(user, token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新成功！"})
}

func GetInformation(c *gin.Context) {
	uid, _ := c.Get("uid")
	c.JSON(http.StatusOK, gin.H{"data": service.GetInformation(uid.(uint))})
}
