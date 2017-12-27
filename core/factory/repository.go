package factory

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/storage"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/personfinance"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/interface/wallet"
	"go2o/core/repository"
)

var (
	Repo *RepoFactory = &RepoFactory{}
)

type RepoFactory struct {
	proMRepo   promodel.IProModelRepo
	valueRepo  valueobject.IValueRepo
	userRepo   user.IUserRepo
	notifyRepo notify.INotifyRepo
	mssRepo    mss.IMssRepo

	expressRepo express.IExpressRepo
	shipRepo    shipment.IShipmentRepo
	memberRepo  member.IMemberRepo
	productRepo product.IProductRepo
	itemWsRepo  item.IItemWholesaleRepo
	catRepo     product.ICategoryRepo
	itemRepo    item.IGoodsItemRepo
	tagSaleRepo item.ISaleLabelRepo
	promRepo    promotion.IPromotionRepo

	shopRepo          shop.IShopRepo
	wholesaleRepo     wholesaler.IWholesaleRepo
	mchRepo           merchant.IMerchantRepo
	cartRepo          cart.ICartRepo
	personFinanceRepo personfinance.IPersonFinanceRepository
	deliveryRepo      delivery.IDeliveryRepo
	contentRepo       content.IContentRepo
	adRepo            ad.IAdRepo
	orderRepo         order.IOrderRepo
	paymentRepo       payment.IPaymentRepo
	asRepo            afterSales.IAfterSalesRepo

	walletRepo wallet.IWalletRepo
}

func (r *RepoFactory) Init(db db.Connector, sto storage.Interface, confPath string) *RepoFactory {
	Repo = r
	orm := db.GetOrm()
	/** Repository **/
	r.proMRepo = repository.NewProModelRepo(db, orm)
	r.valueRepo = repository.NewValueRepo(confPath, db, sto)
	r.userRepo = repository.NewUserRepo(db)
	r.notifyRepo = repository.NewNotifyRepo(db)
	r.mssRepo = repository.NewMssRepo(db, r.notifyRepo, r.valueRepo)
	r.expressRepo = repository.NewExpressRepo(db, r.valueRepo)
	r.shipRepo = repository.NewShipmentRepo(db, r.expressRepo)
	r.memberRepo = repository.NewMemberRepo(sto, db, r.mssRepo, r.valueRepo)
	r.productRepo = repository.NewProductRepo(db, r.proMRepo, r.valueRepo)
	r.itemWsRepo = repository.NewItemWholesaleRepo(db)
	r.catRepo = repository.NewCategoryRepo(db, r.valueRepo, sto)
	r.itemRepo = repository.NewGoodsItemRepo(db, r.catRepo, r.productRepo,
		r.proMRepo, r.itemWsRepo, r.expressRepo, r.valueRepo)
	r.tagSaleRepo = repository.NewTagSaleRepo(db, r.valueRepo)
	r.promRepo = repository.NewPromotionRepo(db, r.itemRepo, r.memberRepo)

	//afterSalesRepo := repository.NewAfterSalesRepo(db)
	r.walletRepo = repository.NewWalletRepo(db)
	r.shopRepo = repository.NewShopRepo(db, sto, r.valueRepo)
	r.wholesaleRepo = repository.NewWholesaleRepo(db)
	r.mchRepo = repository.NewMerchantRepo(db, sto, r.wholesaleRepo,
		r.itemRepo, r.shopRepo, r.userRepo, r.memberRepo, r.mssRepo, r.walletRepo, r.valueRepo)
	r.cartRepo = repository.NewCartRepo(db, r.memberRepo, r.mchRepo, r.itemRepo)
	r.personFinanceRepo = repository.NewPersonFinanceRepository(db, r.memberRepo)
	r.deliveryRepo = repository.NewDeliverRepo(db)
	r.contentRepo = repository.NewContentRepo(db)
	r.adRepo = repository.NewAdvertisementRepo(db, sto)
	r.orderRepo = repository.NewOrderRepo(sto, db, r.mchRepo, nil,
		r.productRepo, r.cartRepo, r.itemRepo, r.promRepo, r.memberRepo,
		r.deliveryRepo, r.expressRepo, r.shipRepo, r.valueRepo)
	r.paymentRepo = repository.NewPaymentRepo(sto, db, r.memberRepo, r.orderRepo, r.valueRepo)
	r.asRepo = repository.NewAfterSalesRepo(db, r.orderRepo, r.memberRepo, r.paymentRepo)

	// 解决依赖
	r.orderRepo.(*repository.OrderRepImpl).SetPaymentRepo(r.paymentRepo)
	// 初始化数据
	r.memberRepo.GetManager().GetAllBuyerGroups()
	return r
}

func (r *RepoFactory) GetProModelRepo() promodel.IProModelRepo {
	return r.proMRepo
}

func (r *RepoFactory) GetValueRepo() valueobject.IValueRepo {
	return r.valueRepo
}
func (r *RepoFactory) GetUserRepo() user.IUserRepo {
	return r.userRepo
}
func (r *RepoFactory) GetNotifyRepo() notify.INotifyRepo {
	return r.notifyRepo
}
func (r *RepoFactory) GetMssRepo() mss.IMssRepo {
	return r.mssRepo
}
func (r *RepoFactory) GetExpressRepo() express.IExpressRepo {
	return r.expressRepo
}
func (r *RepoFactory) GetShipmentRepo() shipment.IShipmentRepo {
	return r.shipRepo
}
func (r *RepoFactory) GetMemberRepo() member.IMemberRepo {
	return r.memberRepo
}
func (r *RepoFactory) GetProductRepo() product.IProductRepo {
	return r.productRepo
}
func (r *RepoFactory) GetItemWholesaleRepo() item.IItemWholesaleRepo {
	return r.itemWsRepo
}
func (r *RepoFactory) GetCategoryRepo() product.ICategoryRepo {
	return r.catRepo
}
func (r *RepoFactory) GetItemRepo() item.IGoodsItemRepo {
	return r.itemRepo
}
func (r *RepoFactory) GetSaleLabelRepo() item.ISaleLabelRepo {
	return r.tagSaleRepo
}
func (r *RepoFactory) GetPromotionRepo() promotion.IPromotionRepo {
	return r.promRepo
}
func (r *RepoFactory) GetShopRepo() shop.IShopRepo {
	return r.shopRepo
}
func (r *RepoFactory) GetWholesaleRepo() wholesaler.IWholesaleRepo {
	return r.wholesaleRepo
}
func (r *RepoFactory) GetMerchantRepo() merchant.IMerchantRepo {
	return r.mchRepo
}
func (r *RepoFactory) GetCartRepo() cart.ICartRepo {
	return r.cartRepo
}
func (r *RepoFactory) GetPersonFinanceRepository() personfinance.IPersonFinanceRepository {
	return r.personFinanceRepo
}
func (r *RepoFactory) GetDeliveryRepo() delivery.IDeliveryRepo {
	return r.deliveryRepo
}
func (r *RepoFactory) GetContentRepo() content.IContentRepo {
	return r.contentRepo
}
func (r *RepoFactory) GetAdRepo() ad.IAdRepo {
	return r.adRepo
}
func (r *RepoFactory) GetOrderRepo() order.IOrderRepo {
	return r.orderRepo
}
func (r *RepoFactory) GetPaymentRepo() payment.IPaymentRepo {
	return r.paymentRepo
}
func (r *RepoFactory) GetAfterSalesRepo() afterSales.IAfterSalesRepo {
	return r.asRepo
}
func (r *RepoFactory) GetWalletRepo() wallet.IWalletRepo {
	return r.walletRepo
}
