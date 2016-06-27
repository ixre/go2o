/**
 * Copyright 2015 @ z3q.net.
 * name : article
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

// 文章
type (
	// 文章
	IArticle interface {
		// 获取领域编号
		GetDomainId() int

		// 获取值
		GetValue() Article

		// 设置值
		SetValue(*Article) error

		// 栏目
		Category() ICategory

		// 保存文章
		Save() (int, error)
	}

	// 栏目
	ICategory interface {
		// 获取领域编号
		GetDomainId() int

		// 获取值
		GetValue() ArticleCategory

		// 获取文章数量
		ArticleNum() int

		// 设置值
		SetValue(*ArticleCategory) error

		// 保存
		Save() (int, error)
	}

	// 文章服务bn
	IArticleManager interface {
		// 获取栏目
		GetCategory(id int) ICategory

		// 获取所有的栏目
		GetAllCategory() []ICategory

		// 创建栏目
		CreateCategory(*ArticleCategory) ICategory

		// 删除栏目
		DelCategory(id int) error

		// 创建文章
		CreateArticle(*Article) IArticle

		// 获取文章
		GetArticle(id int) IArticle

		// 获取文章列表
		GetArticleList(categoryId int, begin, end int) []*Article

		// 删除文章
		DeleteArticle(id int) error
	}

	//栏目
	ArticleCategory struct {
		//编号
		Id int `db:"id" pk:"yes" auto:"yes"`

		//父类编号,如为一级栏目则为0
		ParentId int `db:"parent_id"`

		//名称(唯一)
		Name string `db:"name"`

		// 别名
		Alias string `db:"alias"`

		// 页面标题
		Title string `db:"title"`

		// 关键字
		Keywords string `db:"keywords"`

		// 描述
		Description string `db:"describe"`

		// 排序编号
		SortNumber int `db:"sort_number"`

		// 定位路径（打开栏目页定位到的路径）
		Location string `db:"location"`

		UpdateTime int64 `db:"update_time"`
	}

	//文章
	Article struct {
		//编号
		Id int `db:"id" auto:"yes" pk:"yes"`

		// 栏目编号
		CategoryId int `db:"category_id"`

		//标题
		Title string `db:"title"`

		//小标题
		SmallTitle string `db:"small_title"`

		// 文章附图
		Thumbnail string `db:"thumbnail"`

		//资源地址
		Uri string `db:"-"`

		//重定向URL
		Location string `db:"location"`

		//作者
		PublisherId int `db:"publisher_id"`

		//文档内容
		Content string `db:"content"`

		//标签（关键词）
		Tags string `db:"tags"`

		//显示次数
		ViewCount int `db:"view_count"`

		//排序序号
		SortNumber int `db:"sort_number"`

		//创建时间
		CreateTime int64 `db:"create_time"`

		//最后修改时间
		UpdateTime int64 `db:"update_time"`
	}
)
