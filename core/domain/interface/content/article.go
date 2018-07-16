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
		GetDomainId() int32
		// 获取值
		GetValue() Article
		// 设置值
		SetValue(*Article) error
		// 栏目
		Category() ICategory
		// 保存文章
		Save() (int32, error)
	}

	// 栏目
	ICategory interface {
		// 获取领域编号
		GetDomainId() int32
		// 获取值
		GetValue() ArticleCategory
		// 获取文章数量
		ArticleNum() int
		// 设置值
		SetValue(*ArticleCategory) error
		// 保存
		Save() (int32, error)
	}

	// 文章管理器
	IArticleManager interface {
		// 获取栏目
		GetCategory(id int32) ICategory
		// 根据标识获取文章栏目
		GetCategoryByAlias(alias string) ICategory
		// 获取所有的栏目
		GetAllCategory() []ICategory
		// 创建栏目
		CreateCategory(*ArticleCategory) ICategory
		// 删除栏目
		DelCategory(id int32) error
		// 创建文章
		CreateArticle(*Article) IArticle
		// 获取文章
		GetArticle(id int32) IArticle
		// 获取文章列表
		GetArticleList(categoryId int32, begin, end int) []*Article
		// 删除文章
		DeleteArticle(id int32) error
	}

	//栏目
	ArticleCategory struct {
		//编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		//父类编号,如为一级栏目则为0
		ParentId int32 `db:"parent_id"`
		// 浏览权限
		PermFlag int `db:"perm_flag"`
		// 名称(唯一)
		Name string `db:"name"`
		// 别名
		Alias string `db:"cat_alias"`
		// 排序编号
		SortNum int `db:"sort_num"`
		// 定位路径（打开栏目页定位到的路径）
		Location string `db:"location"`
		// 页面标题
		Title string `db:"title"`
		// 关键字
		Keywords string `db:"keywords"`
		// 描述
		Description string `db:"describe"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	//文章
	Article struct {
		// 编号
		ID int32 `db:"id" auto:"yes" pk:"yes"`
		// 栏目编号
		CatId int32 `db:"cat_id"`
		// 标题
		Title string `db:"title"`
		// 小标题
		SmallTitle string `db:"small_title"`
		// 文章附图
		Thumbnail string `db:"thumbnail"`
		// 重定向URL
		Location string `db:"location"`
		// 优先级,优先级越高，则置顶
		Priority int `db:"priority"`
		// 浏览钥匙
		AccessKey string `db:"access_key"`
		// 作者
		PublisherId int `db:"publisher_id"`
		// 文档内容
		Content string `db:"content"`
		// 标签（关键词）
		Tags string `db:"tags"`
		// 显示次数
		ViewCount int `db:"view_count"`
		// 排序序号
		SortNum int `db:"sort_num"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 最后修改时间
		UpdateTime int64 `db:"update_time"`
	}
)
