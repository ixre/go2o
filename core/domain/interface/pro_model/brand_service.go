package promodel

// 产品品牌
type (
	ProBrand struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 品牌名称
		Name string `db:"name"`
		// 品牌图片
		Image string `db:"image"`
		// 品牌网址
		SiteUrl string `db:"site_url"`
		// 介绍
		Intro string `db:"intro"`
		// 是否审核
		ReviewState int32 `db:"review_state"`
		// 加入时间
		CreateTime int64 `db:"create_time"`
	}

	// 产品模型与品牌关联
	ProModelBrand struct {
		ID       int32 `db:"id" pk:"yes" auto:"yes"`
		BrandId  int32 `db:"brand_id"`
		ProModel int32 `db:"pro_model"`
	}
)

// 品牌服务
type IBrandService interface {
	// 获取品牌
	Get(brandId int32) *ProBrand
	// 保存品牌
	SaveBrand(*ProBrand) (int32, error)
	// 删除品牌
	DeleteBrand(id int32) error
	// 获取所有(已审核的)品牌
	AllBrands() []*ProBrand
	// 获取关联的品牌编号
	Brands(proModel int32) []*ProBrand
}
