/**
 * Copyright 2015 @ S1N1 Team.
 * name : rest_server.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package app
import(
    "github.com/atnet/gof"
    "fmt"
    "strconv"
    "time"
    "go2o/src/core/service/goclient"
    "go2o/src/core/variable"
    "net/http"
    "github.com/atnet/gof/web"
    "go2o/src/app/api"

)

func RunRestApi(app gof.App, port int) {
    fmt.Println("[Started]:Api server running on port [" + strconv.Itoa(port) + "]:")

    //socket client
    time.Sleep(time.Second * 2) //等待启动Socket
    API_DOMAIN = app.Config().GetString(variable.ApiDomain)
    goclient.Configure("tcp", app.Config().GetString(variable.ClientSocketServer), app)

    var in *web.Interceptor= web.NewInterceptor(app,func(ctx *web.Context){
        host := ctx.Request.URL.Host
        // todo: path compare
        if host == API_DOMAIN{
            http.Error(ctx.ResponseWriter,"",http.StatusNotFound)
            return
        }
        //api.HandleApi(ctx)
        api.Handle(ctx)
    })

    //注册路由
    RegisterRoutes(app)

    //启动服务
    err := http.ListenAndServe(":"+strconv.Itoa(port), in)

    if err != nil {
        app.Log().Fatalln("ListenAndServer ", err)
    }
}