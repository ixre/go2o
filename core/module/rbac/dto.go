package rbac

import "github.com/ixre/go2o/core/service/proto"

// 用户初始化数据
type UserInfoResponse struct {
	// 昵称
	Nickname string `json:"nickname"`
	// 头像
	ProfilePhoto string `json:"profilePhoto"`
	// 登录IP
	LoginIp string `json:"loginIp"`
	// 资源Key
	ResourceKeys []string `json:"resourceKeys"`
	// 角色
	Roles []string `json:"roles"`
	// 菜单数据
	MenuData []*proto.SUserMenu `json:"menuData"`
	// 用户设置
	Settings map[string]interface{} `json:"settings"`
}
