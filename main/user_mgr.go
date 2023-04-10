package main

import "study/example06/model"

var (
	mgr  *model.UserMgr
)

// 创建model.UserMgr
func initUserMgr(){
	mgr = model.NewUserMgr(pool)
}