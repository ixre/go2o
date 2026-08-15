//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/pkg/interface/domain/ad"
	afterSales "github.com/ixre/go2o/pkg/interface/domain/aftersales"
	"github.com/ixre/go2o/pkg/interface/domain/approval"
	"github.com/ixre/go2o/pkg/interface/domain/cart"
	"github.com/ixre/go2o/pkg/interface/domain/chat"
	"github.com/ixre/go2o/pkg/interface/domain/content"
	"github.com/ixre/go2o/pkg/interface/domain/delivery"
	"github.com/ixre/go2o/pkg/interface/domain/express"
	"github.com/ixre/go2o/pkg/interface/domain/invoice"
	"github.com/ixre/go2o/pkg/interface/domain/item"
	"github.com/ixre/go2o/pkg/interface/domain/job"
	"github.com/ixre/go2o/pkg/interface/domain/member"
	"github.com/ixre/go2o/pkg/interface/domain/merchant"
	"github.com/ixre/go2o/pkg/interface/domain/merchant/shop"
	"github.com/ixre/go2o/pkg/interface/domain/merchant/staff"
	"github.com/ixre/go2o/pkg/interface/domain/merchant/user"
	"github.com/ixre/go2o/pkg/interface/domain/merchant/wholesaler"
	mss "github.com/ixre/go2o/pkg/interface/domain/message"
	"github.com/ixre/go2o/pkg/interface/domain/order"
	"github.com/ixre/go2o/pkg/interface/domain/payment"
	"github.com/ixre/go2o/pkg/interface/domain/personfinance"
	promodel "github.com/ixre/go2o/pkg/interface/domain/pro_model"
	"github.com/ixre/go2o/pkg/interface/domain/product"
	"github.com/ixre/go2o/pkg/interface/domain/promotion"
	rbac "github.com/ixre/go2o/pkg/interface/domain/rabc"
	"github.com/ixre/go2o/pkg/interface/domain/registry"
	"github.com/ixre/go2o/pkg/interface/domain/shipment"
	"github.com/ixre/go2o/pkg/interface/domain/sys"
	"github.com/ixre/go2o/pkg/interface/domain/valueobject"
	"github.com/ixre/go2o/pkg/interface/domain/wallet"
	"github.com/ixre/go2o/pkg/interface/domain/work/workorder"
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

func GetStationRepo() sys.IStationRepo {
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
func GetContentRepo() content.IArticleRepo {
	panic(wire.Build(provideSets))
}

func GetArticleCategoryRepo() content.IArticleCategoryRepo {
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

func GetPageRepo() content.IPageRepo {
	panic(wire.Build(provideSets))
}

func GetInvoiceTenantRepo() invoice.IInvoiceRepo {
	panic(wire.Build(provideSets))
}

func GetChatRepo() chat.IChatRepository {
	panic(wire.Build(provideSets))
}

func GetWorkorderRepo() workorder.IWorkorderRepo {
	panic(wire.Build(provideSets))
}

func GetApprovalRepo() approval.IApprovalRepository {
	panic(wire.Build(provideSets))
}

func GetRbacRepo() rbac.IRbacRepository {
	panic(wire.Build(provideSets))
}

func GetLogRepo() sys.IApplicationRepository {
	panic(wire.Build(provideSets))
}
