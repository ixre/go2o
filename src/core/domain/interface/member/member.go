/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:49
 * description :
 * history :
 */

package member

import "go2o/src/core/domain/interface/valueobject"

const (
	StateStopped = 0 //已停用
	StateOk      = 1 //正常
	BankNoLock   = 0
	BankLocked   = 1
)

type (
	ValueMember struct {
		Id  int    `db:"id" auto:"yes" pk:"yes"`
		Usr string `db:"usr"`
		Pwd string `db:"Pwd"`
		// 交易密码
		TradePwd string `db:"trade_pwd"`
		// 姓名
		Name string `db:"name"`
		// 经验值
		Exp int `db:"exp"`
		// 等级
		Level int `db:"level"`

		Sex      int    `db:"sex"`
		Avatar   string `db:"avatar"`
		BirthDay string `db:"birthday"`
		Phone    string `db:"phone"`
		Address  string `db:"address"`
		Im       string `db:"im"`
		Email    string `db:"email"`
		// 邀请码
		InvitationCode string `db:"invitation_code"`
		RegFrom        string `db:"reg_from"`
		RegIp          string `db:"reg_ip"`
		State          int    `db:"state"`
		RegTime        int64  `db:"reg_time"`
		Remark         string `db:"remark"` //备注
		Ext1           string `db:"ext_1"`  // 扩展1
		Ext2           string `db:"ext_2"`  // 扩展2
		Ext3           string `db:"ext_3"`  // 扩展3
		Ext4           string `db:"ext_4"`  // 扩展4
		Ext5           string `db:"ext_4"`  // 扩展5
		Ext6           string `db:"ext_4"`  // 扩展6
		LastLoginTime  int64  `db:"last_login_time"`
		UpdateTime     int64  `db:"update_time"`
		DynamicToken   string `db:"-"` // 动态令牌，用于登陆或API调用
		TimeoutTime    int64  `db:"-"` // 超时时间
	}

	// 银行卡信息,因为重要且非频繁更新的数据
	// 所以需要用IsLocked来标记是否锁定
	BankInfo struct {
		MemberId    int    `db:"member_id" pk:"yes"`
		Name        string `db:"name"`
		Account     string `db:"account"`
		AccountName string `db:"account_name"`
		Network     string `db:"network"`
		State       int    `db:"state"`
		IsLocked    int    `db:"is_locked"`
		UpdateTime  int64  `db:"update_time"`
	}

	IMember interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 获取值
		GetValue() ValueMember

		// 邀请管理
		Invitation() IInvitationManager

		// 设置值
		SetValue(*ValueMember) error

		// 获取账户
		GetAccount() IAccount

		// 锁定会员
		Lock() error

		// 解锁会员
		Unlock() error

		// 资料是否完善
		ProfileCompleted() bool

		// 获取提现银行信息
		GetBank() BankInfo

		// 保存提现银行信息,保存后将锁定
		SaveBank(*BankInfo) error

		// 解锁提现银行卡信息
		UnlockBank() error

		//　保存积分记录
		SaveIntegralLog(*IntegralLog) error

		// 增加经验值
		AddExp(exp int) error

		// 获取等级
		GetLevel() *valueobject.MemberLevel

		//　增加积分
		// todo:partnerId 不需要
		AddIntegral(partnerId int, backType int, integral int, log string) error

		// 获取关联的会员
		GetRelation() *MemberRelation

		// 更新会员绑定
		SaveRelation(r *MemberRelation) error

		// 更换用户名
		ChangeUsr(string) error

		// 保存
		Save() (int, error)

		// 修改密码,旧密码可为空
		ModifyPassword(newPwd, oldPwd string) error

		// 修改交易密码，旧密码可为空
		ModifyTradePassword(newPwd, oldPwd string) error

		// 创建配送地址
		CreateDeliver(*DeliverAddress) (IDeliver, error)

		// 获取配送地址
		GetDeliverAddress() []IDeliver

		// 获取配送地址
		GetDeliver(int) IDeliver

		// 删除配送地址
		DeleteDeliver(int) error
	}
)
