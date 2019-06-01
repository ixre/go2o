package repos

import (
	"database/sql"
	"github.com/ixre/gof/db"
	"go2o/core/domain/interface/registry"
	"log"
)

var _ registry.IRegistryRepo = new(registryRepo)

type registryRepo struct {
	conn db.Connector
	data map[string]registry.IRegistry
}

func NewRegistryRepo(conn db.Connector) registry.IRegistryRepo {
	return (&registryRepo{
		conn: conn,
		data: make(map[string]registry.IRegistry),
	}).init()
}

func (r *registryRepo) init() registry.IRegistryRepo {
	// 从数据源加载数据
	list := make([]*registry.Registry, 0)
	err := r.conn.GetOrm().Select(&list, "")
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Registry")
	}
	for _, v := range list {
		r.data[v.Key] = r.Create(v)
	}
	// 合并数据源
	registries := registry.MergeRegistries()
	r.Merge(registries)
	// 清理不再使用的注册表
	r.truncUnused(registries)
	return r
}

func (r *registryRepo) Get(key string) registry.IRegistry {
	return r.data[key]
}

func (r *registryRepo) Remove(key string) error {
	_, err := r.conn.ExecNonQuery("DELETE FROM registry WHERE key=$1", key)
	delete(r.data, key)
	return err
}

func (r *registryRepo) Save(registry registry.IRegistry) (err error) {
	key := registry.Key()
	val := registry.Value()
	_, ok := r.data[key]
	if ok {
		_, _, err = r.conn.GetOrm().Save(key, val)
	} else {
		_, _, err = r.conn.GetOrm().Save(nil, val)
	}
	if err == nil { // 更新缓存
		r.data[key] = registry
	} else if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Registry")
	}
	return err
}

// 合并Registry
func (r *registryRepo) Merge(registries []*registry.Registry) error {
	if registries == nil || len(registries) == 0 {
		return nil
	}
	for _, v := range registries {
		if ir := r.Get(v.Key); ir != nil {
			raw := ir.Value()
			if v.Description != raw.Description || v.DefaultValue != raw.DefaultValue ||
				v.Options != raw.Options {
				// 更新值
				raw.DefaultValue = v.DefaultValue
				raw.Description = v.Description
				raw.Options = v.Options
				// 更新缓存并保存
				ir = r.Create(&raw)
				if err := r.Save(ir); err != nil {
					return err
				}
			}
		} else {
			ir := r.Create(v)
			if err := r.Save(ir); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *registryRepo) Create(v *registry.Registry) registry.IRegistry {
	return registry.NewRegistry(v, r)
}

// 清理不使用的系统键
func (r *registryRepo) truncUnused(registries []*registry.Registry) error {
	exists := true
	for _, ir := range r.data {
		if !ir.IsUser() {
			exists = false
			for _, ir2 := range registries {
				if ir2.Key == ir.Key() {
					exists = true
					break
				}
			}
			if !exists {
				r.Remove(ir.Key())
			}
		}
	}
	return nil
}
