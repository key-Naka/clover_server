package flag_user

import (
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/pwd"
	"fmt"
	"os"

	"golang.org/x/term"
)

type FlagUser struct {
}

func (u *FlagUser) Create() {
	var Role enum.RoleType
	fmt.Println("选择角色  1超级管理员 2普通用户 3访客")
	_, err := fmt.Scanln(&Role)
	if err != nil {
		fmt.Println("输入错误")
		return
	}
	fmt.Println("输入用户名")
	var Username string
	_, err = fmt.Scanln(&Username)
	if err != nil {
		fmt.Println("输入错误")
		return
	}
	var model models.UserModel
	err = global.DB.Take(&model, "username = ?", Username).Error
	if err == nil {
		fmt.Println("用户名已存在")
		return
	}
	fmt.Println("输入密码")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("输入错误")
		return
	}
	fmt.Println("再次输入密码")
	confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("输入错误")
		return
	}
	if string(password) != string(confirmPassword) {
		fmt.Println("两次输入密码不一致")
		return
	}
	hashPassword, err := pwd.GenerateFromPassword(string(password))
	if err != nil {
		fmt.Println("密码加密错误")
		return
	}

	err = global.DB.Create(&models.UserModel{
		Username:       Username,
		Nickname:       Username,
		Password:       hashPassword,
		RegisterSource: enum.RegisterTerminalSourceType,
		Role:           Role,
	}).Error
	if err != nil {
		fmt.Println("创建用户失败")
		return
	}
}
