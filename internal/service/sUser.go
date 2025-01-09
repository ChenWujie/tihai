package service

import (
	"errors"
	"tihai/global"
	"tihai/internal/model"
	"tihai/utils"
	"time"
)

func Register(user *model.User) error {
	if user.Role == "" {
		user.Role = "student"
	}
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPwd
	if err := global.Db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func Login(input *model.User) (token string, err error) {
	var user model.User
	result := global.Db.Where("username = ?", input.Username).First(&user)
	if err = result.Error; err != nil {
		return "", err
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		return "", errors.New("用户名或密码错误")
	}

	token, err = utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func UnLogin(token string) {
	global.RedisDB.Set(token, "blacklist", time.Hour*72+1)
}

func Update(input model.User, token string) error {
	m := make(map[string]interface{})
	if input.Email != "" {
		m["email"] = input.Email
	}
	if input.Nickname != "" {
		m["nickname"] = input.Nickname
	}
	if input.Password != "" {
		hashedPwd, err := utils.HashPassword(input.Password)
		if err != nil {
			return err
		}
		m["password"] = hashedPwd
	}
	if err := global.Db.Model(model.User{}).Where("id = ?", input.ID).Updates(m).Error; err != nil {
		return err
	}
	err := global.RedisDB.Set(token, "blacklist", time.Hour*72+1).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetInformation(uid uint) model.User {
	var user model.User
	global.Db.Select("id,username,role").First(&user, uid)
	return user
}
