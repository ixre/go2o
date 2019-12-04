package api

import (
	"github.com/ixre/gof"
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

func (u utils) success(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, gof.Result{ErrCode: 0, ErrMsg: "success"})
}

func (u utils) error(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusOK, gof.Result{ErrCode: 1, ErrMsg: err.Error()})
}

func (u utils) errorJson(ctx echo.Context, errCode int, err error) error {
	return ctx.JSON(http.StatusOK, gof.Result{ErrCode: errCode, ErrMsg: err.Error()})
}
