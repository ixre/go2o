package io

//　文件系统挂载接口
// 利用本接口可以将文件存放分类归档，同时便于迁移
// IsValid提供了检查挂载点的功能，
// 如果挂在路径不存在或路径不符合规则返回false
// Combine　提供了合并返回完整路径的功能。
type FsMount interface {
	IsValid() (b bool, err error)
	Combine(path string) string
}
