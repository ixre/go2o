/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: rbac.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-27 11:39:12
 * description: RBAC领域模型
 * history:
 */

package rbac

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type (
	// IRbacAggregateRoot 权限聚合根
	IRbacAggregateRoot interface {
		domain.IAggregateRoot
		// GetUser 获取用户
		GetUser(userId int) IRbacUser
	}

	// IRbacUser 用户
	IRbacUser interface {
		domain.IDomain
		// GetValue 获取用户信息
		GetValue() RbacUser
		// GetRoles 获取用户角色
		GetRoles() []int
	}

	// IRbacRepository 权限仓储
	IRbacRepository interface {
		// GetRbacAggregateRoot 获取权限聚合根
		GetRbacAggregateRoot() IRbacAggregateRoot
		// UserRepo 用户仓储
		UserRepo() fw.Repository[RbacUser]
		// UserRoleRepo 用户绑定角色仓储
		UserRoleRepo() fw.Repository[RbacUserRole]
		// LoginLogRepo 登陆日志
		LoginLogRepo() fw.Repository[RbacLoginLog]
	}
)
