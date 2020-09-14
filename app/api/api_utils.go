package api

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"github.com/labstack/echo"
	"go2o/core/service/proto"
	"net/http"
)

type utils struct {
}

func (u utils) SResult(err error) *proto.Result {
	if err != nil {
		return &proto.Result{ErrCode: 1, ErrMsg: err.Error()}
	}
	return &proto.Result{}
}

func (u utils) JSON(ctx echo.Context, ret interface{}) error {
	return ctx.JSON(http.StatusOK, ret)
}


func (u utils) response(err error) *api.Response {
	if err != nil {
		return u.error(err)
	}
	return u.success(nil)
}

func (u utils) errorJson(ctx echo.Context, errCode int, err error) error {
	return ctx.JSON(http.StatusOK, gof.Result{ErrCode: errCode, ErrMsg: err.Error()})
}

func (u utils) success(data interface{}) *api.Response {
	return api.NewResponse(data)
}

func (u utils) error(err error) *api.Response {
	return u.errorWithCode(1, err)
}

func (u utils) errorWithCode(code int, err error) *api.Response {
	return api.ResponseWithCode(code, err.Error())
}
func (u utils) result(r *proto.Result)*api.Response{
	return api.ResponseWithCode(int(r.ErrCode),r.ErrMsg)
}
