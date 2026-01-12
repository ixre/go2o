/**
 * Copyright (C) 2007-2026 fze.NET, All rights reserved.
 *
 * name: repository.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2026-01-12 15:26:02
 * description: 导出OrmMapping函数，用于初始化数据库表结构
 * history:
 */
package repos

import (
	"github.com/ixre/go2o/internal/core/repos"
	"github.com/ixre/go2o/pkg/initial/provide"
)

func OrmMapping() {
	repos.OrmMapping(provide.GetOrmInstance())
}

// notes: 不需要导出，在应用中直接使用inject.GetSystemRepo()获取实例

// var NewSystemRepo = repos.NewSystemRepo
// var NewRegistryRepo = repos.NewRegistryRepo

// var NewProModelRepo = repos.NewProModelRepo
// var NewValueRepo = repos.NewValueRepo
// var NewUserRepo = repos.NewUserRepo
// var NewWalletRepo = repos.NewWalletRepo
// var NewNotifyRepo = repos.NewNotifyRepo
// var NewMssRepo = repos.NewMssRepo
// var NewExpressRepo = repos.NewExpressRepo
// var NewShipmentRepo = repos.NewShipmentRepo
// var NewMemberRepo = repos.NewMemberRepo

// var NewItemWholesaleRepo = repos.NewItemWholesaleRepo
// var NewProductRepo = repos.NewProductRepo
// var NewCategoryRepo = repos.NewCategoryRepo
// var NewShopRepo = repos.NewShopRepo
// var NewGoodsItemRepo = repos.NewGoodsItemRepo
// var NewAfterSalesRepo = repos.NewAfterSalesRepo
// var NewCartRepo = repos.NewCartRepo
// var NewArticleRepo = repos.NewArticleRepo
// var NewMerchantRepo = repos.NewMerchantRepo
// var NewOrderRepo = repos.NewOrderRepo

// var NewPaymentRepo = repos.NewPaymentRepo
// var NewPromotionRepo = repos.NewPromotionRepo
// var NewStationRepo = repos.NewStationRepo
// var NewTagSaleRepo = repos.NewTagSaleRepo
// var NewWholesaleRepo = repos.NewWholesaleRepo
// var NewPersonFinanceRepository = repos.NewPersonFinanceRepository
// var NewDeliverRepo = repos.NewDeliverRepo
// var NewAdvertisementRepo = repos.NewAdvertisementRepo
// var NewJobRepository = repos.NewJobRepository
// var NewStaffRepo = repos.NewStaffRepo

// var NewApprovalRepository = repos.NewApprovalRepository
// var NewInvoiceTenantRepo = repos.NewInvoiceTenantRepo
// var NewPageRepo = repos.NewPageRepo
// var NewArticleCategoryRepo = repos.NewArticleCategoryRepo
// var NewChatRepo = repos.NewChatRepo
// var NewWorkorderRepo = repos.NewWorkorderRepo
// var NewRbacRepo = repos.NewRbacRepo
// var NewSysAppRepo = repos.NewSysAppRepo
