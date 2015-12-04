# Go2o - 前端应用 #

## 目录说明 ##
master ->  管理中心
partner ->  商户管理中心
ucenter ->  会员中心
webui/mobi   -> 移动端前台
webui/weixin -> 微信开放平台
webui/www     ->  PC端应用

## 目录组成 ##
每个端均有控制器(*_c.go文件),路由(route.go)组成。
        #定义控制器
        type HomeController struct{}
        func (this *HomeController) Index(ctx *web.Context){
            //your code
        }
        #定义一个路由表
        var routes *mvc.Route = mvc.NewRoute(nil)
        routes.RegisterController("buy",new(HomeController))
        ## 通过路径 /Home/Index 访问到HomeController的Index方法

可通过定义筛选器，控制器实现mvc.Filter接口，来实现对请求的拦截。可
做会话验证等操作。

