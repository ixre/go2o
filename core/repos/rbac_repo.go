/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: rbac_repo.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-27 11:39:50
 * description: RBAC仓储实现
 * history:
 */
package repos

import (
	rbac "github.com/ixre/go2o/core/domain/interface/rabc"
	rbacImpl "github.com/ixre/go2o/core/domain/rbac"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var _ rbac.IRbacRepository = new(rbacRepoImpl)

type rbacRepoImpl struct {
	userRepo     fw.Repository[rbac.RbacUser]
	userRoleRepo fw.Repository[rbac.RbacUserRole]
}

// GetRbacAggregateRoot implements rbac.IRbacRepository.
func (r *rbacRepoImpl) GetRbacAggregateRoot() rbac.IRbacAggregateRoot {
	// 默认租户:1
	return rbacImpl.NewRbacAggregateRoot(1, r)
}

// UserRepo implements rbac.IRbacRepository.
func (r *rbacRepoImpl) UserRepo() fw.Repository[rbac.RbacUser] {
	return r.userRepo
}

// UserRoleRepo implements rbac.IRbacRepository.
func (r *rbacRepoImpl) UserRoleRepo() fw.Repository[rbac.RbacUserRole] {
	return r.userRoleRepo
}

func NewRbacRepo(o fw.ORM) rbac.IRbacRepository {
	return &rbacRepoImpl{
		userRepo:     &fw.BaseRepository[rbac.RbacUser]{ORM: o},
		userRoleRepo: &fw.BaseRepository[rbac.RbacUserRole]{ORM: o},
	}
}
