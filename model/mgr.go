package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	UserTable = "users"
	RedisPasswd = "admin"
)

type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) (mgr *UserMgr) {
	mgr = &UserMgr{
		pool: pool,
	}
	return
}

// 获取conn并并验证密码
func (p *UserMgr) getConn() redis.Conn {
	conn := p.pool.Get()
	conn.Do("auth",RedisPasswd)
	return conn
}


// 获取用户信息
func (p *UserMgr) getUser(conn redis.Conn,id int)(user *User,err error){
	result, err := redis.String(conn.Do("HGet",UserTable,fmt.Sprintf("%d",id)))
	if err != nil {
		//查账号是否存在
		if err == redis.ErrNil{
			err = ErrUserNotExist
		}
		return
	}

	user = &User{}
	err = json.Unmarshal([]byte(result),user)
	if err != nil {
		return
	}
	return
}

// 登录
func (p *UserMgr) Login(id int,passwd string)(user *User, err error){
	conn:= p.getConn()
	defer conn.Close()
	user, err = p.getUser(conn,id)
	if err != nil{
		return
	}

	// 查账号密码是否匹配
	if user.UserId != id &&user.Passwd!=passwd{
		err =ErrInvalidPasswd
	}
	user.Status = UserStatusOnline
	user.LastLogin = fmt.Sprintf("%v",time.Now())
	return
}

// 注册
func (p *UserMgr) Register(user *User)(err error){
	conn := p.getConn()
	defer conn.Close()
	if user == nil {
		fmt.Println("invalid user")
		err = ErrInvalidParams
		return
	}

	_, err = p.getUser(conn,user.UserId)
	// 返回账号存在错误
	if err == nil {
		err = ErrUserExist
		return
	}

	// getUser其它错误返回
	if err != ErrUserNotExist {
		return
	}

	//用户信息不存在时
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	_, err = conn.Do("HSet",UserTable,fmt.Sprintf("%d",user.UserId),string(data))
	if err != nil{
		return
	}
	return
}
