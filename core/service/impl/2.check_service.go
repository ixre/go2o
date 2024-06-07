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
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
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

// 验证码验证器 todo: 移动到go2o
type CheckCodeVerifier struct {
	store             storage.Interface
	storageKey     string
	expiresSeconds int64
	resendLimit    int64
}

// NewCodeVerifier 创建一个代码验证器, 默认有效期为5分钟, 限制重复发送的时间为2分钟
func NewCodeVerifier(store storage.Interface, storageKey string, minites int, resendLimit int) *CheckCodeVerifier {
	return &CheckCodeVerifier{
		store:             store,
		storageKey:     storageKey,
		expiresSeconds: int64(math.Max(float64(minites), 5)) * 60,
		resendLimit:    int64(math.Max(float64(resendLimit), 60)),
	}
}

func (c *CheckCodeVerifier) getKey(token string) string {
	return fmt.Sprintf("%s-%s", c.storageKey, token)
}

// PrepareToken 1):准备token,30分钟有效, 在获取token时需要验证用户是否能频繁的拿到token计数
func (c *CheckCodeVerifier) PrepareToken(aud string) (string, error) {
	rd := util.RandString(10)
	err := c.store.SetExpire(c.getKey(rd), "0|0|-|0", 1800) // 发送时间|有效时间|验证码|用户编号
	return rd, err
}

// CheckDuration 2):检查短信验证码是否频繁发送
func (c *CheckCodeVerifier) CheckDuration(token string) error {
	if len(token) == 0 {
		return errors.New("token不能为空")
	}
	now := time.Now().Unix()
	s, err := c.store.GetString(c.getKey(token))
	unix, err2 := strconv.Atoi(strings.Split(s, "|")[0])
	if err != nil || err2 != nil {
		return errors.New("操作超时,请重新进入")
	}
	if now-int64(unix) < c.resendLimit {
		return errors.New("请勿在短时间内获取短信验证码")
	}
	return nil
}

// SaveSendData 3):存储验证码,默认5分钟有效, data应为phone和code的组合以保证phone和code是匹配的
func (c *CheckCodeVerifier) SaveSendData(token string, data string, userId int) {
	now := time.Now().Unix()
	expires := now + c.expiresSeconds
	v := fmt.Sprintf("%d|%d|%s|%d", now, expires, data, userId)
	_ = c.store.SetExpire(c.getKey(token), v, 1800)
}

// MatchCode 4):验证验证码是否正确
func (c *CheckCodeVerifier) Valid(token string, data string) (int, error) {
	if len(token) == 0 {
		return 0, errors.New("非法请求")
	}
	s, err := c.store.GetString(c.getKey(token))
	if err != nil {
		return 0, errors.New("验证码已过期")
	}
	arr := strings.Split(s, "|")
	if arr[2] != data {
		return 0, errors.New("验证码不正确")
	}
	now := time.Now().Unix()
	unix, err2 := strconv.Atoi(arr[1])
	if err2 != nil || now-int64(unix) > c.expiresSeconds {
		return 0, errors.New("验证码已过期")
	}
	return strconv.Atoi(arr[3])
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

// CompareCode implements proto.CheckServiceServer.
func (c *checkService) CompareCode(context.Context, *proto.CompareCheckCodeRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// NewCheckService 校验服务
func NewCheckService(repo member.IMemberRepo,
	registryRepo registry.IRegistryRepo,
	store storage.Interface,
) proto.CheckServiceServer {
	s := &checkService{
		repo:         repo,
		registryRepo: registryRepo,
		store:        store,
		CheckCodeVerifier: NewCodeVerifier(store,"sys:go2o:reg:token"),
	}
	return s
}

// SendCode 发送验证码
func (c *checkService) SendCode(_ context.Context, r *proto.SendCheckCodeRequest) (*proto.SendCheckCodeResponse, error) {
// 检查短信验证码是否频繁发送
err := m.CheckDuration(token)
if err != nil {
	return m.ErrorJson(ctx, 7, err)
}
	panic("unimplemented")
}
