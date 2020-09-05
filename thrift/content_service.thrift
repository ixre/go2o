namespace java com.github.jsix.go2o.rpc
namespace netstd com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.content_service
include "ttype.thrift"

/** 内容服务 */
service ContentService{
    /** 获取分页文章 */
    ttype.SPagingResult QueryPagingArticles(1:string cat,2:i32 begin_,3:i32 size)
    /** 获取置顶的文章 */
    list<SArticle> QueryTopArticles(1:string cat)
}



/** 文章 */
struct SArticle{
    /** 编号  */
    1:i32 Id
    /** 栏目编号 */
    2:i32 CatId
    /** 标题 */
    3:string Title
    /** 小标题 */
    4:string SmallTitle
    /** 文章附图 */
    5:string Thumbnail
    /** 重定向URL */
    6:i32 PublisherId
    /** 重定向URL */
    7:string Location
    /** 优先级,优先级越高，则置顶 */
    8:i32 Priority
    /** 浏览钥匙 */
    9:string AccessKey
    /** 文档内容 */
    10:string Content
    /** 标签（关键词） */
    11:string Tags
    /** 显示次数 */
    12:i32 ViewCount
    /** 排序序号 */
    13:i32 SortNum
    /** 创建时间 */
    14:i32 CreateTime
    /** 最后修改时间 */
    15:i32 UpdateTime
}