package api

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"github.com/labstack/echo"
	"go2o/core/service/auto_gen/rpc/ttype"
	"net/http"
)

type utils struct {
}

func (u utils) SResult(err error) *ttype.Result_ {
	if err != nil {
		return &ttype.Result_{ErrCode: 1, ErrMsg: err.Error()}
	}
	return &ttype.Result_{}
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
