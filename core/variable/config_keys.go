/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package variable

const (
	Version = "version"

	//域名
	ServerDomain = "domain"
	ApiDomain    = "api_domain"

	//静态服务器
	StaticServer = "static_server"
	//图片服务器
	ImageServer = "image_server"

	//其他配置
	NoPicPath = "no_pic_path"

	// 上传保存目录
	UploadSaveDir = "upload_save_dir"

	// 支付
	Alipay_Partner = "payment_alipay_partner"
	Alipay_Key     = "payment_alipay_key"
	Alipay_Seller  = "payment_alipay_seller"

	// 邮箱
	SmtpHost    = "smtp_host"
	SmtpPort    = "smtp_port"
	SmtpCreUser = "smtp_user"
	SmtpCrePwd  = "smtp_pwd"
	SmtpFrom    = "smtp_from"
	//是否关闭系统发送邮件队列
	SystemMailQueueOff = "sys_mail_queue_off"
)

var (
	// 域名
	Domain string = "go2o.to2.net"
	// 配置路径
	ConfPath string = "./conf"
)
