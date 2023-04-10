package model

const (
	UserStatusOffline = iota
	UserStatusOnline
)

type User struct {
	UserId int    `json:"user_id"`
	Passwd string `json:"passwd"`
	Nick   string `json:"nick"`
	Sex    string `json:"sex"`
	// 用户头像url
	Header    string `json:"header"`
	LastLogin string `json:"last_login"`
	// 是否在线
	Status int `json:"status"`
}
