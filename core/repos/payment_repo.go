/**
 * Copyright 2015 @ 56x.net.
 * name : payment_repo.go
 * author : jarryliu
 * date : 2016-07-03 12:52
 * description :
 * history :
 */
package repos

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/registry"
	payImpl "github.com/ixre/go2o/core/domain/payment"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

var _ payment.IPaymentRepo = new(paymentRepoImpl)

type paymentRepoImpl struct {
	db.Connector
	Storage storage.Interface
	*payImpl.RepoBase
	_memberRepo      member.IMemberRepo
	_registryRepo    registry.IRegistryRepo
	_divideRepo      fw.Repository[payment.PayDivide]
	_mchRepo         merchant.IMerchantRepo
	_subMchRepo      fw.Repository[payment.PayMerchant]
	_tradeMethodRepo fw.Repository[payment.PayTradeData]
	_subMchMgr       payment.ISubMerchantManager
	_orm             orm.Orm
}

// MerchantRepo implements payment.IPaymentRepo.
func (p *paymentRepoImpl) MerchantRepo() fw.Repository[payment.PayMerchant] {
	return p._subMchRepo
}

// SubMerchantManager implements payment.IPaymentRepo.
func (p *paymentRepoImpl) SubMerchantManager() payment.ISubMerchantManager {
	if p._subMchMgr == nil {
		p._subMchMgr = payImpl.NewSubMerchantManager(p, p._mchRepo, p._memberRepo)
	}
	return p._subMchMgr
}

var payIntegrateAppDaoImplMapped = false

func NewPaymentRepo(sto storage.Interface, o orm.Orm,
	on fw.ORM,
	mmRepo member.IMemberRepo,
	mchRepo merchant.IMerchantRepo,
	registryRepo registry.IRegistryRepo) payment.IPaymentRepo {
	if !payIntegrateAppDaoImplMapped {
		_ = o.Mapping(payment.IntegrateApp{}, "pay_integrate_app")
		payIntegrateAppDaoImplMapped = true
	}
	//todo: 临时取消与orderRepo的循环依赖
	r := &paymentRepoImpl{
		Storage:          sto,
		Connector:        o.Connector(),
		_orm:             o,
		_memberRepo:      mmRepo,
		_mchRepo:         mchRepo,
		_registryRepo:    registryRepo,
		_divideRepo:      &fw.BaseRepository[payment.PayDivide]{ORM: on},
		_subMchRepo:      &fw.BaseRepository[payment.PayMerchant]{ORM: on},
		_tradeMethodRepo: &fw.BaseRepository[payment.PayTradeData]{ORM: on},
	}
	return r
}

// DivideRepo implements payment.IPaymentRepo.
func (p *paymentRepoImpl) DivideRepo() fw.Repository[payment.PayDivide] {
	return p._divideRepo
}

// 根据订单号获取支付单
func (p *paymentRepoImpl) GetPaymentBySalesOrderId(orderId int64) payment.IPaymentOrder {
	e := &payment.Order{}
	if p._orm.GetBy(e, "order_id= $1", orderId) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

// 根据订单号获取支付单
func (p *paymentRepoImpl) GetPaymentOrderByOrderNo(orderType int, orderNo string) payment.IPaymentOrder {
	e := &payment.Order{}
	if p._orm.GetBy(e, "out_order_no = $1 AND order_type= $2",
		orderNo, orderType) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

func (p *paymentRepoImpl) getPaymentOrderCk(id int) string {
	return fmt.Sprintf("go2o:repo:pay:order:%d", id)
}
func (p *paymentRepoImpl) getPaymentOrderCkByNo(orderNO string) string {
	return fmt.Sprintf("go2o:repo:pay:order:%s", orderNO)
}

// DeletePaymentOrder 拆分后删除父支付单
func (p *paymentRepoImpl) DeletePaymentOrder(id int) error {
	key := p.getPaymentOrderCk(id)
	p.Storage.Delete(key)
	return p._orm.DeleteByPk(payment.Order{}, id)
}

// DeletePaymentTradeData 删除支付单的支付数据
func (p *paymentRepoImpl) DeletePaymentTradeData(orderId int) error {
	_, err := p._orm.Delete(payment.PayTradeData{}, "order_id=$1", orderId)
	return err
}

// 根据编号获取支付单
func (p *paymentRepoImpl) GetPaymentOrderById(id int) payment.IPaymentOrder {
	if id <= 0 {
		return nil
	}
	e := &payment.Order{}
	k := p.getPaymentOrderCk(id)
	if err := p.Storage.Get(k, &e); err != nil {
		if p._orm.Get(id, e) != nil {
			return nil
		}
		p.Storage.SetExpire(k, *e, DefaultCacheSeconds)
	}
	return p.CreatePaymentOrder(e)
}

// 根据支付单号获取支付单
func (p *paymentRepoImpl) GetPaymentOrder(paymentNo string) payment.IPaymentOrder {
	k := p.getPaymentOrderCkByNo(paymentNo)
	id, err := p.Storage.GetInt(k)
	if err != nil {
		p.ExecScalar("SELECT id FROM pay_order where trade_no= $1", &id, paymentNo)
		if id == 0 {
			return nil
		}
		p.Storage.SetExpire(k, id, DefaultCacheSeconds*10)
	}
	return p.GetPaymentOrderById(id)
}

// 创建支付单
func (p *paymentRepoImpl) CreatePaymentOrder(
	o *payment.Order) payment.IPaymentOrder {
	return p.RepoBase.CreatePaymentOrder(o, p,
		p._memberRepo, p._registryRepo)
}

// 保存支付单
func (p *paymentRepoImpl) SavePaymentOrder(v *payment.Order) (int, error) {
	stat := v.Status
	if v.Id > 0 {
		stat = p.GetPaymentOrderById(v.Id).Get().Status
	}
	id, err := orm.Save(p._orm, v, v.Id)
	if err == nil {
		v.Id = id
		// 缓存订单
		p.Storage.SetExpire(p.getPaymentOrderCk(id), *v, DefaultCacheSeconds)
		// 缓存订单号与订单的关系
		p.Storage.SetExpire(p.getPaymentOrderCkByNo(v.TradeNo), v.Id, DefaultCacheSeconds*10)
		// 已经更改过状态,且为已成功,则推送到队列中
		if stat != v.Status && v.Status == payment.StateFinished {
			p.notifyPaymentFinish(v.Id)
		}
	}
	return id, err
}

// 通知支付单完成
func (p *paymentRepoImpl) notifyPaymentFinish(paymentOrderId int) error {
	//rc := core.GetRedisConn()
	//defer rc.Close()
	//_, err := rc.Do("RPUSH", variable.KvPaymentOrderFinishQueue, paymentOrderId)
	////log.Println("--  推送支付单成功", paymentOrderId,err)
	//return err
	return nil
}

// CheckTradeNoMatch 检查交易单号是否匹配
func (p *paymentRepoImpl) CheckTradeNoMatch(tradeNo string, id int) bool {
	i := 0
	p.Connector.ExecScalar("SELECT id FROM pay_order WHERE trade_no= $1 AND id<> $2", &i, tradeNo, id)
	return i == 0
}

func (p *paymentRepoImpl) GetTradeChannelItems(tradeNo string) []*payment.PayTradeData {
	return p._tradeMethodRepo.FindList(&fw.QueryOption{
		Limit: 10,
	}, "trade_no= ?", tradeNo)
}

func (p *paymentRepoImpl) SavePaymentTradeChan(tradeNo string, tradeChan *payment.PayTradeData) (int, error) {
	tradeChan.TradeNo = tradeNo
	id, err := orm.Save(p._orm, tradeChan, tradeChan.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PayTradeChan")
	}
	return id, err
}

func (p *paymentRepoImpl) GetMergePayOrders(mergeTradeNo string) []payment.IPaymentOrder {
	var tradeNo = ""
	var tradeNoArr []string
	// 查询支付单号
	p.Connector.Query("SELECT order_trade_no FROM pay_merge_order WHERE merge_trade_no= $1 LIMIT $2",
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&tradeNo)
				tradeNoArr = append(tradeNoArr, "'"+tradeNo+"'")
			}
		}, mergeTradeNo, 10)

	var arr = make([]payment.IPaymentOrder, 0)
	// 查询支付单
	if l := len(tradeNoArr); l > 0 {
		list := make([]*payment.Order, 0)
		p._orm.Select(&list, "trade_no IN ("+strings.Join(tradeNoArr, ",")+
			") LIMIT $1", len(tradeNoArr))
		for _, v := range list {
			arr = append(arr, p.CreatePaymentOrder(v))
		}
	}
	return arr
}

func (p *paymentRepoImpl) ResetMergePaymentOrders(tradeNos []string) error {
	buf := bytes.NewBuffer([]byte("("))
	for i, v := range tradeNos {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("'")
		buf.WriteString(v)
		buf.WriteString("'")
	}
	buf.WriteString(")")
	_, err := p._orm.Delete(&payment.MergeOrder{},
		"order_trade_no in "+buf.String())
	return err
}

func (p *paymentRepoImpl) SaveMergePaymentOrders(mergeTradeNo string, tradeNos []string) error {
	unix := time.Now().Unix()
	orm := p._orm
	for _, v := range tradeNos {
		order := &payment.MergeOrder{
			MergeTradeNo: mergeTradeNo,
			OrderTradeNo: v,
			SubmitTime:   unix,
		}
		if _, _, err := orm.Save(nil, order); err != nil {
			return err
		}
	}
	return nil
}

// FindAllIntegrateApp Select 集成支付应用
func (p *paymentRepoImpl) FindAllIntegrateApp() []*payment.IntegrateApp {
	list := make([]*payment.IntegrateApp, 0)
	err := p._orm.Select(&list, "1=1 ORDER BY sort_number")
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:IntegrateApp")
	}
	return list
}

// SaveIntegrateApp Save 集成支付应用
func (p *paymentRepoImpl) SaveIntegrateApp(v *payment.IntegrateApp) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:IntegrateApp")
	}
	return id, err
}

// DeleteIntegrateApp Delete 集成支付应用
func (p *paymentRepoImpl) DeleteIntegrateApp(primary interface{}) error {
	err := p._orm.DeleteByPk(payment.IntegrateApp{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:IntegrateApp")
	}
	return err
}

// GetAwaitCloseOrders 获取支付超时待关闭的订单
func (p *paymentRepoImpl) GetAwaitCloseOrders(lastId int, size int) []payment.IPaymentOrder {
	list := make([]*payment.Order, 0)
	err := p._orm.Select(&list, `
		status = 1 AND expires_time < $1 
		AND id > $2 ORDER BY id LIMIT $3`,
		time.Now().Unix(),
		lastId,
		size)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PaymentOrder")
	}
	arr := make([]payment.IPaymentOrder, len(list))
	for i, v := range list {
		arr[i] = p.CreatePaymentOrder(v)
	}
	return arr
}
