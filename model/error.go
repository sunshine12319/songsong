package model

import "errors"

var (
	// 用户不存在
	ErrUserNotExist = errors.New("user not exist")
	// 密码不承认
	ErrInvalidPasswd = errors.New("passwd or username not right")
	// 参数不承认
	ErrInvalidParams = errors.New("invalid params")
	//用户存在
	ErrUserExist = errors.New("user exist")
)
