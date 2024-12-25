/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: rbac.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-27 11:38:57
 * description: RBAC领域实现
 * history:
 */

package rbac

import rbac "github.com/ixre/go2o/core/domain/interface/rabc"

var _ rbac.IRbacAggregateRoot = new(rbacAggregateRootImpl)

type rbacAggregateRootImpl struct {
	tenantId int
	repo     rbac.IRbacRepository
}

// GetAggregateRootId implements rbac.IRbacAggregateRoot.
func (r *rbacAggregateRootImpl) GetAggregateRootId() int {
	return r.tenantId
}

// GetUser implements rbac.IRbacAggregateRoot.
func (r *rbacAggregateRootImpl) GetUser(userId int) rbac.IRbacUser {
	v := r.repo.UserRepo().Get(userId)
	if v == nil {
		return nil
	}
	return newRbacUser(v, r.repo)
}

func NewRbacAggregateRoot(tenantId int, repo rbac.IRbacRepository) rbac.IRbacAggregateRoot {
	return &rbacAggregateRootImpl{
		tenantId: tenantId,
		repo:     repo,
	}
}

var _ rbac.IRbacUser = new(rbacUserImpl)

type rbacUserImpl struct {
	value *rbac.RbacUser
	repo  rbac.IRbacRepository
}

// GetDomainId implements rbac.IRbacUser.
func (r *rbacUserImpl) GetDomainId() int {
	return r.value.Id
}

// GetRoles implements rbac.IRbacUser.
func (r *rbacUserImpl) GetRoles() []int {
	panic("unimplemented")
}

// GetValue implements rbac.IRbacUser.
func (r *rbacUserImpl) GetValue() rbac.RbacUser {
	return *r.value
}

func newRbacUser(value *rbac.RbacUser, repo rbac.IRbacRepository) rbac.IRbacUser {
	return &rbacUserImpl{
		value: value,
		repo:  repo,
	}
}
