//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/core/domain/interface/ad"
	afterSales "github.com/ixre/go2o/core/domain/interface/aftersales"
	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/domain/interface/delivery"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/job"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/merchant/user"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/personfinance"
	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/station"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
)

// 解决依赖
//r.orderRepo.(*OrderRepImpl).SetPaymentRepo(r.paymentRepo)
// 初始化数据
//r.memberRepo.GetManager().GetAllBuyerGroups()

func GetProModelRepo() promodel.IProductModelRepo {
	panic(wire.Build(provideSets))
}

// 获取系统仓库
func GetSystemRepo() sys.ISystemRepo {
	panic(wire.Build(provideSets))
}

func GetValueRepo() valueobject.IValueRepo {
	panic(wire.Build(provideSets))

}
func GetUserRepo() user.IUserRepo {
	panic(wire.Build(provideSets))

}
func GetNotifyRepo() mss.INotifyRepo {
	panic(wire.Build(provideSets))

}
func GetMessageRepo() mss.IMessageRepo {
	panic(wire.Build(provideSets))

}
func GetExpressRepo() express.IExpressRepo {
	panic(wire.Build(provideSets))
}

func GetShipmentRepo() shipment.IShipmentRepo {
	panic(wire.Build(provideSets))

}

func GetMemberRepo() member.IMemberRepo {
	panic(wire.Build(provideSets))

}
func GetProductRepo() product.IProductRepo {
	panic(wire.Build(provideSets))

}
func GetItemWholesaleRepo() item.IItemWholesaleRepo {
	panic(wire.Build(provideSets))

}

func GetCategoryRepo() product.ICategoryRepo {
	panic(wire.Build(provideSets))

}
func GetItemRepo() item.IItemRepo {
	panic(wire.Build(provideSets))

}
func GetSaleLabelRepo() item.ISaleLabelRepo {
	panic(wire.Build(provideSets))

}
func GetPromotionRepo() promotion.IPromotionRepo {
	panic(wire.Build(provideSets))

}
func GetShopRepo() shop.IShopRepo {
	panic(wire.Build(provideSets))

}

func GetWholesaleRepo() wholesaler.IWholesaleRepo {
	panic(wire.Build(provideSets))

}

func GetStationRepo() station.IStationRepo {
	panic(wire.Build(provideSets))

}

func GetMerchantRepo() merchant.IMerchantRepo {
	panic(wire.Build(provideSets))

}

func GetCartRepo() cart.ICartRepo {
	panic(wire.Build(provideSets))

}
func GetPersonFinanceRepository() personfinance.IPersonFinanceRepository {
	panic(wire.Build(provideSets))

}
func GetDeliveryRepo() delivery.IDeliveryRepo {
	panic(wire.Build(provideSets))

}
func GetContentRepo() content.IArchiveRepo {
	panic(wire.Build(provideSets))

}
func GetAdRepo() ad.IAdRepo {
	panic(wire.Build(provideSets))

}
func GetOrderRepo() order.IOrderRepo {
	panic(wire.Build(provideSets))

}

func GetPaymentRepo() payment.IPaymentRepo {
	panic(wire.Build(provideSets))

}
func GetAfterSalesRepo() afterSales.IAfterSalesRepo {
	panic(wire.Build(provideSets))

}
func GetWalletRepo() wallet.IWalletRepo {
	panic(wire.Build(provideSets))

}

func GetRegistryRepo() registry.IRegistryRepo {
	panic(wire.Build(provideSets))

}

func GetJobRepo() job.IJobRepo {
	panic(wire.Build(provideSets))
}

func GetStaffRepo() staff.IStaffRepo {
	panic(wire.Build(provideSets))
}
