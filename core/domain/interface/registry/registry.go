package registry

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// 注册表
type Registry struct {
	// 键
	Key string `db:"key" pk:"yes"`
	// 值
	Value string `db:"value"`
	// 默认值
	DefaultValue string `db:"default_value"`
	// 可选值
	Options string `db:"options"`
	// 是否用户定义,0:否,1:是
	UserDefine int16 `db:"user_define"`
	// 描述
	Description string `db:"description"`
}

// 注册项
type IRegistry interface {
	// 获取聚合根编号
	GetAggregateRootId() string
	// 返回键
	Key() string
	// 原始数据
	Value() Registry
	// 是否为用户创建
	IsUser() bool
	// 返回字符值
	StringValue() string
	// 返回Int值
	IntValue() int
	// 返回浮点值
	FloatValue() float64
	// 返回布尔值
	BoolValue() bool
	// 删除项
	Remove() error
	// 重置为默认值
	Reset() error
	// 更新
	Update(value string) error
	// 保存
	Save() error
}

type IRegistryRepo interface {
	// 删除键
	Remove(key string) error
	// 保存键
	Save(registry IRegistry) error
	// 创建
	Create(r *Registry) IRegistry
	// 获取键
	Get(key string) IRegistry
	// 合并数据
	Merge(registries []*Registry) error
	// 搜索注册表
	SearchRegistry(key string) []Registry
}

func KeyFormat(s string) string {
	dst := make([]byte, 0)
	for i, b := range strings.TrimSpace(s) {
		if unicode.IsUpper(b) {
			l := byte(unicode.ToLower(b))
			if i == 0 {
				dst = append(dst, l)
			} else {
				dst = append(dst, byte('_'), l)
			}
		} else {
			dst = append(dst, byte(b))
		}
	}
	return string(dst)
}

var _ IRegistry = new(registryImpl)

type registryImpl struct {
	value *Registry
	repo  IRegistryRepo
}

func NewRegistry(r *Registry, repo IRegistryRepo) IRegistry {
	r.Value = strings.TrimSpace(r.Value)
	return &registryImpl{
		value: r,
		repo:  repo,
	}
}

func (r *registryImpl) IsUser() bool {
	return r.value.UserDefine == 1
}

func (r *registryImpl) GetAggregateRootId() string {
	return r.value.Key
}

func (r *registryImpl) Key() string {
	return r.value.Key
}

func (r *registryImpl) Value() Registry {
	return *r.value
}

func (r *registryImpl) StringValue() string {
	return r.value.Value
}

func (r *registryImpl) IntValue() int {
	v, err := strconv.Atoi(r.value.Value)
	r.panic(err)
	return v
}

func (r *registryImpl) FloatValue() float64 {
	v, err := strconv.ParseFloat(r.value.Value, 64)
	r.panic(err)
	return v
}

func (r *registryImpl) BoolValue() bool {
	v, err := strconv.ParseBool(r.value.Value)
	r.panic(err)
	return v
}

func (r *registryImpl) Remove() error {
	if r.value.UserDefine == 0 {
		return errors.New("registry is not create by user, can't be removed")
	}
	return r.repo.Remove(r.Key())
}

func (r *registryImpl) Reset() error {
	return r.Update(r.value.DefaultValue)
}

func (r *registryImpl) Update(value string) error {
	if r.value.Value != value {
		r.value.Value = value
		return r.Save()
	}
	return nil
}

func (r *registryImpl) Save() error {
	r.value.Key = KeyFormat(r.value.Key)
	if len(r.value.Key) > 45 {
		return errors.New("key length out of 40")
	}
	if len(r.value.Value) > 120 {
		return errors.New("value length out of 120")
	}
	return r.repo.Save(r)
}

func (r *registryImpl) panic(e error) {
	if e != nil {
		panic(fmt.Sprintf("parse registry value fail! key:%s value:%s",
			r.value.Key,
			r.value.Value))
	}
}
