/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package variable

const (
	Version = "version"

	// 经验值对金额的比例
	EXP_BIT = "exp_fee_bit"

	//域名
	ServerDomain = "server_domain"
	ApiDomain    = "api_domain"

	//静态服务器
	StaticServer = "static_server"
	//图片服务器
	ImageServer = "image_server"

	//数据库驱动名称
	DbDriver  = "db_driver"
	DbServer  = "db_server"
	DbPort    = "db_port"
	DbName    = "db_name"
	DbUsr     = "db_usr"
	DbPwd     = "db_pwd"
	DbCharset = "db_charset"

	//redis
	RedisHost     = "redis_host"
	RedisDb       = "redis_db"
	RedisMaxIdle  = "redis_max_idle"
	RedisIdleTout = "redis_idle_timeout"
	RedisPort     = "redis_port"

	//客户端socket server
	ClientSocketServer = "client_socket_server"

	//其他配置
	NoPicPath = "no_pic_path"

<<<<<<< HEAD
	// 上传保存目录
	UploadSaveDir = "upload_save_dir"

=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	//是否关闭系统发送邮件队列
	SystemMailQueueOff = "sys_mail_queue_off"
=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
)
