package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/message/notify"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
)

/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2020-09-05 20:14
 * description :
 * history :
 */

var _ proto.CheckServiceServer = new(checkServiceImpl)

// 校验码信息
type checkData struct {
	// 发送时间
	SendTime int64
	// 过期时间
	ExpiresTime int64
	// 校验码
	CheckCode string
	// 账号
	Account string
	// 用户编号
	UserId int
}

func (c checkData) String() string {
	return fmt.Sprintf("%d|%d|%s|%s|%d",
		c.SendTime,
		c.ExpiresTime,
		c.CheckCode,
		c.Account,
		c.UserId)
}

// 转换校验信息
func parseCheckData(s string) *checkData {
	arr := strings.Split(s, "|")
	if len(arr) != 5 {
		return nil
	}
	sendTime, _ := strconv.ParseInt(arr[0], 10, 64)
	expiresTime, _ := strconv.ParseInt(arr[1], 10, 64)
	userId, _ := strconv.Atoi(arr[4])
	return &checkData{
		SendTime:    sendTime,
		ExpiresTime: expiresTime,
		CheckCode:   arr[2],
		Account:     arr[3],
		UserId:      userId,
	}
}

// 验证码验证器 todo: 移动到go2o
type CheckCodeVerifier struct {
	store          storage.Interface
	storageKey     string
	expiresSeconds int64
	resendLimit    int64
}

// NewCodeVerifier 创建一个代码验证器, 默认有效期为5分钟, 限制重复发送的时间为2分钟
func NewCodeVerifier(store storage.Interface, storageKey string, minites int, resendLimit int) *CheckCodeVerifier {
	return &CheckCodeVerifier{
		store:          store,
		storageKey:     storageKey,
		expiresSeconds: int64(math.Max(float64(minites), 5)) * 60,
		resendLimit:    int64(math.Max(float64(resendLimit), 60)),
	}
}

func (c *CheckCodeVerifier) getKey(token string) string {
	return fmt.Sprintf("%s-%s", c.storageKey, token)
}

// CheckDuration 1):检查短信验证码是否频繁发送
func (c *CheckCodeVerifier) CheckDuration(token string) error {
	if len(token) == 0 {
		return errors.New("token不能为空")
	}
	now := time.Now().Unix()
	s, _ := c.store.GetString(c.getKey(token))
	if len(s) > 0 {
		data := parseCheckData(s)
		if now-int64(data.SendTime) < c.resendLimit {
			return errors.New("请勿在短时间内获取短信验证码")
		}
	}
	return nil
}

// SaveData 3):存储验证码,默认5分钟有效, data应为phone和code的组合以保证phone和code是匹配的
func (c *CheckCodeVerifier) SaveData(token string, data *checkData) {
	now := time.Now().Unix()
	data.SendTime = now
	if data.ExpiresTime-data.SendTime < 60 {
		// 如果过期时间小于60秒,则设置为5分钟
		data.ExpiresTime = now + 5*60
	}
	_ = c.store.SetExpire(c.getKey(token), data.String(), 12*3600)
}

// MatchCode 4):验证验证码是否正确
// [userId]返回用户编号
func (c *CheckCodeVerifier) Validate(token string, receptAccount string, checkCode string) (int, error) {
	if len(token) == 0 {
		return 0, errors.New("非法请求")
	}
	s, err := c.store.GetString(c.getKey(token))
	if err != nil {
		return 0, errors.New("验证码已过期")
	}
	data := parseCheckData(s)
	if data.Account != receptAccount || data.CheckCode != checkCode {
		return 0, errors.New("验证码不正确")
	}
	now := time.Now().Unix()
	if now-int64(data.ExpiresTime) > c.expiresSeconds {
		return 0, errors.New("验证码已过期")
	}
	return data.UserId, nil
}

// Destory 5):销毁本次操作令牌
func (c *CheckCodeVerifier) Destory(token string) {
	c.store.Delete(c.getKey(token))
}

var (
	_checkServiceInstance proto.CheckServiceServer
	_checkServiceOnce     sync.Once
)

// checkServiceImpl 校验服务
type checkServiceImpl struct {
	repo         member.IMemberRepo
	registryRepo registry.IRegistryRepo
	notifyRepo   notify.INotifyRepo
	store        storage.Interface
	*CheckCodeVerifier
	proto.UnimplementedCheckServiceServer
}

// NewCheckService 创建校验服务实现
func NewCheckService(repo member.IMemberRepo,
	notifyRepo notify.INotifyRepo,
	registryRepo registry.IRegistryRepo,
	store storage.Interface,
) proto.CheckServiceServer {
	_checkServiceOnce.Do(func() {
		_checkServiceInstance = &checkServiceImpl{
			repo:              repo,
			registryRepo:      registryRepo,
			store:             store,
			notifyRepo:        notifyRepo,
			CheckCodeVerifier: NewCodeVerifier(store, "go2o:checkcode:token", 0, 0),
		}
	})
	return _checkServiceInstance
}

// SendCode 发送验证码
func (c *checkServiceImpl) SendCode(_ context.Context, r *proto.SendCheckCodeRequest) (*proto.SendCheckCodeResponse, error) {
	if len(r.Token) == 0 {
		return &proto.SendCheckCodeResponse{
			ErrCode: 1001,
			ErrMsg:  "非法请求",
		}, nil
	}
	if len(r.Operation) == 0 {
		return &proto.SendCheckCodeResponse{
			ErrCode: 1002,
			ErrMsg:  "操作名称不能为空",
		}, nil
	}
	if len(r.MsgTemplateId) == 0 {
		return &proto.SendCheckCodeResponse{
			ErrCode: 1003,
			ErrMsg:  "模板不能为空",
		}, nil
	}
	// 检查短信验证码是否频繁发送
	err := c.CheckCodeVerifier.CheckDuration(r.Token)
	if err != nil {
		return &proto.SendCheckCodeResponse{
			ErrCode: 1004,
			ErrMsg:  err.Error(),
		}, nil
	}
	if r.Effective <= 0 {
		// 默认5分钟有效
		r.Effective = 5
	}
	code := domain.NewCheckCode()
	// 发送验证码,如果失败,则输出错误信息
	if err := c.notifyCheckCode(code, r); err != nil {
		log.Println("[ Go2o][ Error]: 发送验证码失败:", err.Error())
	}
	// 保存校验码信息
	unix := time.Now().Unix()
	c.CheckCodeVerifier.SaveData(r.Token,
		&checkData{
			Account:     r.ReceptAccount,
			ExpiresTime: unix + int64(r.Effective*60),
			UserId:      int(r.UserId),
			CheckCode:   code,
		})
	// 返回校验码
	return &proto.SendCheckCodeResponse{
		CheckCode: code,
	}, nil
}

// CompareCode implements proto.CheckServiceServer.
func (c *checkServiceImpl) CompareCode(_ context.Context, r *proto.CompareCheckCodeRequest) (*proto.CompareCheckCodeResponse, error) {
	if len(r.Token) == 0 {
		return &proto.CompareCheckCodeResponse{
			ErrCode: 1001,
			ErrMsg:  "非法请求",
		}, nil
	}
	userId, err := c.CheckCodeVerifier.Validate(r.Token, r.ReceptAccount, r.CheckCode)
	if err != nil {
		return &proto.CompareCheckCodeResponse{
			ErrCode: 1002,
			ErrMsg:  err.Error(),
		}, nil
	}
	if r.ResetIfOk {
		// 如果验证成功,则重置令牌
		c.CheckCodeVerifier.Destory(r.Token)
	}
	return &proto.CompareCheckCodeResponse{
		UserId: int64(userId),
	}, nil
}

func (c *checkServiceImpl) notifyCheckCode(code string, r *proto.SendCheckCodeRequest) error {
	// 创建参数
	data := []string{r.Operation, code, strconv.Itoa(int(r.Effective))}
	// 构造并发送短信
	mg := c.notifyRepo.Manager()
	return mg.SendPhoneMessage(r.ReceptAccount, notify.PhoneMessage(""), data, r.MsgTemplateId)
}
