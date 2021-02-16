package util

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/jwt-api"
	"github.com/labstack/echo/v4"
	"go2o/core/service/proto"
	"net/http"
)

type Utils struct {
}

func (u Utils) SResult(err error) *proto.Result {
	if err != nil {
		return &proto.Result{ErrCode: 1, ErrMsg: err.Error()}
	}
	return &proto.Result{}
}

func (u Utils) Response(err error) *api.Response {
	if err != nil {
		return u.Error(err)
	}
	return u.Success(nil)
}

func (u Utils) ErrorJson(ctx echo.Context, errCode int, err error) error {
	return ctx.JSON(http.StatusOK, gof.Result{ErrCode: errCode, ErrMsg: err.Error()})
}

func (u Utils) Success(data interface{}) *api.Response {
	return api.NewResponse(data)
}

func (u Utils) Error(err error) *api.Response {
	return u.ErrorWithCode(1, err)
}

func (u Utils) ErrorWithCode(code int, err error) *api.Response {
	return api.ResponseWithCode(code, err.Error())
}
func (u Utils) Result(r *proto.Result) *api.Response {
	return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
}
