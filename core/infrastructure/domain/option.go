/**
 * Copyright 2015 @ z3q.net.
 * name : option.go
 * author : jarryliu
 * date : 2016-04-18 13:48
 * description :
 * history :
 */
package domain

import (
	"encoding/json"
	"errors"
	"github.com/ixre/gof"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type (
	IOptionStore interface {
		// the indent of option store
		Indent() string

		// check state
		Stat() error

		// get all options
		All() (keys []string, values []*Option)

		// get option by key
		Get(key string) (value *Option)

		// after set,call flush()
		Set(key string, value *Option)

		// flush to file
		Flush() error

		// destroy and delete file
		Destroy() error
	}

	OptionType int

	Option struct {
		Key   string `json:"key"`   //键
		Type  int    `json:"type"`  //类型
		Must  bool   `json:"must"`  //是否必须填写
		Title string `json:"title"` //标题
		Value string `json:"value"` //值
	}

	OptionStoreWrapper struct {
		_indent string
		_data   map[string]*Option
	}
)

const (
	OptionTypeInt int = iota
	OptionTypeString
	OptionTypeBoolean
)

func BuildOptionsForm(v *[]Option) string {
	return ""
}

var _ IOptionStore = new(OptionStoreWrapper)

func NewOptionStoreWrapper(indent string) IOptionStore {
	return &OptionStoreWrapper{
		_indent: indent,
	}
}

// the indent of option store
func (this *OptionStoreWrapper) Indent() string {
	if len(this._indent) == 0 {
		panic("option store indent must set by implmenter")
	}
	return this._indent
}

func (this *OptionStoreWrapper) getRdKey() string {
	s := strings.Replace(this.Indent(), "/", ":", -1)
	if s[:2] == "./" {
		s = s[2:]
	} else if s[:1] == "." {
		s = s[1:]
	}
	return s
}
func (this *OptionStoreWrapper) Stat() error {
	if this._data == nil {
		return this.load()
	}
	if len(this._data) == 0 {
		return errors.New("empty options!")
	}
	return nil
}

func (this *OptionStoreWrapper) load() error {
	sto := gof.CurrentApp.Storage()
	this._data = make(map[string]*Option)
	rdKey := this.getRdKey()
	if sto.Get(rdKey, &this._data) != nil {
		// 从KV中取得,且KV不存在返回错误
		d, err := ioutil.ReadFile(this.Indent())
		if err == nil {
			err = json.Unmarshal(d, &this._data)
			if err == nil {
				err = sto.SetExpire(rdKey, this._data, 3600) //存储到Kv
			}
		}
		return err
	}
	return nil
}

// get all options
func (this *OptionStoreWrapper) All() (keys []string, values []*Option) {
	if this.Stat() != nil {
		return nil, nil
	}
	vl := make([]*Option, 0)
	kl := make([]string, 0)

	for k, v := range this._data {
		kl = append(kl, k)
		vl = append(vl, v)
	}
	sort.Strings(kl)
	return kl, vl
}

// get option by key
func (this *OptionStoreWrapper) Get(key string) *Option {
	if this.Stat() != nil {
		return nil
	}
	if v, ok := this._data[key]; ok {
		return v
	}
	return nil
}

// after set,call flush()
func (this *OptionStoreWrapper) Set(key string, v *Option) {
	if this._data == nil {
		this._data = make(map[string]*Option)
	}
	this._data[key] = v
	gof.CurrentApp.Storage().Del(this.getRdKey()) // clean cache
}

// flush to file
func (this *OptionStoreWrapper) Flush() (err error) {
	if this._data != nil {
		dir := filepath.Dir(this.Indent()) //检查目录是否存在
		if _, err = os.Stat(dir); err != nil {
			//创建目录
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}
		d, err := json.MarshalIndent(this._data, "", " ")
		if err == nil {
			err = ioutil.WriteFile(this.Indent(), d, os.ModePerm)
		}
		return err
	}
	return nil
}

// destroy and delete file
func (this *OptionStoreWrapper) Destroy() error {
	return os.Remove(this.Indent())
}
