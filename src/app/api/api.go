/**
 * Copyright 2015 @ S1N1 Team.
 * name : apI_portal
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api
import "github.com/atnet/gof/web"

func HandleApi(ctx *web.Context){
    //r := ctx.Request
    ctx.ResponseWriter.Write([]byte("It's working!"))
}