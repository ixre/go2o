package impl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
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

var _ proto.CheckServiceServer = new(checkService)

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
	s, err := c.store.GetString(c.getKey(token))
	if err != nil {
		return errors.New("操作超时,请重新进入")
	}
	data := parseCheckData(s)
	if now-int64(data.SendTime) < c.resendLimit {
		return errors.New("请勿在短时间内获取短信验证码")
	}
	return nil
}

// SaveData 3):存储验证码,默认5分钟有效, data应为phone和code的组合以保证phone和code是匹配的
func (c *CheckCodeVerifier) SaveData(token string, data *checkData) {
	now := time.Now().Unix()
	data.ExpiresTime = now + c.expiresSeconds
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

type checkService struct {
	repo         member.IMemberRepo
	registryRepo registry.IRegistryRepo
	store        storage.Interface
	*CheckCodeVerifier
	serviceUtil
	proto.UnimplementedCheckServiceServer
}

// NewCheckService 校验服务
func NewCheckService(repo member.IMemberRepo,
	registryRepo registry.IRegistryRepo,
	store storage.Interface,
) proto.CheckServiceServer {
	s := &checkService{
		repo:              repo,
		registryRepo:      registryRepo,
		store:             store,
		CheckCodeVerifier: NewCodeVerifier(store, "sys:go2o:reg:token",0,0),
	}
	return s
}

// SendCode 发送验证码
func (c *checkService) SendCode(_ context.Context, r *proto.SendCheckCodeRequest) (*proto.SendCheckCodeResponse, error) {
	if len(r.Token) == 0 {
		return &proto.SendCheckCodeResponse{
			ErrCode: 1002,
			ErrMsg:  "非法请求",
		}, nil
	}
	// 检查短信验证码是否频繁发送
	err := c.CheckCodeVerifier.CheckDuration(r.Token)
	if err != nil {
		return &proto.SendCheckCodeResponse{
			ErrCode: 1003,
			ErrMsg:  "频繁发送",
		}, nil
	}
	code := domain.NewCheckCode()
	// 保存校验码信息
	c.CheckCodeVerifier.SaveData(r.Token,
		&checkData{
			Account:   r.ReceptAccount,
			SendTime:  time.Now().Unix(),
			UserId:    int(r.UserId),
			CheckCode: code,
		})
	// 返回校验码
	return &proto.SendCheckCodeResponse{
		CheckCode: code,
	}, nil
}

// CompareCode implements proto.CheckServiceServer.
func (c *checkService) CompareCode(_ context.Context, r *proto.CompareCheckCodeRequest) (*proto.CompareCheckCodeResponse, error) {
	userId, err := c.CheckCodeVerifier.Validate(r.Token, r.ReceptAccount, r.CheckCode)
	if err != nil {
		return &proto.CompareCheckCodeResponse{
			ErrCode: 1001,
			ErrMsg:  err.Error(),
		}, nil
	}
	return &proto.CompareCheckCodeResponse{
		UserId: int64(userId),
	}, nil
}
