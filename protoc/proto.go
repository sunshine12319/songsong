package protoc

import "study/example06/model"

// 协议
type Message struct {
	Cmd string `json:"cmd"`
	Data string `json:"data"`
}

// 客户端登录信息
type LoginCmd struct{
	Id int `json:"user_id"`
	Passwd string `json:"passwd"`
}

// 客户端注册信息
type RegisterCmd struct {
	User model.User `json:"user"`
}

// 客户端登录返回信息
type LoginCmdRes struct{
	Code int `json:"code"`
	Error string `json:"error"`
}