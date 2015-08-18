/**
 * Copyright 2015 @ z3q.net.
 * name : kv_manager
 * author : jarryliu
 * date : 2015-07-26 22:06
 * description :
 * history :
 */
package partner

const (
	KeyMssTplIdOfProfileComplete string = "mss_profile_complete_mail_tpl"
)

type IKvManager interface {
	// 获取
	Get(k string) string
	// 获取int类型的键值
	GetInt(k string) int
	// 设置
	Set(k, v string)
	// 获取多项
	Gets(k []string) map[string]string
	// 设置多项
	Sets(map[string]string) error
	// 根据关键字获取字典
	GetsByChar(keyword string) map[string]string
}
