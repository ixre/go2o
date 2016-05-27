/**
 * Copyright 2015 @ z3q.net.
 * name : rep
 * author : jarryliu
 * date : 2016-05-24 10:14
 * description :
 * history :
 */
package repository

import(
	"github.com/jsix/gof"
	"sync"
	"encoding/gob"
	"os"
	"path"
)

var(
	mux sync.Mutex
)

// 处理错误
func handleError(err error) error {
	if err != nil && gof.CurrentApp.Debug() {
		gof.CurrentApp.Log().Println("[ Go2o][ Rep][ Error] -", err.Error())
	}
	return err
}

func unMarshalFromFile(file string,dst interface{})error {
	mux.Lock()
	fi, err := os.Open(file)
	if err == nil {
		enc := gob.NewDecoder(fi)
		err = enc.Decode(dst)
	}
	mux.Unlock()
	return err
}


func marshalToFile(file string,src interface{})error{
	//检测目录是否存在,不存在则创建目录
	dir := path.Dir(file)
	if _,err := os.Stat(dir); err == os.ErrNotExist{
		os.MkdirAll(dir,os.ModePerm)
	}
	mux.Lock()
	f, err := os.OpenFile(file,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		os.ModePerm)
	if err == nil{
		enc := gob.NewEncoder(f)
		err = enc.Encode(src)
	}
	mux.Unlock()
	return err
}